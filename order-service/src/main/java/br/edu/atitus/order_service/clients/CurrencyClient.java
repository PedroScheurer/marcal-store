package br.edu.atitus.order_service.clients;

import org.springframework.cloud.openfeign.FeignClient;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestParam;

@FeignClient(name = "currency-service")
public interface CurrencyClient {

    @GetMapping("/currency/convert")
    CurrencyResponse getCurrency(
            @RequestParam("source") String source,
            @RequestParam("target") String target);
}
