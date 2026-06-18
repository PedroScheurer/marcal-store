package dtos

// ProductDTO é o equivalente Go do record Java ProductDTO.
// ConvertedPrice e Environment usam ponteiro (*float64 / *string)
// para poderem ser `null` no JSON, assim como o Double/String
// do Java podem ser null (ex.: quando não há conversão de moeda).
type ProductDTO struct {
	ID                int64    `json:"id"`
	Description       string   `json:"description"`
	Brand             string   `json:"brand"`
	Model             string   `json:"model"`
	Price             float64  `json:"price"`
	Currency          string   `json:"currency"`
	ImageUrl          string   `json:"imageUrl"`
	VideoUrl          string   `json:"videoUrl"`
	Environment       *string  `json:"environment"`
	ConvertedPrice    *float64 `json:"convertedPrice"`
	RequestedCurrency *string  `json:"requestedCurrency"`
}
