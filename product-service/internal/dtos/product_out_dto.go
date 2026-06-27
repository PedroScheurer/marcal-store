package dtos

// ProductOutDTO é a resposta de criação/atualização de um produto,
// sem dados de conversão de moeda.
type ProductOutDTO struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Instructor  string  `json:"instructor"`
	ImageURL    string  `json:"imageUrl"`
	VideoURL    string  `json:"videoUrl"`
	Description string  `json:"description"`
	Workload     int      `json:"workload"`
	Modules      int      `json:"modules"`
	ModuleTitles []string `json:"moduleTitles"`
	Price        float64  `json:"price"`
	Currency     string   `json:"currency"`
}
