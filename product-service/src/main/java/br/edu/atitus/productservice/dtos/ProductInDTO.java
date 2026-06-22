package br.edu.atitus.productservice.dtos;

public record ProductInDTO(
        String description,
        String brand,
        String model,
        String currency,
        Double price,
        Integer stock
) {
}