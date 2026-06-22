package br.edu.atitus.productservice.clients;

import io.github.resilience4j.retry.annotation.Retry;
import org.springframework.cloud.openfeign.FeignClient;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestParam;

@FeignClient(name = "currency-service", path = "/currency", fallbackFactory = CurrencyClientFallbackFactory.class)
public interface CurrencyClient {

    @GetMapping("/convert")
    @Retry(name = "Retry_CurrencyClient_getCurrency")
    CurrencyResponse getCurrency(@RequestParam String source, @RequestParam String target);
}
