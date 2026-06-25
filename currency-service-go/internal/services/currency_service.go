package services

import (
	"context"
	"strings"

	"github.com/PedroScheurer/currency-service/internal/dtos"
	"github.com/PedroScheurer/currency-service/internal/repositories"
)

type CurrencyService struct {
	repository       repositories.CurrencyRepository
	bcbClientService *BCBClientService
	cacheService     *CacheService
	port             string
}

func NewCurrencyService(
	repository repositories.CurrencyRepository,
	bcbClientService *BCBClientService,
	cacheService *CacheService,
	port string,
) *CurrencyService {
	return &CurrencyService{
		repository:       repository,
		bcbClientService: bcbClientService,
		cacheService:     cacheService,
		port:             port,
	}
}

func (s *CurrencyService) findBySourceAndTarget(ctx context.Context, source string, target string) (*dtos.CurrencyDTO, error) {
	source = strings.ToUpper(source)
	target = strings.ToUpper(target)
	environment := `Currency-service running on Port:` + s.port

	if source == target {
		return &dtos.CurrencyDTO{
			Source:         source,
			Target:         target,
			ConversionRate: 1.0,
			Environment:    environment,
		}, nil
	}
	
}
