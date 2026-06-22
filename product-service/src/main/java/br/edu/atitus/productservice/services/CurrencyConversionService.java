package br.edu.atitus.productservice.services;

import br.edu.atitus.productservice.clients.CurrencyClient;
import br.edu.atitus.productservice.clients.CurrencyResponse;
import br.edu.atitus.productservice.dtos.ConversionResult;
import br.edu.atitus.productservice.entities.ProductEntity;
import org.springframework.stereotype.Service;

@Service
public class CurrencyConversionService {

    private static final String CACHE_NAME = "ConvertedValue";
    private static final Double FALLBACK_PRICE = -1.0;

    private final CurrencyClient currencyClient;
    private final CacheService cacheService;

    public CurrencyConversionService(CurrencyClient currencyClient, CacheService cacheService) {
        this.currencyClient = currencyClient;
        this.cacheService = cacheService;
    }

    public ConversionResult convert(ProductEntity product, String targetCurrency, String port) {
        String environment = "Product-service running on Port: " + port;

        if (isSameCurrency(product.getCurrency(), targetCurrency)) {
            return new ConversionResult(product.getPrice(), environment);
        }

        String keyCache = buildCacheKey(product.getCurrency(), targetCurrency);

        Double conversionRate = cacheService.get(CACHE_NAME, keyCache);

        if (conversionRate != null) {
            return new ConversionResult(
                    applyConversion(product.getPrice(), conversionRate),
                    environment + " | Currency in cache"
            );
        }

        CurrencyResponse currency = getCurrencyFromApi(product.getCurrency(), targetCurrency);
        System.out.println(currency);
        if (currency == null) {
            return new ConversionResult(
                    FALLBACK_PRICE,
                    environment + " | " + "Currency fallback"
            );
        }

        cacheService.put(CACHE_NAME, keyCache, currency.conversionRate());

        return new ConversionResult(
                applyConversion(product.getPrice(), currency.conversionRate()),
                environment + " | " + currency.environment()
        );
    }

    private CurrencyResponse getCurrencyFromApi(String currency, String targetCurrency) {
        return currencyClient.getCurrency(currency, targetCurrency);
    }

    private boolean isSameCurrency(String currency, String targetCurrency) {
        return targetCurrency.equalsIgnoreCase(currency);
    }

    private String buildCacheKey(String currency, String targetCurrency) {
        return currency + " - " + targetCurrency;
    }

    private Double applyConversion(Double price, Double conversionRate) {
        return price * conversionRate;
    }
}
