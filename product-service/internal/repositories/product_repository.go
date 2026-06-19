package repositories

import (
	"context"

	"github.com/PedroScheurer/product-service/internal/entities"
)

type ProductRepository interface {
	FindByID(ctx context.Context, id int64) (*entities.ProductEntity, error)
	FindAllPaged(ctx context.Context, page, size int, sortBy, sortDir string) ([]entities.ProductEntity, int64, error)
	ExistsByID(ctx context.Context, id int64) (bool, error)
	Save(ctx context.Context, product *entities.ProductEntity) (*entities.ProductEntity, error)
	DeleteByID(ctx context.Context, id int64) error
}
