# auth-service
O auth-service é um microsserviço de autenticação e gerenciamento de usuários desenvolvido em Java 25 utilizando o ecossistema Spring Boot 4. Ele é responsável por registrar novos usuários, efetuar login seguro por meio de tokens JWT (JSON Web Token), validar regras de negócio e integrar-se a um ecossistema de microsserviços via Netflix Eureka Client.

## 🛠️ Tecnologias e Dependências Principais
A arquitetura do projeto utiliza o gerenciador de dependências Maven e é composta pelas seguintes tecnologias:

Java 25 & Spring Boot 4.0.6: Base do ecossistema e ambiente de execução.

Spring Security: Orquestração e proteção de endpoints, além de gerenciamento do ciclo de vida de autenticação.

JSON Web Tokens (JJWT 0.13.0): Geração e validação de tokens compactos e seguros para autenticação stateless.

Spring Data JPA & PostgreSQL Driver: Abstração da camada de persistência e driver para banco de dados relacional.

Flyway Database Migrations: Evolução controlada e automatizada do esquema do banco de dados (flyway-database-postgresql).

Spring Cloud Netflix Eureka Client: Registro e descoberta dinâmica do microsserviço dentro do cluster.

Spring Boot Actuator: Exposição de métricas operacionais, dados de saúde do app (/health) e informações do sistema.

## 📂 Estrutura de Diretórios do Projeto
O projeto segue um design arquitetural baseado em Camadas Autocontidas, segregando responsabilidades de forma clara e isolada:

```
auth-service/
├── .mvn/
├── src/
│   └── main/
│       └── java/
│           └── br/edu/atitus/authservice/
│               ├── components/          # Utilitários globais (JwtUtil, Validator)
│               ├── configs/             # Classes de configuração (Security, Beans)
│               ├── controllers/         # Endpoints REST expostos (AuthController)
│               ├── dtos/                # Objetos de transferência de dados (Sign-in/up DTOs)
│               ├── entities/            # Entidades mapeadas para o banco (UserEntity)
│               ├── infrastructure/      # Exceções globais e adapters técnicos
│               ├── repositories/        # Interfaces de persistência (UserRepository)
│               └── services/            # Camada de regras de negócio (UserService)
├── pom.xml                              # Definição de dependências e build do Maven
└── README.md                            # Documentação do sistema
```
## 📡 Endpoints da API & Exemplos de cURL
Abaixo estão detalhadas as rotas expostas pela camada de controle para o fluxo de autenticação.

### 1. Registro de Usuário (Sign-up)
Cadastra um novo usuário comum no sistema. A senha é automaticamente criptografada via PasswordEncoder antes de persistir no banco.
```
URL: /auth/signup

Método HTTP: POST

Content-Type: application/json
```
Exemplo de Requisição (cURL):
```
Bash
curl -X POST http://localhost:8080/auth/signup \
-H "Content-Type: application/json" \
-d '{
"name": "Pedro Scheurer",
"email": "pedro@atitus.edu.br",
"password": "senhaSegura123"
}'
Exemplo de Resposta (Mapeamento UserEntity - 201 Created):
```
```
JSON
{
"id": 1,
"name": "Pedro Scheurer",
"email": "pedro@atitus.edu.br",
"type": "Common"
}
```
### 2. Autenticação (Sign-in)
Valida as credenciais enviadas contra o banco de dados. Caso válidas, gera um token JWT contendo as declarações (claims) do usuário.
```
URL: /auth/signin

Método HTTP: POST

Content-Type: application/json
```
```
Exemplo de Requisição (cURL):

Bash
curl -X POST http://localhost:8080/auth/signin \
-H "Content-Type: application/json" \
-d '{
"email": "pedro@atitus.edu.br",
"password": "senhaSegura123"
}'
```

Exemplo de Resposta (SigninResponseDTO - 200 OK):
```
JSON
{
"user": {
"id": 1,
"name": "Pedro Scheurer",
"email": "pedro@atitus.edu.br",
"type": "Common"
},
"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJwZWRyb0BhdGl0dXMuZWR1LmJyIiwiaWQiOjEsInR5cGUiOiJDb21tb24ifQ..."
}
```
### 3. Monitoramento Operacional (Actuator)
Endpoints utilizados para checar a saúde interna do microsserviço, telemetria e coleta de métricas pelo Prometheus.

```
Saúde do App: GET /actuator/health

Métricas do Sistema: GET /actuator/metrics
```
## 🛡️ Regras de Validação de Negócio (UserService)
A camada de serviço executa validações rigorosas antes de persistir qualquer modificação no banco de dados:

Validação de Nome: Rejeita strings nulas ou vazias.

Validação de E-mail: Rejeita formatos fora do padrão corporativo/padrão de internet (validado pela classe interna Validator).

Validação de Senha: Exige uma complexidade mínima de, no mínimo, 6 caracteres.

Unicidade de Identidade:

Ao criar um novo registro, impede cadastros duplicados com o mesmo endereço de e-mail.

Ao atualizar um usuário existente, valida se o e-mail informado já não pertence a outro ID.

## 📝 Documentação da API (Swagger / OpenApi)
A documentação interativa da API é injetada em tempo de execução. Com o microsserviço ativo localmente, você pode acessar a interface gráfica do Swagger, explorar os esquemas e testar os payloads diretamente pelo navegador:

Interface do Swagger UI: http://localhost:8080/swagger-ui/index.html

OpenAPI Specs (JSON): http://localhost:8080/v3/api-docs