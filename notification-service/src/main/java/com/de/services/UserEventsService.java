package com.de.services;

import com.de.exceptions.DecodeException;
import com.de.exceptions.EncodeException;
import com.de.exceptions.InternalServerException;
import com.de.models.DisasterResponse;
import com.de.models.EventDateDto;
import com.de.models.EventLocationDto;
import com.de.models.UserRequest;
import com.de.models.UserResponse;
import com.de.models.kafka.UserEventDto;
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

public class UserEventsService {

    private static final Logger LOGGER = LoggerFactory.getLogger(UserEventsService.class);
    private static final int TIMEOUT = 5000;
    private static final Integer COMMIT_INTERVAL_MS = 1000;

    private static final String GROUP_ID_PATTERN = "Group %s";
    private static final String CLIENT_ID = "notification-user-service";

    private final String topic;
    private final String bootstrapServers;
    private final String userServiceEndpoint;
    private final WebClient webClient;
    private final ObjectMapper objectMapper;
    private final Scheduler scheduler;

    public UserEventsService(String topic,
                             String bootstrapServers,
                             String userServiceEndpoint,
                             WebClient webClient,
                             ObjectMapper objectMapper) {
        this.topic = Objects.requireNonNull(topic);
        this.bootstrapServers = Objects.requireNonNull(bootstrapServers);
        this.userServiceEndpoint = Objects.requireNonNull(userServiceEndpoint);
        this.webClient = Objects.requireNonNull(webClient);
        this.objectMapper = Objects.requireNonNull(objectMapper);
        this.scheduler = Schedulers.newSingle("single-scheduler", true);
    }

    public Flux<UserEventDto> streamUserEvents(String username) {
        return KafkaReceiver.create(KafkaUtils.receiverOptions(bootstrapServers, Collections.singletonList(topic),
                CLIENT_ID, String.format(GROUP_ID_PATTERN, username), Duration.ofMillis(COMMIT_INTERVAL_MS)))
                .receive()
                .publishOn(scheduler)
                .map(ConsumerRecord::value)
                .map(this::parseUserEvent)
                .onErrorContinue(this::logError)
                .filter(userEventDto -> userEventDto.getName().equals(username))
                .filter(this::validateUserEventDto);
    }

    public Mono<UserResponse> getUserEvents(String username, float lat, float lon, float coordinateGap,
                                               long timestamp, long timeGap) {
        return webClient.post()
                .uri(userServiceEndpoint)
                .body(Mono.just(prepareUserRequest(username, lat, lon, coordinateGap, timestamp, timeGap)),
                        String.class)
                .exchange()
                .timeout(Duration.ofMillis(TIMEOUT),
                        Mono.just(ClientResponse.create(HttpStatus.REQUEST_TIMEOUT).build()))
                .flatMap(this::handleUserServiceResponse);
    }

    private String prepareUserRequest(String username, float lat, float lon, float coordinateGap, long timestamp,
                                      long timeGap) {
        final UserRequest userRequest = UserRequest.builder()
                .username(username)
                .location(EventLocationDto.builder().lat(lat).lon(lon).range(coordinateGap).build())
                .date(EventDateDto.builder().timestamp(timestamp).range(timeGap).build()).build();
        try {
            return objectMapper.writeValueAsString(userRequest);
        } catch (JsonProcessingException e) {
            throw new EncodeException(String.format("Failed to encode user request with a reason: %s",
                    e.getMessage()));
        }
    }

    private Mono<UserResponse> handleUserServiceResponse(ClientResponse clientResponse) {
        final HttpStatus responseStatus = clientResponse.statusCode();
        return responseStatus == HttpStatus.OK
                ? clientResponse.bodyToMono(UserResponse.class)
                : Mono.error(new InternalServerException(String.format("Failed to request user service for details,"
                + " response status is : %s", responseStatus)));
    }

    private UserEventDto parseUserEvent(String rowUserEvent) {
        try {
            return objectMapper.readValue(rowUserEvent, UserEventDto.class);
        } catch (JsonProcessingException e) {
            throw new DecodeException(String.format("Failed to decode user event %s with a reason: %s", rowUserEvent,
                    e.getMessage()));
        }
    }

    private boolean validateUserEventDto(UserEventDto userEvent) {
        final Float lat = userEvent.getLat();
        if (lat == null || lat < 0) {
            return false;
        }

        final Float lon = userEvent.getLon();
        if (lon == null || lon < 0) {
            return false;
        }

        final Long date = userEvent.getDate();
        return date != null && date >= 0;
    }

    private Publisher<? extends UserEventDto> logError(Throwable throwable, Object o) {
        LOGGER.warn("User events consumer failed with a reason: {} ", throwable.getMessage());
        return Mono.empty();
    }
}
