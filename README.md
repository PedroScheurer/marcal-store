Este repositório centraliza e orquestra uma arquitetura distribuída de microsserviços voltada para um e-commerce de cursos. O ecossistema combina o poder do ecossistema Spring Cloud (Java) com a alta performance reativa de componentes em Go 1.24, integrando comunicação assíncrona, tolerância a falhas, segurança centralizada e persistência poliglota isolada por serviço.

## 🗺️ Arquitetura do Ecossistema
O fluxo de dados no cluster opera sob o padrão de API Gateway e Service Discovery, estruturado da seguinte forma:

```
                               [ CLIENTES EXTERNOS ]
                                         │
                                         ▼ (Java - 8765)
                                ┌─────────────────┐
                                │ gateway-service │
                                └────────┬────────┘
                                         │
         ┌───────────────────────────────┼───────────────────────────────┐
         ▼ (lb://auth-service)           ▼ (lb://product-service)        ▼ (lb://order-service)
┌──────────────────┐           ┌──────────────────┐           ┌──────────────────┐
│   auth-service   │           │ product-service  │           │  order-service   │
│   (Java - 8900)  │           │   (Go - 8082)    │           │   (Java - 8200)  │
└──────────────────┘           └────────┬─────────┘           └────────┬─────────┘
                                        │                              │
                                        │                              │ (OpenFeign)
                                        │     ┌────────────────────────┘
                                        ▼     ▼
                               ┌──────────────────┐
                               │ currency-service │
                               │   (Java - 8081)  │
                               └──────────────────┘
```
### 📡 Como os Serviços se Comunicam
Borda Unificada (API Gateway): O gateway-service é a única porta aberta para o mundo exterior. Ele intercepta as chamadas, valida os claims do token JWT (usando JwtUtil) e despacha as requisições para a malha interna usando nomes lógicos (lb://).

Service Discovery (Netflix Eureka): O discovery-service atua como o catálogo telefônico do sistema. Todos os microsserviços (incluindo o product-service em Go) se registram e enviam batimentos cardíacos (heartbeats) para ele.

Comunicação Inter-Serviços (HTTP/REST Declarativo):

O order-service consome dados do product-service e currency-service via interfaces Spring Cloud OpenFeign.

O product-service (Go) consome taxas cambiais do currency-service usando um cliente HTTP customizado com políticas nativas de Retry e Circuit Breaker (sony/gobreaker).

## 📦 Divisão e Portas dos Módulos
O parque tecnológico é composto por infraestruturas de suporte e serviços de domínio distribuídos nas seguintes portas locais:

### Componentes de Infraestrutura e Suporte
discovery-service (Porta 8761): Servidor central do Netflix Eureka para registro e descoberta de instâncias.

config-service (Porta 8888): Servidor de configuração centralizado Spring Cloud Config.

postgres (Porta 5433 externo / 5432 interno): Instância isolada do PostgreSQL 16 com base de dados particionada para os microsserviços.

pgadmin (Porta 5050): Interface gráfica para administração das tabelas SQL.

### Serviços de Negócio (Core System)
gateway-service (Porta 8765): Ponto de entrada reativo (WebFlux) com filtro de segurança JWT.

auth-service (Porta 8900): Emissor de credenciais e gestão de usuários (db_user).

product-service (Porta 8082): Catálogo escrito em Go 1.24 com cache LRU em memória e migrações automatizadas (bd_product).

currency-service (Porta 8081): Conversor de cotações monetárias em tempo real (db_currency).

order-service (Porta 8200): Processador de pedidos com conversão cambial histórica de compras (db_order).

greeting-service (Porta 8080): Serviço demonstrativo para validação de entrega de propriedades via Config Server.

## 🛡️ Padrões de Resiliência e Segurança Aplicados
Segurança Centralizada (Edge Security): A validação e descriptografia do token JWT ocorre estritamente no Gateway. Serviços internos recebem os dados do usuário já mastigados via Headers HTTP customizados (X-User-Id, X-User-Email, X-User-Type), eliminando redundância de código de segurança nos microsserviços de negócio.

Isolamento de Estado (Database-per-Service): Nenhum microsserviço acessa a tabela de outro. Embora compartilhem o mesmo container Docker postgres, cada um possui seu banco de dados lógico estritamente isolado (db_user, db_currency, db_order, bd_product).

Circuit Breaker & Fallback: Implementado na comunicação crítica de cotações. Se o currency-service falhar ou apresentar alta latência, o circuito abre e o sistema entra em modo de degradação suave (fallback), servindo valores locais em cache para não interromper a experiência do usuário.

## 🚀 Como Inicializar o Ecossistema Completo
### Pré-requisitos
Docker e Docker Compose instalados na máquina.

Ter executado o build/empacotamento dos artefatos de cada pasta (gerando os arquivos .jar ou binários do Go executados dentro dos respectivos Dockerfile).

### Inicialização via Linha de Comando
Na raiz do projeto (onde se encontra o arquivo docker-compose.yml), execute o comando abaixo para subir toda a malha de microsserviços respeitando a ordem correta de inicialização (healthchecks):

Bash
docker compose up -d --build
### Verificação da Saúde do Cluster
Após alguns instantes, você poderá auditar o estado de registro das instâncias acessando diretamente o dashboard do Eureka:

URL do Eureka: http://localhost:8761

Para conferir os logs ou o status detalhado de saúde via terminal, utilize:

Bash
docker compose ps
docker compose logs -f gateway-service
