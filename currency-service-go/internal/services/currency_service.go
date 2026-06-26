package services

import (
	"context"
	"fmt"
	"time"

	"github.com/PedroScheurer/currency-service/internal/apperrors"
	"github.com/PedroScheurer/currency-service/internal/dtos"
	"github.com/PedroScheurer/currency-service/internal/repositories"
)

const cacheName = "ConversionRateValue"

type CurrencyService struct {
	repository repositories.CurrencyRepository
	bcbClient  BCBClient
	cache      *CacheService
	port       string
}

func NewCurrencyService(
	repository repositories.CurrencyRepository,
	bcbClient BCBClient,
	cache *CacheService,
	port string,
) *CurrencyService {
	return &CurrencyService{
		repository: repository,
		bcbClient:  bcbClient,
		cache:      cache,
		port:       port,
	}
}

func (s *CurrencyService) FindBySourceAndTarget(ctx context.Context, source, target string) (*dtos.CurrencyDTO, error) {
	source = toUpper(source)
	target = toUpper(target)

	environment := "Currency-service running on Port: " + s.port

	// Mesma moeda: taxa 1:1
	if source == target {
		return &dtos.CurrencyDTO{
			SourceCurrency: source,
			TargetCurrency: target,
			ConversionRate: 1.0,
			Environment:    environment,
		}, nil
	}

	// Verifica cache
	cacheKey := buildCacheKey(source, target)
	if rate, found := s.cache.Get(cacheName, cacheKey); found {
		return &dtos.CurrencyDTO{
			SourceCurrency: source,
			TargetCurrency: target,
			ConversionRate: rate,
			Environment:    "Currency-service in cache",
		}, nil
	}

	// Busca taxas no BCB usando data de hoje
	today := time.Now().Format("01-02-2006") // MM-DD-YYYY — formato da API do BCB

	sourceRate, err := s.getRate(ctx, source, today)
	if err != nil {
		return nil, fmt.Errorf("get source rate: %w", err)
	}

	targetRate, err := s.getRate(ctx, target, today)
	if err != nil {
		return nil, fmt.Errorf("get target rate: %w", err)
	}

	// Se ambas as taxas falharam no BCB, usa fallback do banco
	if sourceRate == 0 && targetRate == 0 {
		return s.fallback(ctx, source, target)
	}

	// Calcula a taxa de conversão: sourceRate / targetRate
	// (equivalente ao BigDecimal.divide do Java)
	conversionRate := round2(sourceRate / targetRate)

	s.cache.Put(cacheName, cacheKey, conversionRate)

	return &dtos.CurrencyDTO{
		SourceCurrency: source,
		TargetCurrency: target,
		ConversionRate: conversionRate,
		Environment:    environment + " | Banco Central do Brasil",
	}, nil
}

// getRate retorna a taxa da moeda em relação ao BRL.
// BRL retorna 1.0 diretamente (moeda base).
// Outras moedas consultam o BCB.
// Retorna 0 quando o BCB não tem cotação (sem erro — será tratado como fallback).
func (s *CurrencyService) getRate(ctx context.Context, currency, date string) (float64, error) {
	if currency == "BRL" {
		return 1.0, nil
	}

	rate, err := s.bcbClient.GetConversionRate(ctx,
		currency, date)
	if err != nil {
		return 0, err
	}

	return rate, nil
}

// fallback busca no banco quando o BCB não retorna cotação,
// equivalente ao getCurrencyToFallback do Java.
func (s *CurrencyService) fallback(ctx context.Context, source, target string) (*dtos.CurrencyDTO, error) {
	currency, err := s.repository.FindBySourceCurrencyAndTargetCurrency(ctx, source, target)
	if err != nil {
		return nil, fmt.Errorf("fallback query: %w", err)
	}
	if currency == nil {
		return nil, apperrors.NewCurrencyNotFoundError(
			fmt.Sprintf("Currency not found: %s -> %s", source, target),
		)
	}

	return &dtos.CurrencyDTO{
		SourceCurrency: currency.SourceCurrency,
		TargetCurrency: currency.TargetCurrency,
		ConversionRate: currency.ConversionRate,
		Environment:    "Currency-service fallback running on Port: " + s.port,
	}, nil
}

func buildCacheKey(source, target string) string {
	return source + " - " + target
}

func toUpper(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if 'a' <= c && c <= 'z' {
			c -= 'a' - 'A'
		}
		result[i] = c
	}
	return string(result)
}

func round2(f float64) float64 {
	// Arredonda pra 2 casas decimais, equivalente ao
	// BigDecimal.divide(targetRate, 2, RoundingMode.HALF_UP) do Java.
	return float64(int(f*100+0.5)) / 100
}
