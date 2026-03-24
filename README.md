# x-clone

A conceptual clone of X.com, but limited to text-based tweeting. Implements the microservices architecture on the backend.

This monorepo consists of a Go‑based backend and a frontend. The repository is split into `backend`, `frontend` and `infra` folders with minimal scaffolding for development.

## Repository layout

- **backend/** – Go services for the `post` and `user` domains. Each contains its own `main.go`, `go.mod` and Dockerfile.
- **frontend/** – client application (details TBD).
- **infra/** – infrastructure configuration (Docker Compose, Terraform, etc.).

## Getting started

1. Install Go and Node.js (if working on the frontend).
2. Build and run the backend services from `backend/post` and `backend/user`.
3. See each subdirectory README for service‑specific instructions.

> _This README will be expanded with more detailed setup, testing, and deployment information as the project matures._
