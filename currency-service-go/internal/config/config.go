package config

import (
	"fmt"
	"os"
)

// Config concentra as configurações lidas do ambiente, equivalente
// (por agora) à parte de server.port e spring.datasource do
// application.yaml. As seções de Eureka e Resilience4j serão
// adicionadas quando implementarmos o client e o registro no Eureka.
type Config struct {
	ServerPort      string
	ApplicationName string

	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
	DBSSLMode  string

	EurekaURL string
	HostName  string
}

// Load lê a configuração das variáveis de ambiente, com defaults
// equivalentes aos valores fixos do application.yaml
// (server.port: 8082, spring.datasource.url: .../bd_product).
func Load() Config {
	return Config{
		ServerPort:      getEnv("SERVER_PORT", "8081"),
		ApplicationName: getEnv("APP_NAME", "currency-service"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBName:     getEnv("DB_NAME", "db_currency"),
		DBUser:     getEnv("POSTGRES_USER", ""),
		DBPassword: getEnv("POSTGRES_PASSWORD", ""),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		HostName:  getEnv("HOST_NAME", "localhost"),
		EurekaURL: getEnv("EUREKA_URL", "http://localhost:8761/eureka"),
	}
}

// PostgresDSN monta a connection string no formato esperado pelo
// driver lib/pq, equivalente à
// spring.datasource.url: jdbc:postgresql://localhost:5432/bd_product
// combinada com username/password.
func (c Config) PostgresDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBName, c.DBUser, c.DBPassword, c.DBSSLMode,
	)
}

func (c Config) PostgresMigrationUrl() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName, c.DBSSLMode)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
