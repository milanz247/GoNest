# GoMVC

A clean, lightweight Golang MVC web framework built directly on the standard
`net/http` package — no Gin/Echo/Fiber. Uses GORM + MySQL for persistence,
`html/template` for views, and optional HTMX support for partial-page updates.

## Tech Stack

- Go (standard `net/http`)
- GORM + MySQL
- Go `html/template`
- HTML/CSS, optional HTMX
- `.env`-based configuration

## Architecture

```
HTTP Request
  → Route
  → Middleware
  → Controller
  → DTO Validation
  → Service
  → GORM Model
  → MySQL
  → Service Result
  → Controller
  → HTML View or HTTP Response
```

Layers, and what belongs in each:

| Layer | Responsibility |
|---|---|
| **Controller** | HTTP only: read params, bind DTOs, call services, render/redirect/JSON. No GORM, no business logic. |
| **Service** | Business logic, GORM queries, transactions, model↔DTO mapping. No repository layer. |
| **Model** | GORM structs mapping to DB tables. No HTTP types. |
| **DTO** | Request/response shapes with validation tags. No GORM fields. |
| **View** | `html/template` layouts, pages, partials. |
| **Routes** | Route registration/grouping. |
| **Middleware** | Cross-cutting request handling (logging, recovery, security headers, etc). |
| **Config/Bootstrap** | Typed env config, DB connection, app wiring at startup. |

## Project Structure

```
myapp/
├── app/
│   ├── controllers/     HTTP-facing request handlers
│   ├── services/        Business logic + GORM access
│   ├── models/          GORM models
│   ├── dto/              Request/response DTOs
│   ├── middleware/       App-specific middleware
│   └── validation/       DTO validation helpers
├── bootstrap/            Wires config + logger + routes (+ DB later) into the app
├── config/                Typed configuration (App, Database) loaded from .env
├── framework/
│   ├── http/              Router, route, context, response (Phase 2)
│   ├── view/               Template renderer (Phase 5)
│   ├── database/           GORM connection manager (Phase 2)
│   ├── middleware/          Reusable framework middleware (Phase 2)
│   └── application.go       Core Application lifecycle: start, signal handling, graceful shutdown
├── routes/                Route definitions (Phase 2+)
├── resources/views/        HTML templates (layouts, pages, partials)
├── public/                 Static assets (css/js/images)
├── storage/                Logs and uploads
├── database/               Migrations and seeders
├── .env / .env.example
└── main.go                 Entry point: bootstrap.Boot() → app.Run()
```

## Configuration

Copy `.env.example` to `.env` and adjust as needed:

```env
APP_NAME=GoMVC
APP_ENV=local
APP_DEBUG=true
APP_HOST=127.0.0.1
APP_PORT=8080

DB_DRIVER=mysql
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=gomvc
DB_USER=root
DB_PASSWORD=
DB_CHARSET=utf8mb4
DB_PARSE_TIME=true
DB_LOCATION=Local
```

`config.Load` reads `.env` first (without overriding real environment
variables), then builds a typed `Config` struct and validates that required
fields are present. `DatabaseConfig.String()` masks the password so it's
never accidentally logged.

## Running

```bash
go run .
```

The server starts at `http://APP_HOST:APP_PORT` (default
`http://127.0.0.1:8080`).

### Health check

```bash
curl -i http://127.0.0.1:8080/health
```

```
HTTP/1.1 200 OK
Content-Type: application/json

{"status":"ok","time":"2026-07-24T18:01:18Z"}
```

### Graceful shutdown

`framework.Application.Run()` listens for `SIGINT`/`SIGTERM`. On receipt, it
stops accepting new connections and drains in-flight requests via
`http.Server.Shutdown` (bounded by a 10s timeout) before exiting cleanly.

```bash
# Ctrl+C, or:
kill -TERM <pid>
```

## Roadmap

- **Phase 1 (done):** module init, folder structure, typed env config,
  minimal `net/http` server, `/health`, graceful shutdown.
- **Phase 2:** custom router (`framework/http`) with params, groups, named
  routes, static file serving, 404/500 handlers; custom request `Context`;
  core middleware (logger, recovery, request ID, security headers, body size
  limit); GORM + MySQL connection setup.
- **Phase 3:** `User` model, base model fields.
- **Phase 4:** DTOs + validation system, `UserService`.
- **Phase 5:** `html/template` renderer (layouts/partials), `UserController`,
  full CRUD views, optional HTMX fragment responses.

## Development Rules

- No repository layer between services and GORM.
- No business logic in controllers.
- No binding HTML forms directly into GORM models — always go through a DTO.
- Constructor injection only; no global singletons, no DI framework.
- Passwords are hashed with bcrypt and never logged or returned to views.
