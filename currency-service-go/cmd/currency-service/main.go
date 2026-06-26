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

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/PedroScheurer/currency-service/internal/clients"
	"github.com/PedroScheurer/currency-service/internal/config"
	"github.com/PedroScheurer/currency-service/internal/controllers"
	"github.com/PedroScheurer/currency-service/internal/infra"
	"github.com/PedroScheurer/currency-service/internal/repositories"
	"github.com/PedroScheurer/currency-service/internal/services"
)

func main() {
	cfg := config.Load()

	db, err := sqlx.Connect("postgres", cfg.PostgresDSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := infra.RunMigrations(cfg.PostgresMigrationUrl()); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
	log.Println("migrations applied successfully")

	// Camada de dados
	currencyRepository := repositories.NewCurrencyRepository(db)

	// Camada de serviços
	cacheService := services.NewCacheService(500, 15*time.Second)
	bcbClient := services.NewHttpBCBClient(5 * time.Second)
	currencyService := services.NewCurrencyService(currencyRepository, bcbClient, cacheService, cfg.ServerPort)

	// Camada de controllers
	currencyController := controllers.NewCurrencyController(currencyService)

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	currencyController.RegisterRoutes(router)

	router.Route("/management", func(r chi.Router) {
		r.Handle("/health", infra.NewHealthHandler(db))
		r.Get("/info", infra.NewInfoHandler())
		r.Handle("/metrics", promhttp.Handler())
	})

	// Graceful shutdown + Eureka
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	eurekaClient := clients.NewEurekaClient(cfg.EurekaURL, "currency-service", cfg.HostName, cfg.ServerPort)

	if err := eurekaClient.Register(ctx); err != nil {
		log.Printf("failed to register on eureka, will retry in background: %v", err)
	}
	go eurekaClient.StartHeartbeatLoop(ctx, 30*time.Second)

	addr := ":" + cfg.ServerPort
	server := &http.Server{Addr: addr, Handler: router}

	go func() {
		log.Printf("currency-service listening on %s", addr)
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

	log.Println("currency-service stopped gracefully")
}
