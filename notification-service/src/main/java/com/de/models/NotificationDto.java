package com.de.models;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Value;

@Value
@AllArgsConstructor
@Builder
public class NotificationDto {
    String disasterDescription;

    String userEventDescription;

    Float lon;

    Float lat;

    Long date;
}