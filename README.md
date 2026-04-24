# Mock API

API to create and manage dynamic mocks for testing. Developers define input/output instructions and the system automatically creates mock endpoints.

## Features

- **Create mocks** via REST API
- **Static responses** (fixed)
- **Dynamic responses** with templating (`{{input.field}}`, `{{faker.field}}`)
- **Fake data** (name, email, uuid, phone, etc.)
- **Simulated delay** for testing timeouts
- **Swagger/OpenAPI** documentation

## Quick Start

### Docker (Recommended)

```bash
docker-compose up -d
```

### Local

```bash
go run ./cmd/server
# or
go build -o bin/mocker ./cmd/server && ./bin/mocker
```

## URLs

| Service | URL |
|---------|-----|
| Frontend | http://localhost |
| API | http://localhost:8080 |
| Swagger UI | http://localhost:8080/docs |
| Health Check | http://localhost:8080/health |

## Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/v1/mocks` | Create mock |
| `GET` | `/api/v1/mocks` | List mocks |
| `GET` | `/api/v1/mocks/:id` | Get mock details |
| `PUT` | `/api/v1/mocks/:id` | Update mock |
| `DELETE` | `/api/v1/mocks/:id` | Delete mock |
| `POST` | `/api/v1/mocks/:id/test` | Test mock |
| `POST` | `/api/v1/mocks/:id/activate` | Activate mock |
| `POST` | `/api/v1/mocks/:id/deactivate` | Deactivate mock |

## Examples

### 1. Mock with fake data

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

### 2. Static mock

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
      "body": {"id": "123", "name": "John Doe", "created": true}
    }
  }'
```

### 3. Mock with delay

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

## Faker Types

| Type | Example |
|------|---------|
| `name` | "John Doe" |
| `first_name` | "John" |
| `last_name` | "Doe" |
| `email` | "user@gmail.com" |
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
| `credit_card_number` | "4532 1234 5678 9010" |
| `word` | "hello" |
| `sentence` | "This is a test." |

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_HOST` | `0.0.0.0` | Server host |
| `SERVER_PORT` | `8080` | Server port |
| `DATABASE_PATH` | `./mocker.db` | SQLite database path |
| `DATABASE_MAX_OPEN_CONNS` | `25` | Max open connections |
| `DATABASE_MAX_IDLE_CONNS` | `5` | Max idle connections |

## Docker

```bash
# Build and run
docker-compose up -d

# Logs
docker logs -f mocker

# Stop
docker-compose down

# Rebuild
docker-compose build --no-cache
```

## Project Structure

```
mocker/
├── cmd/server/      # Entry point
├── internal/
│   ├── config/      # Configuration
│   ├── models/      # Data models
│   ├── handlers/    # HTTP handlers
│   ├── storage/     # SQLite repository
│   ├── generator/   # Generation engine
│   └── middleware/  # HTTP middleware
├── migrations/      # SQL migrations
├── frontend/        # React frontend
├── Dockerfile
├── docker-compose.yml
└── README.md
```

## License

MIT