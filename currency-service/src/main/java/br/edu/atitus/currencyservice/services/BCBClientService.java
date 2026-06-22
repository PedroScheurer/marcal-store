package br.edu.atitus.currencyservice.services;

import br.edu.atitus.currencyservice.clients.BCBClient;
import br.edu.atitus.currencyservice.clients.BCBResponse;
import br.edu.atitus.currencyservice.dtos.CurrencyDTO;
import org.springframework.stereotype.Service;

import java.math.BigDecimal;
import java.util.HashMap;
import java.util.Map;

@Service
public class BCBClientService {

    private final BCBClient bcbClient;

    public BCBClientService(BCBClient bcbClient) {
        this.bcbClient = bcbClient;
    }

    public CurrencyDTO getCurrency(String moeda, String dataCotacao, String port){
        String environment = "Currency-service running on Port: " + port;

        BCBResponse bcbResponse = getConversionRateFromBCBApi(moeda,dataCotacao);
        System.out.println(bcbResponse);
        if (isValueInvalid(bcbResponse)) {
            return null;
        }

        BigDecimal conversionRate = BigDecimal.valueOf(
                bcbResponse.value().getLast().cotacaoCompra()
        );

        return new CurrencyDTO(
                moeda,
                "BRL",
                conversionRate,
                environment);
    }

    private boolean isValueInvalid(BCBResponse bcbResponse) {
        return bcbResponse == null || bcbResponse.value() == null || bcbResponse.value().isEmpty();
    }

    private BCBResponse getConversionRateFromBCBApi(String moeda, String dataCotacao) {
        return bcbClient.getConversionRate(moeda ,dataCotacao);
    }
}
