package main

import (
	"log"
	"net/http"
	"time"

	"github.com/PedroScheurer/product-service/internal/infra"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/PedroScheurer/product-service/internal/config"
	"github.com/PedroScheurer/product-service/internal/controllers"
	"github.com/PedroScheurer/product-service/internal/repositories"
	"github.com/PedroScheurer/product-service/internal/services"
)

// main é o equivalente Go da classe ProductServiceApplication
// (@SpringBootApplication + SpringApplication.run). Aqui montamos
// manualmente o grafo de dependências que o Spring monta via injeção
// de dependência automática (controllers -> services -> repositories),
// e iniciamos o servidor HTTP.
func main() {
	cfg := config.Load()

	db, err := sqlx.Connect("postgres", cfg.PostgresDSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// --- Camada de dados ---
	productRepository := repositories.NewProductRepository(db)

	// --- Camada de serviços ---
	// Equivalente ao spring.cache.caffeine.spec=maximumSize=500,expireAfterWrite=15s
	cacheService := services.NewCacheService(500, 15*time.Second)

	currencyServiceURL := "http://localhost:8081"
	currencyClient := services.NewHTTPCurrencyClient(currencyServiceURL, 5*time.Second)

	currencyConversionService := services.NewCurrencyConversionService(currencyClient, cacheService)
	productService := services.NewProductService(productRepository, currencyConversionService, cfg.ServerPort)
	wsProductService := services.NewWsProductService(productRepository)

	// --- Camada de controllers ---
	productController := controllers.NewProductController(productService)
	wsProductController := controllers.NewWsProductController(wsProductService)

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	productController.RegisterRoutes(router)
	wsProductController.RegisterRoutes(router)

	router.Route("/management", func(r chi.Router) {
		r.Handle("/health", infra.NewHealthHandler(db))
		r.Get("/info", infra.NewInfoHandler())
		r.Handle("/metrics", promhttp.Handler())
	})

	addr := ":" + cfg.ServerPort
	log.Printf("product-service listening on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
