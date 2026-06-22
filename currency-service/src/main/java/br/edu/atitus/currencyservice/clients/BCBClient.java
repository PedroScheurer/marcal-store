package br.edu.atitus.currencyservice.clients;

import io.github.resilience4j.retry.annotation.Retry;
import org.springframework.cloud.openfeign.FeignClient;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;


@FeignClient(name = "bcbclient",
        url = "https://olinda.bcb.gov.br/olinda/servico/PTAX/versao/v1/odata",
        fallback = BCBClientFallback.class)
public interface BCBClient {

    @GetMapping("/CotacaoMoedaDia(moeda=@moeda,dataCotacao=@dataCotacao)?@moeda='{moeda}'&@dataCotacao='{dataCotacao}'&$format=json")
    @Retry(name = "Retry_BCBClient_getConversionRate")
    BCBResponse getConversionRate(@PathVariable String moeda,
                           @PathVariable String dataCotacao);

}
