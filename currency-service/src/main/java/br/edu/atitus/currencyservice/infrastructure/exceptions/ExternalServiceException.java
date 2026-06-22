package br.edu.atitus.currencyservice.infrastructure.exceptions;

public class ExternalServiceException extends RuntimeException{
    public ExternalServiceException(String message){
        super(message);
    }
}
