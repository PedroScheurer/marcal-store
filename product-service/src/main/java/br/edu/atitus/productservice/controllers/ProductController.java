package br.edu.atitus.productservice.controllers;

import br.edu.atitus.productservice.dtos.ProductDTO;
import br.edu.atitus.productservice.services.ProductService;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.domain.Sort;
import org.springframework.data.web.PageableDefault;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/products")
public class ProductController {

    private final ProductService service;

    public ProductController(ProductService service) {
        this.service = service;
    }

    @GetMapping(path = {"{idProduct}"}, params = {"targetCurrency"})
    public ResponseEntity<ProductDTO> findProduct(@PathVariable Long idProduct,
                                                 @RequestParam String targetCurrency){

        ProductDTO productDTO = service.findById(idProduct, targetCurrency);

        return ResponseEntity.ok().body(productDTO);
    }

    @GetMapping("/noconverter/{idProduct}")
    public ResponseEntity<ProductDTO> getProductNoConverter(@PathVariable Long idProduct) {
        ProductDTO dto = service.findProductNoConversion(idProduct);

        return ResponseEntity.ok(dto);
    }

    @GetMapping
    public ResponseEntity<Page<ProductDTO>> getAllProducts(
            @RequestParam String targetCurrency,
            @PageableDefault(
                    page = 0,
                    size = 5,
                    sort = "description",
                    direction = Sort.Direction.ASC
            ) Pageable pageable) {

        Page<ProductDTO> productDTOs = service.findProductsPaged(targetCurrency, pageable);


        return ResponseEntity.ok(productDTOs);
    }
}
