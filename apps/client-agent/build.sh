#!/bin/bash
# Cross-platform build script for PHD Client Agent

set -e

VERSION="1.0.0"
BINARY_NAME="phd-client-agent"
BUILD_DIR="build"
MAIN_PATH="./cmd/agent"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}╔═══════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║     PHD Client Agent - Build Script           ║${NC}"
echo -e "${GREEN}║                 Version ${VERSION}                  ║${NC}"
echo -e "${GREEN}╚═══════════════════════════════════════════════╝${NC}"
echo ""

# Clean previous builds
echo -e "${YELLOW}Cleaning previous builds...${NC}"
rm -rf ${BUILD_DIR}
mkdir -p ${BUILD_DIR}

# Install dependencies
echo -e "${YELLOW}Installing dependencies...${NC}"
go mod download
go mod tidy

# Build function
build_binary() {
    local GOOS=$1
    local GOARCH=$2
    local OUTPUT=$3

    echo -e "${YELLOW}Building for ${GOOS}/${GOARCH}...${NC}"

    GOOS=${GOOS} GOARCH=${GOARCH} go build \
        -ldflags="-s -w -X main.version=${VERSION}" \
        -o ${OUTPUT} \
        ${MAIN_PATH}

    if [ $? -eq 0 ]; then
        SIZE=$(du -h ${OUTPUT} | cut -f1)
        echo -e "${GREEN}✓ Built ${OUTPUT} (${SIZE})${NC}"
    else
        echo -e "${RED}✗ Failed to build ${OUTPUT}${NC}"
        exit 1
    fi
}

# Build for all platforms
echo -e "${YELLOW}Building binaries...${NC}"
echo ""

# Linux
build_binary "linux" "amd64" "${BUILD_DIR}/${BINARY_NAME}-linux-amd64"
build_binary "linux" "arm64" "${BUILD_DIR}/${BINARY_NAME}-linux-arm64"

# macOS
build_binary "darwin" "amd64" "${BUILD_DIR}/${BINARY_NAME}-darwin-amd64"
build_binary "darwin" "arm64" "${BUILD_DIR}/${BINARY_NAME}-darwin-arm64"

# Windows
build_binary "windows" "amd64" "${BUILD_DIR}/${BINARY_NAME}-windows-amd64.exe"

echo ""
echo -e "${GREEN}╔═══════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║           Build Complete!                     ║${NC}"
echo -e "${GREEN}╚═══════════════════════════════════════════════╝${NC}"
echo ""
echo -e "${YELLOW}Binaries available in ${BUILD_DIR}/:${NC}"
ls -lh ${BUILD_DIR}/

echo ""
echo -e "${GREEN}Distribution guide:${NC}"
echo -e "  • Linux (Ubuntu/Debian):   ${BUILD_DIR}/${BINARY_NAME}-linux-amd64"
echo -e "  • macOS (Intel):           ${BUILD_DIR}/${BINARY_NAME}-darwin-amd64"
echo -e "  • macOS (Apple Silicon):   ${BUILD_DIR}/${BINARY_NAME}-darwin-arm64"
echo -e "  • Windows:                 ${BUILD_DIR}/${BINARY_NAME}-windows-amd64.exe"
