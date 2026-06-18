package dtos

// ProductInDTO é o equivalente Go do record Java ProductInDTO,
// usado como payload de entrada para criação de produto (POST).
type ProductInDTO struct {
	Description string  `json:"description"`
	Brand       string  `json:"brand"`
	Model       string  `json:"model"`
	Currency    string  `json:"currency"`
	Price       float64 `json:"price"`
	ImageUrl    string  `json:"imageUrl"`
	VideoUrl    string  `json:"videoUrl"`
}
