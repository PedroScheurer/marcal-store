# ─── Estágio 1: build ────────────────────────────────────────────
FROM golang:1.26.4 AS builder

WORKDIR /app

# Copia só os arquivos de dependência primeiro (melhor uso de cache Docker)
COPY go.mod go.sum ./
RUN go mod download

# Agora copia o restante do código
COPY . .

# Compila o binário estático pra Linux
RUN CGO_ENABLED=0 GOOS=linux go build -o currency-service-go ./cmd/currency-service

# ─── Estágio 2: imagem final ─────────────────────────────────────
FROM alpine:3.20

# Certificados TLS (necessário pra chamadas HTTPS futuras)
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copia só o binário do estágio anterior
COPY --from=builder /app/currency-service-go .

# Copia as migrations (o binário lê esses arquivos em runtime)
COPY --from=builder /app/db_migration ./db_migration

EXPOSE 8081

CMD ["./currency-service-go"]