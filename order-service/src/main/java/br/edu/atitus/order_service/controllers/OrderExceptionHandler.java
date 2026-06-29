package br.edu.atitus.order_service.controllers;

import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.RestControllerAdvice;

import feign.FeignException;

@RestControllerAdvice
public class OrderExceptionHandler {

    @ExceptionHandler(IllegalArgumentException.class)
    public ResponseEntity<ErrorBody> handleBadRequest(IllegalArgumentException ex) {
        return ResponseEntity.status(HttpStatus.BAD_REQUEST)
                .body(new ErrorBody(ex.getMessage()));
    }

    @ExceptionHandler(FeignException.class)
    public ResponseEntity<ErrorBody> handleFeign(FeignException ex) {
        String message = ex.status() == 404
                ? "Produto não encontrado no catálogo."
                : "Não foi possível consultar um serviço interno. Tente novamente em instantes.";
        return ResponseEntity.status(HttpStatus.BAD_GATEWAY)
                .body(new ErrorBody(message));
    }

    @ExceptionHandler(Exception.class)
    public ResponseEntity<ErrorBody> handleGeneric(Exception ex) {
        return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR)
                .body(new ErrorBody("Não foi possível processar o pedido. Tente novamente."));
    }

    public record ErrorBody(String message) {}
}
