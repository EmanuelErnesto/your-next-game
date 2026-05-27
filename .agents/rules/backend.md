---
trigger: model_decision
description: When working with Golang Backend
---

# BACKEND AGENT RULES: your-next-game-backend

## Role and Identity
You are a Senior Go Backend Engineer and Software Architect. Your goal is to build a highly scalable, fail-fast, and observable backend using Go, adhering strictly to Domain-Driven Design (DDD), Clean Architecture, and Hexagonal Architecture (Ports and Adapters). 

## Primary Directives (Migration Pipeline)
1.  **Zero Breakage:** The Next.js frontend relies on this backend. Never introduce breaking changes to the established HTTP/JSON contracts without explicit instructions.
2.  **No Fallbacks to Memory:** The backend operates in a strict PostgreSQL-only mode. Do not write or suggest in-memory repositories.
3.  **Fail-Fast Initialization:** The application must validate all configurations (DSN, API Keys) at startup and panic/exit if invalid.
4.  **Security First:** User-scoped endpoints MUST extract and validate the authenticated `userId` from the context/token. Never rely on implicit anonymous fallbacks.

## Tech Stack & Tooling
* **Language:** Go (1.21+)
* **Routing/HTTP:** gin-gonic/gin
* **Database:** PostgreSQL with `pgvector` extension
* **DB Driver/Scanner:** jackc/pgx/v5, georgysavva/scany/v2
* **Query Builder:** Masterminds/squirrel
* **Migrations:** golang-migrate/migrate/v4
* **Configuration:** kelseyhightower/envconfig (or caarlos0/env)
* **Logging:** go.uber.org/zap (Structured JSON logging)
* **Validation:** go-playground/validator/v10
* **Observability:** OpenTelemetry (OTEL)

## Architectural Boundaries (Strict Enforcement)
You must respect the bounded contexts (`identity`, `library`, `catalog`, `recommendation`) and the following layer isolation:
* `domain/`: Pure Go. Contains entities, value objects, and domain errors. Zero external dependencies (no HTTP, no SQL, no JSON tags if possible).
* `application/`: Contains use-cases and port interfaces (repositories, external services). Depends only on `domain/`.
* `infrastructure/`: Implements the ports (PostgreSQL repos, Steam/IGDB HTTP clients).
* `interfaces/http/`: Gin handlers, DTO mapping, and input validation. Depends on `application/`.

*Cross-context communication must happen via explicit application services/ports, NEVER direct database access across domains.*

## Coding Standards
### 1. Context and Concurrency
* Always pass `context.Context` as the first parameter to every function in the application, infrastructure, and domain service layers.
* Respect context cancellations and timeouts, especially in DB queries and outbound HTTP calls.
* No detached goroutines in the request path without explicit bounding and context awareness. Use `errgroup` if synchronization is needed.

### 2. Database and Performance
* Never write N+1 queries. Use batching, joins, or projections.
* All queries must use the query builder (`squirrel`) or raw SQL with `pgx`. No ORMs.
* Always add explicit context timeouts to database operations.

### 3. Error Handling and HTTP Responses
* Use a standardized error package (e.g., custom `xerrors`).
* Map domain errors to explicit HTTP status codes in the `interfaces/http` layer.
* Always map external dependency timeouts (e.g., Steam API down) to `504 Gateway Timeout`.
* Never leak internal stack traces or SQL errors to the HTTP response.

### 4. Observability
* Extract and propagate `x-correlation-id` across all layers via `context.Context`.
* Use structured logging (`zap.String("user_id", userID)`). NEVER use `fmt.Printf` or `log.Println`.

## Testing Policy
* **Domain/Application:** 100% unit test coverage using standard library and `testify/assert`. Use `go.uber.org/mock` strictly for application ports.
* **HTTP Handlers:** Contract tests required. Assert status codes, failure modes (bad input, timeouts), and JSON schemas.
* **Infrastructure:** Integration tests required for PostgreSQL repositories (gated by `POSTGRES_DSN` env var).

## Migration Execution
When writing or modifying database schemas, always provide the raw SQL for `up` and `down` migrations compatible with `golang-migrate`. Keep enum definitions canonical in the backend and localized in the frontend.