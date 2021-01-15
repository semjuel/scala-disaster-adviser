package com.de.services;

import com.de.exceptions.DecodeException;
import com.de.exceptions.EncodeException;
import com.de.exceptions.InternalServerException;
import com.de.models.DisasterDto;
import com.de.models.DisasterRequest;
import com.de.models.DisasterResponse;
import com.de.models.EventDateDto;
import com.de.models.EventLocationDto;
import com.de.models.kafka.DisasterEventDto;
import com.de.utils.KafkaUtils;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.kafka.clients.consumer.ConsumerRecord;
import org.reactivestreams.Publisher;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.http.HttpStatus;
import org.springframework.web.reactive.function.client.ClientResponse;
import org.springframework.web.reactive.function.client.WebClient;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;
import reactor.core.scheduler.Scheduler;
import reactor.core.scheduler.Schedulers;
import reactor.kafka.receiver.KafkaReceiver;

import java.time.Duration;
import java.util.Collections;
import java.util.Objects;

public class DisasterEventsService {

    private static final Logger LOGGER = LoggerFactory.getLogger(DisasterEventsService.class);
    private static final int TIMEOUT = 5000;
    private static final int COMMIT_INTERVAL_MS = 1000;
    private static final String GROUP_ID_PATTERN = "Group %s";
    private static final String CLIENT_ID = "notification-disaster-service";

    private final String topic;
    private final String bootstrapServers;
    private final String disastersEndpoint;
    private final WebClient webClient;
    private final ObjectMapper objectMapper;
    private final Scheduler scheduler;

    public DisasterEventsService(String topic,
                                 String bootstrapServers,
                                 String disastersEndpoint,
                                 WebClient webClient,
                                 ObjectMapper objectMapper) {
        this.topic = Objects.requireNonNull(topic);
        this.bootstrapServers = Objects.requireNonNull(bootstrapServers);
        this.disastersEndpoint = Objects.requireNonNull(disastersEndpoint);
        this.webClient = Objects.requireNonNull(webClient);
        this.objectMapper = Objects.requireNonNull(objectMapper);
        this.scheduler = Schedulers.newSingle("single-scheduler");
    }

    public Flux<DisasterEventDto> streamDisasterEvents(String userId) {
        return KafkaReceiver.create(KafkaUtils.receiverOptions(bootstrapServers, Collections.singletonList(topic),
                CLIENT_ID, String.format(GROUP_ID_PATTERN, userId), Duration.ofMillis(COMMIT_INTERVAL_MS)))
                .receive()
                .publishOn(scheduler)
                .map(ConsumerRecord::value)
                .map(this::parseDisasterEvent)
                .onErrorResume(this::logError)
                .filter(this::validateDisasterEventDto);
    }

    public Mono<DisasterResponse> getDisasterEvents(float lat, float lon, float coordinateGap, long timestamp,
                                                    long timeGap) {
        return webClient.post()
                .uri(disastersEndpoint)
                .body(Mono.just(prepareDisasterRequest(lat, lon, coordinateGap, timestamp, timeGap)),
                        String.class)
                .exchange()
                .timeout(Duration.ofMillis(TIMEOUT),
                        Mono.just(ClientResponse.create(HttpStatus.REQUEST_TIMEOUT).build()))
                .flatMap(this::handleDisasterServiceResponse);
    }

    private String prepareDisasterRequest(float lat, float lon, float coordinateGap, long timestamp, long timeGap) {
        final DisasterRequest disasterRequest = DisasterRequest.builder()
                .location(EventLocationDto.builder().lat(lat).lon(lon).range(coordinateGap).build())
                .date(EventDateDto.builder().timestamp(timestamp).range(timeGap).build()).build();
        try {
            return objectMapper.writeValueAsString(disasterRequest);
        } catch (JsonProcessingException e) {
            throw new EncodeException(String.format("Failed to encode disaster request with a reason: %s",
                    e.getMessage()));
        }
    }

    private Mono<DisasterResponse> handleDisasterServiceResponse(ClientResponse clientResponse) {
        final HttpStatus responseStatus = clientResponse.statusCode();
        return responseStatus == HttpStatus.OK
                ? clientResponse.bodyToMono(DisasterResponse.class)
                : Mono.error(new InternalServerException(String.format("Failed to request disaster service for details,"
                + " response status is : %s", responseStatus)));
    }

    private DisasterEventDto parseDisasterEvent(String rowDisasterEvent) {
        try {
            return objectMapper.readValue(rowDisasterEvent, DisasterEventDto.class);
        } catch (JsonProcessingException e) {
            throw new DecodeException(String.format("Failed to decode user event %s with a reason: %s",
                    rowDisasterEvent, e.getMessage()));
        }
    }

    private boolean validateDisasterEventDto(DisasterEventDto disasterEventDto) {
        final DisasterDto disaster = disasterEventDto.getDisaster();
        if (disaster == null) {
            return false;
        }
        final Float lat = disaster.getLat();
        if (lat == null || lat < 0) {
            return false;
        }

        final Float lon = disaster.getLon();
        if (lon == null || lon < 0) {
            return false;
        }

        final Long date = disaster.getDate();
        return date != null && date >= 0;
    }

    private Publisher<? extends DisasterEventDto> logError(Throwable throwable) {
        LOGGER.warn("Disaster events consumer failed with a reason: {} ", throwable.getMessage());
        return Flux.empty();
    }
}
