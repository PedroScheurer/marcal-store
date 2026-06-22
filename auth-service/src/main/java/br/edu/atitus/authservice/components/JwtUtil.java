package br.edu.atitus.authservice.components;

import java.util.Date;
import java.util.HashMap;
import java.util.Map;

import javax.crypto.SecretKey;

import br.edu.atitus.authservice.entities.UserType;
import io.jsonwebtoken.Claims;
import io.jsonwebtoken.Jwts;
import io.jsonwebtoken.security.Keys;
import jakarta.servlet.http.HttpServletRequest;

public class JwtUtil {

    private static final String SECRET_KEY = "chaveSuperSecretaParaJWTdeExemplo!@#123"; // Chave secreta (use uma mais segura)
    private static final long EXPIRATION_TIME = 1000 * 60 * 60 * 24; // 1000 milisegundos * 60 segundos * 60 minutos * 24 horas

    private static final String AUTHORIZATION_HEADER = "Authorization";
    private static final String BEARER_PREFIX = "Bearer ";

    private static final String CLAIM_ID = "id";
    private static final String CLAIM_EMAIL = "email";
    private static final String CLAIM_TYPE = "type";

    private JwtUtil() {
        throw new UnsupportedOperationException("Utility class");
    }

    private static SecretKey getSigningKey() {
        return Keys.hmacShaKeyFor(SECRET_KEY.getBytes());
    }

    public static String generateToken(String email, Long id, UserType type) {
        return Jwts.builder()
                .claims(buildClaims(email, id, type))
                .issuedAt(new Date())
                .expiration(createExpirationDate())
                .signWith(getSigningKey())
                .compact();
    }

    public static Claims validateToken(String token) {
        try {
            return Jwts.parser()
                    .verifyWith(getSigningKey())
                    .build()
                    .parseSignedClaims(token)
                    .getPayload();
        } catch (Exception e) {
            return null;
        }
    }

    public static String getJwtFromRequest(HttpServletRequest request) {
        String authorizationHeader = request.getHeader(AUTHORIZATION_HEADER);

        if (!hasBearerToken(authorizationHeader)) {
            return null;
        }

        return authorizationHeader.substring(BEARER_PREFIX.length());
    }

    private static Map<String, Object> buildClaims(
            String email,
            Long id,
            UserType type
    ) {
        Map<String, Object> claims = new HashMap<>();
        claims.put(CLAIM_ID, id);
        claims.put(CLAIM_EMAIL, email);
        claims.put(CLAIM_TYPE, type.ordinal());

        return claims;
    }

    private static Date createExpirationDate() {
        return new Date(System.currentTimeMillis() + EXPIRATION_TIME);
    }

    private static boolean hasBearerToken(String authorizationHeader) {
        return authorizationHeader != null
                && authorizationHeader.startsWith(BEARER_PREFIX);
    }
}