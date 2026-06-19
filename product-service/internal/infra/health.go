package infra

import (
	"context"
	"net/http"
	"time"

	"github.com/hellofresh/health-go/v5"
	"github.com/jmoiron/sqlx"
)

// NewHealthHandler configura os health checks equivalentes ao Spring Actuator.
func NewHealthHandler(db *sqlx.DB) http.Handler {
	// 1. Cria a instância base do Health (similar ao endpoint /actuator/health)
	h, _ := health.New(
		health.WithComponent(health.Component{
			Name:    "product-service",
			Version: "1.0.0",
		}),
		health.WithSystemInfo(), // Inclui uso de memória e tempo de atividade (uptime)
	)

	// 2. Registra o check do Banco de Dados (Equivalente ao db health indicator do Actuator)
	// Ele vai rodar um Ping() no banco a cada requisição ou usar o cache configurado.
	_ = h.Register(health.Config{
		Name:      "postgres",
		Timeout:   2 * time.Second,
		SkipOnErr: false, // Se o banco falhar, o status geral do app vira DOWN (503)

		// Em vez de usar postgres.New(postgres.Config{...}), usamos a função nativa do Go:
		Check: func(ctx context.Context) error {
			// O db.PingContext garante que o banco está vivo e responde dentro do timeout do contexto
			return db.PingContext(ctx)
		},
	})

	// 3. Exemplo de um Custom Check (Caso queira validar outra dependência na mão)
	// No Java você faria isso implementando a interface `HealthIndicator`.
	_ = h.Register(health.Config{
		Name:      "custom-cache-check",
		Timeout:   1 * time.Second,
		SkipOnErr: true, // Se o cache falhar, não derruba o app inteiro (continua UP)
		Check: func(ctx context.Context) error {
			// Aqui você colocaria uma lógica rápida, ex: testar se a memória está estourada
			return nil
		},
	})

	// Retorna o Handler HTTP nativo pronto para ser usado no Chi Router
	return h.Handler()
}
