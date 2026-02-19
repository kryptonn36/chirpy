# Chirpy

A lightweight Twitter-like backend API written in Go with PostgreSQL.

Chirpy provides:
- user registration and login
- JWT-based auth with refresh tokens
- chirp (post) creation, listing, filtering, and deletion
- a webhook-driven `Chirpy Red` upgrade flow
- basic admin endpoints for health, metrics, and local reset

---

## Tech Stack

- Go `1.24.4`
- PostgreSQL
- SQLC-generated query layer (`internal/database`)
- Goose migrations (`sql/schema`)
- Argon2id password hashing
- JWT auth (`HS256`)

---

## Project Structure

- `main.go` — server startup and route registration
- `handler_*.go` — HTTP handlers for users, auth, chirps, and admin routes
- `internal/auth` — hashing, JWT, bearer token, API key helpers
- `internal/database` — generated SQLC models and queries
- `sql/schema` — Goose database migrations
- `sql/queries` — SQLC query definitions
- `Makefile` — migration helper commands

---

## Prerequisites

Install:
- Go `1.24+`
- PostgreSQL
- Goose CLI (`go install github.com/pressly/goose/v3/cmd/goose@latest`)

Create a PostgreSQL database, then set environment variables in `.env`.

Example `.env`:

```env
DB_URL=postgres://postgres:postgres@localhost:5432/chirpy?sslmode=disable
PLATFORM=dev
SECRET=your-jwt-secret
POLKA_KEY=your-polka-api-key
```

### Environment Variables

- `DB_URL` (required): PostgreSQL connection string
- `PLATFORM` (required): must be `dev` to allow `/admin/reset`
- `SECRET` (required): JWT signing secret
- `POLKA_KEY` (required only for webhook auth): expected API key for `/api/polka/webhooks`

---

## Setup and Run

1. Install dependencies:

```bash
go mod download
```

2. Run database migrations:

```bash
make migrate-up
```

3. Start the server:

```bash
go run .
```

Server runs on:
- `http://localhost:8080`

Static app files are served at:
- `http://localhost:8080/app/`

---

## API Overview

Base URL: `http://localhost:8080`

### Health & Admin

- `GET /api/healthz` — readiness check
- `GET /admin/metrics` — returns HTML with file-server hit count
- `POST /admin/reset` — resets metrics and deletes users (only when `PLATFORM=dev`)

### Users & Auth

- `POST /api/users` — register user
- `POST /api/login` — login, returns access + refresh token
- `POST /api/refresh` — get new access token using refresh token
- `POST /api/revoke` — revoke refresh token
- `PUT /api/users` — update authenticated user email/password

### Chirps

- `POST /api/chirps` — create chirp (auth required)
- `GET /api/chirps` — list chirps
  - optional `author_id=<uuid>` filter
  - optional `sort=desc` for reverse chronological order
- `GET /api/chirps/{chirpID}` — fetch one chirp
- `DELETE /api/chirps/{chirpID}` — delete chirp (owner only)

### Webhooks

- `POST /api/polka/webhooks`
  - expects `Authorization: ApiKey <POLKA_KEY>`
  - processes `event == "user.upgraded"`
  - upgrades user to `is_chirpy_red = true`

---

## Auth Conventions

### Bearer Tokens

Protected endpoints expect:

```http
Authorization: Bearer <token>
```

Used for:
- access token on protected user/chirp routes
- refresh token for `/api/refresh` and `/api/revoke`

### API Key

Webhook endpoint expects:

```http
Authorization: ApiKey <key>
```

---

## Example Requests

### Register

```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@example.com","password":"secret123"}'
```

### Login

```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@example.com","password":"secret123"}'
```

### Create Chirp

```bash
curl -X POST http://localhost:8080/api/chirps \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{"body":"hello chirpy"}'
```

### Refresh Access Token

```bash
curl -X POST http://localhost:8080/api/refresh \
  -H "Authorization: Bearer <refresh_token>"
```

### Revoke Refresh Token

```bash
curl -X POST http://localhost:8080/api/revoke \
  -H "Authorization: Bearer <refresh_token>"
```

### List Chirps by Author

```bash
curl "http://localhost:8080/api/chirps?author_id=<user_uuid>"
```

### Delete Chirp

```bash
curl -X DELETE http://localhost:8080/api/chirps/<chirp_id> \
  -H "Authorization: Bearer <access_token>"
```

---

## Behavior Notes

- Chirp body max length is `140` characters.
- Profanity filter replaces these words with `****` (case-insensitive):
  - `kerfuffle`
  - `sharbert`
  - `fornax`
- Deleting a user cascades deletes to chirps and refresh tokens.

---

## Development Notes

### Migrations

- Apply: `make migrate-up`
- Rollback one step: `make migrate-down`

### SQLC

If you change SQL in `sql/queries`, regenerate code with:

```bash
sqlc generate
```

(`sqlc.yaml` is already configured in this repository.)

---

## Quick Smoke Test Flow

1. `POST /api/users`
2. `POST /api/login`
3. `POST /api/chirps` (with access token)
4. `GET /api/chirps`
5. `POST /api/refresh` (with refresh token)
6. `POST /api/revoke` (with refresh token)

This verifies the core user-auth-chirp lifecycle end to end.
