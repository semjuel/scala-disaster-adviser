package com.de.models.kafka;

import com.de.models.DisasterDto;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@Builder
@AllArgsConstructor
@NoArgsConstructor
public class DisasterEventDto {

    DisasterDto disaster;
}
