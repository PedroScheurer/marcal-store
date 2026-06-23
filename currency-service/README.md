# currency-service
O currency-service é um microsserviço de alta performance desenvolvido em Java 25 com Spring Boot 4. A sua principal responsabilidade é calcular taxas de conversão monetária entre diferentes moedas. O serviço utiliza uma estratégia híbrida: consome a API oficial do Banco Central do Brasil (BCB) em tempo real, gerencia um cache em memória de alta performance com Caffeine e possui um mecanismo de tolerância a falhas utilizando banco de dados local (PostgreSQL) como fallback, orquestrado por padrões de Circuit Breaker e Retry com o Resilience4j.

## 🛠️ Tecnologias e Dependências Principais
Java 25 & Spring Boot 4.0.5: Core do microsserviço.

Spring Cloud OpenFeign: Cliente HTTP declarativo utilizado para consumir a API externa do Banco Central do Brasil.

Resilience4j (Circuit Breaker & Retry): Implementação de resiliência para isolar falhas na API externa do BCB.

Caffeine Cache: Mecanismo de cache local síncrono para otimizar o tempo de resposta e poupar requisições externas.

Spring Data JPA & PostgreSQL: Camada de persistência relacional para armazenamento das cotações de contingência (fallback).

Flyway Migrations: Controle de versão automatizado do esquema de banco de dados (db_migration).

Spring Cloud Netflix Eureka Client: Registro e descoberta dinâmica do microsserviço no ecossistema de nuvem.

Spring Boot Actuator: Monitoramento operacional completo exposto via HTTP.

## 📂 Estrutura de Diretórios do Projeto
Conforme mapeado na árvore de componentes (visível em image_4284a9.png), o projeto adota o seguinte design de pacotes:

```
currency-service/
├── .mvn/
├── src/
│   └── main/
│       ├── java/
│       │   └── br/edu/atitus/currencyservice/
│       │       ├── clients/         # Interfaces OpenFeign (BCBClient)
│       │       ├── controllers/     # Endpoints expostos (CurrencyController)
│       │       ├── dtos/            # Records e estruturas de payload (CurrencyDTO)
│       │       ├── entities/        # Modelos relacionais (CurrencyEntity)
│       │       ├── infrastructure/  # Middlewares e Exceções customizadas
│       │       ├── repositories/    # Interfaces JPA (CurrencyRepository)
│       │       └── services/        # Orquestração e regras (CurrencyService, CacheService)
│       │       └── CurrencyServiceApplication.java
│       └── resources/
│           ├── db_migration/        # Scripts SQL gerenciados pelo Flyway
│           └── application.yml      # Configurações de propriedades do ambiente
```
## 📡 Endpoints da API & Exemplos de cURL
O microsserviço roda por padrão na porta 8081.

### 1. Conversão de Moedas
Realiza o cálculo da taxa de conversão cambial entre uma moeda de origem (source) e uma moeda de destino (target).
```
URL: /currency/convert

Método HTTP: GET

Parâmetros obrigatórios: source (Moeda Origem), target (Moeda Destino)
```
Exemplo de Requisição (cURL):
```
Bash
curl -X GET "http://localhost:8081/currency/convert?source=USD&target=BRL" \
-H "Accept: application/json"
```
Exemplo de Resposta (200 OK - Via Banco Central ou Cache):

```
JSON
{
"sourceCurrency": "USD",
"targetCurrency": "BRL",
"conversionRate": 5.25,
"environment": "Currency-service running on Port: 8081 | Banco Central do Brasil"
}
```

Exemplo de Resposta (200 OK - Cache Ativo):
```
JSON
{
"sourceCurrency": "USD",
"targetCurrency": "BRL",
"conversionRate": 5.25,
"environment": "Currency-service in cache"
}
```

Exemplo de Resposta (200 OK - Fallback acionado devido à indisponibilidade do BCB):
```
JSON
{
"sourceCurrency": "USD",
"targetCurrency": "BRL",
"conversionRate": 5.10,
"environment": "Currency-service fallback running on Port: 8081"
}
```
## ⚡ Resiliência, Cache e Elasticidade
### 🧠 Estratégia de Cache (Caffeine)
Para evitar sobrecarregar a API do Banco Central, o microsserviço armazena as cotações calculadas em memória com a seguinte política:

Tamanho Máximo: 500 registros.

Tempo de Expiração: 15 segundos após a escrita (expireAfterWrite=15s).

### 🛡️ Tolerância a Falhas (Resilience4j)
Se as requisições para o Banco Central começarem a falhar ou apresentarem lentidão (Timeout de Conexão: 1s / Leitura: 3s), o sistema se protege de forma autônoma:

Mecanismo de Retry (Retry_BCBClient_getConversionRate): Realiza até 3 tentativas automáticas com recuo exponencial multiplicador de 2.0 antes de desistir da chamada externa.

Circuit Breaker (BCBClientgetConversionRateStringString): Se a taxa de falha atingir 50% em uma janela de 10 chamadas, o circuito se abre por 30 segundos, redirecionando instantaneamente todas as requisições subsequentes para a tabela de contingência no PostgreSQL (CurrencyEntity) sem gastar tempo tentando acessar a rede.

## 📊 Monitoramento e Saúde (Actuator)
O ecossistema expõe dados operacionais vitais para ferramentas de telemetria, disponíveis nos seguintes endpoints:

Métricas Gerais: GET http://localhost:8081/actuator/metrics

Painel de Saúde Completo: GET http://localhost:8081/actuator/health

Eventos do Circuit Breaker: GET http://localhost:8081/actuator/circuit-breaker-events

## ⚙️ Como Executar o Projeto Localmente
Certifique-se de possuir o Java 25 instalado.

Garanta que uma instância do PostgreSQL esteja rodando na porta 5432 com o banco db_currency.

Configure as variáveis de ambiente necessárias para o banco:

```
export POSTGRES_USER=seu_usuario
export POSTGRES_PASSWORD=sua_senha
```
4. Execute o microsserviço através do Maven:
```
./mvnw spring-boot:run
``` 