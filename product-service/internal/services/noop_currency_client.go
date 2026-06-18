package services

// noopCurrencyClient é uma implementação provisória de CurrencyClient,
// usada apenas para o serviço compilar e rodar antes de implementarmos
// o client HTTP real (que vai chamar o currency-service, com descoberta
// via Eureka, retry e circuit breaker — próxima etapa da conversão).
//
// Ela sempre retorna (nil, nil), o que faz a CurrencyConversionService
// cair no caminho de fallback (mesmo comportamento do
// CurrencyClientFallbackFactory do Java quando a causa não é um
// FeignException.NotFound: `return null;`).
type noopCurrencyClient struct{}

// NewNoopCurrencyClient cria o client placeholder. Substituir por um
// client HTTP real assim que implementarmos a etapa de integração com
// o currency-service.
func NewNoopCurrencyClient() CurrencyClient {
	return &noopCurrencyClient{}
}

func (n *noopCurrencyClient) GetCurrency(source, target string) (*CurrencyResponse, error) {
	return nil, nil
}
