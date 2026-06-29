package br.edu.atitus.order_service.clients;

import org.springframework.cloud.openfeign.FeignClient;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestParam;

@FeignClient(name = "product-service")
public interface ProductClient {

    @GetMapping("/products/noconverter/{id}")
    ProductResponse getProductById(@PathVariable("id") Long id);

    @GetMapping("/products/{id}")
    ProductResponse getProductByIdWithCurrency(
            @PathVariable("id") Long id,
            @RequestParam("targetCurrency") String targetCurrency);
}
