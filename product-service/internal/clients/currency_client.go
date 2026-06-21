package clients

import "context"

// CurrencyResponse é o equivalente Go do record Java CurrencyResponse,
// que representa o corpo de resposta do endpoint /currency/convert
// do currency-service.
type CurrencyResponse struct {
	SourceCurrency string  `json:"sourceCurrency"`
	TargetCurrency string  `json:"targetCurrency"`
	ConversionRate float64 `json:"conversionRate"`
	Environment    string  `json:"environment"`
}

// CurrencyClient é o equivalente Go da interface Feign CurrencyClient.
//
//	@FeignClient(name = "currency-service", path = "/currency", fallbackFactory = ...)
//	public interface CurrencyClient {
//	    @GetMapping("/convert")
//	    @Retry(name = "Retry_CurrencyClient_getCurrency")
//	    CurrencyResponse getCurrency(@RequestParam String source, @RequestParam String target);
//	}
//
// Por enquanto esta é só a interface (contrato) usada pela
// CurrencyConversionService — a implementação concreta, que fará a
// chamada HTTP de fato ao currency-service (incluindo descoberta via
// Eureka, retry e circuit breaker, equivalentes ao Resilience4j do
// application.yaml), será adicionada na próxima etapa.
//
// GetCurrency deve retornar (nil, nil) quando o fallback é acionado e
// decide "engolir" o erro (equivalente ao `return null;` dentro do
// fallback do Java), e (nil, error) quando deve propagar uma falha real
// (equivalente ao `throw new ExternalServiceException(...)`).
type CurrencyClient interface {
	GetCurrency(ctx context.Context, source, target string) (*CurrencyResponse, error)
}
