package services

import (
	"strings"
	"testing"
	"time"

	"github.com/PedroScheurer/product-service/internal/entities"
)

func TestConvert_SameCurrency_ReturnsOriginalPrice(t *testing.T) {
	cache := NewCacheService(500, 15*time.Second)
	client := NewNoopCurrencyClient()
	svc := NewCurrencyConversionService(client, cache)

	product := &entities.ProductEntity{Price: 100, Currency: "BRL"}

	result := svc.Convert(product, "brl", "8082")

	if result.ConvertedPrice != 100 {
		t.Fatalf("expected 100, got %v", result.ConvertedPrice)
	}
}

func TestConvert_DifferentCurrency_NoopClient_ReturnsFallback(t *testing.T) {
	cache := NewCacheService(500, 15*time.Second)
	client := NewNoopCurrencyClient()
	svc := NewCurrencyConversionService(client, cache)

	product := &entities.ProductEntity{Price: 100, Currency: "BRL"}

	result := svc.Convert(product, "USD", "8082")

	if result.ConvertedPrice != fallbackPrice {
		t.Fatalf("expected fallback price %v, got %v", fallbackPrice, result.ConvertedPrice)
	}
}

type fakeCurrencyClient struct {
	rate float64
}

func (f *fakeCurrencyClient) GetCurrency(source, target string) (*CurrencyResponse, error) {
	return &CurrencyResponse{
		SourceCurrency: source,
		TargetCurrency: target,
		ConversionRate: f.rate,
		Environment:    "currency-service running on Port: 9091",
	}, nil
}

func TestConvert_DifferentCurrency_UsesClientThenCache(t *testing.T) {
	cache := NewCacheService(500, 15*time.Second)
	client := &fakeCurrencyClient{rate: 5}
	svc := NewCurrencyConversionService(client, cache)

	product := &entities.ProductEntity{Price: 10, Currency: "BRL"}

	first := svc.Convert(product, "USD", "8082")
	if first.ConvertedPrice != 50 {
		t.Fatalf("expected 50, got %v", first.ConvertedPrice)
	}

	// Segunda chamada deve usar o cache (o client não seria mais necessário,
	// mas como o fake sempre retorna o mesmo rate, validamos pela mensagem
	// de environment, que deve indicar "Currency in cache").
	second := svc.Convert(product, "USD", "8082")
	if second.ConvertedPrice != 50 {
		t.Fatalf("expected 50 from cache, got %v", second.ConvertedPrice)
	}
	if !strings.Contains(second.Environment, "Currency in cache") {
		t.Fatalf("expected environment to mention cache, got %q", second.Environment)
	}
}
