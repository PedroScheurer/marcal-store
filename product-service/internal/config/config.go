package config

import (
	"fmt"
	"os"
	"strconv"
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

	UploadDir   string
	MaxVideoMB  int64
	MaxImageMB  int64
}

// Load lê a configuração das variáveis de ambiente, com defaults
// equivalentes aos valores fixos do application.yaml
// (server.port: 8082, spring.datasource.url: .../bd_product).
func Load() Config {
	return Config{
		ServerPort:      getEnv("SERVER_PORT", "8082"),
		ApplicationName: getEnv("APP_NAME", "product-service"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBName:     getEnv("DB_NAME", "bd_product"),
		DBUser:     getEnv("POSTGRES_USER", ""),
		DBPassword: getEnv("POSTGRES_PASSWORD", ""),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		HostName:  getEnv("HOST_NAME", "localhost"),
		EurekaURL: getEnv("EUREKA_URL", "http://localhost:8761/eureka"),

		UploadDir:  getEnv("UPLOAD_DIR", "./uploads"),
		MaxVideoMB: getEnvInt64("MAX_VIDEO_MB", 100),
		MaxImageMB: getEnvInt64("MAX_IMAGE_MB", 10),
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

func getEnvInt64(key string, fallback int64) int64 {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}
	v, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return fallback
	}
	return v
}
