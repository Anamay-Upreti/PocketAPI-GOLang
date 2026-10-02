# PocketAPI

A modular REST API for managing personal notes, user accounts, and authentication. PocketAPI is built with Go and Fiber, uses PostgreSQL for persistence, and supports JWT-based access and refresh tokens.

> **Project status:** This project is under active development. The API structure and core authentication and notes workflows are in place, while additional tests, documentation, and production hardening can be added over time.

## Features

- User registration and login
- JWT access and refresh token authentication
- Refresh-token rotation and logout support
- Logout from all devices
- Protected user profile endpoint
- Avatar upload with image validation and a 2 MB size limit
- Create, read, update, and delete notes
- Note pagination and search
- PostgreSQL persistence with GORM
- Automatic database migrations for users, notes, and refresh tokens
- Request IDs, request logging, panic recovery, and CORS middleware
- Docker Compose setup for PostgreSQL and pgAdmin

## Tech stack

- **Language:** Go 1.27.1+
- **Web framework:** Fiber v2
- **Database:** PostgreSQL 17
- **ORM:** GORM
- **Authentication:** JWT access and refresh tokens
- **Validation:** go-playground/validator
- **Local tooling:** Docker Compose and pgAdmin

## Project structure

```text
.
├── cmd/server/              # Application entrypoint
├── internal/
│   ├── config/              # Environment-based configuration
│   ├── database/            # PostgreSQL connection and migrations
│   ├── handlers/            # HTTP request handlers
│   ├── middleware/          # Authentication, logging, CORS, and request middleware
│   ├── models/              # Database models
│   ├── repositories/        # Data-access layer
│   ├── routes/              # API route registration
│   ├── services/            # Business logic
│   ├── utils/               # Shared response and helper functions
│   └── validators/           # Request validation
├── test/                    # Authentication and note tests
├── docker-compose.yml       # PostgreSQL and pgAdmin services
├── Dockerfile               # Container image definition
├── go.mod
└── README.md
```

## Requirements

Install the following before running the project locally:

- Go 1.27.1 or later
- PostgreSQL 17 or a compatible PostgreSQL version
- Docker and Docker Compose, if using the provided database setup

## Quick start

### 1. Clone the repository

```bash
git clone https://github.com/Anamay-Upreti/PocketAPI-GOLang.git
cd PocketAPI-GOLang
```

### 2. Start PostgreSQL with Docker Compose

```bash
docker compose up -d postgres
```

The Compose configuration starts PostgreSQL with these development credentials:

| Setting | Value |
|---|---|
| Host | `localhost` |
| Port | `5432` |
| Database | `pocketapi` |
| Username | `postgres` |
| Password | `password` |

To start pgAdmin as well:

```bash
docker compose up -d
```

pgAdmin is available at `http://localhost:5050`.

> The credentials in `docker-compose.yml` are intended for local development only. Use secure, externally managed credentials in production.

### 3. Configure environment variables

Create a `.env` file in the project root:

```env
APP_NAME=PocketAPI
PORT=8000
ENV=development
JWT_SECRET=replace-this-with-a-long-random-secret
DATABASE_URL=host=localhost user=postgres password=password dbname=pocketapi port=5432 sslmode=disable
```

`DATABASE_URL` uses the PostgreSQL connection string format supported by the application.

### 4. Download dependencies and run the API

```bash
go mod download
go run ./cmd/server
```

The API starts at:

```text
http://localhost:8000
```

On startup, PocketAPI connects to PostgreSQL and automatically migrates the `users`, `notes`, and `refresh_tokens` tables.

## Health checks

Check that the server is running:

```bash
curl http://localhost:8000/
```

Check the API and database status:

```bash
curl http://localhost:8000/health
```

## API overview

All application routes are prefixed with `/api`. Protected endpoints require a valid access token, usually sent as a bearer token:

```http
Authorization: Bearer <access-token>
```

### Authentication

| Method | Endpoint | Auth required | Description |
|---|---|---:|---|
| `POST` | `/api/auth/register` | No | Register a new user |
| `POST` | `/api/auth/login` | No | Log in and receive access and refresh tokens |
| `POST` | `/api/auth/refresh` | No | Refresh the access token |
| `POST` | `/api/auth/logout` | No | Revoke a refresh token |

Register example:

```bash
curl -X POST http://localhost:8000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Ada Lovelace",
    "email": "ada@example.com",
    "password": "strong-password"
  }'
```

Login example:

```bash
curl -X POST http://localhost:8000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "ada@example.com",
    "password": "strong-password"
  }'
```

### Users

| Method | Endpoint | Auth required | Description |
|---|---|---:|---|
| `GET` | `/api/users/me` | Yes | Get the authenticated user's profile |
| `POST` | `/api/users/avatar` | Yes | Upload an avatar using multipart field `avatar` |
| `POST` | `/api/users/logout-all` | Yes | Revoke sessions on all devices |

Avatar upload example:

```bash
curl -X POST http://localhost:8000/api/users/avatar \
  -H "Authorization: Bearer <access-token>" \
  -F "avatar=@/path/to/avatar.jpg"
```

Uploaded avatars are served from `/uploads/<filename>`.

### Notes

| Method | Endpoint | Auth required | Description |
|---|---|---:|---|
| `POST` | `/api/notes/` | Yes | Create a note |
| `GET` | `/api/notes/` | Yes | List notes with pagination and search |
| `GET` | `/api/notes/:id` | Yes | Get one note |
| `PUT` | `/api/notes/:id` | Yes | Update a note |
| `DELETE` | `/api/notes/:id` | Yes | Delete a note |

Create a note:

```bash
curl -X POST http://localhost:8000/api/notes/ \
  -H "Authorization: Bearer <access-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "My first note",
    "content": "This note is stored by PocketAPI."
  }'
```

List notes with pagination and search:

```bash
curl "http://localhost:8000/api/notes/?page=1&limit=10&search=project" \
  -H "Authorization: Bearer <access-token>"
```

## Configuration reference

| Variable | Default | Description |
|---|---|---|
| `APP_NAME` | `PocketAPI` | Application name |
| `PORT` | `8000` | HTTP server port |
| `ENV` | `development` | Runtime environment |
| `JWT_SECRET` | `secret` | Secret used for JWT operations; override this outside local development |
| `DATABASE_URL` | Empty | PostgreSQL connection string |

Do not commit `.env` files or production secrets to source control.

## Running tests

Run the full Go test suite with:

```bash
go test ./...
```

For additional output:

```bash
go test -v ./...
```

## Docker notes

The current Compose setup provides PostgreSQL and pgAdmin for local development. The application can also be started directly with Go using `go run ./cmd/server`.

Before using the `Dockerfile` for deployment, verify that the image includes all required build and runtime steps for your target environment.

## Development workflow

1. Start PostgreSQL with Docker Compose.
2. Create or update the local `.env` file.
3. Run the API with `go run ./cmd/server`.
4. Use the health endpoints to verify connectivity.
5. Authenticate and use the returned access token for protected routes.
6. Run `go test ./...` before submitting changes.

## Security considerations

- Replace the default `JWT_SECRET` with a long, random secret.
- Do not use the sample PostgreSQL credentials in production.
- Keep `.env` files and uploaded files out of version control.
- Use HTTPS in production.
- Restrict CORS origins instead of allowing broad development defaults.
- Review token expiry, refresh-token storage, file validation, and rate limiting before exposing the API publicly.

## License

No license has been selected for this repository yet. Add a license file if you plan to distribute or accept external contributions.

## Author

Created and maintained by [Anamay Upreti](https://github.com/Anamay-Upreti).
