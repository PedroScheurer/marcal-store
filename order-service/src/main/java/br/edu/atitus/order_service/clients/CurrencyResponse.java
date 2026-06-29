package br.edu.atitus.order_service.clients;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;

@JsonIgnoreProperties(ignoreUnknown = true)
public record CurrencyResponse(
        String sourceCurrency,
        String targetCurrency,
        double conversionRate,
        String environment
) {}
