package br.edu.atitus.gatewayservice.filters;

import br.edu.atitus.gatewayservice.components.JwtUtil;
import io.jsonwebtoken.Claims;
import org.apache.http.HttpHeaders;
import org.springframework.cloud.gateway.filter.GatewayFilterChain;
import org.springframework.cloud.gateway.filter.GlobalFilter;
import org.springframework.core.Ordered;
import org.springframework.http.HttpStatus;
import org.springframework.http.server.reactive.ServerHttpRequest;
import org.springframework.stereotype.Component;
import org.springframework.web.server.ServerWebExchange;
import reactor.core.publisher.Mono;

import java.util.List;

@Component
public class AuthFilter implements GlobalFilter, Ordered {
    private static final List<String> PROTECTED_ROUTES = List.of("/ws/");
    private static final String BEARER_PREFIX = "Bearer ";


    @Override
    public Mono<Void> filter(ServerWebExchange exchange, GatewayFilterChain chain) {
        ServerHttpRequest request = exchange.getRequest();
        String path = request.getURI().getPath();

        if(!isProtectedRoute(path)){
            return chain.filter(exchange);
        }

        String authHeader = request.getHeaders().getFirst(HttpHeaders.AUTHORIZATION);

        if(!hasBearerToken(authHeader)){
            exchange.getResponse().setStatusCode(HttpStatus.UNAUTHORIZED);
            return exchange.getResponse().setComplete();
        }

        String jwt = getJwtFromAuthHeader(authHeader);
        Claims payload = JwtUtil.validateToken(jwt);

        if(payload == null){
            exchange.getResponse().setStatusCode(HttpStatus.UNAUTHORIZED);
            return exchange.getResponse().setComplete();
        }

        ServerHttpRequest mutatedRequest = request.mutate()
                .header("X-User-Id", String.valueOf(payload.get("id", Long.class)))
                .header("X-User-Type", String.valueOf(payload.get("type",Integer.class)))
                .header("X-User-Email", payload.get("email", String.class))
                .build();

        return chain.filter(
                exchange.mutate()
                        .request(mutatedRequest)
                        .build()
        );
    }

    private boolean isProtectedRoute(String path) {
        return PROTECTED_ROUTES.stream().anyMatch(path::startsWith);
    }

    private String getJwtFromAuthHeader(String authHeader) {
        return authHeader.substring(BEARER_PREFIX.length());
    }

    private static boolean hasBearerToken(String authorizationHeader) {
        return authorizationHeader != null
                && authorizationHeader.startsWith(BEARER_PREFIX);
    }

    @Override
    public int getOrder() {
        return -1;
    }
}
