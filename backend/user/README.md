# like-x backend: User Service

A Go microservice, uses Fiber (HTTP), pgx (PostgreSQL), and a clean repository-service-transport layered architecture.

## Features

- `POST /v1/users/sign-up`
- Validation for `full_name` and `password`
- PostgreSQL persistence via `pgx` and `pgxpool`
- Structured packages:
  - `model` (domain models)
  - `repository` (DB repository)
  - `service` (business logic)
  - `transport/http` (HTTP handlers)

## Requirements

- Go 1.26 (or latest supported)
- PostgreSQL (DB setup via environment variables)

## Environment variables

- `POSTGRES_USER`
- `POSTGRES_PASSWORD`
- `POSTGRES_HOST`
- `POSTGRES_PORT`
- `POSTGRES_DB_NAME`
- `POSTGRES_SSL_MODE` (e.g. `disable`)
- `HTTP_SERVER_PORT` (e.g. `3001`)

## Run

```bash
cd backend/user
go run .
```

If using `.env` file, existing code loads it via `github.com/joho/godotenv`.

## API

Request:

```bash
curl -X POST http://localhost:3001/v1/users/sign-up \
  -H 'Content-Type: application/json' \
  -d '{"full_name":"John Doe","password":"secret123"}'
```

Response:

```json
{"id":"<public-id>"}
```

## Tests

```bash
go test -cover ./... -count=1 -v
```
