# AGENTS.md - NOFX Project Documentation

## Project Name

**NOFX** - Agentic Trading OS: AI-Powered Decentralized Crypto Trading Platform

## Overview

**NOFX** is an innovative open-source AI trading system that enables users to run multiple AI models (DeepSeek, Qwen, GPT, Claude, Gemini, Grok, Kimi) to trade cryptocurrency futures automatically across multiple exchanges (Binance, Bybit, OKX, Hyperliquid, Aster DEX, Lighter).

The platform features a web-based configuration interface that eliminates the need for manual JSON editing, a real-time dashboard for monitoring performance, and a unique AI competition mode where multiple AI traders compete simultaneously to find optimal trading strategies. Users can configure strategies through a visual Strategy Studio with support for various technical indicators, risk controls, and coin sources.

This system is experimental and recommended for learning, research, or small-scale testing due to the inherent risks of automated AI trading.

## Technology Stack

### Backend

- **Language**: Go 1.25.0
- **API Framework**: Gin-gonic (HTTP server & routing)
- **Database**: SQLite (modernc.org/sqlite)
- **Authentication**: JWT (golang-jwt/jwt/v5)
- **Crypto**: AES-256 & RSA encryption (golang.org/x/crypto)
- **WebSocket**: gorilla/websocket for real-time updates
- **Exchange APIs**:
  - Binance: go-binance/v2
  - Bybit: bybit.go.api
  - OKX: Built-in implementation
  - Hyperliquid: go-hyperliquid
  - Lighter: Custom implementation
- **Technical Analysis**: TA-Lib integration
- **Logging**: RS/zerolog, sirupsen/logrus
- **Environment**: godotenv for .env support

### Frontend

- **Language**: TypeScript 5.8+
- **Framework**: React 18.3+
- **Build Tool**: Vite 6.0+
- **Styling**: Tailwind CSS 3.4+
- **UI Components**: Radix UI
- **State Management**: Zustand 5.0+
- **HTTP Client**: Axios 1.13+
- **Data Fetching**: SWR 2.2+
- **Charts**: Recharts 2.15+
- **Router**: React Router DOM 7.9+
- **Linting**: ESLint 9.39+ with TypeScript support
- **Code Formatting**: Prettier 3.6+
- **Testing**: Vitest 2.1+
- **Dev Tools**: Husky + lint-staged for git hooks

### DevOps & Infrastructure

- **Containerization**: Docker & Docker Compose
- **Reverse Proxy**: Nginx (for frontend serving)
- **Encryption**: Web Crypto API for browser-side encryption
- **Security**: HTTPS support with Cloudflare integration

## Project Structure

```
nofx/
├── main.go                      # Application entry point
├── go.mod, go.sum              # Go module dependencies
├── Makefile                     # Build, test, and dev commands
├── docker-compose.yml           # Development Docker setup
├── docker-compose.prod.yml      # Production Docker setup
├── .env.example                 # Environment variables template
│
├── api/                         # REST API handlers & endpoints
│   ├── server.go               # Gin server setup & routing
│   ├── strategy.go             # Strategy configuration endpoints
│   ├── backtest.go             # Backtesting logic
│   ├── debate.go               # AI debate/competition endpoints
│   ├── crypto_handler.go       # Encryption/decryption handlers
│   └── utils.go                # Utility functions
│
├── auth/                        # Authentication module
│   └── auth.go                 # JWT & OTP authentication logic
│
├── config/                      # Configuration management
│   └── config.go               # .env parsing & config initialization
│
├── store/                       # Database layer (SQLite)
│   ├── store.go                # Core store interface
│   ├── user.go                 # User data persistence
│   ├── trader.go               # Trader configuration storage
│   ├── strategy.go             # Strategy persistence
│   ├── exchange.go             # Exchange credentials storage
│   ├── ai_model.go             # AI model configuration storage
│   ├── position.go             # Position data
│   ├── equity.go               # Equity tracking
│   ├── backtest.go             # Backtest results
│   ├── debate.go               # Debate/competition data
│   └── decision.go             # AI decision logs
│
├── trader/                      # Exchange trading implementations
│   ├── interface.go            # Trader interface definition
│   ├── binance_futures.go      # Binance Futures trader
│   ├── bybit_trader.go         # Bybit perpetual trader
│   ├── okx_trader.go           # OKX trader
│   ├── hyperliquid_trader.go   # Hyperliquid DEX trader
│   ├── aster_trader.go         # Aster DEX trader
│   ├── lighter_trader.go       # Lighter DEX trader (v1)
│   ├── lighter_trader_v2.go    # Lighter DEX trader (v2)
│   ├── auto_trader.go          # Universal trader wrapper
│   ├── position_sync.go        # Position synchronization
│   ├── position_rebuild.go     # Position recovery logic
│   └── helpers.go              # Shared trader utilities
│
├── market/                      # Market data & analysis
│   ├── api_client.go           # Exchange market data API
│   ├── websocket_client.go     # WebSocket market streams
│   ├── combined_streams.go     # Multi-stream handling
│   ├── monitor.go              # Real-time price monitoring
│   ├── historical.go           # Historical data fetching
│   ├── timeframe.go            # Timeframe management
│   ├── data.go                 # OHLCV data structures
│   └── types.go                # Market data types
│
├── manager/                     # Trader lifecycle management
│   └── trader_manager.go       # Start/stop/monitor traders
│
├── logger/                      # Logging module
│   └── Unified logging interface
│
├── crypto/                      # Cryptographic utilities
│   └── AES & RSA encryption implementation
│
├── hook/                        # Webhook handlers
│   └── External notification integration
│
├── pool/                        # Resource pooling
│   └── Connection and resource management
│
├── decision/                    # AI decision recording
│   └── Chain of Thought logging
│
├── debate/                      # AI competition module
│   └── Multi-AI comparison and leaderboard
│
├── backtest/                    # Backtesting engine
│   └── Strategy historical analysis
│
├── migrations/                  # Database schema migrations
│   └── Database initialization scripts
│
├── docker/                      # Docker configurations
│   ├── Dockerfile.backend      # Backend container image
│   └── Dockerfile.frontend     # Frontend container image
│
├── nginx/                       # Nginx reverse proxy config
│   └── Proxy & static file serving
│
├── web/                         # React/TypeScript frontend application
│   ├── src/
│   │   ├── pages/              # Page components
│   │   │   ├── LandingPage.tsx
│   │   │   ├── StrategyStudioPage.tsx
│   │   │   ├── DebateArenaPage.tsx
│   │   │   └── FAQPage.tsx
│   │   ├── components/         # Reusable UI components
│   │   ├── hooks/              # React custom hooks
│   │   ├── stores/             # Zustand state management
│   │   ├── contexts/           # React context providers
│   │   ├── utils/              # Utility functions
│   │   ├── lib/                # Library helpers
│   │   ├── constants/          # Application constants
│   │   ├── i18n/               # Internationalization (i18n)
│   │   ├── types.ts            # TypeScript type definitions
│   │   ├── App.tsx             # Root component
│   │   └── main.tsx            # Entry point
│   ├── package.json            # Frontend dependencies
│   ├── vite.config.ts          # Vite build configuration
│   ├── tsconfig.json           # TypeScript configuration
│   ├── tailwind.config.js      # Tailwind CSS configuration
│   ├── eslint.config.js        # ESLint rules
│   └── vitest.config.ts        # Testing configuration
│
├── docs/                        # Documentation
│   ├── api/                    # API documentation
│   ├── architecture/           # Architecture guides
│   ├── getting-started/        # Quick start guides
│   ├── guides/                 # How-to guides
│   ├── i18n/                   # Internationalization docs
│   ├── prompt-guide.md         # AI prompt engineering guide
│   ├── pnl.md                  # P&L calculation documentation
│   ├── MIGRATION_GUIDE.md      # Version migration guide
│   └── roadmap/                # Project roadmap
│
├── scripts/                     # Utility scripts
│   └── Automation & deployment scripts
│
├── install.sh                   # One-click installation script
├── start.sh                     # Application startup script
├── README.md                    # Main project documentation
├── README.ja.md                 # Japanese documentation
├── CONTRIBUTING.md             # Contributing guidelines
├── CODE_OF_CONDUCT.md         # Community code of conduct
├── SECURITY.md                 # Security policy
├── LICENSE                     # AGPL-3.0 license
├── CHANGELOG.md                # Version history
└── DISCLAIMER.md               # Risk disclaimer
```

## Key Features

- **Multi-AI Support**: Deploy and switch between DeepSeek, Qwen, GPT, Claude, Gemini, Grok, Kimi instantly
- **Multi-Exchange Trading**: Trade on 6 exchanges simultaneously (Binance, Bybit, OKX, Hyperliquid, Aster, Lighter)
- **Web-Based Configuration**: Visual UI for all settings—no JSON file editing required
- **AI Competition Mode**: Multiple AI traders compete in real-time with live leaderboard tracking
- **Real-Time Dashboard**: TradingView-style candlestick charts, position management, Chain of Thought AI decision logs
- **Strategy Studio**: Visual strategy builder with:
  - Coin sources (static lists, AI500 pool, OI Top)
  - Technical indicators (EMA, MACD, RSI, ATR, Volume, OI, Funding Rate)
  - Risk controls (leverage limits, position sizing, margin usage)
  - Real-time AI testing with prompt preview
- **Backtesting Engine**: Historical strategy performance analysis
- **Real-Time P&L Tracking**: Equity curves and performance metrics
- **Debate Arena**: AI agents discuss and analyze trading decisions
- **Encryption**: Optional transport encryption with Web Crypto API
- **Multi-Language Support**: English, Chinese, Japanese, Korean, Russian, Ukrainian, Vietnamese

## Getting Started

### Prerequisites

#### System Requirements

- **macOS/Linux**: Bash shell
- **Windows**: WSL2 or manual installation
- **RAM**: 2GB minimum (4GB+ recommended for running multiple AI traders)
- **Disk**: 500MB+ for database and logs

#### Required Software

- **Go 1.21+**
- **Node.js 18+** (with npm)
- **TA-Lib** (technical analysis library)
- **Docker & Docker Compose** (for containerized deployment)

### Installation

#### Option 1: One-Click Install (Recommended - Linux/macOS)

```bash
curl -fsSL https://raw.githubusercontent.com/NoFxAiOS/nofx/main/install.sh | bash
```

Then access the web interface at: **http://localhost:3000**

#### Option 2: Docker Compose (All Platforms)

```bash
# Download Docker Compose file
curl -O https://raw.githubusercontent.com/NoFxAiOS/nofx/main/docker-compose.prod.yml

# Start services
docker compose -f docker-compose.prod.yml up -d

# View logs
docker compose -f docker-compose.prod.yml logs -f

# Stop services
docker compose -f docker-compose.prod.yml down

# Update to latest version
docker compose -f docker-compose.prod.yml pull && docker compose -f docker-compose.prod.yml up -d
```

Access web interface at: **http://localhost:3000**

#### Option 3: Manual Development Installation

##### macOS

```bash
# Install TA-Lib
brew install ta-lib

# Clone repository
git clone https://github.com/NoFxAiOS/nofx.git
cd nofx

# Install backend dependencies
go mod download

# Install frontend dependencies
cd web
npm install
cd ..

# Build backend
go build -o nofx

# Start backend (in terminal 1)
./nofx

# Start frontend (in terminal 2)
cd web
npm run dev
```

##### Ubuntu/Debian

```bash
# Install TA-Lib
sudo apt-get install libta-lib0-dev

# Clone repository
git clone https://github.com/NoFxAiOS/nofx.git
cd nofx

# Install backend dependencies
go mod download

# Install frontend dependencies
cd web
npm install
cd ..

# Build backend
go build -o nofx

# Start backend (in terminal 1)
./nofx

# Start frontend (in terminal 2)
cd web
npm run dev
```

##### Windows (WSL2)

```bash
# Inside WSL2 Ubuntu environment, follow Ubuntu/Debian instructions above
```

### Initial Configuration

After starting NOFX, complete the web interface setup:

1. **Configure AI Models** - Add API keys for your chosen AI providers
   - Visit OpenAI, Anthropic, Deepseek, Qwen, Gemini, etc. to get keys
2. **Configure Exchanges** - Set up exchange API credentials
   - Create API keys with appropriate permissions on each exchange
3. **Create Strategy** - Build trading strategy in Strategy Studio
   - Define coin sources, technical indicators, and risk parameters
4. **Create Trader** - Combine AI model + Exchange + Strategy
   - Name your trader and link components
5. **Start Trading** - Launch traders and monitor performance
   - View real-time P&L and AI decision logs

### Usage Examples

#### Using Docker (Easiest)

```bash
# Development with file watching
docker compose up -d

# View real-time logs
docker compose logs -f nofx

# Access dashboard
# Open http://localhost:3000 in browser
```

#### Using Make Commands

```bash
# Start backend
make run

# Start frontend (new terminal)
make run-frontend

# Run tests
make test

# Build for production
make build
make build-frontend

# Clean build artifacts
make clean
```

#### Using npm (Frontend)

```bash
# Development mode with HMR
cd web
npm run dev

# Production build
npm run build

# Preview production build
npm run preview

# Run tests
npm run test

# Lint code
npm run lint

# Format code
npm run format
```

## Development

### Available Scripts

#### Go Backend

```bash
# Run backend in development
go run main.go

# Build binary
go build -o nofx

# Run tests with coverage
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Format code
go fmt ./...

# Lint (requires golangci-lint)
golangci-lint run
```

#### TypeScript/React Frontend

```bash
cd web

# Development with hot reload
npm run dev

# Build for production
npm run build

# Run tests
npm run test

# Lint code
npm run lint

# Fix linting issues
npm run lint:fix

# Format code
npm run format

# Check formatting
npm run format:check
```

#### Make Commands (Convenient Wrapper)

```bash
make help          # Show all available commands
make test          # Run all tests
make test-backend  # Backend tests only
make test-frontend # Frontend tests only
make test-coverage # Generate coverage report
make build         # Build backend
make build-frontend # Build frontend
make clean         # Clean build artifacts
make run           # Run backend
make run-frontend  # Run frontend dev server
make fmt           # Format Go code
make lint          # Lint Go code
make docker-build  # Build Docker images
make docker-up     # Start Docker containers
make docker-down   # Stop Docker containers
make docker-logs   # View Docker logs
make deps          # Download Go dependencies
make deps-update   # Update Go dependencies
make deps-frontend # Install frontend dependencies
```

### Development Workflow

1. **Setup**: Clone repo and install dependencies

   ```bash
   git clone https://github.com/NoFxAiOS/nofx.git
   cd nofx
   make deps
   make deps-frontend
   ```

2. **Create .env file**: Copy from `.env.example` and configure

   ```bash
   cp .env.example .env
   # Edit .env with your settings
   ```

3. **Development**: Run backend and frontend in separate terminals

   ```bash
   # Terminal 1: Backend
   make run

   # Terminal 2: Frontend
   make run-frontend
   ```

4. **Testing**: Write tests alongside code

   ```bash
   # Run tests
   make test

   # Generate coverage
   make test-coverage
   ```

5. **Code Quality**:

   ```bash
   # Format Go code
   make fmt

   # Lint Go code
   make lint

   # Format TypeScript/React
   cd web && npm run lint:fix && npm run format
   ```

6. **Submit Changes**:
   - Create feature branch: `git checkout -b feature/your-feature`
   - Commit with clear messages (see Commit Guidelines)
   - Push and create Pull Request
   - Reference related issues in PR description

### Testing Strategy

**Backend Testing**:

- Unit tests for each Go package
- Integration tests for trader implementations
- Race condition tests for concurrent operations
- Mock external API calls

**Frontend Testing**:

- Component tests with React Testing Library
- Vitest runner for fast execution
- ESLint for code quality
- Prettier for code formatting

Run tests before committing:

```bash
make test
```

## Configuration

### Environment Variables (.env)

**Server Configuration**

```bash
NOFX_BACKEND_PORT=8080        # Backend API port
NOFX_FRONTEND_PORT=3000        # Frontend web server port
NOFX_TIMEZONE=Asia/Shanghai    # Application timezone
```

**Authentication (Required)**

```bash
JWT_SECRET=<32+ character random string>  # JWT signing secret
# Generate: openssl rand -base64 32
```

**Encryption Keys (Required)**

```bash
DATA_ENCRYPTION_KEY=<32-byte base64 encoded>  # AES-256 key for database
RSA_PRIVATE_KEY=<PEM format RSA key>         # Client-server encryption
# Generate: openssl rand -base64 32
# Generate: openssl genrsa 2048
```

**Security Options**

```bash
TRANSPORT_ENCRYPTION=false    # Enable browser-side encryption of API keys
# true = HTTPS required, false = HTTP/IP access allowed
```

**Optional Services**

```bash
TELEGRAM_BOT_TOKEN=<token>    # Telegram notifications
TELEGRAM_CHAT_ID=<chat_id>    # Where to send notifications
```

### Database Configuration

- **Location**: `data/data.db` (SQLite)
- **Auto-migration**: Schemas are created automatically on startup
- **Volume mounting**: In Docker, mounted to `/app/data`

### Logging Configuration

- **Backend**: RS/Zerolog with structured logging
- **Frontend**: Browser console and network requests
- **Log level**: Configured via code (Info by default)

## Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────┐
│         Web Browser (React + TypeScript)        │
│  - Strategy Studio                              │
│  - Real-time Dashboard                          │
│  - AI Competition Leaderboard                   │
└─────────────┬─────────────────────────────────┘
              │ HTTP/WebSocket
              ↓
┌─────────────────────────────────────────────────┐
│         Nginx Reverse Proxy                     │
│  - Static file serving                          │
│  - API proxying                                 │
│  - SSL/TLS termination (optional)               │
└─────────────┬─────────────────────────────────┘
              │ HTTP
              ↓
┌─────────────────────────────────────────────────┐
│      Go Backend (REST API + WebSocket)          │
│                                                 │
│  ┌─────────────────────────────────────────┐  │
│  │  API Handler Layer (api/)               │  │
│  │  - REST endpoints                       │  │
│  │  - WebSocket connections                │  │
│  └─────────────────────────────────────────┘  │
│                    ↓                            │
│  ┌─────────────────────────────────────────┐  │
│  │  Business Logic Layer                   │  │
│  │  - Trader Manager (manager/)            │  │
│  │  - Strategy Engine                      │  │
│  │  - Backtest Engine (backtest/)          │  │
│  │  - Debate/Competition (debate/)         │  │
│  │  - AI Decision Logging (decision/)      │  │
│  └─────────────────────────────────────────┘  │
│                    ↓                            │
│  ┌─────────────────────────────────────────┐  │
│  │  Exchange Layer (trader/)               │  │
│  │  - Binance Futures                      │  │
│  │  - Bybit Perpetuals                     │  │
│  │  - OKX Futures                          │  │
│  │  - Hyperliquid DEX                      │  │
│  │  - Aster DEX                            │  │
│  │  - Lighter DEX                          │  │
│  │  - Position sync & recovery             │  │
│  └─────────────────────────────────────────┘  │
│                    ↓                            │
│  ┌─────────────────────────────────────────┐  │
│  │  Market Data Layer (market/)            │  │
│  │  - Real-time price feeds                │  │
│  │  - WebSocket streams                    │  │
│  │  - Historical data                      │  │
│  │  - Technical indicators (TA-Lib)        │  │
│  └─────────────────────────────────────────┘  │
│                    ↓                            │
│  ┌─────────────────────────────────────────┐  │
│  │  Data Layer (store/)                    │  │
│  │  - SQLite database                      │  │
│  │  - User management                      │  │
│  │  - Trader configuration                 │  │
│  │  - Strategy persistence                 │  │
│  │  - Position tracking                    │  │
│  │  - Equity curve                         │  │
│  └─────────────────────────────────────────┘  │
│                    ↓                            │
│  ┌─────────────────────────────────────────┐  │
│  │  Cross-Cutting Concerns                 │  │
│  │  - Auth & encryption (auth/, crypto/)   │  │
│  │  - Logging (logger/)                    │  │
│  │  - Configuration (config/)              │  │
│  │  - Webhooks (hook/)                     │  │
│  └─────────────────────────────────────────┘  │
└─────────────────────────────────────────────────┘
                    ↓
        ┌──────────────────────────┐
        │  External Exchange APIs  │
        │  - WebSocket feeds       │
        │  - REST endpoints        │
        │  - Order execution       │
        └──────────────────────────┘
                    ↓
        ┌──────────────────────────┐
        │   External AI APIs       │
        │  - OpenAI/Claude/etc     │
        │  - Deepseek/Qwen/etc     │
        └──────────────────────────┘
```

### Key Design Patterns

1. **Trader Interface Pattern**: Exchange implementations share common interface
2. **Manager Pattern**: TraderManager handles trader lifecycle (start, stop, monitor)
3. **Store Pattern**: Database abstraction layer for all persistence
4. **State Management**: Zustand for frontend state, no Redux complexity
5. **Component Composition**: Modular React components with clear boundaries
6. **Error Handling**: Centralized error responses with meaningful messages

### Data Flow

1. **Trading Loop**:

   - Market data fetched via WebSocket/API
   - Technical indicators calculated (TA-Lib)
   - AI model receives prompt with market data
   - AI generates trading signal
   - Trader executes order on exchange
   - Position updated in database
   - P&L calculated and displayed

2. **Competition Mode**:
   - Multiple traders run simultaneously
   - Results compared in real-time
   - Leaderboard updated continuously
   - Decision logs stored for analysis

## Contributing

### How to Contribute

Contributions are welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for details.

**Contribution Types**:

- **Code commits** - Bug fixes, features, improvements
- **Bug reports** - Detailed issue descriptions
- **Feature suggestions** - New functionality ideas
- **Documentation** - README, guides, translations
- **Testing** - Test cases, edge cases
- **Community** - Help others, answer questions

### Development Requirements

- Follow [Coding Standards](CONTRIBUTING.md#coding-standards)
- Write tests for new code
- Commit messages follow [Guidelines](CONTRIBUTING.md#commit-message-guidelines)
- Create PR with issue references
- Pass CI/CD checks

### Contributor Airdrop Program

All contributions are tracked. When NOFX generates revenue, contributors receive airdrops:

| Contribution Type             |    Weight    |
| ----------------------------- | :----------: |
| **Pinned Issue PRs**          | ⭐⭐⭐⭐⭐⭐ |
| **Code Commits (Merged PRs)** |  ⭐⭐⭐⭐⭐  |
| **Bug Fixes**                 |   ⭐⭐⭐⭐   |
| **Feature Suggestions**       |    ⭐⭐⭐    |
| **Bug Reports**               |     ⭐⭐     |
| **Documentation**             |     ⭐⭐     |

**PRs resolving Pinned Issues receive the HIGHEST rewards!**

## Additional Resources

### Documentation Structure

- **[Architecture Overview](docs/architecture/README.md)** - Complete system design
- **[Strategy Module](docs/architecture/STRATEGY_MODULE.md)** - Coin selection, indicators, AI prompts
- **[Backtest Module](docs/architecture/BACKTEST_MODULE.md)** - Historical simulation, metrics
- **[Debate Module](docs/architecture/DEBATE_MODULE.md)** - Multi-AI debate mechanism
- **[FAQ](docs/faq/README.md)** - Common questions and troubleshooting
- **[Security Policy](SECURITY.md)** - Vulnerability reporting guidelines
- **[Code of Conduct](CODE_OF_CONDUCT.md)** - Community guidelines

### External Resources

- **[Prompt Engineering Guide](docs/prompt-guide.md)** - AI prompt optimization
- **[P&L Documentation](docs/pnl.md)** - P&L calculation methodology
- **[Migration Guide](docs/MIGRATION_GUIDE.md)** - Version upgrade instructions
- **[Project Roadmap](docs/roadmap/README.md)** - Future development plans

## Common Deployment Patterns

### Pattern 1: Local Development (Single Machine)

```bash
# All-in-one development environment
make run &              # Terminal 1: Backend on :8080
cd web && npm run dev  # Terminal 2: Frontend on :3000
# Access: http://localhost:3000
```

### Pattern 2: Docker Compose (Recommended for Servers)

```bash
# Production-ready with Nginx + Docker
docker compose -f docker-compose.prod.yml up -d
# Access: http://server-ip:3000 or https://yourdomain.com
```

### Pattern 3: Advanced (HTTPS with Cloudflare)

```bash
# Secure production deployment
# 1. Add domain to Cloudflare
# 2. Set DNS A record to server IP
# 3. Enable "Proxied" in Cloudflare DNS
# 4. Set .env: TRANSPORT_ENCRYPTION=true
# 5. Access via https://yourdomain.com
```

## Troubleshooting Guide

### Backend Issues

**TA-Lib not found**

```bash
# macOS
brew install ta-lib

# Ubuntu/Debian
sudo apt-get install libta-lib0-dev

# WSL2
wsl sudo apt-get install libta-lib0-dev
```

**Port already in use**

```bash
# Find process using port 8080
lsof -i :8080
# Kill it or use different port in .env
export NOFX_BACKEND_PORT=8081
```

**Database locked error**

```bash
# Ensure only one instance is running
pkill -f "go run main.go"
rm -f data/data.db
# Restart application
```

### Frontend Issues

**Node modules not found**

```bash
cd web
rm -rf node_modules package-lock.json
npm install
npm run dev
```

**WebSocket connection fails**

```bash
# Ensure backend is running on http://localhost:8080
# Check firewall settings
# Verify CORS in api/server.go if customizing ports
```

**Hot reload not working**

```bash
# Vite cache issue
cd web
rm -rf .vite
npm run dev
```

### Docker Issues

**Container fails to start**

```bash
# Check logs
docker compose logs -f nofx

# Rebuild images
docker compose build --no-cache
docker compose up -d
```

**Port conflicts**

```bash
# Edit docker-compose.yml to use different ports
# ports:
#   - "3001:3000"   # Changed from 3000:3000
```

## Performance Optimization Tips

### Backend Optimization

- Monitor memory usage: `watch -n 1 'docker stats'`
- Enable database query caching for historical data
- Use WebSocket for real-time updates instead of polling
- Increase timeout for slow API connections: adjust in `.env`

### Frontend Optimization

- Use React DevTools to identify re-renders
- Lazy load chart components for better performance
- Zustand store is already optimized (no Redux overhead)
- Run `npm run build` and test production build locally

### Database Optimization

- Backup regularly: `cp data/data.db data/data.db.backup`
- Monitor database size: `du -h data/data.db`
- Clean old backtest records periodically
- Use indexes for frequently queried fields

## License

This project is licensed under the **GNU Affero General Public License v3.0 (AGPL-3.0)**

- **Source**: See [LICENSE](LICENSE) file
- **Summary**:
  - Free for private use
  - Free for non-commercial use
  - If modified and deployed as a service, modifications must be shared
  - Suitable for learning and research
- **Legal**: You MUST provide access to source code if you deploy this as a service

## Project Links

- **GitHub**: [NoFxAiOS/nofx](https://github.com/NoFxAiOS/nofx)
- **Issues & Roadmap**: [GitHub Issues](https://github.com/NoFxAiOS/nofx/issues)
- **Developer Community**: [Telegram Group](https://t.me/nofx_dev_community)
- **Twitter**: [@nofx_official](https://x.com/nofx_official)
- **Backed by**: [Amber.ac](https://amber.ac)

## Disclaimer

⚠️ **Risk Warning**: This system is experimental. AI auto-trading carries significant financial risks. Use for:

- Learning and research purposes
- Testing with small amounts
- Paper trading before real money

Not recommended for production trading without thorough testing and risk management.

---

**Last Updated**: December 2024  
**Current Version**: 3.0.0  
**Latest Changes**: Web-based configuration, SQLite database integration, multi-exchange support  
**Maintainers**: NOFX Core Team  
**Status**: Active Development
