package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/PedroScheurer/product-service/internal/clients"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/PedroScheurer/product-service/internal/config"
	"github.com/PedroScheurer/product-service/internal/controllers"
	"github.com/PedroScheurer/product-service/internal/infra"
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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, err := sqlx.Connect("postgres", cfg.PostgresDSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// --- Executar Migrations ---
	if err := infra.RunMigrations(cfg.PostgresMigrationUrl()); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	// --- Camada de dados ---
	productRepository := repositories.NewProductRepository(db)

	// --- Camada de serviços ---
	// Equivalente ao spring.cache.caffeine.spec=maximumSize=500,expireAfterWrite=15s
	cacheService := services.NewCacheService(500, 15*time.Second)

	eurekaURL := "http://localhost:8761/eureka" // ajuste se já tiver isso no cfg
	eurekaClient := clients.NewEurekaClient(eurekaURL, cfg.ApplicationName, cfg.HostName, cfg.ServerPort)

	discovery := clients.NewServiceDiscovery(eurekaURL, &http.Client{Timeout: 5 * time.Second}, 20*time.Second)
	currencyClient := clients.NewHTTPCurrencyClient(discovery, 5*time.Second)

	if err := eurekaClient.Register(ctx); err != nil {
		log.Printf("failed to register on eureka: %v", err)
	}

	go eurekaClient.StartHeartbeatLoop(ctx, 30*time.Second)

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
	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		log.Printf("product-service listening on %s", addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := eurekaClient.Deregister(shutdownCtx); err != nil {
		log.Printf("eureka deregister failed: %v", err)
	}

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown failed: %v", err)
	}

	log.Println("server stopped gracefully")
}
