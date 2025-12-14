# PHD - Device Control System

Hệ thống giám sát và điều khiển máy tính từ xa sử dụng Hedera Blockchain.

## Tổng quan

Bot theo dõi máy tính: nhân viên khi được cấp máy sẽ phải cài đặt app này. App có khả năng thay đổi hình nền/hình nền khóa khi nhận lệnh từ server. Áp dụng blockchain để lưu trữ lệnh trên smart contract Hedera.

### Kiến trúc

- **Smart Contract (Solidity)**: Lưu trữ và phát lệnh trên Hedera blockchain
- **Backend (NestJS)**: API server để quản lý lệnh và tương tác với blockchain
- **Client**: Desktop application lắng nghe và thực thi lệnh
- **Shared Library**: Types và utilities dùng chung

### Luồng hoạt động

1. Admin gọi hàm `Trigger()` trên smart contract
2. Smart contract emit event `CommandTriggered`
3. Client apps query interval để kiểm tra event
4. Khi phát hiện event mới, client gọi `GetFunction()` để lấy command
5. Command có thể là script hoặc URL:
   - **SCRIPT**: Thực thi script trực tiếp
   - **URL**: Curl URL rồi set desktop screen

## Cấu trúc Monorepo

```
phd-monorepo/
├── apps/
│   ├── backend/              # NestJS API server
│   └── smart-contract/       # Foundry smart contract
├── libs/
│   └── shared/              # Shared TypeScript library
├── nx.json                  # Nx configuration
├── package.json             # Root package.json
└── pnpm-workspace.yaml      # pnpm workspace config
```

## Prerequisites

### Option 1: Docker (Recommended)
- Docker & Docker Compose

### Option 2: Local Development
- Node.js >= 18
- pnpm >= 8
- PostgreSQL >= 14
- Foundry (for smart contract development)

## Installation

### Using Docker (Recommended) 🐳

```bash
# 1. Copy environment file
cp .env.docker .env

# 2. Edit .env with your blockchain config
# Change JWT_SECRET in production!

# 3. Start all services (PostgreSQL + Backend)
docker-compose up -d

# 4. View logs
docker-compose logs -f
```

**✅ That's it!** Services will:
- Start PostgreSQL on port 5432
- Run database migrations automatically
- Seed default admin account
- Start backend on port 3000

**Access**:
- API: http://localhost:3000/api
- Swagger Docs: http://localhost:3000/api/docs
- Default Admin: `username: admin`, `password: Admin@123`

📖 **Full Docker guide**: See [README.docker.md](README.docker.md)

### Local Development

```bash
# Install dependencies
pnpm install

# Install Foundry (if not installed)
curl -L https://foundry.paradigm.xyz | bash
foundryup

# Setup database
createdb phd_admin

# Copy and configure .env
cp .env.example .env

# Run migrations
pnpm migration:run

# Seed admin account
pnpm seed

# Start backend
pnpm dev:backend
```

## Quick Start

**📖 Xem [GETTING_STARTED.md](GETTING_STARTED.md) để có hướng dẫn chi tiết từng bước chạy backend.**

**🐳 Dùng Docker**: Xem [README.docker.md](README.docker.md)

Workflow nhanh:
1. Build contract: `pnpm contract:build`
2. Export ABI: `pnpm contract:export-abi`
3. Deploy contract: `pnpm contract:deploy:testnet`
4. Configure `.env` với CONTRACT_ADDRESS
5. Run backend: `pnpm dev:backend` hoặc `docker-compose up`

## Configuration

### Backend Configuration

Copy `.env.example` to `.env` and configure:

```bash
# Blockchain
BLOCKCHAIN_NETWORK=testnet
CONTRACT_ADDRESS=your_contract_address
PRIVATE_KEY=your_private_key

# Server
PORT=3000
NODE_ENV=development
```

### Smart Contract Configuration

Copy `apps/smart-contract/.env.example` to `apps/smart-contract/.env`:

```bash
PRIVATE_KEY=your_deployment_private_key
RPC_URL=https://testnet.hashio.io/api
```

## Development

### Start Backend

```bash
pnpm dev:backend
```

### Test Smart Contract

```bash
pnpm test:contract

# Watch mode
pnpm test:contract:watch
```

### Build All Projects

```bash
pnpm build
```

## Smart Contract

### Build Contract

```bash
pnpm contract:build
```

### Deploy Contract

```bash
# Deploy to local network
pnpm contract:deploy:local

# Deploy to Hedera testnet
pnpm contract:deploy:testnet
```

### Verify Contract

Sau khi deploy, verify contract trên block explorer:

```bash
# Thêm CONTRACT_ADDRESS vào .env trước
# Verify on testnet
pnpm contract:verify

# Verify on mainnet
pnpm contract:verify:mainnet

# Xem thông tin verification
pnpm contract:verify:info
```

Hoặc sử dụng script trực tiếp:
```bash
cd apps/smart-contract
./verify.sh <contract_address> testnet
```

Chi tiết xem [VERIFY.md](apps/smart-contract/VERIFY.md)

### Smart Contract Functions

- `Trigger(commandType, data)`: Admin triggers new command
- `GetFunction()`: Get current command details
- `getLatestCommandId()`: Get latest command ID
- `transferAdmin(newAdmin)`: Transfer admin role

## Backend API

### API Documentation

Full interactive API documentation available at: **http://localhost:3000/api/docs** (Swagger UI)

### Authentication Endpoints

#### Login
```bash
POST /api/auth/login
{
  "username": "admin",
  "password": "Admin@123"
}

Response:
{
  "accessToken": "eyJhbGciOiJIUzI1...",
  "admin": {
    "id": "uuid",
    "username": "admin",
    "email": "admin@phd.local"
  }
}
```

#### Get Profile (Protected)
```bash
GET /api/auth/profile
Headers: Authorization: Bearer <token>
```

### Scripts Management (All Protected)

#### Create Script
```bash
POST /api/scripts
Headers: Authorization: Bearer <token>
{
  "name": "Hello World Script",
  "description": "A simple script",
  "jsonData": {
    "commandType": 0,
    "data": "console.log('Hello')",
    "metadata": {}
  }
}
```

#### List Scripts
```bash
GET /api/scripts?page=1&limit=20&search=hello
Headers: Authorization: Bearer <token>
```

#### Get Script
```bash
GET /api/scripts/:id
Headers: Authorization: Bearer <token>
```

#### Update Script
```bash
PATCH /api/scripts/:id
Headers: Authorization: Bearer <token>
{
  "name": "Updated Name"
}
```

#### Delete Script
```bash
DELETE /api/scripts/:id
Headers: Authorization: Bearer <token>
```

#### Trigger Script on Blockchain
```bash
POST /api/scripts/:id/trigger
Headers: Authorization: Bearer <token>

Response:
{
  "success": true,
  "transactionHash": "0x123...",
  "script": { "id": "...", "name": "..." }
}
```

### Events (All Protected)

#### List Events
```bash
GET /api/events?page=1&limit=20&commandType=0
Headers: Authorization: Bearer <token>
```

#### Get Event
```bash
GET /api/events/:id
Headers: Authorization: Bearer <token>
```

#### Get Statistics
```bash
GET /api/events/stats
Headers: Authorization: Bearer <token>
```

### Legacy Commands Endpoints

#### Trigger Command (Protected)
```bash
POST /api/commands/trigger
Headers: Authorization: Bearer <token>
{
  "commandType": 0,  # 0=SCRIPT, 1=URL
  "data": "script or url"
}
```

#### Get Current Command (Public)
```bash
GET /api/commands/current
```

#### Get Latest Command ID (Public)
```bash
GET /api/commands/latest-id
```

#### Get Command History (Public)
```bash
GET /api/commands/history
```

#### Get Blockchain Info (Public)
```bash
GET /api/commands/info
```

## Scripts

### Development

| Command | Description |
|---------|-------------|
| `pnpm dev` | Start backend in dev mode |
| `pnpm dev:backend` | Start backend only |
| `pnpm build` | Build all projects |
| `pnpm build:backend` | Build backend only |
| `pnpm test` | Run all tests |
| `pnpm lint` | Lint all projects |

### Database

| Command | Description |
|---------|-------------|
| `pnpm migration:run` | Run database migrations |
| `pnpm migration:revert` | Revert last migration |
| `pnpm migration:generate` | Generate new migration |
| `pnpm seed` | Seed default admin account |

### Smart Contract

| Command | Description |
|---------|-------------|
| `pnpm contract:build` | Build smart contract |
| `pnpm contract:test` | Test smart contract |
| `pnpm contract:deploy:testnet` | Deploy to Hedera testnet |
| `pnpm contract:verify` | Verify contract on testnet |
| `pnpm contract:verify:mainnet` | Verify contract on mainnet |

### Docker

| Command | Description |
|---------|-------------|
| `docker-compose up -d` | Start services in background |
| `docker-compose logs -f` | View logs |
| `docker-compose down` | Stop services |
| `docker-compose down -v` | Stop and remove volumes |

### Other

| Command | Description |
|---------|-------------|
| `pnpm graph` | View dependency graph |

## Technology Stack

- **Monorepo**: Nx
- **Backend**: NestJS, TypeScript
- **Database**: PostgreSQL, TypeORM
- **Authentication**: JWT, Passport, bcrypt
- **API Documentation**: Swagger/OpenAPI
- **Validation**: class-validator, class-transformer
- **Blockchain**: Hedera, Solidity, Foundry
- **Smart Contract**: ethers.js v6
- **Package Manager**: pnpm
- **Containerization**: Docker, Docker Compose

## Project Structure Details

### apps/backend
NestJS application với các modules:
- **DatabaseModule**: PostgreSQL connection & configuration
- **AuthModule**: JWT authentication & authorization
- **AdminModule**: Admin user management
- **ScriptsModule**: Script CRUD & blockchain triggering
- **EventsModule**: Blockchain event logging
- **BlockchainModule**: Smart contract interaction
- **CommandsModule**: Legacy command endpoints

### apps/smart-contract
Foundry project với:
- **DeviceControl.sol**: Main smart contract
- **Deploy.s.sol**: Deployment script
- **Tests**: Comprehensive test suite

### libs/shared
Shared library chứa:
- **Types**: Command, Blockchain, API types
- **Constants**: App configuration và constants

## License

MIT
