package com.de.configuration;

import com.de.services.DisasterEventsService;
import com.de.services.NotificationService;
import com.de.services.UserEventsService;
import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.MapperFeature;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.PropertyNamingStrategy;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.web.reactive.function.client.WebClient;

@Configuration
public class WebConfiguration {

    @Bean
    ObjectMapper jacksonMapper() {
        final ObjectMapper mapper = new ObjectMapper();
        mapper.configure(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES, false);
        mapper.configure(MapperFeature.DEFAULT_VIEW_INCLUSION, true);
        mapper.setPropertyNamingStrategy(PropertyNamingStrategy.SNAKE_CASE);

        return mapper;
    }


    @Bean
    WebClient webClient() {
        return WebClient.create();
    }

    @Bean
    UserEventsService userEventsService(@Value("${kafka.user-events.topic}") String topic,
                                        @Value("${kafka.bootstrap-servers}") String bootstrapServers,
                                        @Value("${user-service.endpoint}") String userServiceEndpoint,
                                        WebClient webClient,
                                        ObjectMapper jacksonMapper) {
        return new UserEventsService(
                topic,
                bootstrapServers,
                userServiceEndpoint,
                webClient,
                jacksonMapper);
    }

    @Bean
    DisasterEventsService disasterEventsService(
            @Value("${kafka.disaster-events.topic}") String topic,
            @Value("${kafka.bootstrap-servers}") String bootstrapServers,
            @Value("${disaster-service.endpoint}") String disasterServiceEndpoint,
            WebClient webClient,
            ObjectMapper jacksonMapper) {

        return new DisasterEventsService(
                topic,
                bootstrapServers,
                disasterServiceEndpoint,
                webClient,
                jacksonMapper);
    }

    @Bean
    NotificationService notificationService(DisasterEventsService disasterEventsService,
                                            UserEventsService userEventsService,
                                            @Value("${notification-service.coordinates-gap}") Float coordinatesGap,
                                            @Value("${notification-service.timestamp-gap}") Long timestampGap) {
        return new NotificationService(disasterEventsService, userEventsService, coordinatesGap, timestampGap);
    }
}
