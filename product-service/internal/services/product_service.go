package services

import (
	"context"
	"fmt"

	"github.com/PedroScheurer/product-service/internal/apperrors"
	"github.com/PedroScheurer/product-service/internal/dtos"
	"github.com/PedroScheurer/product-service/internal/entities"
	"github.com/PedroScheurer/product-service/internal/repositories"
)

type ProductService struct {
	repository                repositories.ProductRepository
	currencyConversionService *CurrencyConversionService
	port                      string
}

func NewProductService(
	repository repositories.ProductRepository,
	currencyConversionService *CurrencyConversionService,
	port string,
) *ProductService {
	return &ProductService{
		repository:                repository,
		currencyConversionService: currencyConversionService,
		port:                      port,
	}
}

func (s *ProductService) FindByID(ctx context.Context, id int64, targetCurrency string) (*dtos.ProductDTO, error) {
	product, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("find product: %w", err)
	}
	if product == nil {
		return nil, apperrors.NewProductNotFoundError("Produto não encontrado")
	}

	conversionResult := s.currencyConversionService.Convert(product, targetCurrency, s.port)

	return toProductDTO(product, &conversionResult, &targetCurrency), nil
}

func (s *ProductService) FindProductNoConversion(ctx context.Context, id int64) (*dtos.ProductDTO, error) {
	product, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("find product: %w", err)
	}
	if product == nil {
		return nil, apperrors.NewProductNotFoundError(fmt.Sprintf("Produto não encontrado com o ID: %d", id))
	}

	environment := "Product-service running on port: " + s.port

	return toProductDTO(product, &dtos.ConversionResult{Environment: environment}, nil), nil
}

func (s *ProductService) FindProductsPaged(ctx context.Context, targetCurrency string, page, size int, sortBy, sortDir string) (*dtos.Page[dtos.ProductDTO], error) {
	products, total, err := s.repository.FindAllPaged(ctx, page, size, sortBy, sortDir)
	if err != nil {
		return nil, fmt.Errorf("find products paged: %w", err)
	}

	content := make([]dtos.ProductDTO, 0, len(products))
	for i := range products {
		product := products[i]
		conversionResult := s.currencyConversionService.Convert(&product, targetCurrency, s.port)
		content = append(content, *toProductDTO(&product, &conversionResult, &targetCurrency))
	}

	return dtos.NewPage(content, page, size, total), nil
}

func toProductDTO(product *entities.ProductEntity, conversionResult *dtos.ConversionResult, requestedCurrency *string) *dtos.ProductDTO {
	var convertedPrice *float64
	var environment *string

	if conversionResult != nil {
		environment = &conversionResult.Environment
		if requestedCurrency != nil {
			price := conversionResult.ConvertedPrice
			convertedPrice = &price
		}
	}

	return &dtos.ProductDTO{
		ID:                product.ID,
		Name:              product.Name,
		Instructor:        product.Instructor,
		ImageURL:          product.ImageURL,
		VideoURL:          product.VideoURL,
		Description:       product.Description,
		Workload:          product.Workload,
		Modules:           product.Modules,
		Price:             product.Price,
		Currency:          product.Currency,
		ConvertedPrice:    convertedPrice,
		RequestedCurrency: requestedCurrency,
		Environment:       environment,
	}
}
