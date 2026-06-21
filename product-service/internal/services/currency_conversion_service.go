package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/PedroScheurer/product-service/internal/clients"
	"github.com/PedroScheurer/product-service/internal/dtos"
	"github.com/PedroScheurer/product-service/internal/entities"
)

const (
	convertedValueCacheName = "ConvertedValue"
	fallbackPrice           = -1.0
)

// CurrencyConversionService é o equivalente Go da classe Java
// CurrencyConversionService: decide se precisa converter a moeda do
// produto, consultando o cache antes de chamar o currency-service via
// CurrencyClient, e atualizando o cache com o resultado.
type CurrencyConversionService struct {
	currencyClient clients.CurrencyClient
	cacheService   *CacheService
}

func NewCurrencyConversionService(currencyClient clients.CurrencyClient, cacheService *CacheService) *CurrencyConversionService {
	return &CurrencyConversionService{
		currencyClient: currencyClient,
		cacheService:   cacheService,
	}
}

// Convert é o equivalente a CurrencyConversionService.convert(product, targetCurrency, port).
func (s *CurrencyConversionService) Convert(ctx context.Context, product *entities.ProductEntity, targetCurrency, port string) dtos.ConversionResult {
	environment := "Product-service running on Port: " + port

	if isSameCurrency(product.Currency, targetCurrency) {
		return dtos.ConversionResult{
			ConvertedPrice: product.Price,
			Environment:    environment,
		}
	}

	keyCache := buildCacheKey(product.Currency, targetCurrency)

	if conversionRate, found := s.cacheService.Get(convertedValueCacheName, keyCache); found {
		return dtos.ConversionResult{
			ConvertedPrice: applyConversion(product.Price, conversionRate),
			Environment:    environment + " | Currency in cache",
		}
	}

	currency, err := s.currencyClient.GetCurrency(ctx, product.Currency, targetCurrency)
	if err != nil || currency == nil {
		return dtos.ConversionResult{
			ConvertedPrice: fallbackPrice,
			Environment:    environment + " | Currency fallback",
		}
	}

	s.cacheService.Put(convertedValueCacheName, keyCache, currency.ConversionRate)

	return dtos.ConversionResult{
		ConvertedPrice: applyConversion(product.Price, currency.ConversionRate),
		Environment:    fmt.Sprintf("%s | %s", environment, currency.Environment),
	}
}

func isSameCurrency(currency, targetCurrency string) bool {
	// Equivalente a targetCurrency.equalsIgnoreCase(currency) no Java.
	return strings.EqualFold(currency, targetCurrency)
}

func buildCacheKey(currency, targetCurrency string) string {
	return currency + " - " + targetCurrency
}

func applyConversion(price, conversionRate float64) float64 {
	return price * conversionRate
}
