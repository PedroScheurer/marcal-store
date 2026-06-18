package dtos

// ProductUpdateInDTO representa o payload de entrada do PUT /ws/products/{idProduct}.
//
// No Java original, o WsProductController.putProduct reaproveita o record
// ProductDTO (de saída, com campos de conversão de moeda) como tipo do
// @RequestBody, mas o WsProductService.alterProduct só lê os campos
// description, brand, model, price e currency dele. Em Go, separamos isso
// em um DTO de entrada próprio, mais explícito sobre o que a API espera
// receber no corpo da requisição.
type ProductUpdateInDTO struct {
	Description string  `json:"description"`
	Brand       string  `json:"brand"`
	Model       string  `json:"model"`
	Price       float64 `json:"price"`
	Currency    string  `json:"currency"`
	ImageUrl    string  `json:"imageUrl"`
	VideoUrl    string  `json:"videoUrl"`
}
