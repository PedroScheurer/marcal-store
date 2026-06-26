package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/PedroScheurer/currency-service/internal/entities"
	"github.com/jmoiron/sqlx"
)

type PostgresCurrencyRepository struct {
	db *sqlx.DB
}

func NewCurrencyRepository(db *sqlx.DB) CurrencyRepository {
	return &PostgresCurrencyRepository{db: db}
}

const selectColumns = `id, source_currency, target_currency, conversion_rate`

func (r *PostgresCurrencyRepository) FindByID(ctx context.Context, id int64) (*entities.CurrencyEntity, error) {
	var currency entities.CurrencyEntity

	query := `SELECT ` + selectColumns + ` FROM tb_currency WHERE id = $1`

	err := r.db.GetContext(ctx, &currency, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find currency by id %d: %w", id, err)
	}

	return &currency, nil
}

func (r *PostgresCurrencyRepository) FindBySourceCurrencyAndTargetCurrency(ctx context.Context, sourceCurrency string, targetCurrency string) (*entities.CurrencyEntity, error) {
	var currency entities.CurrencyEntity

	query := `SELECT ` + selectColumns + `FROM tb_currency WHERE source_currency = $1 AND target_currency = $2`

	err := r.db.GetContext(ctx, &currency, query, sourceCurrency, targetCurrency)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find currency by source %s and target %s: %w", sourceCurrency, targetCurrency, err)
	}
	return &currency, nil
}
