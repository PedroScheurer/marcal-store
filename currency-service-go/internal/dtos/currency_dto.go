package dtos

type CurrencyDTO struct {
	Source         string
	Target         string
	ConversionRate float64
	Environment    string
}
