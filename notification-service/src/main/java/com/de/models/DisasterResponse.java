package com.de.models;

import lombok.Data;
import lombok.NoArgsConstructor;

import java.util.List;

@Data
@NoArgsConstructor
public class DisasterResponse {

    List<DisasterDto> disasters;

    Boolean isHot;
}
