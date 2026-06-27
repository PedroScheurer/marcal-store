package br.edu.atitus.gatewayservice.configs;

import org.springframework.cloud.gateway.route.RouteLocator;
import org.springframework.cloud.gateway.route.builder.RouteLocatorBuilder;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class GatewayConfig {

    @Bean
    RouteLocator getGatewayRoutes(RouteLocatorBuilder builder){
        return builder.routes()
                .route(predicateSpec -> predicateSpec
                        .path("/products/**")
                        .uri("lb://product-service"))
                .route(predicateSpec -> predicateSpec
                        .path("/ws/products/**")
                        .uri("lb://product-service"))
                .route(predicateSpec -> predicateSpec
                        .path("/media/**")
                        .uri("lb://product-service"))
                .route(predicateSpec -> predicateSpec
                        .path("/currency/**")
                        .uri("lb://currency-service"))
                .route(predicateSpec -> predicateSpec
                        .path("/ws/currency/**")
                        .uri("lb://currency-service"))
                .route(p -> p
                        .path("/ws/orders/**")
                        .uri("lb://order-service"))
                .route(predicateSpec -> predicateSpec
                        .path("/auth/**")
                        .uri("lb://auth-service"))
                .build();
    }
}
