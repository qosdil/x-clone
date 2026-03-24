# x-clone Backend (post)

This directory contains the backend service for the `post` component of the x-clone project.

## Overview

- **Language:** Go 1.26
- **Framework:** Fiber v3
- **Entry point:** `main.go`
- **Database:** PostgreSQL (via pgx connection pool)
- **Authentication:** Header-based (Auth-User-ID header)
- **Dependencies:** managed by `go.mod`
- **Database schema:** `schema.sql`

### Features

- Post like endpoint with validation
- Database connection pooling
- Environment-based configuration
- Asynchronous post liking (optional)

## Database Schema

The service uses two main tables:

### posts
- `id` (serial, primary key)
- `public_id` (text, unique)
- `user_id` (integer)
- `post` (varchar(255))
- `created_at` (timestamp)
- `updated_at` (timestamp)

### post_likes
- `id` (serial, primary key)
- `post_id` (integer, foreign key to posts.id)
- `user_id` (integer)
- `created_at` (timestamp)
- Unique constraint on (post_id, user_id) to prevent duplicate likes

## API Endpoints

### POST /v1/posts/:public_id/like

Likes a post identified by its public ID.

**Authentication:** Required (via `Auth-User-ID` header)

**Parameters:**
- `public_id` (path): The public ID of the post to like

**Response Codes:**
- `200 OK`: Post liked successfully (synchronous)
- `202 Accepted`: Post like queued for processing (asynchronous)
- `400 Bad Request`: Invalid post ID
- `403 Forbidden`: Cannot like own post
- `404 Not Found`: Post not found
- `409 Conflict`: Post already liked by user
- `500 Internal Server Error`: Database error

**Notes:**
- Users cannot like their own posts
- Duplicate likes are prevented
- Asynchronous mode can be enabled via `ASYNC_POST_LIKE=true` environment variable

## Getting Started

### Prerequisites

- Go 1.26+
- PostgreSQL
- Environment variables configured (see below)

### Environment Variables

Create a `.env` file based on `.env.example`:

- `ASYNC_POST_LIKE`: Enable asynchronous post liking (default: false)
- `DEBUG`: Enable debug mode with pprof middleware (default: false)
- `X_CLONE_HTTP_SERVER_PORT`: HTTP server port (default: 3000)
- `X_CLONE_POSTGRES_HOST`: PostgreSQL host
- `X_CLONE_POSTGRES_PORT`: PostgreSQL port
- `X_CLONE_POSTGRES_SSL_MODE`: SSL mode for PostgreSQL connection
- `X_CLONE_POSTGRES_USER`: PostgreSQL username
- `X_CLONE_POSTGRES_PASSWORD`: PostgreSQL password
- `X_CLONE_POSTGRES_DB_NAME`: PostgreSQL database name

### Setup

1. Install dependencies:
   ```sh
   go mod download
   ```

2. Create a `.env` file based on `.env.example`

3. Initialize the database:
   ```sh
   psql -U postgres -h localhost -d x_clone_post -f schema.sql
   ```

4. Build the service:
   ```sh
   go build -o tmp/bin/post ./...
   ```

5. Run the service:
   ```sh
   ./tmp/bin/post
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

## Simple, Local Load Testing

If not present, add a Post with ID 1 and User ID 1:
```sql
INSERT INTO "public"."posts" ("id", "public_id", "post", "user_id") VALUES (1, 'abcd1234', 'Hi.', 1)
```

If not empty, reset `post_likes` table:
```sql
TRUNCATE TABLE post_likes RESTART IDENTITY;
```

Run the HTTP server with constrained resources for a lightweight load test:

```sh
GOMAXPROCS=0.05 GOMEMLIMIT=128MiB go run main.go
```

### Install vegeta

Let's use `vegeta` for this load testing. Install via Homebrew (macOS/Linux), Go install, or download a binary:

```sh
# Homebrew (recommended on macOS)
brew install vegeta

# Go install (with Go 1.20+)
go install github.com/tsenart/vegeta/v12@latest

# Verify installation
vegeta -version
```

Using vegeta, test the load with 1,000 RPS for 30 seconds, start with User ID 2:

```sh
seq 2 30001 \
  | awk '{ printf "POST http://localhost:3000/v1/posts/abcd1234/like\nAuth-User-ID: %s\n\n", $1 }' \
  | vegeta attack -rate=1000 -duration=30s -timeout=60s \
  | vegeta report
```
