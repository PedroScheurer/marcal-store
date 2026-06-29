package br.edu.atitus.order_service.clients;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;

@JsonIgnoreProperties(ignoreUnknown = true)
public record ProductResponse(
        Long id,
        String name,
        String instructor,
        String description,
        String imageUrl,
        String videoUrl,
        double price,
        String currency,
        String environment,
        Double convertedPrice
) {}
