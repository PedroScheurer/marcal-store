package services

import (
	"context"
	"fmt"

	"github.com/PedroScheurer/product-service/internal/apperrors"
	"github.com/PedroScheurer/product-service/internal/dtos"
	"github.com/PedroScheurer/product-service/internal/entities"
	"github.com/PedroScheurer/product-service/internal/repositories"
)

// adminType é o equivalente Go de WsProductService.ADMIN_TYPE (Integer = 0).
const adminType = 0

// WsProductService é o equivalente Go da classe Java WsProductService:
// operações de escrita (criar, alterar, excluir) com checagem de
// permissão baseada no header X-User-Type.
type WsProductService struct {
	repository repositories.ProductRepository
}

func NewWsProductService(repository repositories.ProductRepository) *WsProductService {
	return &WsProductService{repository: repository}
}

// CreateProduct é o equivalente a WsProductService.createProduct(dto, userType).
//
// Atenção: no Java original a checagem está invertida em relação ao que o
// nome ADMIN_TYPE sugere — `if (userType.equals(ADMIN_TYPE)) throw ...`
// bloqueia justamente o usuário do tipo "admin" (0) de criar produtos,
// e libera qualquer outro tipo. Mantemos esse mesmo comportamento aqui
// por fidelidade ao original, mas vale confirmar com o time se essa
// regra é intencional ou um bug histórico.
func (s *WsProductService) CreateProduct(ctx context.Context, dto dtos.ProductInDTO, userType int) (*dtos.ProductOutDTO, error) {
	if userType != adminType {
		return nil, apperrors.NewAuthenticationError("Usuário sem Permissão!")
	}

	product := &entities.ProductEntity{
		Description: dto.Description,
		Brand:       dto.Brand,
		Model:       dto.Model,
		Currency:    dto.Currency,
		Price:       dto.Price,
		ImageUrl:    dto.ImageUrl,
		VideoUrl:    dto.VideoUrl,
	}

	newProduct, err := s.repository.Save(ctx, product)
	if err != nil {
		return nil, fmt.Errorf("save product: %w", err)
	}

	return toProductOutDTO(newProduct), nil
}

// AlterProduct é o equivalente a WsProductService.alterProduct(idProduct, dto, userType).
// Aqui a checagem segue a lógica "normal": só ADMIN_TYPE (0) pode alterar.
func (s *WsProductService) AlterProduct(ctx context.Context, id int64, dto dtos.ProductUpdateInDTO, userType int) (*dtos.ProductOutDTO, error) {
	if userType != adminType {
		return nil, apperrors.NewAuthenticationError("Usuário sem Permissão!")
	}

	product, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("find product: %w", err)
	}
	if product == nil {
		return nil, apperrors.NewProductNotFoundError(fmt.Sprintf("Produto não encontrado com o ID: %d", id))
	}

	product.Description = dto.Description
	product.Brand = dto.Brand
	product.Model = dto.Model
	product.Price = dto.Price
	product.Currency = dto.Currency
	product.ImageUrl = dto.ImageUrl
	product.VideoUrl = dto.VideoUrl

	updatedProduct, err := s.repository.Save(ctx, product)
	if err != nil {
		return nil, fmt.Errorf("save product: %w", err)
	}

	return toProductOutDTO(updatedProduct), nil
}

// DeleteProduct é o equivalente a WsProductService.deleteProduct(idProduct, userType).
func (s *WsProductService) DeleteProduct(ctx context.Context, id int64, userType int) error {
	if userType != adminType {
		return apperrors.NewAuthenticationError("Usuário sem Permissão!")
	}

	exists, err := s.repository.ExistsByID(ctx, id)
	if err != nil {
		return fmt.Errorf("check product exists: %w", err)
	}
	if !exists {
		return apperrors.NewProductNotFoundError(fmt.Sprintf("Produto não encontrado com o ID: %d", id))
	}

	if err := s.repository.DeleteByID(ctx, id); err != nil {
		return fmt.Errorf("delete product: %w", err)
	}

	return nil
}

func toProductOutDTO(product *entities.ProductEntity) *dtos.ProductOutDTO {
	return &dtos.ProductOutDTO{
		ID:          product.ID,
		Description: product.Description,
		Brand:       product.Brand,
		Model:       product.Model,
		Price:       product.Price,
		Currency:    product.Currency,
		ImageUrl:    product.ImageUrl,
		VideoUrl:    product.VideoUrl,
	}
}
