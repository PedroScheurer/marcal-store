# order-service
O order-service é um microsserviço crítico do ecossistema, desenvolvido em Java 21 e baseado no Spring Boot 3.5.0. Sua principal atribuição é gerenciar o ciclo de vida e a persistência de pedidos de compras. Ele consome dados assincronamente de outros microsserviços do cluster por meio de clientes OpenFeign para calcular precificações em tempo real com conversões monetárias dinâmicas.

## 🛠️ Tecnologias e Dependências Principais
Java 21 & Spring Boot 3.5.0: Ambientes de execução e core do microsserviço.

Spring Cloud OpenFeign: Cliente HTTP declarativo usado para comunicação síncrona inter-service (consultas ao product-service e currency-service).

Spring Data JPA & PostgreSQL: Framework de persistência relacional e driver de banco de dados para armazenamento de ordens e seus respectivos itens.

Flyway Core: Evolução incremental controlada do esquema do banco de dados relacional através de scripts versionados (db_migration).

Spring Cloud Netflix Eureka Client: Mecanismo de registro e autodescoberta do microsserviço no cluster na porta 8200.

Spring Boot Actuator: Monitoramento em tempo real do estado de saúde e telemetria da aplicação.

## 📂 Estrutura de Diretórios do Projeto
Seguindo o mapeamento de pacotes visualizado na estrutura da aplicação (visível em image_4536e3.png), o projeto é estruturado da seguinte forma:

Plaintext
order-service/
├── .mvn/
├── src/
│   └── main/
│       ├── java/
│       │   └── br/edu/atitus/order_service/
│       │       ├── clients/         # Interfaces OpenFeign (ProductClient, CurrencyClient)
│       │       ├── controllers/     # Endpoints HTTP expostos (OrderController)
│       │       ├── dtos/            # Registros de transporte de dados (OrderDTO, OrderItemDTO)
│       │       ├── entities/        # Entidades JPA (OrderEntity, OrderItemEntity)
│       │       ├── repositories/    # Camada de acesso a dados (OrderRepository)
│       │       └── services/        # Regras e cálculos de precificação (OrderService)
│       │       └── OrderServiceApplication.java
│       └── resources/
│           ├── db_migration/        # Scripts SQL de migração da base de dados
│           └── application.properties # Parâmetros do banco de dados, Eureka e Actuator
## 📡 Endpoints da API & Exemplos de cURL
O microsserviço roda nativamente na porta 8200. Os cabeçalhos de contexto de usuário (X-User-*) são injetados de forma transparente pelo API Gateway após a validação do Token JWT.

### 1. Criação de Novo Pedido (Cálculo em Tempo Real)
Cria um pedido associado ao usuário conectado. Para cada item enviado, o serviço consulta o catálogo de produtos e converte o preço para a moeda padrão (USD).

URL: /ws/orders

Método HTTP: POST

Content-Type: application/json

Exemplo de Requisição (cURL):

Bash
curl -X POST http://localhost:8200/ws/orders \
-H "Content-Type: application/json" \
-H "X-User-Id: 1024" \
-H "X-User-Email: usuario@atitus.edu.br" \
-H "X-User-Type: 1" \
-d '{
"items": [
{ "productId": 15, "quantity": 2 },
{ "productId": 22, "quantity": 1 }
]
}'
### 2. Listagem Paginada de Pedidos com Conversão Dinâmica
Retorna o histórico de compras do usuário conectado. O grande diferencial deste endpoint é a conversão em tempo de execução de todo o montante histórico para uma moeda alvo informada no parâmetro (targetCurrency).

URL: /ws/orders

Método HTTP: GET

Parâmetros de URL: targetCurrency (ex: BRL, EUR), page, size

Exemplo de Requisição (cURL):

Bash
curl -X GET "http://localhost:8200/ws/orders?targetCurrency=BRL&page=0&size=5" \
-H "X-User-Id: 1024" \
-H "X-User-Email: usuario@atitus.edu.br" \
-H "X-User-Type: 1"
## 🧠 Fluxo de Regras de Negócio (OrderService)
O processamento interno do serviço orquestra as seguintes etapas ao gerar uma ordem:

Captura do Contexto: O ID do cliente é amarrado à ordem de forma automática através dos cabeçalhos repassados pelo Gateway.

Consulta de Catálogo (product-service): O microsserviço bate no catálogo remoto para capturar a descrição, integridade e o valor real do preço do produto no momento exato da compra (priceAtPurchase).

Conversão Cambial Dinâmica (currency-service): O serviço calcula dinamicamente o valor acumulado e os preços convertidos par a par utilizando as cotações mais recentes fornecidas pelo microsserviço financeiro.

Fechamento e Persistência: Salva os totais originais e os totais convertidos no PostgreSQL para fins de auditoria financeira histórica.

## 📈 Monitoramento e Gestão (Actuator)
O monitoramento operacional está completamente aberto para raspagem de logs e auditoria interna através dos endpoints expostos:

Saúde das Conexões e Probes: GET http://localhost:8200/actuator/health

Métricas Gerais do Spring: GET http://localhost:8200/actuator/metrics

## ⚙️ Inicialização Local
Instale o Java 21.

Crie uma base PostgreSQL local chamada db_order com credenciais padrão postgres/postgres.

Certifique-se de que o Eureka Server (discovery-service) esteja operacional na porta 8761.

Inicie o microsserviço:

Bash
./mvnw spring-boot:run