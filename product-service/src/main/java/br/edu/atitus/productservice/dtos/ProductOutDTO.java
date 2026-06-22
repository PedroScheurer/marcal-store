package br.edu.atitus.productservice.dtos;

public record ProductOutDTO(Long id,
                            String description,
                            String brand,
                            String model,
                            Double price,
                            String currency,
                            Integer stock) {
}
