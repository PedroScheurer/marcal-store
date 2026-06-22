package br.edu.atitus.productservice.services;

import br.edu.atitus.productservice.dtos.ConversionResult;
import br.edu.atitus.productservice.dtos.ProductDTO;
import br.edu.atitus.productservice.entities.ProductEntity;
import br.edu.atitus.productservice.infrastructure.exceptions.ProductNotFoundException;
import br.edu.atitus.productservice.repositories.ProductRepository;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.stereotype.Service;

@Service
public class ProductService {

    private final ProductRepository repository;
    private final CurrencyConversionService currencyConversionService;

    @Value("${server.port}")
    private String port;

    public ProductService(ProductRepository repository, CurrencyConversionService currencyConversionService) {
        this.repository = repository;
        this.currencyConversionService = currencyConversionService;
    }

    public ProductDTO findById(Long id, String targetCurrency) {
        ProductEntity product = repository.findById(id)
                .orElseThrow(() -> new ProductNotFoundException("Produto não encontrado"));

        ConversionResult conversionResult = currencyConversionService.convert(product, targetCurrency, port);

        return new ProductDTO(product.getId(),
                product.getDescription(),
                product.getBrand(),
                product.getModel(),
                product.getPrice(),
                product.getCurrency(),
                product.getStock(),
                conversionResult.environment(),
                conversionResult.convertedPrice(),
                targetCurrency);
    }

    public ProductDTO findProductNoConversion(Long idProduct) {
        ProductEntity product = repository.findById(idProduct)
                .orElseThrow(() -> new ProductNotFoundException("Produto não encontrado com o ID: " + idProduct));

        return new ProductDTO(
                product.getId(),
                product.getDescription(),
                product.getBrand(),
                product.getModel(),
                product.getPrice(),
                product.getCurrency(),
                product.getStock(),
                "Product-service running on port: " + port,
                null,
                null
        );
    }

    public Page<ProductDTO> findProductsPaged(String targetCurrency, Pageable pageable) {
        Page<ProductEntity> products = repository.findAll(pageable);

        return products.map(product -> {
            ConversionResult conversionResult = currencyConversionService.convert(product, targetCurrency, port);

            return new ProductDTO(product.getId(),
                    product.getDescription(),
                    product.getBrand(),
                    product.getModel(),
                    product.getPrice(),
                    product.getCurrency(),
                    product.getStock(),
                    conversionResult.environment(),
                    conversionResult.convertedPrice(),
                    targetCurrency);
        });
    }
}
