package dtos

type CurrencyDTO struct {
	SourceCurrency string  `json:"sourceCurrency"`
	TargetCurrency string  `json:"targetCurrency"`
	ConversionRate float64 `json:"conversionRate"`
	Environment    string  `json:"environment"`
}
