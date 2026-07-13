# API Clientes

[![Go](https://img.shields.io/badge/Go-1.21-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-13+-4169E1?logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Architecture](https://img.shields.io/badge/Architecture-Layered-0A0A0A)](#arquitetura)

API REST desenvolvida em Go para gerenciamento de clientes, produtos e pedidos, com foco em boas praticas de arquitetura, validacao de dados e consistencia transacional no controle de estoque.

## Visao de produto

Este projeto simula um backend de e-commerce em escala inicial, cobrindo operacoes centrais:

- onboarding de clientes
- catalogo de produtos com estoque
- criacao de pedidos com validacao de regras de negocio
- persistencia em PostgreSQL com migracoes versionadas

## Highlights tecnicos

- Arquitetura em camadas (handler -> repository -> database)
- Roteamento com Chi e middlewares de observabilidade (request id, logger, recoverer)
- Password hashing com bcrypt
- Controle de estoque no fluxo de pedidos
- Migracoes SQL automaticas ao iniciar a aplicacao
- Configuracao orientada a ambiente (.env para dev local)

## Demo rapida

Base URL local: http://localhost:8080

- Health check: GET /health
- Criar cliente: POST /customers
- Login cliente: POST /customers/login
- Criar produto: POST /products
- Criar pedido: POST /orders

Colecao pronta de testes HTTP: examples.http

## Arquitetura

```mermaid
flowchart LR
    A[Cliente HTTP] --> B[Router Chi]
    B --> C[Handlers]
    C --> D[Repositories]
    D --> E[(PostgreSQL)]
    F[Migrations Embedded] --> E
```

## Estrutura do projeto

```text
cmd/main.go                # bootstrap da aplicacao
internal/config            # configuracao por variaveis de ambiente
internal/db                # conexao com banco
internal/migrate           # executor de migracoes + SQL embedado
internal/router            # definicao de rotas
internal/handler           # camada HTTP e validacoes de entrada
internal/repository        # acesso a dados e regras de persistencia
internal/model             # entidades e payloads
examples.http              # cenarios de teste prontos
```

## Stack

- Go 1.21
- Chi v5
- PostgreSQL
- bcrypt
- godotenv

## Como executar

### 1) Pre-requisitos

- Go 1.21+
- PostgreSQL 13+

### 2) Configuracao

Use .env.example como base:

```env
PORT=8080
DATABASE_URL=postgres://postgres:postgres@localhost:5432/api_clientes?sslmode=disable
```

Defaults caso variaveis nao sejam definidas:

- PORT=8080
- DATABASE_URL=postgres://postgres:postgres@localhost:5432/api_clientes?sslmode=disable

### 3) Banco de dados

```bash
createdb api_clientes
```

### 4) Subir aplicacao

```bash
go mod download
go run ./cmd
```

Saida esperada:

```text
API rodando em http://localhost:8080
```

## Endpoints

### Health

- GET /health

### Customers

- POST /customers
- POST /customers/login
- GET /customers/{id}

### Products

- POST /products
- GET /products
- GET /products/{id}
- PUT /products/{id}

### Orders

- POST /orders
- GET /orders?limit=10&offset=0
- GET /orders/{id}

## Regras de negocio

- Cliente: nome, email e senha obrigatorios
- Produto: preco e estoque nao podem ser negativos
- Pedido: exige customer_id e ao menos um item
- Pedido: responde 422 quando nao ha estoque suficiente

## Migracoes

As migracoes ficam em internal/migrate/sql e sao aplicadas automaticamente no startup.

- Tabela de controle: schema_migrations
- Extensao habilitada: pgcrypto (gen_random_uuid)
- Migracao inicial: 000001_init.up.sql

## Exemplo de request

```bash
curl -X POST http://localhost:8080/customers \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Joao Silva",
    "email": "joao@example.com",
    "password": "senha@123"
  }'
```

## Formato de erro

```json
{
  "error": "mensagem de erro"
}
```

## Possiveis proximos passos

- Autenticacao com JWT e middleware de autorizacao
- Documentacao OpenAPI/Swagger
- Testes automatizados (unitarios e integracao)
- Containerizacao com Docker + docker compose
- CI com lint, testes e build
