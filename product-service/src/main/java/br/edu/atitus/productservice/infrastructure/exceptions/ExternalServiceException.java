package br.edu.atitus.productservice.infrastructure.exceptions;

public class ExternalServiceException extends RuntimeException{
    public ExternalServiceException(String message){
        super(message);
    }
}
