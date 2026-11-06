# FlowTrade

A brokerage platform built as microservices in Go, where users will be able to view assets, build a portfolio, place buy and sell orders, and track performance.

**Status: in progress.** The API gateway and account service exist today. Orders, matching and notifications are planned. See [Current status](#current-status) for what works and what doesn't yet.

## Architecture

```text
Client
  |  HTTP / JSON
  v
API Gateway (Go, Gin)
  |  gRPC (Protocol Buffers)
  v
Account Service (Go) ---> PostgreSQL
```

Planned services, not built yet: orders, matching and notifications.

### Design decisions

- **API gateway:** clients talk to one HTTP entry point, which validates JWTs and calls internal services over gRPC.
- **Hexagonal architecture:** the account service keeps its domain and use cases in `internal/core`, behind interfaces in `internal/core/ports`. PostgreSQL, gRPC and security code live in `internal/adapters`.
- **Ledger for balances:** each balance change is designed as a `ledger_entries` row that records the type, amount and balance after the change.
- **Decimal money:** account amounts use `shopspring/decimal` in Go and `NUMERIC(20, 8)` in PostgreSQL, not floats.

## Current status

| Area | Status |
|---|---|
| API gateway with JWT validation | Built |
| Gateway to account service over gRPC | Built |
| Register, login, refresh, logout, logout all | Built |
| Get, update and delete the current user | Built |
| PostgreSQL migrations: users, refresh tokens, accounts, ledger, assets | Built |
| Account use cases: create, deposit, withdraw, check balance | Written, not exposed through gRPC or HTTP yet |
| Orders, matching, notifications | Planned |
| Dockerfiles and Compose for the services | Started. `docker-compose.yml` only defines infrastructure, and references services that don't exist yet |
| CI and Kubernetes | Placeholder files only |
| Automated tests | Not written yet |

## Tech stack

- **Language:** Go 1.25
- **HTTP:** Gin
- **Internal communication:** gRPC, Protocol Buffers
- **Database:** PostgreSQL 17
- **Auth:** JWT, refresh tokens, bcrypt

## Project structure

```text
.
├── proto/                     # gRPC contracts and generated code
├── services/
│   ├── api-gateway/           # HTTP entry point
│   └── account-service/       # Users, auth, accounts
│       ├── internal/core/     # Domain, use cases and ports
│       ├── internal/adapters/ # gRPC, PostgreSQL, security
│       └── migrations/        # SQL migrations
├── infra/k8s/                 # Placeholder
└── docker-compose.yml         # Infrastructure, in progress
```

## Running locally

You need Go 1.25 and Docker.

**1. Start PostgreSQL**

```bash
docker run -d --name flowtrade-account-db -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=account_db -p 5432:5432 postgres:17-alpine
```

**2. Apply the migrations**

Run each `.up.sql` file in `services/account-service/migrations` in numeric order, for example with [golang-migrate](https://github.com/golang-migrate/migrate).

**3. Start the account service**

```bash
cd services/account-service
JWT_SECRET=change-me go run .
```

It listens for gRPC on port `50051`.

**4. Start the API gateway, in another terminal**

```bash
cd services/api-gateway
JWT_SECRET=change-me ACCOUNT_SERVICE_ADDR=localhost:50051 go run .
```

It listens for HTTP on port `8080`. Both services must use the same `JWT_SECRET`.

## Configuration

Both services read environment variables or a local `.env` file.

| Service | Variable | Default |
|---|---|---|
| Account service | `GRPC_PORT` | `50051` |
| Account service | `DB_HOST`, `DB_PORT` | `localhost`, `5432` |
| Account service | `DB_USER`, `DB_PASS`, `DB_NAME` | `postgres`, `postgres`, `account_db` |
| Account service | `DB_SSLMODE` | `disable` |
| Account service | `JWT_TOKEN_DURATION`, `REFRESH_TOKEN_DURATION` | `15m`, `168h` |
| API gateway | `PORT` | `8080` |
| API gateway | `ACCOUNT_SERVICE_ADDR` | `0.0.0.0` |
| Both | `JWT_SECRET` | Development value only. Always set your own |

## API

| Method | Path | Auth | Description |
|---|---|---|---|
| `POST` | `/users/` | No | Register |
| `POST` | `/auth/login` | No | Log in |
| `POST` | `/auth/refresh` | Yes | Refresh the session |
| `POST` | `/auth/logout` | Yes | Log out of this session |
| `POST` | `/auth/logout-all` | Yes | Log out of all sessions |
| `GET` | `/users/me` | Yes | Get the current user |
| `PUT` | `/users/me` | Yes | Update the current user |
| `DELETE` | `/users/me` | Yes | Delete the current user |

## Roadmap

1. Automated tests and CI
2. Expose accounts and ledger operations through gRPC and the gateway
3. Orders service
4. Matching and settlement
5. Event-driven communication between services
6. Notifications
7. Containers and Kubernetes for every service
