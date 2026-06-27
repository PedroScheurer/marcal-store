# Bancos de dados — database-per-service

Cada microsserviço com persistência possui **sua própria instância PostgreSQL** no Docker.

| Container           | Banco        | Serviço          | Porta host |
|---------------------|--------------|------------------|------------|
| `postgres-auth`     | `db_user`    | auth-service     | 5434       |
| `postgres-product`  | `bd_product` | product-service  | 5435       |
| `postgres-currency` | `db_currency`| currency-service | 5436       |
| `postgres-order`    | `db_order`   | order-service    | 5437       |

**Credenciais (todos):** usuário `admin`, senha `admin`

## Migrações

| Serviço          | Ferramenta        | Local das migrations                          |
|------------------|-------------------|-----------------------------------------------|
| auth-service     | Flyway            | `auth-service/src/main/resources/db_migration` |
| currency-service | Flyway            | `currency-service/src/main/resources/db_migration` |
| order-service    | Flyway            | `order-service/src/main/resources/db_migration` |
| product-service  | golang-migrate    | `product-service/db_migration`                |

O schema é criado automaticamente na subida de cada serviço.

## Acesso local (fora do Docker)

```
jdbc:postgresql://localhost:5434/db_user      # auth
jdbc:postgresql://localhost:5435/bd_product # product
jdbc:postgresql://localhost:5436/db_currency # currency
jdbc:postgresql://localhost:5437/db_order   # order
```

PgAdmin: http://localhost:5050 (`admin@admin.com` / `admin`) — servidores pré-configurados.

## Admin (login no app)

| Campo | Valor |
|-------|--------|
| E-mail | `admin@admin.dev` |
| Senha | `admin123` |

Conta criada pelo Flyway (`V2__PopulateTableTbUser.sql`), com `type = 0` (Admin).
Após subir o auth-service, a migration `V3__ResetAdminPassword.sql` garante essa senha em bancos já existentes.

Login via gateway:

```bash
curl -X POST http://localhost:8765/auth/signin \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@admin.dev","password":"admin123"}'
```

## Testes manuais rápidos

```bash
# Subir tudo
docker compose up -d --build

# Gateway (ponto de entrada do app)
curl http://localhost:8765/actuator/health
curl "http://localhost:8765/products?targetCurrency=BRL&page=0&size=5"

# Cadastro + login
curl -X POST http://localhost:8765/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"name":"Teste","email":"teste@test.com","password":"senha123"}'

curl -X POST http://localhost:8765/auth/signin \
  -H "Content-Type: application/json" \
  -d '{"email":"teste@test.com","password":"senha123"}'
```

**Frontend:** `cd marcal-store-frontend && npx expo start --port 8083`

| URL | Descrição |
|-----|-----------|
| http://localhost:8765 | API Gateway |
| http://localhost:8761 | Eureka Dashboard |
| http://localhost:5050 | PgAdmin |
