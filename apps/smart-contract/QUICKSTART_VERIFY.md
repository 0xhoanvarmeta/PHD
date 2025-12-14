# Quick Start: Contract Verification

Hướng dẫn nhanh để verify smart contract sau khi deploy.

## Step 1: Deploy Contract

```bash
# Deploy to Hedera testnet
pnpm contract:deploy:testnet
```

Sau khi deploy xong, bạn sẽ thấy output như:
```
DeviceControl deployed at: 0x1234567890abcdef...
Admin address: 0xabcdef1234567890...
```

## Step 2: Configure Environment

Copy contract address và thêm vào `.env`:

```bash
cd apps/smart-contract
cp .env.example .env
```

Edit `.env`:
```bash
CONTRACT_ADDRESS=0x1234567890abcdef...
PRIVATE_KEY=your_private_key
ETHERSCAN_API_KEY=your_api_key  # Optional
```

## Step 3: Verify Contract

### Option A: Using pnpm (Recommended)

```bash
# From root directory
pnpm contract:verify
```

### Option B: Using script directly

```bash
cd apps/smart-contract
./verify.sh 0x1234567890abcdef... testnet
```

### Option C: Using Nx

```bash
nx verify smart-contract
```

## Step 4: Check Verification

Sau khi verify thành công, mở browser:

**Hedera Testnet:**
```
https://hashscan.io/testnet/contract/0x1234567890abcdef...
```

Bạn sẽ thấy:
- ✓ Tab "Contract" với source code
- ✓ Tab "Read Contract" để đọc state
- ✓ Tab "Write Contract" để gọi functions
- ✓ Verified checkmark badge

## Troubleshooting

### Error: CONTRACT_ADDRESS not set
```bash
# Set in .env file
echo "CONTRACT_ADDRESS=0x1234..." >> .env
```

### Error: Verification timeout
```bash
# Try manual verification with --watch flag
forge verify-contract \
  0x1234... \
  src/DeviceControl.sol:DeviceControl \
  --chain-id 296 \
  --verifier-url https://server-verify.hashscan.io \
  --watch
```

### Success message
```
✓ Contract verified successfully!
View on explorer: https://hashscan.io/testnet/contract/0x1234...
```

## Full Workflow Example

```bash
# 1. Deploy
pnpm contract:deploy:testnet

# Output:
# DeviceControl deployed at: 0xABC123...
# Admin address: 0xDEF456...

# 2. Add to .env
echo "CONTRACT_ADDRESS=0xABC123..." >> apps/smart-contract/.env

# 3. Verify
pnpm contract:verify

# Output:
# === Contract Verification ===
# Contract Address: 0xABC123...
# Network: testnet
# Chain ID: 296
# Starting verification...
# ✓ Contract verified successfully!

# 4. Check on explorer
open https://hashscan.io/testnet/contract/0xABC123...
```

## What's Next?

Sau khi verify contract:
1. Update `CONTRACT_ADDRESS` trong root `.env` cho backend
2. Start backend: `pnpm dev:backend`
3. Test API endpoints
4. Deploy client applications

Happy coding! 🚀
