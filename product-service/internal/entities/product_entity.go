package entities

// ProductEntity é o equivalente Go da entidade JPA ProductEntity.
// Mapeia a tabela tb_product. As tags `db` são usadas pelo sqlx
// para fazer o scan das linhas do banco para a struct.
type ProductEntity struct {
	ID          int64   `db:"id"`
	Name        string  `db:"name"`
	Instructor  string  `db:"instructor"`
	ImageURL    string  `db:"image_url"`
	VideoURL    string  `db:"video_url"`
	Description string  `db:"description"`
	Workload    int     `db:"workload"`
	Modules     int     `db:"modules"`
	Price       float64 `db:"price"`
	Currency    string  `db:"currency"`
}
