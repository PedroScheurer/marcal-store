package services

import (
	"context"
	"fmt"

	"github.com/PedroScheurer/product-service/internal/apperrors"
	"github.com/PedroScheurer/product-service/internal/dtos"
	"github.com/PedroScheurer/product-service/internal/entities"
	"github.com/PedroScheurer/product-service/internal/repositories"
)

// ProductService é o equivalente Go da classe Java ProductService.
// Cobre as operações de leitura usadas pelo ProductController:
// busca por id com conversão de moeda, busca sem conversão e
// listagem paginada com conversão.
type ProductService struct {
	repository                repositories.ProductRepository
	currencyConversionService *CurrencyConversionService

	// port é o equivalente ao @Value("${server.port}") private String port
	// do Java: a porta em que este serviço está rodando, usada apenas para
	// compor a string de "environment" nas respostas.
	port string
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

// FindByID é o equivalente a ProductService.findById(id, targetCurrency).
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

// FindProductNoConversion é o equivalente a ProductService.findProductNoConversion(idProduct).
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

// FindProductsPaged é o equivalente a ProductService.findProductsPaged(targetCurrency, pageable).
//
// page é 0-based (igual ao Pageable do Spring), size é o tamanho da página,
// sortBy/sortDir equivalem ao Sort embutido no Pageable
// (default: "description" ASC, vide @PageableDefault no controller Java).
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

// toProductDTO monta o ProductDTO de saída a partir da entidade e do
// resultado de conversão, equivalente ao `new ProductDTO(...)` repetido
// três vezes no Java original.
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
		Description:       product.Description,
		Brand:             product.Brand,
		Model:             product.Model,
		Price:             product.Price,
		Currency:          product.Currency,
		ImageUrl:          product.ImageUrl,
		VideoUrl:          product.VideoUrl,
		Environment:       environment,
		ConvertedPrice:    convertedPrice,
		RequestedCurrency: requestedCurrency,
	}
}
