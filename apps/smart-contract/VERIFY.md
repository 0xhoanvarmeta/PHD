# Smart Contract Verification Guide

Hướng dẫn verify smart contract DeviceControl trên Hedera block explorer.

## Prerequisites

1. Contract đã được deploy lên network (testnet hoặc mainnet)
2. Có địa chỉ contract address
3. Đã cấu hình `ETHERSCAN_API_KEY` trong file `.env` (optional nhưng recommended)

## Configuration

Thêm contract address vào file `.env`:

```bash
CONTRACT_ADDRESS=0x1234567890abcdef...
ETHERSCAN_API_KEY=your_api_key_here
```

## Verification Methods

### Method 1: Using pnpm scripts (Recommended)

#### Verify on Hedera Testnet
```bash
# Set CONTRACT_ADDRESS in .env first
pnpm contract:verify
```

#### Verify on Hedera Mainnet
```bash
pnpm contract:verify:mainnet
```

#### View verification info
```bash
pnpm contract:verify:info
```

### Method 2: Using bash script directly

```bash
cd apps/smart-contract
./verify.sh <contract_address> testnet
```

For mainnet:
```bash
./verify.sh <contract_address> mainnet
```

### Method 3: Using Nx directly

```bash
nx verify smart-contract
# or
nx verify:mainnet smart-contract
```

### Method 4: Manual verification with forge

#### Hedera Testnet
```bash
forge verify-contract \
  <contract_address> \
  src/DeviceControl.sol:DeviceControl \
  --chain-id 296 \
  --verifier-url https://server-verify.hashscan.io \
  --etherscan-api-key $ETHERSCAN_API_KEY \
  --watch
```

#### Hedera Mainnet
```bash
forge verify-contract \
  <contract_address> \
  src/DeviceControl.sol:DeviceControl \
  --chain-id 295 \
  --verifier-url https://server-verify.hashscan.io \
  --etherscan-api-key $ETHERSCAN_API_KEY \
  --watch
```

## Network Information

### Hedera Testnet
- Chain ID: `296`
- Verifier URL: `https://server-verify.hashscan.io`
- Explorer: `https://hashscan.io/testnet`
- RPC URL: `https://testnet.hashio.io/api`

### Hedera Mainnet
- Chain ID: `295`
- Verifier URL: `https://server-verify.hashscan.io`
- Explorer: `https://hashscan.io/mainnet`
- RPC URL: `https://mainnet.hashio.io/api`

## Troubleshooting

### Verification fails with "Contract source code already verified"
Contract đã được verify rồi. Bạn có thể xem trên explorer:
- Testnet: `https://hashscan.io/testnet/contract/<contract_address>`
- Mainnet: `https://hashscan.io/mainnet/contract/<contract_address>`

### Verification fails with "Invalid API Key"
1. Kiểm tra `ETHERSCAN_API_KEY` trong file `.env`
2. Đảm bảo API key hợp lệ cho Hedera network
3. Thử verify mà không cần API key (có thể sẽ chậm hơn)

### Verification timeout
Thêm flag `--watch` để theo dõi progress:
```bash
forge verify-contract ... --watch
```

### Contract bytecode doesn't match
1. Đảm bảo bạn đang verify đúng contract source code
2. Kiểm tra compiler version và optimization settings trong `foundry.toml`:
   - Solidity version: `0.8.28`
   - Optimizer: `true`
   - Optimizer runs: `200`

## Example Workflow

1. Deploy contract:
```bash
pnpm contract:deploy:testnet
```

2. Copy contract address từ output

3. Thêm vào `.env`:
```bash
CONTRACT_ADDRESS=0x1234...
```

4. Verify contract:
```bash
pnpm contract:verify
```

5. Kiểm tra trên explorer:
```
https://hashscan.io/testnet/contract/0x1234...
```

## Verification Status

Sau khi verify thành công, bạn sẽ thấy:
- ✓ Source code được hiển thị trên explorer
- ✓ ABI có thể download
- ✓ Read/Write functions có thể tương tác trực tiếp trên explorer
- ✓ Contract được đánh dấu "verified" với checkmark

## Additional Resources

- [Foundry Verification Docs](https://book.getfoundry.sh/reference/forge/forge-verify-contract)
- [Hedera Hashscan](https://hashscan.io)
- [Hedera Docs](https://docs.hedera.com)
