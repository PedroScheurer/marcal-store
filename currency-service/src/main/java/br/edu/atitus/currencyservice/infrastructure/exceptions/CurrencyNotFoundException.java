package br.edu.atitus.currencyservice.infrastructure.exceptions;

public class CurrencyNotFoundException extends RuntimeException{
    public CurrencyNotFoundException(String message){
        super(message);
    }
}
