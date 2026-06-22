package br.edu.atitus.gatewayservice.infrastructure.exceptions;

public class TokenExpiredException extends RuntimeException {
    public TokenExpiredException(String message, Throwable cause) {
        super(message, cause);
    }
}