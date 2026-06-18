package dtos

// ProductOutDTO é o equivalente Go do record Java ProductOutDTO,
// usado como resposta de criação/alteração de produto (sem dados de conversão).
type ProductOutDTO struct {
	ID          int64   `json:"id"`
	Description string  `json:"description"`
	Brand       string  `json:"brand"`
	Model       string  `json:"model"`
	Price       float64 `json:"price"`
	Currency    string  `json:"currency"`
	ImageUrl    string  `json:"imageUrl"`
	VideoUrl    string  `json:"videoUrl"`
}
