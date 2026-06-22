package br.edu.atitus.productservice.services;

import br.edu.atitus.productservice.dtos.ProductDTO;
import br.edu.atitus.productservice.dtos.ProductInDTO;
import br.edu.atitus.productservice.dtos.ProductOutDTO;
import br.edu.atitus.productservice.entities.ProductEntity;
import br.edu.atitus.productservice.infrastructure.exceptions.ProductNotFoundException;
import br.edu.atitus.productservice.repositories.ProductRepository;
import org.springframework.beans.BeanUtils;
import org.springframework.stereotype.Service;

import javax.naming.AuthenticationException;

@Service
public class WsProductService {

    private final static Integer ADMIN_TYPE = 0;

    private final ProductRepository repository;

    public WsProductService(ProductRepository repository) {
        this.repository = repository;
    }

    public ProductOutDTO createProduct(ProductInDTO dto, Integer userType) throws AuthenticationException {
        if (userType.equals(ADMIN_TYPE))
            throw new AuthenticationException("Usuário sem Permissão!");

        ProductEntity product = convertDTO2Entity(dto);
        ProductEntity newProduct = repository.save(product);
        return new ProductOutDTO(newProduct.getId(),
                newProduct.getDescription(),
                newProduct.getBrand(),
                newProduct.getModel(),
                newProduct.getPrice(),
                newProduct.getCurrency(),
                newProduct.getStock());
    }


    public ProductOutDTO alterProduct(Long idProduct, ProductDTO dto, Integer userType) throws AuthenticationException {
        if (!userType.equals(ADMIN_TYPE)) {
            throw new AuthenticationException("Usuário sem Permissão!");
        }

        ProductEntity product = repository.findById(idProduct)
                .orElseThrow(() -> new ProductNotFoundException("Produto não encontrado com o ID: " + idProduct));

        product.setDescription(dto.description());
        product.setBrand(dto.brand());
        product.setModel(dto.model());
        product.setPrice(dto.price());
        product.setCurrency(dto.currency());

        ProductEntity updatedEntity = repository.save(product);

        return new ProductOutDTO(
                updatedEntity.getId(),
                updatedEntity.getDescription(),
                updatedEntity.getBrand(),
                updatedEntity.getModel(),
                updatedEntity.getPrice(),
                updatedEntity.getCurrency(),
                updatedEntity.getStock()
        );
    }

    public void deleteProduct(Long idProduct, Integer userType) throws AuthenticationException {
        if (!userType.equals(ADMIN_TYPE)) {
            throw new AuthenticationException("Usuário sem Permissão!");
        }

        if (!repository.existsById(idProduct)) {
            throw new ProductNotFoundException("Produto não encontrado com o ID: " + idProduct);
        }
        repository.deleteById(idProduct);
    }

    private ProductEntity convertDTO2Entity(ProductInDTO dto){
        var entity = new ProductEntity();
        BeanUtils.copyProperties(dto, entity);
        return entity;
    }


}
