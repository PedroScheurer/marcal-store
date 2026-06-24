# product-service (Go Port)
O product-service é o microsserviço de catálogo e gerenciamento de produtos do ecossistema, reescrito em Go 1.24 para atingir performance extrema, baixíssima pegada de memória e tempos de inicialização instantâneos.

A aplicação espelha rigorosamente as regras de negócio e contratos de API de sua contraparte original em Java/Spring Boot. Ela se conecta dinamicamente ao ecossistema micro-controlado por meio de um cliente nativo do Netflix Eureka, consome taxas de câmbio dinâmicas do currency-service, expõe telemetria para o Prometheus e implementa estratégias robustas de tolerância a falhas utilizando uma variação do algoritmo de Circuit Breaker (através do gobreaker).

## 🛠️ Tecnologias e Dependências Principais
Go 1.24.0: Core da aplicação utilizando os recursos mais recentes de concorrência e tipagem genérica da linguagem.

Go-Chi Router (v5): Roteador HTTP idiomático, leve e extremamente rápido, ideal para a construção de REST APIs modulares.

Sqlx & Lib/PQ: Extensão da biblioteca padrão database/sql para mapeamento relacional fluído e de alta performance acoplado ao PostgreSQL.

Go-Migrate: Mecanismo de migração de banco de dados nativo executado de forma programática durante o bootstrap da aplicação (db_migration).

Sony/Gobreaker & Avast/Retry-Go: Combinação cirúrgica para resiliência na comunicação de rede, implementando padrões estáveis de Circuit Breaker e Retry com recuo exponencial.

HelloFresh/Health-Go: Engine para monitoramento de saúde de dependências físicas (como conectividade de banco de dados).

Prometheus Client Golang: Instrumentação nativa e exposição de métricas de runtime do Go no padrão OpenMetrics.

## 📂 Estrutura de Diretórios do Projeto
Mapeado de acordo com a arquitetura limpa de pacotes padrão Go e a árvore exibida em image_453ac4.png:

```
product-service/
├── cmd/
│   └── product-service/
│       └── main.go              # Ponto de entrada (Montagem manual do grafo de injeção de dependência)
├── internal/
│   ├── apperrors/               # Centralizador de handlers globais e mapeamento de Status HTTP
│   ├── clients/                 # Clientes HTTP (OpenFeign Replacement para Eureka e CurrencyClient)
│   ├── config/                  # Estruturas de leitura de variáveis de ambiente do SO (Load)
│   ├── controllers/             # Camada de Handlers HTTP (ProductController, WsProductController)
│   ├── dtos/                    # Estruturas de dados e Records Go (ProductDTO, Page[T])
│   ├── entities/                # Structs de mapeamento de tabelas relacionais do banco
│   ├── infra/                   # Inicializadores de infra (Migrations, HealthCheckers)
│   ├── repositories/            # Camada de persistência SQL pura
│   └── services/                # Regras de negócio e Cache LRU In-Memory concorrente
├── go.mod                       # Definição de dependências do módulo
└── go.sum                       # Somas de verificação de integridade
```

## 📡 Endpoints da API & Exemplos de cURL
A aplicação roda nativamente na porta 8082 (ou na variável de ambiente SERVER_PORT).

### 1. Consulta de Produto com Conversão Monetária Online
Recupera os detalhes do produto e injeta o preço convertido com base na moeda requisitada.
```
URL: /products/{idProduct}

Método HTTP: GET

Query Params Obrigatórios: targetCurrency
```

Exemplo de Requisição (cURL):
```
Bash
curl -X GET "http://localhost:8082/products/1?targetCurrency=BRL" \
-H "Accept: application/json"
Exemplo de Resposta (200 OK):

JSON
{
"id": 1,
"name": "Arquitetura de Microsserviços com Go",
"instructor": "Pedro Konig Scheurer",
"imageUrl": "http://...",
"videoUrl": "http://...",
"description": "Curso avançado",
"workload": 40,
"modules": 5,
"price": 10.00,
"currency": "USD",
"convertedPrice": 52.50,
"requestedCurrency": "BRL",
"environment": "Product-service running on Port: 8082 | Currency-service running on Port: 8081 | Banco Central do Brasil"
}
```
### 2. Criação de Produto (Área Administrativa)
Garante que apenas usuários do tipo ADMIN (código 0) possam persistir novos produtos. Os cabeçalhos X-User-* são repassados e validados na borda pelo API Gateway.
```
URL: /ws/products

Método HTTP: POST
```

Exemplo de Requisição (cURL):
```
Bash
curl -X POST http://localhost:8082/ws/products \
-H "Content-Type: application/json" \
-H "X-User-Id: 55" \
-H "X-User-Email: admin@atitus.edu.br" \
-H "X-User-Type: 0" \
-d '{
"name": "Go para Iniciantes",
"instructor": "Pedro Scheurer",
"description": "Do zero ao profissional",
"workload": 20,
"modules": 3,
"price": 100.0,
"currency": "BRL"
}'
```
## 🧠 Engenharia de Resiliência & Mecanismos Customizados
### 🧊 Cache LRU Concorrente (CacheService)
Para emular perfeitamente a especificação do Caffeine Cache utilizado na stack Java (maximumSize=500,expireAfterWrite=15s), foi codificado um serviço de cache nativo em Go utilizando uma lista duplamente encadeada (container/list) protegida por travas de exclusão mútua (sync.Mutex).

Eviction automatizada baseada em LRU (Least Recently Used) quando o tamanho ultrapassa 500 itens.

Controle de TTL estrito por entrada para expiração de valores cambiais voláteis após 15 segundos.

### 🔌 Roteamento e Paginação Idiomáticos (Spring Data Emulation)
Como o Go-Chi não possui resoluções de assinaturas automáticas de paginação como o Spring Data (Pageable), o pacote controllers estende um parser customizado (parsePageable) que extrai de forma transparente parâmetros como ?page=0&size=5&sort=price,desc e encapsula os dados em um wrapper genérico Page[T] compatível com o formato JSON enviado pelo ecossistema Java original.

## 📊 Observabilidade e Monitoramento
A porta de gerenciamento disponibiliza as seguintes interfaces de telemetria expostas sob o path /management:
```
HealthCheck Unificado: GET http://localhost:8082/management/health

Métricas do Prometheus (Runtime Go & GC): GET http://localhost:8082/management/metrics

Informações do Sistema: GET http://localhost:8082/management/info
```
## ⚙️ Como Executar a Aplicação
Tenha o Go 1.24+ configurado em sua máquina.

Certifique-se de expor as variáveis de ambiente necessárias ou use as configurações padrão:

```
export DB_HOST=localhost
export POSTGRES_USER=postgres
export POSTGRES_PASSWORD=postgres
```
```
Execute o bootstrap do serviço (as migrações do banco via Flyway equivalente serão aplicadas instantaneamente):

go run cmd/product-service/main.go
```

## 📡 Descoberta de Serviços e Clientes Reativos (/internal/clients)
Para operar de forma transparente em uma arquitetura de microsserviços originalmente desenhada para Java/Spring, a camada de comunicação foi reescrita utilizando padrões nativos em Go, implementando o ciclo de vida completo de Descoberta de Serviços e Resiliência de Rede:

```
[ product-service (Go) ]
│
├───► (1) ResolveURL() ──► [ Service Discovery (Cache Local TTL 20s) ] ──► [ Eureka Server ]
│
└───► (2) GetCurrency() ──► [ Circuit Breaker & Retry ] ──► HTTP GET ──► [ currency-service ]
```
### 1. Registro e Ciclo de Vida no Eureka (EurekaClient)
Ao contrário do Java, onde a anotação @EnableDiscoveryClient oculta a complexidade, em Go o ciclo de vida REST do Netflix Eureka é gerenciado de forma explícita através de goroutines:

Bootstrap (Register): Dispara um payload POST /eureka/apps/PRODUCT-SERVICE mapeando os metadados da instância, URLs de saúde e portas no padrão esperado pelo Spring Cloud.

Renovação de Lease (StartHeartbeatLoop): Uma goroutine executa requisições PUT a cada 30 segundos. Caso o Eureka responda 404 Not Found (indicação de que a instância sofreu eviction por timeout), o cliente reconstrói o estado e efetua o re-registro de forma autônoma.

Shutdown Gracioso (Deregister): Intercepta sinais do sistema operacional (SIGTERM e os.Interrupt), disparando um DELETE para remover imediatamente o nó do balanceador, evitando o roteamento de requisições fantasmas (zombies).

### 2. Balanceamento Reativo e Resolução Dinâmica (ServiceDiscovery)
Para emular o comportamento do Spring Cloud LoadBalancer, foi codificado um componente de descoberta com cache síncrono local por TTL (20 segundos) protegido por sync.Mutex:

Evita sobrecarregar o Eureka com chamadas HTTP a cada requisição de negócio.

Estratégia de Resiliência Passiva: Se o servidor do Eureka estiver temporariamente fora do ar, o ServiceDiscovery intercepta o erro de rede e continua servindo os endereços antigos já cacheados (mesmo que expirados), preferindo uma rota potencialmente desatualizada a quebrar o fluxo completo do cliente.

### 3. Resiliência de Rede Baseada em Contratos (HTTPCurrencyClient)
A implementação substitui os componentes declarativos do @FeignClient e @Retry/@CircuitBreaker do Resilience4j acoplando duas bibliotecas consolidadas do Go:

#### 🔄 Política de Retry (avast/retry-go)
Envolve a chamada de rede executando até 3 tentativas com recuo exponencial e fator multiplicador (Exponential Backoff Delay).

Filtro de Exceções: Avalia o erro de forma inteligente. Erros de infraestrutura (como Timeout ou Connection Refused) disparam novas tentativas, enquanto erros de negócio (404 Not Found da moeda) abortam o laço imediatamente, poupando recursos de rede.

#### 🔌 Política de Circuit Breaker (sony/gobreaker)
Usa uma janela móvel para monitorar as falhas de rede do currency-service. O circuito é configurado de forma equivalente às definições do seu application.yml original:

Métricas: Exige um mínimo de 5 chamadas para análise; se a taxa de falha for superior a 50%, o circuito muda para o estado Aberto.

Tratamento de Erros e Fallback: Conforme o contrato especificado para interrupções, se o circuito estiver em estado Aberto (gobreaker.ErrOpenState), a falha é suprimida retornando (nil, nil), acionando de forma limpa a camada de banco de dados local do seu CurrencyConversionService (fluxo de contingência).
