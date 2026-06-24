# gateway-service
- O gateway-service é a camada de borda (API Gateway) unificada do ecossistema de microsserviços, desenvolvida em Java 25 e alimentada por Spring Boot 4.0.6. 
- Utilizando o Spring Cloud Gateway (WebFlux), ele atua como o ponto único de entrada para todas as requisições de clientes externos, sendo responsável pelo roteamento dinâmico inteligente, balanceamento de carga reativo e validação centralizada de segurança através de tokens JWT.

## 🛠️ Tecnologias e Dependências Principais
<table class="tech-table">
    <thead>
        <tr>
            <th>Tecnologia / Dependência</th>
            <th>Descrição e Função no Ecossistema</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td><strong>Java 25 & Spring Boot 4.0.6</strong></td>
            <td>Base reativa de altíssima performance para o ciclo de vida da aplicação.</td>
        </tr>
        <tr>
            <td><strong>Spring Cloud Gateway Server (WebFlux)</strong></td>
            <td>Mecanismo assíncrono e não-bloqueante para roteamento de fluxos baseado em predicados e filtros.</td>
        </tr>
        <tr>
            <td><strong>JSON Web Tokens (JJWT 0.13.0)</strong></td>
            <td>Componente integrado para decodificação, checagem de integridade e assinatura criptográfica de claims (<code>jjwt-api</code>, <code>jjwt-impl</code>, <code>jjwt-jackson</code>).</td>
        </tr>
        <tr>
            <td><strong>Spring Cloud Netflix Eureka Client</strong></td>
            <td>Integração ao servidor de <em>Service Discovery</em> para permitir a resolução de rotas via nomes lógicos de serviços com balanceamento de carga ativo (<code>lb://</code>).</td>
        </tr>
        <tr>
            <td><strong>Spring Boot Actuator</strong></td>
            <td>Exposição de métricas internas, status das rotas e monitoramento de saúde em tempo real.</td>
        </tr>
    </tbody>
</table>

## 📂 Estrutura de Diretórios do Projeto
Baseado na organização estrutural visível na árvore do projeto, o microsserviço divide suas responsabilidades da seguinte forma:
```
gateway-service/
├── .mvn/
├── src/
│   └── main/
│       ├── java/
│       │   └── br/edu/atitus/gatewayservice/
│       │       ├── components/          # Utilitários de segurança e criptografia (JwtUtil)
│       │       ├── configs/             # Definições programáticas de rotas (GatewayConfig)
│       │       ├── filters/             # Filtros globais e customizados para interceptação
│       │       ├── infrastructure/      # Exceções globais (TokenExpiredException, etc.)
│       │       └── GatewayServiceApplication.java
│       └── resources/
│           └── application.yml          # Propriedades de rede, Eureka e Actuator
```

## 🛣️ Malha de Roteamento Dinâmico
O gateway intercepta as requisições na porta 8765 e despacha-as dinamicamente para as instâncias ativas no Netflix Eureka utilizando o prefixo balanceado lb://. 
Abaixo estão as regras explícitas de encaminhamento configuradas na classe GatewayConfig:
<table class="routing-table">
    <thead>
        <tr>
            <th>Caminho da Requisição (Path)</th>
            <th>Microsserviço Destino (URI)</th>
            <th>Descrição do Fluxo</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td><code>/auth/**</code></td>
            <td><code>lb://auth-service</code></td>
            <td>Fluxos de autenticação, login e cadastro de usuários.</td>
        </tr>
        <tr>
            <td><code>/products/**</code></td>
            <td><code>lb://product-service</code></td>
            <td>Gerenciamento e listagem pública de catálogo de produtos.</td>
        </tr>
        <tr>
            <td><code>/ws/products/**</code></td>
            <td><code>lb://product-service</code></td>
            <td>WebSockets e comunicação estendida de produtos administrativamente.</td>
        </tr>
        <tr>
            <td><code>/currency/**</code></td>
            <td><code>lb://currency-service</code></td>
            <td>Conversão de moedas e consultas de taxas ao Banco Central (BCB).</td>
        </tr>
        <tr>
            <td><code>/ws/currency/**</code></td>
            <td><code>lb://currency-service</code></td>
            <td>Canal reativo / WebSocket de cotações monetárias.</td>
        </tr>
        <tr>
            <td><code>/ws/orders/**</code></td>
            <td><code>lb://order-service</code></td>
            <td>Processamento, persistência e eventos de pedidos/ordens de compra.</td>
        </tr>
    </tbody>
</table>

## 🛡️ Segurança Centralizada (JwtUtil)
O Gateway impede que requisições maliciosas ou sem autenticação cheguem aos microsserviços internos. <br/> 
Através da classe utilitária JwtUtil, o gateway descriptografa e valida a assinatura criptográfica dos tokens utilizando o algoritmo HMAC-SHA com uma chave privada configurada localmente.
O processo lança exceções de infraestrutura customizadas quando anomalias são encontradas:
- TokenExpiredException: Disparada automaticamente quando o tempo de vida do token expirou (ExpiredJwtException).
- InvalidTokenException: Disparada quando o payload foi violado, a assinatura é incompatível ou o formato é malformado (JwtException).
## 📊 Monitoramento Operacional e Telemetria

O ecossistema expõe dados para ferramentas de monitoramento e auditoria na raiz do barramento. As propriedades do arquivo application.yml configuram a exposição dos seguintes endpoints do Actuator:

<table class="telemetry-table">
    <thead>
        <tr>
            <th>Métrica / Monitoramento</th>
            <th>Método HTTP</th>
            <th>Endpoint do Actuator</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td><strong>Saúde Completa da Malha (Probes Liveness/Readiness)</strong></td>
            <td><mark style="background-color: #e2f0d9; color: #385723; padding: 2px 6px; border-radius: 4px;">GET</mark></td>
            <td><code>http://localhost:8765/actuator/health</code></td>
        </tr>
        <tr>
            <td><strong>Status dos Circuit Breakers</strong></td>
            <td><mark style="background-color: #e2f0d9; color: #385723; padding: 2px 6px; border-radius: 4px;">GET</mark></td>
            <td><code>http://localhost:8765/actuator/circuit-breakers</code></td>
        </tr>
        <tr>
            <td><strong>Métricas Gerais de Tráfego</strong></td>
            <td><mark style="background-color: #e2f0d9; color: #385723; padding: 2px 6px; border-radius: 4px;">GET</mark></td>
            <td><code>http://localhost:8765/actuator/metrics</code></td>
        </tr>
    </tbody>
</table>

## ⚙️ Requisitos para Inicialização Local
Certifique-se de utilizar o Java 25.Garanta que o servidor de registro Eureka Server (discovery-service) esteja ativo na porta 8761 antes de subir este serviço, permitindo a correta amarração das rotas dinâmicas.
Suba a aplicação utilizando o wrapper do Maven:
```
./mvnw spring-boot:run
```
