package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/PedroScheurer/product-service/internal/entities"
)

type postgresProductRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) ProductRepository {
	return &postgresProductRepository{db: db}
}

// allowedSortColumns é a whitelist de colunas permitidas no ORDER BY,
// evitando SQL injection via parâmetro de ordenação.
var allowedSortColumns = map[string]string{
	"id":          "id",
	"name":        "name",
	"instructor":  "instructor",
	"description": "description",
	"workload":    "workload",
	"modules":     "modules",
	"price":       "price",
	"currency":    "currency",
}

const selectColumns = `id, name, instructor, image_url, video_url, description, workload, modules, module_titles, price, currency`

func (r *postgresProductRepository) FindByID(ctx context.Context, id int64) (*entities.ProductEntity, error) {
	var product entities.ProductEntity

	query := `SELECT ` + selectColumns + ` FROM tb_product WHERE id = $1`

	err := r.db.GetContext(ctx, &product, query, id)
	if errors.Is(err, sql.ErrNoRows) {
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
		column = "name"
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
	query := fmt.Sprintf(
		`SELECT `+selectColumns+` FROM tb_product ORDER BY %s %s LIMIT $1 OFFSET $2`,
		column, direction,
	)

	var products []entities.ProductEntity
	if err := r.db.SelectContext(ctx, &products, query, size, offset); err != nil {
		return nil, 0, fmt.Errorf("find products paged: %w", err)
	}

	return products, total, nil
}

func (r *postgresProductRepository) ExistsByID(ctx context.Context, id int64) (bool, error) {
	var exists bool
	err := r.db.GetContext(ctx, &exists, `SELECT EXISTS(SELECT 1 FROM tb_product WHERE id = $1)`, id)
	if err != nil {
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
	query := `
		INSERT INTO tb_product (name, instructor, image_url, video_url, description, workload, modules, module_titles, price, currency)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id`

	err := r.db.QueryRowContext(ctx, query,
		product.Name, product.Instructor, product.ImageURL, product.VideoURL,
		product.Description, product.Workload, product.Modules, product.ModuleTitles, product.Price, product.Currency,
	).Scan(&product.ID)
	if err != nil {
		return nil, fmt.Errorf("insert product: %w", err)
	}

	return product, nil
}

func (r *postgresProductRepository) update(ctx context.Context, product *entities.ProductEntity) (*entities.ProductEntity, error) {
	query := `
		UPDATE tb_product
		SET name = $1, instructor = $2, image_url = $3, video_url = $4,
		    description = $5, workload = $6, modules = $7, module_titles = $8, price = $9, currency = $10
		WHERE id = $11`

	_, err := r.db.ExecContext(ctx, query,
		product.Name, product.Instructor, product.ImageURL, product.VideoURL,
		product.Description, product.Workload, product.Modules, product.ModuleTitles, product.Price, product.Currency,
		product.ID,
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
