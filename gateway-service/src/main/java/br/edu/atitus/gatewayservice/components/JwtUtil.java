package br.edu.atitus.gatewayservice.components;

import br.edu.atitus.gatewayservice.infrastructure.exceptions.InvalidTokenException;
import br.edu.atitus.gatewayservice.infrastructure.exceptions.TokenExpiredException;
import io.jsonwebtoken.Claims;
import io.jsonwebtoken.ExpiredJwtException;
import io.jsonwebtoken.JwtException;
import io.jsonwebtoken.Jwts;
import io.jsonwebtoken.security.Keys;

import javax.crypto.SecretKey;

public class JwtUtil {

    private static final String SECRET_KEY = "chaveSuperSecretaParaJWTdeExemplo!@#123"; // Chave secreta (use uma mais segura)

    private JwtUtil() {
        throw new UnsupportedOperationException("Utility class");
    }

    private static SecretKey getSigningKey() {
        return Keys.hmacShaKeyFor(SECRET_KEY.getBytes());
    }

    public static Claims validateToken(String token) {
        try {
            return Jwts.parser()
                    .verifyWith(getSigningKey())
                    .build()
                    .parseSignedClaims(token)
                    .getPayload();
        } catch (ExpiredJwtException e) {
            throw new TokenExpiredException("Token expirado", e);
        } catch (JwtException e) {
            throw new InvalidTokenException("Token inválido", e);
        }
    }
}