package services

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/PedroScheurer/product-service/internal/clients"
	"github.com/PedroScheurer/product-service/internal/entities"
)

// fakeCurrencyClient é um CurrencyClient de teste: nunca faz chamada HTTP
// de verdade, permitindo testar a lógica de CurrencyConversionService
// (cache, fallback, mesma moeda) isoladamente, sem depender de rede,
// Eureka ou currency-service reais.
type fakeCurrencyClient struct {
	rate float64
	err  error
}

func (f *fakeCurrencyClient) GetCurrency(ctx context.Context, source, target string) (*clients.CurrencyResponse, error) {
	if f.err != nil {
		return nil, f.err
	}
	return &clients.CurrencyResponse{
		SourceCurrency: source,
		TargetCurrency: target,
		ConversionRate: f.rate,
		Environment:    "currency-service running on Port: 9091",
	}, nil
}

func TestConvert_SameCurrency_ReturnsOriginalPrice(t *testing.T) {
	cache := NewCacheService(500, 15*time.Second)
	client := &fakeCurrencyClient{}
	svc := NewCurrencyConversionService(client, cache)

	product := &entities.ProductEntity{Price: 100, Currency: "BRL"}

	result := svc.Convert(context.Background(), product, "brl", "8082")

	if result.ConvertedPrice != 100 {
		t.Fatalf("expected 100, got %v", result.ConvertedPrice)
	}
}

func TestConvert_DifferentCurrency_ClientError_ReturnsFallback(t *testing.T) {
	cache := NewCacheService(500, 15*time.Second)
	client := &fakeCurrencyClient{err: errClientUnavailable}
	svc := NewCurrencyConversionService(client, cache)

	product := &entities.ProductEntity{Price: 100, Currency: "BRL"}

	result := svc.Convert(context.Background(), product, "USD", "8082")

	if result.ConvertedPrice != fallbackPrice {
		t.Fatalf("expected fallback price %v, got %v", fallbackPrice, result.ConvertedPrice)
	}
	if !strings.Contains(result.Environment, "fallback") {
		t.Fatalf("expected environment to mention fallback, got %q", result.Environment)
	}
}

func TestConvert_DifferentCurrency_UsesClientThenCache(t *testing.T) {
	cache := NewCacheService(500, 15*time.Second)
	client := &fakeCurrencyClient{rate: 5}
	svc := NewCurrencyConversionService(client, cache)

	product := &entities.ProductEntity{Price: 10, Currency: "BRL"}

	first := svc.Convert(context.Background(), product, "USD", "8082")
	if first.ConvertedPrice != 50 {
		t.Fatalf("expected 50, got %v", first.ConvertedPrice)
	}

	// Segunda chamada deve usar o cache (o client não seria mais necessário,
	// mas como o fake sempre retorna o mesmo rate, validamos pela mensagem
	// de environment, que deve indicar "Currency in cache").
	second := svc.Convert(context.Background(), product, "USD", "8082")
	if second.ConvertedPrice != 50 {
		t.Fatalf("expected 50 from cache, got %v", second.ConvertedPrice)
	}
	if !strings.Contains(second.Environment, "Currency in cache") {
		t.Fatalf("expected environment to mention cache, got %q", second.Environment)
	}
}

// errClientUnavailable simula uma falha de chamada ao currency-service
// (ex.: timeout, conexão recusada, circuit breaker aberto).
var errClientUnavailable = &fakeClientError{"currency-service unavailable"}

type fakeClientError struct{ msg string }

func (e *fakeClientError) Error() string { return e.msg }
