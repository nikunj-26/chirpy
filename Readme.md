# Chirpy

<p align="center">
  <strong>A small, friendly social API for sharing chirps.</strong><br>
  Built with Go, PostgreSQL, JWTs, and a little moderation magic.
</p>

## About

Chirpy is a lightweight REST API for creating and sharing short messages called **chirps**. It includes:

- User registration and profile updates
- Secure password hashing with Argon2id
- JWT access tokens and long-lived refresh tokens
- Chirp creation with a 140-character limit
- Automatic filtering of selected profane words
- Author filtering and chronological sorting
- Owner-only chirp deletion
- Chirpy Red upgrades through a Polka webhook
- Static file serving and simple admin metrics

## Requirements

- Go **1.27.1** or newer
- PostgreSQL
- [Goose](https://github.com/pressly/goose) for database migrations

## Getting started

### 1. Clone the project

```bash
git clone https://github.com/nikunj-26/chirpy.git
cd chirpy
```

### 2. Create a PostgreSQL database

Chirpy uses PostgreSQL and the `gen_random_uuid()` function. Enable the `pgcrypto` extension in your database:

```sql
CREATE DATABASE chirpy;
\c chirpy
CREATE EXTENSION IF NOT EXISTS pgcrypto;
```

### 3. Configure the environment

Create a `.env` file in the project root:

```env
DB_URL=postgres://postgres:postgres@localhost:5432/chirpy?sslmode=disable
PLATFORM=dev
JWT_SECRET=replace-this-with-a-long-random-secret
POLKA_KEY=replace-this-with-your-polka-api-key
```

`PLATFORM=dev` enables the development-only database reset endpoint.

### 4. Run migrations

```bash
goose -dir sql/schema postgres "$DB_URL" up
```

### 5. Start Chirpy

```bash
go run .
```

The server listens on **http://localhost:8080**.

## API at a glance

### Health and administration

| Method | Endpoint | Auth | Description |
| --- | --- | --- | --- |
| `GET` | `/api/healthz` | — | Readiness check |
| `GET` | `/admin/metrics` | — | View file-server hit count |
| `POST` | `/admin/reset` | — | Reset metrics and database; development only |

### Users and authentication

| Method | Endpoint | Auth | Description |
| --- | --- | --- | --- |
| `POST` | `/api/users` | — | Create a user |
| `PUT` | `/api/users` | Bearer JWT | Update the current user's email and password |
| `POST` | `/api/login` | — | Log in and receive access and refresh tokens |
| `POST` | `/api/refresh` | Bearer refresh token | Create a new access token |
| `POST` | `/api/revoke` | Bearer refresh token | Revoke a refresh token |

### Chirps

| Method | Endpoint | Auth | Description |
| --- | --- | --- | --- |
| `POST` | `/api/chirps` | Bearer JWT | Create a chirp |
| `GET` | `/api/chirps` | — | List chirps |
| `GET` | `/api/chirps/{chirpID}` | — | Get one chirp |
| `DELETE` | `/api/chirps/{chirpID}` | Bearer JWT | Delete your own chirp |

List chirps by author or newest first:

```text
GET /api/chirps?author_id=<user-id>
GET /api/chirps?sort=desc
```

### Polka integration

| Method | Endpoint | Auth | Description |
| --- | --- | --- | --- |
| `POST` | `/api/polka/webhooks` | `ApiKey` | Upgrade a user after `user.upgraded` |

## Example workflow

Create a user:

```bash
curl -X POST http://localhost:8080/api/users \
  -H 'Content-Type: application/json' \
  -d '{"email":"you@example.com","password":"correct horse battery staple"}'
```

Log in:

```bash
curl -X POST http://localhost:8080/api/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"you@example.com","password":"correct horse battery staple"}'
```

Create a chirp with the returned access token:

```bash
curl -X POST http://localhost:8080/api/chirps \
  -H "Authorization: Bearer <access-token>" \
  -H 'Content-Type: application/json' \
  -d '{"body":"Hello from Chirpy!"}'
```

Chirps may contain up to 140 characters. The words `kerfuffle`, `sharbert`, and `fornax` are replaced with `****`.

## Authentication

- Access tokens are JWTs and expire after **1 hour**.
- Refresh tokens expire after **60 days** and can be revoked.
- Passwords are stored as Argon2id hashes; plaintext passwords are never returned by the API.
- Protected requests use the standard header:

```text
Authorization: Bearer <token>
```

Polka webhooks use:

```text
Authorization: ApiKey <polka-key>
```

## Project layout

```text
.
├── main.go                 # HTTP server and route registration
├── handler_*.go            # API handlers
├── internal/auth/          # Password, JWT, refresh-token, and API-key helpers
├── internal/database/      # sqlc-generated PostgreSQL access
├── sql/schema/             # Goose database migrations
├── sql/queries/            # SQL queries used by sqlc
└── index.html              # Static application page
```

## Development

Regenerate database access code after changing SQL queries or schema:

```bash
sqlc generate
```

Run the test suite:

```bash
go test ./...
```

## License

This project is available for learning and personal use.
