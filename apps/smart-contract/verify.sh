#!/bin/bash

# Script to verify DeviceControl contract on Hedera block explorer
# Usage: ./verify.sh <contract_address> [network]
# Example: ./verify.sh 0x1234... testnet

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if contract address is provided
if [ -z "$1" ]; then
    echo -e "${RED}Error: Contract address is required${NC}"
    echo "Usage: ./verify.sh <contract_address> [network]"
    echo "Example: ./verify.sh 0x1234567890abcdef... testnet"
    exit 1
fi

CONTRACT_ADDRESS=$1
NETWORK=${2:-testnet}

# Load environment variables
if [ -f .env ]; then
    source .env
fi

# Network configurations
case $NETWORK in
    testnet)
        CHAIN_ID=296
        VERIFIER_URL="https://server-verify.hashscan.io"
        EXPLORER_URL="https://hashscan.io/testnet"
        echo -e "${GREEN}Using Hedera Testnet${NC}"
        ;;
    mainnet)
        CHAIN_ID=295
        VERIFIER_URL="https://server-verify.hashscan.io"
        EXPLORER_URL="https://hashscan.io/mainnet"
        echo -e "${GREEN}Using Hedera Mainnet${NC}"
        ;;
    local)
        echo -e "${YELLOW}Local network does not support verification${NC}"
        exit 0
        ;;
    *)
        echo -e "${RED}Error: Unknown network '$NETWORK'${NC}"
        echo "Supported networks: testnet, mainnet, local"
        exit 1
        ;;
esac

# Check if ETHERSCAN_API_KEY is set
if [ -z "$ETHERSCAN_API_KEY" ]; then
    echo -e "${YELLOW}Warning: ETHERSCAN_API_KEY not set in .env file${NC}"
    echo "Verification may not work without API key"
fi

echo ""
echo -e "${GREEN}=== Contract Verification ===${NC}"
echo "Contract Address: $CONTRACT_ADDRESS"
echo "Network: $NETWORK"
echo "Chain ID: $CHAIN_ID"
echo "Verifier URL: $VERIFIER_URL"
echo ""

# Run verification
echo -e "${GREEN}Starting verification...${NC}"
forge verify-contract \
    $CONTRACT_ADDRESS \
    src/DeviceControl.sol:DeviceControl \
    --chain-id $CHAIN_ID \
    --verifier-url $VERIFIER_URL \
    ${ETHERSCAN_API_KEY:+--etherscan-api-key $ETHERSCAN_API_KEY} \
    --watch

if [ $? -eq 0 ]; then
    echo ""
    echo -e "${GREEN}✓ Contract verified successfully!${NC}"
    echo -e "${GREEN}View on explorer: $EXPLORER_URL/contract/$CONTRACT_ADDRESS${NC}"
else
    echo ""
    echo -e "${RED}✗ Verification failed${NC}"
    exit 1
fi
