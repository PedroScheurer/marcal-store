package br.edu.atitus.currencyservice.dtos;

import java.math.BigDecimal;

public record CurrencyDTO(
                          String sourceCurrency,
                          String targetCurrency,
                          BigDecimal conversionRate,
                          String environment) {
}
