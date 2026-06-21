package clients

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/sony/gobreaker"
)

// HTTPCurrencyClient é a implementação real da interface CurrencyClient.
// Faz a ponte via HTTP com o microsserviço "currency-service".
type HTTPCurrencyClient struct {
	discovery  *ServiceDiscovery
	httpClient *http.Client
	cb         *gobreaker.CircuitBreaker
}

// NewHTTPCurrencyClient inicializa o cliente HTTP configurando o Circuit Breaker.
// A baseURL pode ser injetada via variável de ambiente ou resolvida dinamicamente (ex: Eureka).
func NewHTTPCurrencyClient(discovery *ServiceDiscovery, timeout time.Duration) *HTTPCurrencyClient {
	// Configuração do Circuit Breaker equivalente ao Resilience4j do Java
	cbSettings := gobreaker.Settings{
		Name:        "CurrencyClientCircuitBreaker",
		MaxRequests: 3,
		Interval:    10 * time.Second,
		Timeout:     15 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 5 && failureRatio > 0.5
		},
		IsSuccessful: func(err error) bool {
			// 404 é erro de negócio (moeda/par não suportado), não falha de
			// infraestrutura — não deve contar para abrir o circuito.
			if errors.Is(err, ErrCurrencyNotFound) {
				return true
			}
			return err == nil
		},
	}

	return &HTTPCurrencyClient{
		discovery: discovery,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		cb: gobreaker.NewCircuitBreaker(cbSettings),
	}
}

// GetCurrency executa a chamada HTTP protegida por uma política de Retry e Circuit Breaker.
func (c *HTTPCurrencyClient) GetCurrency(ctx context.Context, source, target string) (*CurrencyResponse, error) {
	var response *CurrencyResponse

	// 1. Aplica a política de Retry (Equivalente ao @Retry(name = "Retry_CurrencyClient_getCurrency"))
	err := retry.Do(
		func() error {
			// 2. Aplica o Circuit Breaker envolvendo a chamada de rede real
			res, cbErr := c.cb.Execute(func() (interface{}, error) {
				return c.doHTTPRequest(ctx, source, target)
			})

			if cbErr != nil {
				return cbErr
			}

			response = res.(*CurrencyResponse)
			return nil
		},
		retry.Context(ctx),
		retry.Attempts(3),                   // 3 tentativas no total
		retry.Delay(100*time.Millisecond),   // Delay inicial entre retries
		retry.DelayType(retry.BackOffDelay), // Backoff exponencial para poupar o servidor destino
		retry.LastErrorOnly(true),
		retry.RetryIf(func(err error) bool {
			// Erro de negócio (moeda/par não suportado) não se resolve
			// tentando de novo — só erros de infraestrutura justificam retry.
			return !errors.Is(err, ErrCurrencyNotFound)
		}),
	)

	if err != nil {
		// --- Lógica de Fallback Alinhada ao seu Contrato ---
		// Se o circuito estiver Aberto (gobreaker.ErrOpenState), o Javadoc diz que devemos "engolir"
		// o erro retornando (nil, nil) para ativar o fluxo de fallback da service.
		if errors.Is(err, gobreaker.ErrOpenState) {
			return nil, nil
		}

		// Se for um erro real ou falhas consecutivas exauridas, propagamos a falha
		return nil, fmt.Errorf("currency client failure after retries: %w", err)
	}

	return response, nil
}

// doHTTPRequest realiza o envio do pacote HTTP de fato e decodifica o payload JSON.
func (c *HTTPCurrencyClient) doHTTPRequest(ctx context.Context, source, target string) (*CurrencyResponse, error) {
	// 1. Constrói e sanitiza a URL com Query Parameters (?source=USD&target=BRL)
	baseURL, err := c.discovery.ResolveURL(ctx, "currency-service")
	if err != nil {
		return nil, fmt.Errorf("resolve currency-service: %w", err)
	}

	endpoint, err := url.Parse(fmt.Sprintf("%s/currency/convert", baseURL))
	if err != nil {
		return nil, fmt.Errorf("invalid base url: %w", err)
	}

	queryParams := url.Values{}
	queryParams.Add("source", source)
	queryParams.Add("target", target)
	endpoint.RawQuery = queryParams.Encode()

	// 2. Cria a requisição HTTP GET
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 3. Dispara a chamada
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err // Falhas de I/O de rede (Timeout, Connection Refused) caem aqui e disparam o Retry
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Printf("warning: failed to close response body: %v", closeErr)
		}
	}()

	// 4. Validação do Status Code de Resposta
	var ErrCurrencyNotFound = errors.New("currency mapping not found")
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("%w: %s to %s", ErrCurrencyNotFound, source, target)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("external currency service returned status: %d", resp.StatusCode)
	}

	// 5. Unmarshal do JSON
	var currencyResp CurrencyResponse
	if err := json.NewDecoder(resp.Body).Decode(&currencyResp); err != nil {
		return nil, fmt.Errorf("failed to decode currency body: %w", err)
	}

	return &currencyResp, nil
}
