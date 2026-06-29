# Config Server — repositório de propriedades

Arquivos de configuração centralizada consumidos pelo **config-service** (Spring Cloud Config).

## Estrutura

Cada microsserviço que usa o Config Server possui um arquivo `{application-name}.properties` nesta pasta.

| Arquivo | Serviço |
|---------|---------|
| `greeting-service.properties` | greeting-service (base) |
| `greeting-service-*.properties` | Profiles de idioma (en, es, fr, it) |

## Fonte no runtime

O config-service busca estas propriedades no GitHub:

- Repositório: `https://github.com/PedroScheurer/marcal-store.git`
- Branch: `main`
- Pasta: `configs/`

Em ambiente Docker, a pasta local `./configs` também é montada como fallback (profile `composite`), garantindo subida do cluster mesmo sem acesso ao GitHub.

## Alterar configuração

1. Edite o arquivo correspondente nesta pasta.
2. Commit e push para `main` no repositório backend.
3. Reinicie o config-service (ou use `/actuator/refresh` nos clientes que suportarem).
