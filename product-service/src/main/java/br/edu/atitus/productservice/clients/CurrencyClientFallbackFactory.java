package br.edu.atitus.productservice.clients;

import br.edu.atitus.productservice.infrastructure.exceptions.ExternalServiceException;
import feign.FeignException;
import org.springframework.cloud.openfeign.FallbackFactory;
import org.springframework.stereotype.Component;

@Component
public class CurrencyClientFallbackFactory
        implements FallbackFactory<CurrencyClient> {

    @Override
    public CurrencyClient create(Throwable cause) {

        return (currency, targetCurrency) -> {

            if (cause instanceof FeignException.NotFound) {

                throw new ExternalServiceException(
                        "Currency not found"
                );
            }

            return null;
        };
    }
}