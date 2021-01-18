package com.de.controllers;

import com.de.models.NotificationDto;
import com.de.services.DisasterEventsService;
import com.de.services.NotificationService;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.MediaType;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;
import reactor.core.publisher.Flux;

import java.util.Objects;

@RestController
public class NotificationController {

    private static final Logger LOGGER = LoggerFactory.getLogger(NotificationController.class);


    private final NotificationService notificationService;

    @Autowired
    NotificationController(NotificationService notificationService) {
        this.notificationService = Objects.requireNonNull(notificationService);
    }

    @GetMapping(path = "/disaster-notifications", produces = MediaType.TEXT_EVENT_STREAM_VALUE)
    public Flux<NotificationDto> streamDisasterNotifications(@RequestParam String userName) {
        return notificationService.getNotifications(userName)
                .onErrorContinue((throwable, failedEvent) -> LOGGER.error("Event notification failed for {} with a"
                        + " reason {}", failedEvent, throwable.getMessage()));
    }
}
