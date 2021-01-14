package com.de.models;

import lombok.Builder;
import lombok.Data;

@Data
@Builder
public class DisasterRequest {

    EventDateDto date;

    EventLocationDto location;
}
