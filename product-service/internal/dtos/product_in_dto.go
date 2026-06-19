package dtos

// ProductInDTO é o payload de entrada para criação de um produto (POST).
type ProductInDTO struct {
	Name        string  `json:"name"`
	Instructor  string  `json:"instructor"`
	ImageURL    string  `json:"imageUrl"`
	VideoURL    string  `json:"videoUrl"`
	Description string  `json:"description"`
	Workload    int     `json:"workload"`
	Modules     int     `json:"modules"`
	Price       float64 `json:"price"`
	Currency    string  `json:"currency"`
}
