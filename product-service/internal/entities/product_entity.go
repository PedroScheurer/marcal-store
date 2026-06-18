package entities

// ProductEntity é o equivalente Go da entidade JPA ProductEntity.
// Mapeia a tabela tb_product. As tags `db` são usadas pelo sqlx
// para fazer o scan das linhas do banco para a struct.
type ProductEntity struct {
	ID          int64   `db:"id"`
	Description string  `db:"description"`
	Brand       string  `db:"brand"`
	Model       string  `db:"model"`
	Currency    string  `db:"currency"`
	Price       float64 `db:"price"`
	ImageUrl    string  `db:"image_url"`
	VideoUrl    string  `db:"video_url"`
}
