package repositories

import (
	"context"

	"github.com/PedroScheurer/product-service/internal/entities"
)

// ProductRepository é o equivalente Go da interface Java
// ProductRepository extends JpaRepository<ProductEntity, Long>.
//
// O Spring Data JPA gera a implementação de findById, findAll(Pageable),
// existsById, save e deleteById automaticamente a partir da interface.
// Em Go não há geração automática, então declaramos a interface aqui
// (para permitir mocks em testes) e implementamos manualmente em
// postgres_product_repository.go.
type ProductRepository interface {
	// FindByID busca um produto pelo ID.
	// Retorna (nil, nil) quando não encontrado — equivalente ao
	// Optional<ProductEntity> vazio do Java; quem chama decide se
	// isso é um erro (ver services, que traduz para ProductNotFoundError).
	FindByID(ctx context.Context, id int64) (*entities.ProductEntity, error)

	// FindAllPaged busca produtos paginados e ordenados.
	// Equivalente a repository.findAll(Pageable) do Spring Data.
	// Retorna os produtos da página solicitada e o total de registros
	// (necessário para montar a resposta de paginação, como o Page<T> do Java).
	FindAllPaged(ctx context.Context, page, size int, sortBy, sortDir string) ([]entities.ProductEntity, int64, error)

	// ExistsByID equivalente a repository.existsById(id).
	ExistsByID(ctx context.Context, id int64) (bool, error)

	// Save insere ou atualiza um produto.
	// Se product.ID == 0, insere e retorna o produto com o ID gerado
	// (equivalente ao INSERT do Java quando a entidade é nova).
	// Se product.ID != 0, atualiza o registro existente.
	Save(ctx context.Context, product *entities.ProductEntity) (*entities.ProductEntity, error)

	// DeleteByID equivalente a repository.deleteById(id).
	DeleteByID(ctx context.Context, id int64) error
}
