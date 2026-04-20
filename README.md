# Mock API

API para criar e gerenciar mocks dinâmicos para testes. O desenvolvedor define instructions de input/output e o sistema cria endpoints mock automaticamente.

## Funcionalidades

- **Criar mocks** via API REST
- **Respostas estáticas** (fixed)
- **Respostas dinâmicas** com templating (`{{input.field}}`, `{{faker.field}}`)
- **Dados fake** (nome, email, uuid, telefone, etc.)
- **Delay simulado** para testar timeouts
- **Swagger/OpenAPI** documentação

## Quick Start

### Docker (Recomendado)

```bash
# Clone, build e rode
docker-compose up -d

# Acesse
curl http://localhost:8080/health   # Health check
curl http://localhost:8080/docs     # Swagger UI
```

### Local

```bash
# Rode localmente
go run ./cmd/server
# Ou build e rode
go build -o bin/mocker ./cmd/server && ./bin/mocker
```

Acesse: http://localhost:8080

## Documentação

| Recurso | URL |
|--------|-----|
| **Swagger UI** | http://localhost:8080/docs |
| **OpenAPI Spec** | http://localhost:8080/docs/swagger.json |
| **Health Check** | http://localhost:8080/health |

## Endpoints

### Gerenciamento de Mocks

| Método | Endpoint | Descrição |
|--------|----------|------------|
| `POST` | `/api/v1/mocks` | Criar mock |
| `GET` | `/api/v1/mocks` | Listar mocks |
| `GET` | `/api/v1/mocks/:id` | Detalhes do mock |
| `PUT` | `/api/v1/mocks/:id` | Atualizar mock |
| `DELETE` | `/api/v1/mocks/:id` | Deletar mock |
| `POST` | `/api/v1/mocks/:id/test` | Testar mock |
| `POST` | `/api/v1/mocks/:id/activate` | Ativar mock |
| `POST` | `/api/v1/mocks/:id/deactivate` | Desativar mock |

### Mock Server (endpoints dos mocks)

Os mocks ficam acessíveis diretamente nos paths configurados:

```
GET    /api/users     → retorna dados fake
POST   /api/users    → retorna resposta estática
GET    /api/products → retorna lista de produtos
```

## Exemplos

### 1. Criar mock com dados fake

```bash
curl -X POST http://localhost:8080/api/v1/mocks \
  -H "Content-Type: application/json" \
  -d '{
    "name": "get-users",
    "path": "/api/users",
    "method": "GET",
    "responseConfig": {
      "type": "faker",
      "statusCode": 200,
      "faker": {
        "fields": [
          {"name": "id", "type": "uuid"},
          {"name": "name", "type": "name"},
          {"name": "email", "type": "email"}
        ]
      }
    }
  }'
```

**Acessar:**
```bash
curl http://localhost:8080/api/users
# {"id":"abc-123","name":"John Doe","email":"john@example.com"}
```

### 2. Criar mock estático

```bash
curl -X POST http://localhost:8080/api/v1/mocks \
  -H "Content-Type: application/json" \
  -d '{
    "name": "create-user",
    "path": "/api/users",
    "method": "POST",
    "responseConfig": {
      "type": "static",
      "statusCode": 201,
      "body": {
        "id": "123",
        "name": "João Silva",
        "created": true
      }
    }
  }'
```

### 3. Criar mock com delay

```bash
curl -X POST http://localhost:8080/api/v1/mocks \
  -H "Content-Type: application/json" \
  -d '{
    "name": "slow",
    "path": "/api/slow",
    "method": "GET",
    "responseConfig": {
      "type": "static",
      "statusCode": 200,
      "delayMs": 2000,
      "body": {"message": "Delayed response"}
    }
  }'
```

### 4. Dynamic template

```bash
curl -X POST http://localhost:8080/api/v1/mocks \
  -H "Content-Type: application/json" \
  -d '{
    "name": "echo",
    "path": "/api/echo",
    "method": "POST",
    "responseConfig": {
      "type": "dynamic",
      "statusCode": 200,
      "body": {
        "echoed": "{{input.message}}",
        "timestamp": "{{timestamp}}",
        "uuid": "{{uuid}}"
      }
    }
  }'
```

**Acessar:**
```bash
curl -X POST http://localhost:8080/api/echo \
  -H "Content-Type: application/json" \
  -d '{"message": "Hello"}'
# {"echoed":"Hello","timestamp":"1700000000","uuid":"abc-123"}
```

## Tipos de Faker

| Tipo | Exemplo |
|------|---------|
| `name` | "John Doe" |
| `first_name` | "John" |
| `last_name` | "Doe" |
| `email` | "user123@gmail.com" |
| `phone` | "+1-555-123-4567" |
| `uuid` | "550e8400-e29b..." |
| `date` | "2024-01-15" |
| `datetime` | "2024-01-15T10:30:00Z" |
| `address` | "123 Main St" |
| `city` | "New York" |
| `country` | "United States" |
| `state` | "NY" |
| `zipcode` | "10001" |
| `url` | "https://example.com" |
| `ip` | "192.168.1.1" |
| `ipv6` | "2001:0db8:..." |
| `credit_card_number` | "4532 1234 5678 9010" |
| `credit_card_type` | "Visa" |
| `word` | "hello" |
| `sentence` | "This is a test." |
| `paragraph` | "Lorem ipsum..." |

## Variáveis de Ambiente

| Variável | Default | Descrição |
|----------|---------|------------|
| `SERVER_HOST` | `0.0.0.0` |.Host do servidor |
| `SERVER_PORT` | `8080` | Porta do servidor |
| `SERVER_READ_TIMEOUT` | `15s` | Read timeout |
| `SERVER_WRITE_TIMEOUT` | `15s` | Write timeout |
| `DATABASE_PATH` | `./mocker.db` | Caminho do banco SQLite |
| `DATABASE_MAX_OPEN_CONNS` | `25` | Conexões abertas |
| `DATABASE_MAX_IDLE_CONNS` | `5` | Conexões ociosas |

## Docker

```bash
# Build e run
docker-compose up -d

# Logs
docker-compose logs -f mocker

# Stop
docker-compose down

# Rebuild
docker-compose build --no-cache
```

## Desenvolvimento

```bash
# Setup
go mod download

# Run local
go run ./cmd/server

# Testes
go test ./...

# Build
go build -o bin/mocker ./cmd/server
```

## Estrutura do Projeto

```
/home/vinicius/mocker/
├── cmd/server/          # Entry point
├── internal/
│   ├── config/         # Configuração
│   ├── models/         # Modelos de dados
│   ├── handlers/       # HTTP handlers
│   ├── storage/       # SQLite repository
│   ├── generator/     # Motor de geração
│   ├── validator/     # Validação
│   └── middleware/   # Middleware HTTP
├── migrations/         # SQL migrations
├── docs/              # Documentação
└── docker-compose.yml # Docker config
```

## Endpoints Importantes

### Mock Server

```bash
# GET users (faker)
curl http://localhost:8080/api/users

# POST users (static)
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Test"}'

# GET products
curl http://localhost:8080/api/products
```

### Gerenciamento

```bash
# Listar todos mocks
curl http://localhost:8080/api/v1/mocks

# Ver detalhes de um mock
curl http://localhost:8080/api/v1/mocks/{id}

# Testar mock
curl -X POST http://localhost:8080/api/v1/mocks/{id}/test \
  -H "Content-Type: application/json" \
  -d '{"input": {}}'

# Ativar/Desativar mock
curl -X POST http://localhost:8080/api/v1/mocks/{id}/activate
curl -X POST http://localhost:8080/api/v1/mocks/{id}/deactivate
```

## Estrutura do Projeto

```
/home/vinicius/mocker/
├── cmd/server/           # Entry point
├── internal/
│   ├── config/         # Configuração
│   ├── models/         # Modelos de dados
│   ├── handlers/       # HTTP handlers
│   ├── storage/       # SQLite repository
│   ├── generator/     # Motor de geração
│   ├── validator/     # Validação
│   └── middleware/   # Middleware HTTP
├── docs/
│   ├── openapi.yaml   # OpenAPI spec
│   └── swagger-ui/    # Swagger UI assets
├── docker-compose.yml
├── Dockerfile
└── README.md
```

## API Integrations

| Integração | URL |
|------------|-----|
| Swagger UI | http://localhost:8080/docs |
| OpenAPI JSON | http://localhost:8080/docs/swagger.json |
| ReDoc | http://localhost:8080/docs/redoc.html |

## Licença

MIT