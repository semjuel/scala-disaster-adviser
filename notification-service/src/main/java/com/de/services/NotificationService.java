package com.de.services;

import com.de.models.DisasterDto;
import com.de.models.DisasterResponse;
import com.de.models.NotificationDto;
import com.de.models.UserResponse;
import com.de.models.kafka.DisasterEventDto;
import com.de.models.kafka.UserEventDto;
import org.apache.commons.collections4.CollectionUtils;
import org.springframework.util.StringUtils;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;

import java.util.Collections;
import java.util.Objects;
import java.util.Set;
import java.util.stream.Collectors;

public class NotificationService {

    private final DisasterEventsService disasterEventsService;
    private final UserEventsService userEventsService;
    private final float coordinatesGap;
    private final long timestampGap;

    public NotificationService(DisasterEventsService disasterEventsService,
                               UserEventsService userEventsService,
                               Float coordinatesGap,
                               Long timestampGap) {
        this.disasterEventsService = Objects.requireNonNull(disasterEventsService);
        this.userEventsService = Objects.requireNonNull(userEventsService);
        this.coordinatesGap = Objects.requireNonNull(coordinatesGap);
        this.timestampGap = Objects.requireNonNull(timestampGap);
    }

    public Flux<NotificationDto> getNotifications(String userName) {
        return Flux.merge(streamOnUserEventNotifications(userName), streamOnDisasterEventNotifications(userName));
    }

    private Flux<NotificationDto> streamOnUserEventNotifications(String userName) {
        return userEventsService.streamUserEvents(userName)
                .flatMap(userEventDto -> disasterEventsService.getDisasterEvents(userEventDto.getLat(),
                        userEventDto.getLon(), coordinatesGap, userEventDto.getDate(), timestampGap)
                        .map(disasterResponse -> makeOnUserEventNotification(disasterResponse, userEventDto)));
    }

    private Flux<NotificationDto> streamOnDisasterEventNotifications(String userName) {
        return disasterEventsService.streamDisasterEvents(userName)
                .flatMap(disasterEventDto -> getUserEvents(userName, disasterEventDto)
                        .map(userEventResponse -> makeOnDisasterEventNotification(userEventResponse, disasterEventDto)))
                .flatMap(Flux::fromIterable);
    }

    private Mono<UserResponse> getUserEvents(String userName, DisasterEventDto disasterEventDto) {
        final DisasterDto disaster = disasterEventDto.getDisaster();
        return userEventsService.getUserEvents(userName, disaster.getLat(),
                disaster.getLon(), coordinatesGap, disaster.getDate(), timestampGap);
    }

    private NotificationDto makeOnUserEventNotification(DisasterResponse disasterResponse, UserEventDto userEventDto) {
        return CollectionUtils.isNotEmpty(disasterResponse.getDisasters()) || disasterResponse.getIsHot()
                ? NotificationDto.builder()
                .disasterDescription(makeDisasterDescription(disasterResponse))
                .userEventDescription(userEventDto.getDescription())
                .lat(userEventDto.getLat())
                .lon(userEventDto.getLon())
                .date(userEventDto.getDate())
                .build()
                : null;
    }

    private String makeDisasterDescription(DisasterResponse disasterResponse) {
        String disasters = disasterResponse.getDisasters()
                .stream()
                .map(DisasterDto::getDescription)
                .collect(Collectors.joining(","));
        disasters = StringUtils.isEmpty(disasters)
                ? "Next disaster probably could be near your event location: "
                : disasters;
        return disasterResponse.getIsHot() != null && disasterResponse.getIsHot()
                ? disasters + " Location often has disasters withing you event time"
                : disasters;
    }

    private Set<NotificationDto> makeOnDisasterEventNotification(UserResponse userResponse,
                                                                 DisasterEventDto disasterEvent) {
        return CollectionUtils.isNotEmpty(userResponse.getEvents())
                ? userResponse.getEvents().stream()
                .map(userEventDto -> NotificationDto.builder()
                        .lon(userEventDto.getLon())
                        .lat(userEventDto.getLat())
                        .disasterDescription(disasterEvent.getDisaster().getDescription())
                        .userEventDescription(userEventDto.getDescription())
                        .date(userEventDto.getDate())
                        .build())
                .collect(Collectors.toSet())
                : Collections.emptySet();
    }
}
