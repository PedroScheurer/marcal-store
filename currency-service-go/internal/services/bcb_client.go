package services

import "context"

type BCBClient interface {
	GetConversionRate(ctx context.Context, target, currencyDate string) (float64, error)
}
