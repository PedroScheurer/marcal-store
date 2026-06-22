package br.edu.atitus.productservice.controllers;

import br.edu.atitus.productservice.repositories.ProductRepository;
import io.restassured.RestAssured;
import io.restassured.http.ContentType;
import org.hamcrest.Matchers;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.web.server.LocalServerPort;
import org.springframework.http.HttpStatus;

@SpringBootTest(webEnvironment = SpringBootTest.WebEnvironment.RANDOM_PORT)
class ProductControllerTest {
    @LocalServerPort
    private int port;

    @Autowired
    private ProductRepository repository;

    @BeforeEach
    void setup() {
        RestAssured.port = port;
        RestAssured.basePath = "/products";
    }

    @Test
    public void shouldReturnProductWithNullConvertedPriceWhenCurrencyIsInvalid() {
        Integer id = 1;

        RestAssured
                .given()
                .pathParam("id", id)
                .queryParam("targetCurrency", "XXX")
                .accept(ContentType.JSON)
                .when()
                .get("/{id}")
                .then()
                .statusCode(HttpStatus.OK.value())
                .body("id", Matchers.equalTo(id))
                .body("convertedPrice", Matchers.equalTo(null))
                .body("requestedCurrency", Matchers.equalTo("XXX"));

    }

}