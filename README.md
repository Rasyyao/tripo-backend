# Tripo Backend

Backend API for Tripo, a trip and itinerary planning app. Built with Go and [Fiber v3](https://github.com/gofiber/fiber), backed by PostgreSQL.

## Tech Stack

- Go 1.25+
- Fiber v3 (HTTP framework)
- PostgreSQL 16 (via Docker Compose)
- pgAdmin 4 (database UI)

## Project Structure

```
cmd/api/main.go          # Application entry point
internal/
  config/                # Environment-based configuration
  routes/                # Route registration and dependency wiring
  controllers/           # HTTP handlers and request DTOs
  services/              # Business logic
  repositories/          # Data access layer
  models/                # Domain models (user, itinerary, place, chat)
db/migrations/           # SQL schema (auto-applied on first DB start)
docker-compose.yml       # PostgreSQL + pgAdmin
```

The code follows a layered flow: `routes -> controllers -> services -> repositories`.

## Getting Started

### Prerequisites

- Go 1.25 or newer
- Docker and Docker Compose

### Setup

```bash
# 1. Configure environment
cp .env.example .env

# 2. Start PostgreSQL and pgAdmin
docker compose up -d

# 3. Run the API
go run ./cmd/api
```

The API listens on `http://localhost:3000` by default.

### Environment Variables

| Variable                 | Default              | Description                      |
| ------------------------ | -------------------- | -------------------------------- |
| `APP_PORT`               | `3000`               | Port the API listens on          |
| `POSTGRES_USER`          | `tripo`              | Database user                    |
| `POSTGRES_PASSWORD`      | `tripo_password`     | Database password                |
| `POSTGRES_DB`            | `tripo`              | Database name                    |
| `POSTGRES_PORT`          | `5432`               | Host port for PostgreSQL         |
| `PGADMIN_DEFAULT_EMAIL`  | `admin@tripo.local`  | pgAdmin login email              |
| `PGADMIN_DEFAULT_PASSWORD` | `admin`            | pgAdmin login password           |
| `PGADMIN_PORT`           | `5050`               | Host port for pgAdmin            |

Change the default credentials before using this outside local development.

## Database

Migrations in `db/migrations/` run automatically the first time the PostgreSQL container initializes its data volume. To re-apply them from scratch:

```bash
docker compose down -v && docker compose up -d
```

Tables: `users`, `itineraries`, `itinerary_days`, `itinerary_stops`, `places`, `route_legs`, `chat_conversations`, `chat_messages`, `memories`, `activities`, `itinerary_shares`.

pgAdmin is available at `http://localhost:5050` with the `tripo` server pre-registered.

## API

Base path: `/api/v1`

| Method | Endpoint      | Description     |
| ------ | ------------- | --------------- |
| POST   | `/users`      | Create a user   |
| GET    | `/users`      | List users      |
| GET    | `/users/:id`  | Get a user      |
| PUT    | `/users/:id`  | Update a user   |
| DELETE | `/users/:id`  | Delete a user   |

## Development

```bash
go build ./...   # compile
go vet ./...     # static checks
go test ./...    # run tests
```
