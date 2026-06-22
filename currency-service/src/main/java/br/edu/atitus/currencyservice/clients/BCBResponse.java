package br.edu.atitus.currencyservice.clients;

import br.edu.atitus.currencyservice.dtos.CurrencyQuoteResponse;
import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.List;

public record BCBResponse(

        @JsonProperty("@odata.context")
        String context,

        List<CurrencyQuoteResponse> value
) {
}