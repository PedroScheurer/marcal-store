package entities

type CurrencyEntity struct {
	ID             int64   `db:"id"`
	SourceCurrency string  `db:"source_currency"`
	TargetCurrency string  `db:"target_currency"`
	ConversionRate float64 `db:"conversion_rate"`
}
