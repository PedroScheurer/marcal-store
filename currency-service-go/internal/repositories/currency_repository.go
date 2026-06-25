package repositories

import (
	"context"

	"github.com/PedroScheurer/currency-service/internal/entities"
)

type CurrencyRepository interface {
	FindByID(ctx context.Context, id int64) (*entities.CurrencyEntity, error)
	FindBySourceCurrencyAndTargetCurrency(ctx context.Context, sourceCurrency string, targetCurrency string) (*entities.CurrencyEntity, error)
}
