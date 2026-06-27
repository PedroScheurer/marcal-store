package services

import (
	"context"
	"fmt"

	"github.com/PedroScheurer/product-service/internal/apperrors"
	"github.com/PedroScheurer/product-service/internal/dtos"
	"github.com/PedroScheurer/product-service/internal/entities"
	"github.com/PedroScheurer/product-service/internal/modules"
	"github.com/PedroScheurer/product-service/internal/repositories"
)

const adminType = 0

type WsProductService struct {
	repository repositories.ProductRepository
}

func NewWsProductService(repository repositories.ProductRepository) *WsProductService {
	return &WsProductService{repository: repository}
}

func (s *WsProductService) CreateProduct(ctx context.Context, dto dtos.ProductInDTO, userType int) (*dtos.ProductOutDTO, error) {
	if userType != adminType {
		return nil, apperrors.NewAuthenticationError("Usuário sem Permissão!")
	}

	product := &entities.ProductEntity{
		Name:         dto.Name,
		Instructor:   dto.Instructor,
		ImageURL:     dto.ImageURL,
		VideoURL:     dto.VideoURL,
		Description:  dto.Description,
		Workload:     dto.Workload,
		Modules:      modules.ResolveCount(dto.ModuleTitles, dto.Modules),
		ModuleTitles: modules.Encode(dto.ModuleTitles),
		Price:        dto.Price,
		Currency:     dto.Currency,
	}

	newProduct, err := s.repository.Save(ctx, product)
	if err != nil {
		return nil, fmt.Errorf("save product: %w", err)
	}

	return toProductOutDTO(newProduct), nil
}

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

	product.Name = dto.Name
	product.Instructor = dto.Instructor
	product.ImageURL = dto.ImageURL
	product.VideoURL = dto.VideoURL
	product.Description = dto.Description
	product.Workload = dto.Workload
	product.Modules = modules.ResolveCount(dto.ModuleTitles, dto.Modules)
	product.ModuleTitles = modules.Encode(dto.ModuleTitles)
	product.Price = dto.Price
	product.Currency = dto.Currency

	updatedProduct, err := s.repository.Save(ctx, product)
	if err != nil {
		return nil, fmt.Errorf("save product: %w", err)
	}

	return toProductOutDTO(updatedProduct), nil
}

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
		ID:           product.ID,
		Name:         product.Name,
		Instructor:   product.Instructor,
		ImageURL:     product.ImageURL,
		VideoURL:     product.VideoURL,
		Description:  product.Description,
		Workload:     product.Workload,
		Modules:      product.Modules,
		ModuleTitles: modules.Decode(product.ModuleTitles),
		Price:        product.Price,
		Currency:     product.Currency,
	}
}
