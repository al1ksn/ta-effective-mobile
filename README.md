# Subscriptions Management Service

REST API service for managing user subscriptions, built with Go.

## Features

- Full CRUD operations for subscriptions
- Total cost calculation with flexible period and user/service filters
- Automatic database migrations on startup
- Swagger UI for interactive API exploration
- Dockerized setup with PostgreSQL

## Tech Stack

- **Go 1.26** — language
- **Chi** — HTTP router
- **PostgreSQL 16** — database
- **pgx v5** — PostgreSQL driver
- **golang-migrate** — database migrations
- **Swagger / swaggo** — API documentation

## Getting Started

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/) and Docker Compose

### Run with Docker

```bash
cp .env.example .env
docker-compose up
```

The API will be available at `http://localhost:8080`.  
Swagger UI: `http://localhost:8080/swagger/`

### Run locally

```bash
# Install dependencies
go mod download

# Configure environment
cp .env.example .env
# Edit .env with your PostgreSQL credentials

# Start the server
go run ./cmd/api
```

## Configuration

| Variable      | Default        | Description           |
|---------------|----------------|-----------------------|
| `DB_HOST`     | `localhost`    | PostgreSQL host       |
| `DB_PORT`     | `5432`         | PostgreSQL port       |
| `DB_USER`     | `postgres`     | Database user         |
| `DB_PASSWORD` | `postgres`     | Database password     |
| `DB_NAME`     | `subscriptions`| Database name         |
| `SERVER_PORT` | `8080`         | HTTP server port      |

## API Reference

Base path: `/api/v1`

| Method   | Endpoint                  | Description                        |
|----------|---------------------------|------------------------------------|
| `POST`   | `/subscriptions`          | Create a subscription              |
| `GET`    | `/subscriptions`          | List all subscriptions             |
| `GET`    | `/subscriptions/{id}`     | Get subscription by ID             |
| `PUT`    | `/subscriptions/{id}`     | Update a subscription              |
| `DELETE` | `/subscriptions/{id}`     | Delete a subscription              |
| `GET`    | `/subscriptions/total`    | Calculate total cost with filters  |

### Create / Update subscription

```json
{
  "service_name": "Netflix",
  "price": 990,
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "start_date": "01-2026",
  "end_date": "12-2026"
}
```

> Dates use `MM-YYYY` format. `end_date` is optional.

### Total cost query parameters

| Parameter      | Type   | Description                        |
|----------------|--------|------------------------------------|
| `from`         | string | Period start in `MM-YYYY` format   |
| `to`           | string | Period end in `MM-YYYY` format     |
| `user_id`      | UUID   | Filter by user                     |
| `service_name` | string | Filter by service name             |

### Response example

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "service_name": "Netflix",
  "price": 990,
  "user_id": "550e8400-e29b-41d4-a716-446655440001",
  "start_date": "01-2026",
  "end_date": "12-2026",
  "created_at": "2026-01-01T00:00:00Z",
  "updated_at": "2026-01-01T00:00:00Z"
}
```

## Database Schema

```sql
CREATE TABLE subscriptions (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    service_name VARCHAR(255) NOT NULL,
    price        INTEGER      NOT NULL CHECK (price > 0),
    user_id      UUID         NOT NULL,
    start_date   DATE         NOT NULL,
    end_date     DATE,
    created_at   TIMESTAMPTZ  DEFAULT NOW(),
    updated_at   TIMESTAMPTZ  DEFAULT NOW()
);
```

## Project Structure

```
.
├── cmd/api/          # Entry point
├── internal/
│   ├── config/       # Environment configuration
│   ├── handler/      # HTTP handlers
│   ├── model/        # Data models
│   └── repository/   # Database layer
├── migrations/       # SQL migrations
├── docs/             # Generated Swagger docs
├── Dockerfile
└── docker-compose.yml
```

## License

MIT © 2026 Alexander Kalugin
