# Getting Started

Hướng dẫn từng bước để chạy dự án PHD Device Control System.

## Prerequisites

Đảm bảo bạn đã cài đặt:
- Node.js >= 18
- pnpm >= 8
- Foundry (cho smart contract)

## Installation

```bash
# Clone repository (nếu cần)
cd /path/to/PHD

# Install dependencies
pnpm install

# Install Foundry (nếu chưa có)
curl -L https://foundry.paradigm.xyz | bash
foundryup
```

## Complete Workflow

### Step 1: Build Smart Contract & Export ABI

```bash
# Build contract
pnpm contract:build

# Export ABI to shared library
pnpm contract:export-abi
```

Sau bước này, ABI sẽ được export vào `libs/shared/src/lib/abis/` và backend có thể sử dụng.

### Step 2: Deploy Smart Contract (Optional - For Testing)

#### Deploy to Local Network

```bash
# Terminal 1: Start local anvil node
anvil

# Terminal 2: Deploy contract
pnpm contract:deploy:local
```

#### Deploy to Hedera Testnet

```bash
# Configure .env in apps/smart-contract
cd apps/smart-contract
cp .env.example .env

# Edit .env với private key của bạn
# PRIVATE_KEY=your_private_key_here

# Deploy
pnpm contract:deploy:testnet

# Lưu lại contract address từ output
# DeviceControl deployed at: 0xABC123...
```

### Step 3: Verify Smart Contract (After Deployment)

```bash
# Thêm CONTRACT_ADDRESS vào apps/smart-contract/.env
echo "CONTRACT_ADDRESS=0xABC123..." >> apps/smart-contract/.env

# Verify trên testnet
pnpm contract:verify
```

### Step 4: Configure Backend

```bash
# Copy .env.example to .env
cp .env.example .env

# Edit .env và điền thông tin:
```

`.env`:
```bash
# Blockchain Configuration
BLOCKCHAIN_NETWORK=testnet          # hoặc mainnet, local
CONTRACT_ADDRESS=0xABC123...        # Địa chỉ contract vừa deploy
PRIVATE_KEY=your_private_key        # Private key của admin (optional)

# Backend Configuration
PORT=3000
NODE_ENV=development

# API Configuration
API_PREFIX=api
CORS_ORIGIN=*

# Logging
LOG_LEVEL=debug
```

**Lưu ý:**
- `PRIVATE_KEY` chỉ cần nếu bạn muốn backend có thể trigger commands
- Nếu không có `PRIVATE_KEY`, backend chỉ có thể đọc data từ contract

### Step 5: Run Backend

```bash
# Start backend in development mode
pnpm dev:backend
```

Backend sẽ chạy tại `http://localhost:3000`

### Step 6: Test API

#### Check blockchain info
```bash
curl http://localhost:3000/api/commands/info
```

Response:
```json
{
  "success": true,
  "data": {
    "network": "connected",
    "contractAddress": "0xABC123...",
    "walletAddress": "0xDEF456...",
    "adminAddress": "0xDEF456...",
    "historyCount": 0
  }
}
```

#### Get current command
```bash
curl http://localhost:3000/api/commands/current
```

#### Trigger new command (requires PRIVATE_KEY)
```bash
curl -X POST http://localhost:3000/api/commands/trigger \
  -H "Content-Type: application/json" \
  -d '{
    "commandType": 0,
    "data": "set-wallpaper https://example.com/image.jpg"
  }'
```

## Development Workflow

### When you modify smart contract:

```bash
# 1. Make changes to contract
# 2. Test changes
pnpm contract:test

# 3. Build and export ABI
pnpm contract:build
pnpm contract:export-abi

# 4. Rebuild backend (if needed)
pnpm build:backend
```

### When you modify backend:

```bash
# Backend auto-reloads in dev mode
pnpm dev:backend
```

## Project Structure

```
PHD/
├── apps/
│   ├── backend/                 # NestJS API
│   │   └── src/
│   │       ├── app/
│   │       │   ├── blockchain/  # Blockchain service (uses ABI)
│   │       │   └── commands/    # Commands API
│   │       └── main.ts
│   └── smart-contract/          # Foundry contracts
│       ├── src/DeviceControl.sol
│       ├── export-abi.sh        # Export ABI script
│       └── verify.sh            # Verify script
├── libs/
│   └── shared/                  # Shared library
│       └── src/
│           ├── lib/
│           │   ├── types/
│           │   ├── constants/
│           │   └── abis/        # Contract ABIs (auto-generated)
│           └── index.ts
├── .env                         # Backend config
└── package.json
```

## Common Issues

### Backend cannot connect to contract
- Kiểm tra `CONTRACT_ADDRESS` trong `.env`
- Kiểm tra `BLOCKCHAIN_NETWORK` đúng với network đã deploy
- Kiểm tra RPC URL trong `libs/shared/src/lib/types/blockchain.types.ts`

### Cannot trigger commands
- Cần có `PRIVATE_KEY` trong `.env`
- Wallet phải là admin của contract
- Đảm bảo có đủ gas

### ABI not found
- Chạy `pnpm contract:export-abi` sau khi build contract
- Kiểm tra `libs/shared/src/lib/abis/DeviceControl.json` tồn tại

## Next Steps

1. Develop client application để listen events
2. Implement authentication cho API
3. Add database để lưu command history
4. Deploy lên production

## Useful Commands

| Command | Description |
|---------|-------------|
| `pnpm dev:backend` | Start backend dev server |
| `pnpm contract:build` | Build smart contract |
| `pnpm contract:export-abi` | Export ABI to shared lib |
| `pnpm contract:test` | Test smart contract |
| `pnpm contract:deploy:testnet` | Deploy to Hedera testnet |
| `pnpm contract:verify` | Verify contract |
| `pnpm build` | Build all projects |
| `pnpm test` | Run all tests |

Happy coding! 🚀
