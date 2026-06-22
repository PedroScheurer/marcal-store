package br.edu.atitus.currencyservice.services;

import br.edu.atitus.currencyservice.dtos.CurrencyDTO;
import br.edu.atitus.currencyservice.entities.CurrencyEntity;
import br.edu.atitus.currencyservice.infrastructure.exceptions.CurrencyNotFoundException;
import br.edu.atitus.currencyservice.repositories.CurrencyRepository;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;

import java.math.BigDecimal;
import java.math.RoundingMode;

@Service
public class CurrencyService {

    @Value("${server.port}")
    private String port;
    private static final String CACHE_NAME = "ConversionRateValue";

    private final CurrencyRepository repository;
    private final BCBClientService bcbClientService;
    private final CacheService cacheService;

    public CurrencyService(CurrencyRepository repository, BCBClientService bcbClientService, CacheService cacheService) {
        this.repository = repository;
        this.bcbClientService = bcbClientService;
        this.cacheService = cacheService;
    }

    public CurrencyDTO findBySourceAndTarget(String source, String target) {
        source = source.toUpperCase();
        target = target.toUpperCase();
        String environment = "Currency-service running on Port: " + port;

        if (source.equals(target)) {
            return new CurrencyDTO(source, target, BigDecimal.ONE, environment);
        }

        String keyCache = buildCacheKey(source, target);
        BigDecimal conversionRateInCache = cacheService.get(CACHE_NAME, keyCache);

        if(conversionRateInCache != null){
            return new CurrencyDTO(source,target,conversionRateInCache,
                    "Currency-service in cache");
        }

        BigDecimal sourceRate = getRate(source);
        BigDecimal targetRate = getRate(target);

        if(sourceRate == null && targetRate == null){
            return getCurrencyToFallback(source, target);
        }

        BigDecimal conversionRate = sourceRate.divide(targetRate, 2, RoundingMode.HALF_UP);

        cacheService.put(CACHE_NAME, keyCache, conversionRate);

        return new CurrencyDTO(
                source,
                target,
                conversionRate,
                "Currency-service running on Port: " + port + " | Banco Central do Brasil");
    }

    private BigDecimal getRate(String currency) {
        if (!isCurrencyBRL(currency)) {
            CurrencyDTO currencyFromBCB = bcbClientService.getCurrency(currency, "05-21-2026", port);

            if(currencyFromBCB != null){
                return currencyFromBCB.conversionRate();
            }

            return null;
        }
        return BigDecimal.ONE;
    }

    private boolean isCurrencyBRL(String currency) {
        return currency.equals("BRL");
    }

    private CurrencyDTO getCurrencyToFallback(String source, String target) {
        CurrencyEntity currency = repository.findBySourceCurrencyAndTargetCurrency(source, target)
                .orElseThrow(() -> new CurrencyNotFoundException("Currency not found."));

        return new CurrencyDTO(
                currency.getSourceCurrency(),
                currency.getTargetCurrency(),
                currency.getConversionRate(),
                "Currency-service fallback running on Port: " + port
        );
    }

    private String buildCacheKey(String source, String target) {
        return source + " - " + target;
    }
}
