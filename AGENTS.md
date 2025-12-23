# AGENTS.md

## Project Name

**NOFX** — Agentic Trading OS for AI-powered crypto futures trading (multi-model agents, multi-exchange, web UI, backtesting, and debate arena).

## Overview

NOFX is a full-stack automated trading platform focused on cryptocurrency perpetual/futures trading. It combines market data ingestion, strategy configuration, AI-model-driven decision making, execution across multiple exchanges, and persistent tracking of positions/P&L.

The system runs a Go backend (HTTP API + trading engine) backed by a local SQLite database, and a React/TypeScript single-page app that provides a configuration UI (AI models, exchanges, strategies, traders) plus real-time monitoring dashboards, competition views, and backtest workflows.

Security-sensitive values (exchange API keys, model API keys) are designed to be encrypted at rest (AES) and can optionally be encrypted in transit from browser to backend (RSA/WebCrypto). This project is experimental and intended for learning/research or careful small-scale testing—automated trading carries significant financial risk.

## Technology Stack

- **Language/Runtime**
  - Backend: Go (module/CI targets Go 1.25)
  - Frontend: TypeScript + React (Node.js 18+; Docker build uses Node 20)
- **Framework(s)**
  - Backend HTTP API: Gin (`github.com/gin-gonic/gin`)
  - Frontend: React + Vite (SPA)
- **Key Dependencies**
  - Database: SQLite (`modernc.org/sqlite`)
  - Auth: JWT (`github.com/golang-jwt/jwt/v5`) + OTP (`github.com/pquerna/otp`)
  - Realtime: WebSocket (`github.com/gorilla/websocket`)
  - Exchanges: Binance (`github.com/adshao/go-binance/v2`), Bybit, OKX, Hyperliquid, Aster, Lighter (mix of SDK + custom)
  - AI clients: `mcp/` package (DeepSeek, Qwen, OpenAI, Claude, Gemini, Grok, Kimi)
  - Frontend UI/state: TailwindCSS, Radix UI, Zustand, SWR, Recharts
- **Build Tools**
  - Backend: `go build`, `go test`, `Makefile`
  - Frontend: Vite, TypeScript, ESLint, Prettier, Vitest
  - Deployment: Docker + Docker Compose, Nginx

## Project Structure

```
nofx/
├── main.go                  # Backend entrypoint (initializes config, DB, crypto, managers, API server)
├── go.mod / go.sum          # Go module dependencies
├── Makefile                 # Dev/test/build helpers
├── .env.example             # Environment template (ports, encryption keys, transport encryption)
├── docker-compose.yml        # Local/dev compose (builds images from source)
├── docker-compose.prod.yml   # Production compose (uses prebuilt images)
├── docker/                   # Dockerfiles for backend/frontend
├── nginx/nginx.conf          # Nginx config (serves SPA + proxies /api to backend)
│
├── api/                      # Gin HTTP API (routes under /api)
├── auth/                     # JWT + OTP auth helpers
├── config/                   # Global config loader from env
├── crypto/                   # AES/RSA encryption for storage + optional transport encryption
├── store/                    # SQLite persistence layer (users, traders, models, exchanges, strategies, backtests, equity, etc.)
├── manager/                  # Trader lifecycle orchestration
├── trader/                   # Exchange adapters + AutoTrader execution engine
├── market/                   # Market data (WS + REST), historical data utilities
├── decision/                 # Strategy/decision engine (AI prompt + parse + execution hooks)
├── debate/                   # Multi-agent debate engine
├── backtest/                 # Backtest engine (runs, caching, persistence, metrics)
├── mcp/                      # AI provider client implementations and request/response handling
├── provider/                 # Data providers (e.g., coin lists/sources)
│
├── web/                      # Frontend React SPA
│   ├── package.json          # Frontend deps + scripts
│   ├── vite.config.ts        # Dev server + /api proxy
│   └── src/                  # UI pages/components/stores/lib
│
├── docs/                     # Documentation center and guides
├── scripts/                  # Helper scripts (PR tooling, encryption migration, etc.)
└── .github/workflows/         # CI workflows (Go tests/coverage, docker healthcheck, etc.)
```

## Key Features

- Multi-AI model support (switch providers/models)
- Multi-exchange execution (CEX/DEX perps)
- Web-based configuration (no manual JSON editing expected)
- Strategy Studio concepts (coin sources, indicators, prompts, risk controls)
- Backtesting (historical simulation + metrics + resume/checkpoint)
- Debate Arena (multi-agent roles + consensus workflow)
- Real-time monitoring (positions, equity, decisions)
- Encryption for stored secrets; optional browser→server transport encryption

## Getting Started

### Prerequisites

Required (development / building from source):

- Go (CI uses Go 1.25)
- Node.js 18+ (frontend dev); npm
- TA-Lib (required for builds that rely on technical indicator bindings)
- Docker + Docker Compose (recommended for deployment)

Optional:

- Python (used in CI helper scripts)

### Installation

#### Option A: Production (recommended) — Docker Compose (prebuilt images)

```bash
# Download compose file
curl -O https://raw.githubusercontent.com/NoFxAiOS/nofx/main/docker-compose.prod.yml

# Start
docker compose -f docker-compose.prod.yml up -d

# Logs
docker compose -f docker-compose.prod.yml logs -f
```

#### Option B: Local/dev with Docker Compose (build from this repo)

```bash
# Ensure .env exists
cp .env.example .env

# Build + start
docker compose up -d --build

# Logs
docker compose logs -f
```

#### Option C: Manual development (backend + frontend)

```bash
# Backend deps
go mod download

# Frontend deps
cd web && npm install
```

If TA-Lib is required locally:

```bash
# macOS
brew install ta-lib

# Ubuntu/Debian
sudo apt-get install libta-lib0-dev
```

### Usage

Defaults:

- Frontend: `http://localhost:3000`
- Backend API: `http://localhost:8080` (health: `/api/health`)

Run backend locally:

```bash
# Ensure required env vars exist (see Configuration)
cp .env.example .env

# Run
go run main.go
```

Run frontend locally (separate terminal):

```bash
cd web
npm run dev
```

After the UI loads, configure:

1. AI models (provider + API keys)
2. Exchanges (API keys/credentials)
3. Strategies
4. Traders

## Development

### Available Scripts

Backend (Makefile):

- `make test` / `make test-backend` / `make test-frontend`
- `make build` / `make build-frontend`
- `make run` / `make run-frontend`
- `make fmt` / `make lint`
- `make docker-build` / `make docker-up` / `make docker-down` / `make docker-logs`

Frontend (`web/package.json`):

- `npm run dev` — Vite dev server
- `npm run build` — typecheck + production build
- `npm run preview` — preview build
- `npm run test` — Vitest
- `npm run lint` / `npm run lint:fix`
- `npm run format` / `npm run format:check`

Git hooks:

- `.husky/pre-commit` runs `cd web && npx lint-staged` (root hook)
- `web/.husky/pre-commit` runs `npm test` when Husky is installed from `web/package.json`

### Coding Conventions

- Go: run `go fmt ./...` before committing; keep error handling explicit and code idiomatic.
- TypeScript/React: follow Prettier in `web/.prettierrc.json` (no semicolons, single quotes, ~80 cols); run `npm run lint` and `npm run format`.
- Safety: never commit secrets (`.env`, API keys) or local data (`data/data.db`).

### PR Workflow

- Use conventional commits (`feat|fix|docs|style|refactor|perf|test|chore|ci|security`) and keep PRs focused.
- Before opening a PR: `make test` (or `go test ./...` and `cd web && npm test`).
- Helper scripts: `scripts/pr-check.sh` (analyze PR quality) and `scripts/pr-fix.sh` (guided fixes).

### Troubleshooting

- Primary guide: `docs/guides/TROUBLESHOOTING.md`
- FAQ: `docs/guides/faq.en.md` and localized docs under `docs/i18n/`
- Docker quickstart script (local repo): `start.sh`
- Encryption details: `scripts/ENCRYPTION_README.md`
- Architecture index: `docs/architecture/README.md`

### Development Workflow

1. Copy `.env.example` to `.env` and generate real secrets (or use `start.sh` which can generate keys).
2. Run backend + frontend locally (`make run` + `make run-frontend`) or use Docker Compose.
3. Run tests before pushing:

```bash
make test
```

## Configuration

### Important Files

- `.env.example` / `.env` — runtime secrets and ports
- `config/config.go` — backend env configuration (JWT secret, registration flags, limits)
- `docker-compose.yml`, `docker-compose.prod.yml` — service wiring and port mapping
- `nginx/nginx.conf` — SPA serving + `/api/` proxy to backend
- `web/vite.config.ts` — frontend dev server + `/api` proxy

### Environment Variables

Core (required to run backend safely):

- `JWT_SECRET` — JWT signing secret (set a strong value)
- `DATA_ENCRYPTION_KEY` — AES key material for encrypting secrets at rest
- `RSA_PRIVATE_KEY` — RSA private key PEM (store as single-line with `\n` escapes)

Ports and networking:

- `NOFX_BACKEND_PORT` — host port mapping for Docker backend (`docker-compose*.yml`)
- `NOFX_FRONTEND_PORT` — host port mapping for Docker frontend (`docker-compose*.yml`)
- `API_SERVER_PORT` — backend listen port when running _without_ Docker (defaults to `8080`)
- `TZ` — timezone used inside containers (compose defaults to `Asia/Shanghai`)

Security/behavior flags:

- `TRANSPORT_ENCRYPTION` — when `true`, frontend encrypts sensitive payloads before sending (requires `https://` or `http://localhost`)
- `REGISTRATION_ENABLED` — enable/disable new user registration
- `MAX_USERS` — max number of users (0 = unlimited, default 10)
- `EXPERIENCE_IMPROVEMENT` — set to `false` to disable anonymous usage statistics

AI/backtest runtime:

- `AI_MAX_TOKENS` — default max tokens used by `mcp/` clients (also set in compose)
- `DEEPSEEK_API_KEY` — used by `main.go` for creating the shared MCP client in some flows

Testing-only keys:

- `COINANK_API_KEY` — used by `provider/coinank/*_test.go` (tests skip if unset)

## Architecture

At a high level:

- The React SPA calls the Go backend under `/api`.
- The backend exposes endpoints for authentication, configuration (models/exchanges/strategies/traders), monitoring, debate sessions, and backtesting.
- Trading state and history are persisted to SQLite (`data/data.db` by default).
- Sensitive values are encrypted for storage (AES); optional end-to-end transport encryption uses RSA/WebCrypto.

Key runtime flow (simplified):

1. Market data collection (`market/`) supplies prices/klines.
2. Decision engine (`decision/`) constructs prompts and parses AI outputs.
3. AI provider calls go through `mcp/` clients.
4. Execution goes through exchange adapters (`trader/`), orchestrated by `manager/`.
5. Results (positions, equity, decisions) are stored in `store/` and displayed by the UI.

## Contributing

- Follow `CONTRIBUTING.md` (coding standards, PR guidelines, conventional commit messages).
- Do not commit secrets (`.env`, API keys, `data/data.db`).
- Run formatting and tests before opening a PR: `go fmt ./...` and `make test`.

## License

This project is licensed under the **GNU Affero General Public License v3.0 (AGPL-3.0)**.

See `LICENSE` in the repository root.

---
