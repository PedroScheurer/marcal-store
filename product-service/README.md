# product-service (Go) — Etapa 1: Controllers, Services, Repositories

Conversão do `product-service` Java/Spring Boot para Go. Esta primeira etapa
cobre a camada de domínio (controllers, services, repositories), conforme
acordado. O client HTTP para o `currency-service`, o registro no Eureka e a
resiliência (retry/circuit breaker) ficam para as próximas etapas.

## Stack escolhida

- Roteamento HTTP: [chi](https://github.com/go-chi/chi) (v5.1.0, compatível com Go 1.22)
- Acesso a dados: `database/sql` + [sqlx](https://github.com/jmoiron/sqlx) (v1.3.5)
- Driver Postgres: [lib/pq](https://github.com/lib/pq)
- Go 1.22.2 (mesma versão disponível no ambiente onde isso foi escrito e testado)

> Nota sobre versões: o ambiente onde fiz a conversão só tinha acesso a
> `archive.ubuntu.com` (apt) e GitHub direto (sem `proxy.golang.org`), então
> fixei `chi v5.1.0` (a v5.3.0 mais recente exige Go 1.23+) e `sqlx v1.3.5`
> (versões mais novas trazem o driver `go-sql-driver/mysql` como dependência
> indireta, que por sua vez depende de `filippo.io/edwards25519`, inacessível
> nesse ambiente). No seu ambiente local, com acesso normal ao proxy do Go,
> você pode tentar atualizar essas versões com `go get -u` se quiser, mas elas
> funcionam perfeitamente como estão.

## Estrutura

```
cmd/product-service/main.go     equivalente ao ProductServiceApplication (monta tudo e inicia o servidor)
internal/config/                carrega configuração do ambiente (porta, dados de conexão do Postgres)
internal/controllers/           ProductController e WsProductController (rotas chi)
internal/services/              ProductService, WsProductService, CurrencyConversionService, CacheService
internal/repositories/          ProductRepository (interface) + implementação Postgres com sqlx
internal/entities/              ProductEntity (struct com tags `db`)
internal/dtos/                  ProductDTO, ProductInDTO, ProductOutDTO, ProductUpdateInDTO, Page[T], ErrorResponse
internal/apperrors/             erros customizados + tradução para respostas HTTP (equivalente ao @ControllerAdvice)
db_migration/                   scripts SQL originais (ainda não plugados a uma ferramenta de migração)
```

## Mapeamento Java → Go (visão rápida)

| Java | Go |
|---|---|
| `@RestController` + `@GetMapping`/`@PostMapping`/... | `chi.Router` + handlers em `internal/controllers` |
| `@Service` | struct simples em `internal/services`, injetada via construtor `New...` |
| `JpaRepository<ProductEntity, Long>` | interface `ProductRepository` + impl. `postgresProductRepository` com sqlx |
| `Optional<ProductEntity>` vazio | `(*entities.ProductEntity)(nil)` retornado pelo repository |
| `@ControllerAdvice` / `GlobalExceptionHandler` | `apperrors.WriteErrorResponse`, chamado explicitamente em cada handler |
| Exceções customizadas (`ProductNotFoundException` etc.) | tipos de erro em `internal/apperrors`, usando `errors.As` para despachar |
| Cache Caffeine (`spring.cache.caffeine.spec`) | `CacheService` próprio, com TTL + eviction LRU em memória |
| `Page<T>` do Spring Data | `dtos.Page[T]` (genérico), com os mesmos campos (`content`, `totalElements`, etc.) |
| `@RequestHeader` obrigatório | leitura manual do header + erro 401 se ausente/inválido |
| `@Value("${server.port}")` | `config.Config.ServerPort`, lido de env var com fallback `8082` |

## O que ainda falta (próximas etapas, por sua escolha)

1. **Client HTTP do currency-service**: hoje existe só a interface
   `services.CurrencyClient` e uma implementação `noopCurrencyClient`
   (sempre cai no fallback). Falta o client real fazendo a chamada HTTP.
2. **Registro no Eureka**: o serviço ainda não se registra em nenhum
   service discovery.
3. **Resiliência**: retry e circuit breaker (equivalentes ao
   `@Retry`/`@CircuitBreaker` do Resilience4j no `application.yaml`) serão
   adicionados quando implementarmos o client.
4. **Migrations**: os arquivos `.sql` foram copiados para `db_migration/`,
   mas ainda não há uma ferramenta (ex.: `golang-migrate`/`goose`) plugada
   para rodá-los automaticamente como o Flyway fazia no Java.
5. **Validação de payload de entrada**: o Java usa Bean Validation
   (`@Valid` + anotações) nos DTOs de entrada — isso ainda não foi
   replicado nos `ProductInDTO`/`ProductUpdateInDTO` em Go.

## Pontos de atenção / decisões a confirmar com você

- **`ErrorResponse.status`**: no Java o campo do record está escrito como
  `staus` (erro de digitação). Corrigi para `status` no JSON de saída do Go.
  Se algum client já depende do nome com erro de digitação, me avise para
  eu reverter.
- **Bug aparente em `WsProductService.createProduct`**: a checagem original
  em Java é `if (userType.equals(ADMIN_TYPE)) throw new AuthenticationException(...)`,
  ou seja, bloqueia justamente o usuário do tipo "admin" (0) de criar
  produtos, e libera qualquer outro tipo — o oposto do que os métodos de
  alterar/excluir fazem (que exigem `ADMIN_TYPE`). Mantive esse comportamento
  idêntico no Go (veja o comentário em `ws_product_service.go`), mas vale
  confirmar se é intencional.
- **DTO de entrada do PUT**: o Java reaproveita o record `ProductDTO` (de
  saída, com campos de conversão de moeda) como `@RequestBody` do PUT, mas só
  lê 5 campos dele. Criei um `ProductUpdateInDTO` próprio em Go, mais
  explícito sobre o que a API espera no corpo da requisição.

## Como rodar (assumindo Postgres já rodando)

```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_NAME=bd_product
export POSTGRES_USER=seu_usuario
export POSTGRES_PASSWORD=sua_senha
export SERVER_PORT=8082

go run ./cmd/product-service
```

## Testes

```bash
go test ./...
```

Hoje só há testes unitários para `CurrencyConversionService` (lógica pura,
sem banco). Testes de integração para repository/controllers podem ser
adicionados numa próxima etapa, se quiser.
