package br.edu.atitus.currencyservice.clients;

import org.springframework.stereotype.Component;

import java.util.Map;

@Component
public class BCBClientFallback implements BCBClient {
    @Override
    public BCBResponse getConversionRate(String moeda, String dataCotacao) {
        return null;
    }
}
