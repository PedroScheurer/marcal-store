package br.edu.atitus.authservice.infrastructure.exceptions;

public class ExternalServiceException extends RuntimeException{
    public ExternalServiceException(String message){
        super(message);
    }
}
