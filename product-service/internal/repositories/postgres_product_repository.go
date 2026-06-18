package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/PedroScheurer/product-service/internal/entities"
)

// postgresProductRepository é a implementação concreta de ProductRepository
// usando sqlx sobre Postgres. Equivalente, na prática, ao que o Spring Data
// JPA gera por trás da interface ProductRepository.
type postgresProductRepository struct {
	db *sqlx.DB
}

// NewProductRepository cria um ProductRepository baseado em Postgres.
func NewProductRepository(db *sqlx.DB) ProductRepository {
	return &postgresProductRepository{db: db}
}

// allowedSortColumns é a whitelist de colunas que podem ser usadas no
// ORDER BY. O Pageable do Spring aceita o nome do campo Java diretamente
// (ex.: "description") vindo da query string, então sem essa validação
// estaríamos abrindo a porta para SQL injection via parâmetro de sort.
var allowedSortColumns = map[string]string{
	"id":          "id",
	"description": "description",
	"brand":       "brand",
	"model":       "model",
	"currency":    "currency",
	"price":       "price",
}

func (r *postgresProductRepository) FindByID(ctx context.Context, id int64) (*entities.ProductEntity, error) {
	var product entities.ProductEntity

	query := `SELECT id, description, brand, model, currency, price
	          FROM tb_product WHERE id = $1`

	err := r.db.GetContext(ctx, &product, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		// Equivalente ao Optional vazio do Java: quem chama decide o que fazer.
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find product by id %d: %w", id, err)
	}

	return &product, nil
}

func (r *postgresProductRepository) FindAllPaged(ctx context.Context, page, size int, sortBy, sortDir string) ([]entities.ProductEntity, int64, error) {
	column, ok := allowedSortColumns[sortBy]
	if !ok {
		column = "description" // mesmo default usado no @PageableDefault do Java
	}

	direction := "ASC"
	if sortDir == "DESC" || sortDir == "desc" {
		direction = "DESC"
	}

	var total int64
	if err := r.db.GetContext(ctx, &total, `SELECT COUNT(*) FROM tb_product`); err != nil {
		return nil, 0, fmt.Errorf("count products: %w", err)
	}

	offset := page * size

	// column/direction vêm de uma whitelist fixa acima, nunca diretamente
	// do usuário, então a interpolação aqui é segura.
	query := fmt.Sprintf(`SELECT id, description, brand, model, currency, price, image_url, video_url
	          FROM tb_product ORDER BY %s %s LIMIT $1 OFFSET $2`, column, direction)

	var products []entities.ProductEntity
	if err := r.db.SelectContext(ctx, &products, query, size, offset); err != nil {
		return nil, 0, fmt.Errorf("find products paged: %w", err)
	}

	return products, total, nil
}

func (r *postgresProductRepository) ExistsByID(ctx context.Context, id int64) (bool, error) {
	var exists bool

	query := `SELECT EXISTS(SELECT 1 FROM tb_product WHERE id = $1)`

	if err := r.db.GetContext(ctx, &exists, query, id); err != nil {
		return false, fmt.Errorf("check product exists id %d: %w", id, err)
	}

	return exists, nil
}

func (r *postgresProductRepository) Save(ctx context.Context, product *entities.ProductEntity) (*entities.ProductEntity, error) {
	if product.ID == 0 {
		return r.insert(ctx, product)
	}
	return r.update(ctx, product)
}

func (r *postgresProductRepository) insert(ctx context.Context, product *entities.ProductEntity) (*entities.ProductEntity, error) {
	query := `INSERT INTO tb_product (description, brand, model, currency, price, image_url, video_url)
	          VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	err := r.db.QueryRowContext(ctx, query,
		product.Description, product.Brand, product.Model,
		product.Currency, product.Price, product.ImageUrl, product.VideoUrl,
	).Scan(&product.ID)
	if err != nil {
		return nil, fmt.Errorf("insert product: %w", err)
	}

	return product, nil
}

func (r *postgresProductRepository) update(ctx context.Context, product *entities.ProductEntity) (*entities.ProductEntity, error) {
	query := `UPDATE tb_product
	          SET description = $1, brand = $2, model = $3, currency = $4, price = $5, image_url = $6, video_url = $7
	          WHERE id = $8`

	_, err := r.db.ExecContext(ctx, query,
		product.Description, product.Brand, product.Model,
		product.Currency, product.Price, product.ImageUrl, product.VideoUrl, product.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("update product id %d: %w", product.ID, err)
	}

	return product, nil
}

func (r *postgresProductRepository) DeleteByID(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM tb_product WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete product id %d: %w", id, err)
	}
	return nil
}
