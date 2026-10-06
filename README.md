# Go Backend Starter Skeleton

A clean, minimalist Go backend starter skeleton following the **Clean Architecture / Hexagonal (Ports & Adapters)** pattern.

Designed as an unopinionated, lightweight boilerplate for future projects with structured configuration, PostgreSQL persistence, Uber Zap logging, and REST routing.

---

## 🏛 Directory Structure

```
sceleton-go/
├── cmd/                          # CLI commands & application wiring (Cobra)
│   ├── main.go                   # Application entrypoint
│   ├── root.go                   # Root Cobra command
│   ├── serve.go                  # `serve-rest` command (--port, --config)
│   ├── server.go                 # HTTP server bootstrap & graceful shutdown
│   └── api_version_one.go        # Router v1 group (/api/v1) & route registration
│
├── internal/                     # Private application packages
│   ├── config/                   # Environment configuration loader
│   │   └── config.go
│   │
│   ├── domain/                   # Business domain entities & sentinel errors
│   │   ├── error.go              # Standard domain errors & wrapping
│   │   └── user.go               # Sample domain entity
│   │
│   ├── user/                     # Feature module (Ports & Business Logic)
│   │   ├── port.go               # Service and Repository interfaces
│   │   └── service.go            # Domain service implementation
│   │
│   ├── repository/               # Data persistence layer (PostgreSQL)
│   │   └── user_repository.go    # SQL database implementation
│   │
│   ├── infrastructure/           # External infrastructure adapters
│   │   ├── db/postgres/          # sql.DB connection pool
│   │   └── logger/               # Uber Zap logger (dev & prod)
│   │
│   └── rest/                     # Delivery layer (HTTP REST via Gin)
│       ├── middleware/           # Middlewares (Auth, CORS, Logger, RateLimiter)
│       ├── response/             # Standard API response envelopes
│       └── handler/              # HTTP controllers
│           └── user/             # User handler, DTOs & route registration
│
├── migrations/                   # SQL database migration files
│   └── 000001_init_schema.sql
├── .air.toml                     # Air live-reload configuration
├── .env.example                  # Environment variables template
├── docker-compose.yml            # Docker stack (App + PostgreSQL 17)
├── Dockerfile                    # Multi-stage production build
├── Makefile                      # Developer shortcut commands
└── go.mod
```

---

## 🚀 Quick Start

### 1. Configure Environment
```bash
cp .env.example .env
```

### 2. Start PostgreSQL (Docker)
```bash
docker compose up -d postgres
```

### 3. Run the API Server
```bash
# Run server:
go run ./cmd serve-rest

# Or with Makefile:
make run

# Or with live reload:
make dev
```

Server starts at `http://localhost:8080`.
- **Health Check**: `GET http://localhost:8080/api/v1/health`
- **Register User**: `POST http://localhost:8080/api/v1/users/register`

---

## 🛠 Adding a New Feature Module (Step-by-Step)

Follow this pattern whenever you add a new domain entity (e.g. `Product`):

1. **Entity**: Define `internal/domain/product.go`
2. **Ports**: Define `internal/product/port.go` (`Service` & `Repository` interfaces)
3. **Service**: Implement business logic in `internal/product/service.go`
4. **Repository**: Implement database queries in `internal/repository/product_repository.go`
5. **REST Handler**: Add controller in `internal/rest/handler/product/` (`handler.go`, `dto.go`, `routes.go`)
6. **Wire Up**: In `cmd/server.go`, instantiate repo ➔ service ➔ handler, and register in `cmd/api_version_one.go`.
