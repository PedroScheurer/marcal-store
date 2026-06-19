package dtos

// ProductDTO é a resposta de leitura de um produto,
// incluindo os campos opcionais de conversão de moeda.
// ConvertedPrice, RequestedCurrency e Environment usam ponteiro
// para poderem ser null no JSON quando não há conversão.
type ProductDTO struct {
	ID                int64    `json:"id"`
	Name              string   `json:"name"`
	Instructor        string   `json:"instructor"`
	ImageURL          string   `json:"imageUrl"`
	VideoURL          string   `json:"videoUrl"`
	Description       string   `json:"description"`
	Workload          int      `json:"workload"`
	Modules           int      `json:"modules"`
	Price             float64  `json:"price"`
	Currency          string   `json:"currency"`
	ConvertedPrice    *float64 `json:"convertedPrice"`
	RequestedCurrency *string  `json:"requestedCurrency"`
	Environment       *string  `json:"environment"`
}
