package com.de.models;

import com.de.models.kafka.UserEventDto;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.util.List;

@Data
@NoArgsConstructor
public class UserResponse {

    List<UserEventDto> events;
}
