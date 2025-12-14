#!/bin/bash

# Script to export contract ABIs for backend usage
# This script builds the contracts and exports ABIs to shared library

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${GREEN}=== Exporting Contract ABIs ===${NC}"

# Build contracts first
echo "Building contracts..."
forge build

# Create ABI output directory
ABI_DIR="../../libs/shared/src/lib/abis"
mkdir -p $ABI_DIR

# Export DeviceControl ABI
CONTRACT_NAME="DeviceControl"
ABI_FILE="out/${CONTRACT_NAME}.sol/${CONTRACT_NAME}.json"

if [ ! -f "$ABI_FILE" ]; then
    echo -e "${YELLOW}Error: $ABI_FILE not found. Make sure contracts are built.${NC}"
    exit 1
fi

echo "Extracting ABI for $CONTRACT_NAME..."

# Extract ABI using jq (if available) or node
if command -v jq &> /dev/null; then
    cat $ABI_FILE | jq '.abi' > "$ABI_DIR/${CONTRACT_NAME}.json"
else
    # Use node to extract ABI
    node -e "
        const fs = require('fs');
        const contract = JSON.parse(fs.readFileSync('$ABI_FILE', 'utf8'));
        fs.writeFileSync('$ABI_DIR/${CONTRACT_NAME}.json', JSON.stringify(contract.abi, null, 2));
    "
fi

# Create TypeScript export file
cat > "$ABI_DIR/${CONTRACT_NAME}.ts" << 'TSEOF'
import abi from './DeviceControl.json';

export const DeviceControlABI = abi as const;

export default DeviceControlABI;
TSEOF

# Create index.ts to export all ABIs
cat > "$ABI_DIR/index.ts" << 'INDEXEOF'
export { DeviceControlABI } from './DeviceControl';
export { default as DeviceControlABIJson } from './DeviceControl.json';
INDEXEOF

echo -e "${GREEN}✓ ABI exported successfully to libs/shared/src/lib/abis/${NC}"
echo -e "${GREEN}✓ Generated TypeScript exports${NC}"
echo ""
echo "Files created:"
echo "  - libs/shared/src/lib/abis/DeviceControl.json"
echo "  - libs/shared/src/lib/abis/DeviceControl.ts"
echo "  - libs/shared/src/lib/abis/index.ts"
