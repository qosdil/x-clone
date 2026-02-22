# x-clone Backend (post)

This directory contains the backend service for the `post` component of the x-clone project.

## Overview

- **Language:** Go
- **Framework:** Fiber v3
- **Entry point:** `main.go`
- **Database:** PostgreSQL (via pgx connection pool)
- **Authentication:** Header-based (Auth-User-ID header)
- **Dependencies:** managed by `go.mod`
- **Database schema:** `schema.sql`

### Features

- User authentication middleware
- Post like endpoint with validation
- Database connection pooling
- Environment-based configuration

## Getting Started

### Prerequisites

- Go 1.26+
- PostgreSQL
- Environment variables configured (see below)

### Setup

1. Install dependencies:
   ```sh
   go mod download
   ```

2. Create a `.env` file, see `.env.example` as example:
3. Initialize the database:
   ```sh
   psql -U postgres -h localhost -d x_clone -f schema.sql
   ```

4. Build the service:
   ```sh
   go build -o tmp/bin/post ./...
   ```

5. Run the service:
   ```sh
   ./tmp/bin/post
   ```

## API Endpoints

### `POST /v1/posts/:id/like`

Like a post.

**Headers:**
- `Auth-User-ID: <uint>` – User ID performing the action (required; must be between 1-100,000)

**Response:**
- `200 OK` – Like recorded successfully
- `400 Bad Request` – Invalid post ID
- `403 Forbidden` – Invalid user ID or attempting to like own post
- `404 Not Found` – Post does not exist
- `409 Conflict` – User already liked the post
- `500 Internal Server Error` – Database error

## Environment Variables

| Variable | Description | Example |
|---|---|---|
| `X_CLONE_POSTGRES_USER` | PostgreSQL user | postgres |
| `X_CLONE_POSTGRES_PASSWORD` | PostgreSQL password | password |
| `X_CLONE_POSTGRES_HOST` | PostgreSQL host | localhost |
| `X_CLONE_POSTGRES_PORT` | PostgreSQL port | 5432 |
| `X_CLONE_POSTGRES_DB_NAME` | PostgreSQL database name | x_clone |
| `X_CLONE_POSTGRES_SSLMODE` | SSL mode | disable |
| `X_CLONE_HTTP_SERVER_PORT` | HTTP server port | 8080 |

## Project Structure

```
post/
├── main.go          # application entry point
├── go.mod           # Go modules definition
├── schema.sql       # database schema
```

## Docker

A multi-stage `Dockerfile` is provided for building and running the service in a container.

Build the image from the `post` directory:
```sh
docker build -t x-clone-post .
```

Adjust ports and environment variables (see `.env.example`) as needed by the application.

Run the container:
```sh
docker run --rm --env-file .env -p 3000:3000 --name x-clone-post-latest x-clone-post:latest
```
