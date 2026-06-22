package br.edu.atitus.productservice.controllers;

import br.edu.atitus.productservice.dtos.ProductDTO;
import br.edu.atitus.productservice.dtos.ProductInDTO;
import br.edu.atitus.productservice.dtos.ProductOutDTO;
import br.edu.atitus.productservice.entities.ProductEntity;
import br.edu.atitus.productservice.services.WsProductService;
import org.springframework.http.HttpStatus;
import org.springframework.http.HttpStatusCode;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import javax.naming.AuthenticationException;

@RestController
@RequestMapping("/ws/products")
public class WsProductController {

    private final WsProductService service;

    public WsProductController(WsProductService service) {
        this.service = service;
    }

    @PostMapping
    public ResponseEntity<ProductOutDTO> postProduct(
            @RequestBody ProductInDTO dto,
            @RequestHeader("X-User-Id") Long userId,
            @RequestHeader("X-User-Email") String userEmail,
            @RequestHeader("X-User-Type") Integer type) throws AuthenticationException {

        ProductOutDTO productDTO = service.createProduct(dto, type);

        return ResponseEntity.status(201).body(productDTO);
    }

    @PutMapping("/{idProduct}")
    public ResponseEntity<ProductOutDTO> putProduct(
            @PathVariable Long idProduct,
            @RequestBody ProductDTO dto,
            @RequestHeader("X-User-Id") Long userId,
            @RequestHeader("X-User-Email") String userEmail,
            @RequestHeader("X-User-Type") Integer type) throws AuthenticationException {
        ProductOutDTO productDTO = service.alterProduct(idProduct, dto, type);

        return ResponseEntity.ok().body(productDTO);
    }

    @DeleteMapping("/{idProduct}")
    public ResponseEntity<Void> deleteProduct(
            @PathVariable Long idProduct,
            @RequestHeader("X-User-Id") Long userId,
            @RequestHeader("X-User-Email") String userEmail,
            @RequestHeader("X-User-Type") Integer type) throws AuthenticationException {

        service.deleteProduct(idProduct, type);

        return ResponseEntity.status(HttpStatus.NO_CONTENT).build();
    }

    @ExceptionHandler(Exception.class)
    public ResponseEntity<String> handleException(Exception e){
        String message = e.getMessage().replace("/r/n", "");
        return ResponseEntity.badRequest().body(message);
    }
}