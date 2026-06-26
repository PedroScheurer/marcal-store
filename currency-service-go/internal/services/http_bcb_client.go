package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/PedroScheurer/currency-service/internal/clients"
	"github.com/PedroScheurer/currency-service/internal/dtos"
	"github.com/avast/retry-go/v4"
	"github.com/sony/gobreaker"
)

const bcbBaseURL = "https://olinda.bcb.gov.br/olinda/servico/PTAX/versao/v1/odata"

var ErrBCBCurrencyNotFound = errors.New("bcb: no quote found for currency/date")

type HttpBCBClient struct {
	httpClient *http.Client
	cb         *gobreaker.CircuitBreaker
}

func NewHttpBCBClient(timeout time.Duration) *HttpBCBClient {
	cbSettings := getBCBSettings()

	return &HttpBCBClient{
		httpClient: &http.Client{
			Timeout: timeout,
		},
		cb: gobreaker.NewCircuitBreaker(cbSettings),
	}
}

func getBCBSettings() gobreaker.Settings {
	return gobreaker.Settings{
		Name:        "BCBClientCircuitBreaker",
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
			if errors.Is(err, clients.ErrCurrencyNotFound) {
				return true
			}
			return err == nil
		},
	}
}

func (c *HttpBCBClient) GetConversionRate(ctx context.Context, currency, date string) (float64, error) {
	var rate float64

	rate, err := c.getRateWithRetry(ctx, currency, date)

	if err != nil {
		if errors.Is(err, gobreaker.ErrOpenState) {
			return 0, nil
		}
		return 0, fmt.Errorf("bcb client failure after retries: %w", err)
	}

	return rate, nil
}

func (c *HttpBCBClient) getRateWithRetry(ctx context.Context, currency, date string) (float64, error) {
	var rate float64

	err := retry.Do(
		func() error {
			res, cbErr := c.cb.Execute(func() (interface{}, error) {
				return c.doRequest(currency, date)
			})
			if cbErr != nil {
				return cbErr
			}

			// Type assertion seguro
			v, ok := res.(float64)
			if !ok {
				return fmt.Errorf("unexpected response type: %T", res)
			}

			rate = v
			return nil
		},
		retry.Context(ctx), // Idealmente usamos o context que vem de cima em vez do Background()
		retry.Attempts(3),
		retry.Delay(500*time.Millisecond),
		retry.DelayType(retry.BackOffDelay),
		retry.LastErrorOnly(true),
		retry.RetryIf(func(err error) bool {
			return !errors.Is(err, ErrBCBCurrencyNotFound)
		}),
	)

	if err != nil {
		return 0, err
	}

	return rate, nil
}

func (c *HttpBCBClient) doRequest(currency, date string) (float64, error) {
	// Monta a URL exatamente igual ao @GetMapping do Feign Java:
	// /CotacaoMoedaDia(moeda=@moeda,dataCotacao=@dataCotacao)?@moeda='{moeda}'&@dataCotacao='{dataCotacao}'&$format=json
	path := fmt.Sprintf(
		"/CotacaoMoedaDia(moeda=@moeda,dataCotacao=@dataCotacao)?@moeda='%s'&@dataCotacao='%s'&$format=json",
		url.QueryEscape(currency),
		url.QueryEscape(date),
	)

	req, err := http.NewRequest(http.MethodGet, bcbBaseURL+path, nil)
	if err != nil {
		return 0, fmt.Errorf("create bcb request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("call bcb: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("bcb returned status %d", resp.StatusCode)
	}

	var bcbResp dtos.BCBResponse
	if err := json.NewDecoder(resp.Body).Decode(&bcbResp); err != nil {
		return 0, fmt.Errorf("decode bcb response: %w", err)
	}

	if len(bcbResp.Value) == 0 {
		return 0, ErrBCBCurrencyNotFound
	}

	return bcbResp.Value[0].CotacaoCompra, nil
}
