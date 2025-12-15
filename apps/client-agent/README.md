# PHD Client Agent - Blockchain Command Executor

Cross-platform client agent that listens to blockchain events and executes commands on local machines.

## Features

✅ **Cross-Platform** - Runs on macOS, Ubuntu/Linux, and Windows
✅ **Blockchain Integration** - Direct polling of Hedera smart contract
✅ **Event-Driven** - Listens for `CommandTriggered` events
✅ **Dual Execution Modes**:
  - `SCRIPT`: Execute script content directly
  - `URL`: Fetch script from URL and execute
✅ **Standalone Binary** - No runtime dependencies (Node.js, Python, etc.)
✅ **Lightweight** - Binary size ~15-20MB
✅ **Auto Retry** - Configurable retry attempts with backoff
✅ **Comprehensive Logging** - File and console output

---

## Architecture

```
Smart Contract (Hedera)
    ↓ (emit CommandTriggered event)
Client Agent (polling every 5s)
    ↓ (detect new event)
Fetch Command Details
    ↓
┌─────────────────────────────────┐
│  CommandType.SCRIPT → Execute   │
│  CommandType.URL → Fetch & Run  │
└─────────────────────────────────┘
    ↓
Execute on Local Machine
    ↓
Log Result
```

---

## Prerequisites

- **Go 1.21+** (for building from source)
- **Smart Contract** deployed on Hedera (testnet/mainnet)
- **Network Access** to Hedera RPC endpoint

---

## Installation

### Option 1: Download Pre-built Binary

Download the binary for your platform from releases:

- **macOS (Intel)**: `phd-client-agent-darwin-amd64`
- **macOS (Apple Silicon)**: `phd-client-agent-darwin-arm64`
- **Linux (Ubuntu/Debian)**: `phd-client-agent-linux-amd64`
- **Windows**: `phd-client-agent-windows-amd64.exe`

### Option 2: Build from Source

```bash
# Clone repository
cd apps/client-agent

# Install dependencies
make deps

# Build for current platform
make build

# Or build for all platforms
make build-all

# Or use build script
./build.sh
```

---

## Configuration

Create a `.env` file or set environment variables:

```bash
# Copy example config
cp .env.example .env

# Edit configuration
nano .env
```

### Required Configuration

```env
# Blockchain Configuration
BLOCKCHAIN_NETWORK=testnet
CONTRACT_ADDRESS=0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb
RPC_URL=https://testnet.hashio.io/api

# Client Agent Configuration
CLIENT_ID=my-client-001
POLLING_INTERVAL=5000
EXECUTION_TIMEOUT=30000
MAX_RETRY_ATTEMPTS=3

# Logging
LOG_LEVEL=info
LOG_FILE=client-agent.log
```

### Configuration Options

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `BLOCKCHAIN_NETWORK` | Network name (testnet/mainnet/local) | testnet | No |
| `CONTRACT_ADDRESS` | Smart contract address | - | **Yes** |
| `RPC_URL` | Hedera RPC endpoint | https://testnet.hashio.io/api | **Yes** |
| `CLIENT_ID` | Unique client identifier | auto-generated UUID | No |
| `POLLING_INTERVAL` | Event polling interval (ms) | 5000 | No |
| `EXECUTION_TIMEOUT` | Script execution timeout (ms) | 30000 | No |
| `MAX_RETRY_ATTEMPTS` | Max retry on failure | 3 | No |
| `LOG_LEVEL` | Log level (debug/info/warn/error) | info | No |
| `LOG_FILE` | Log file path | client-agent.log | No |

---

## Usage

### Running the Agent

#### macOS/Linux

```bash
# Make executable
chmod +x phd-client-agent-darwin-amd64

# Run
./phd-client-agent-darwin-amd64
```

#### Windows

```powershell
# Run in PowerShell or CMD
.\phd-client-agent-windows-amd64.exe
```

### Running as Background Service

#### Linux (systemd)

Create service file `/etc/systemd/system/phd-client-agent.service`:

```ini
[Unit]
Description=PHD Client Agent
After=network.target

[Service]
Type=simple
User=your-user
WorkingDirectory=/opt/phd-client-agent
EnvironmentFile=/opt/phd-client-agent/.env
ExecStart=/opt/phd-client-agent/phd-client-agent-linux-amd64
Restart=on-failure
RestartSec=10

[Install]
WantedBy=multi-user.target
```

Enable and start:

```bash
sudo systemctl daemon-reload
sudo systemctl enable phd-client-agent
sudo systemctl start phd-client-agent
sudo systemctl status phd-client-agent
```

#### macOS (launchd)

Create plist file `~/Library/LaunchAgents/com.phd.client-agent.plist`:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.phd.client-agent</string>
    <key>ProgramArguments</key>
    <array>
        <string>/usr/local/bin/phd-client-agent</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardErrorPath</key>
    <string>/tmp/phd-client-agent.err</string>
    <key>StandardOutPath</key>
    <string>/tmp/phd-client-agent.out</string>
</dict>
</plist>
```

Load service:

```bash
launchctl load ~/Library/LaunchAgents/com.phd.client-agent.plist
launchctl start com.phd.client-agent
```

#### Windows (NSSM)

Using [NSSM (Non-Sucking Service Manager)](https://nssm.cc/):

```powershell
# Install NSSM
choco install nssm

# Install service
nssm install PHDClientAgent "C:\path\to\phd-client-agent-windows-amd64.exe"
nssm set PHDClientAgent AppDirectory "C:\path\to"
nssm set PHDClientAgent Start SERVICE_AUTO_START

# Start service
nssm start PHDClientAgent
```

---

## How It Works

### 1. Event Polling

The agent polls the smart contract every 5 seconds (configurable) for new `CommandTriggered` events:

```go
// Event signature
event CommandTriggered(
    uint256 indexed commandId,
    uint256 timestamp,
    uint8 commandType
)
```

### 2. Fetch Command Details

When a new event is detected, the agent calls `getCommand(commandId)` to fetch full command details:

```solidity
function getCommand(uint256 commandId) external view returns (
    uint256 id,
    CommandType commandType,
    string data,
    uint256 timestamp,
    address triggeredBy
)
```

### 3. Execute Command

Based on `commandType`:

#### CommandType.SCRIPT (0)

Execute script content directly:

```bash
# Example: Simple bash script
echo "Hello from blockchain!"
ls -la
```

#### CommandType.URL (1)

Fetch script from URL and execute:

```bash
# Example URL: https://example.com/scripts/backup.sh
curl -s https://example.com/scripts/backup.sh | bash
```

### 4. Cross-Platform Execution

| Platform | Shell | Script Extension |
|----------|-------|------------------|
| **macOS** | `/bin/bash` | `.sh` |
| **Linux** | `/bin/bash` | `.sh` |
| **Windows** | `powershell` | `.ps1` |

Scripts are saved to temp directory, executed, then cleaned up.

---

## Security Considerations

⚠️ **IMPORTANT: This agent executes arbitrary code from blockchain!**

### Recommended Security Measures

1. **Whitelist Contract Addresses**
   - Only connect to trusted smart contracts
   - Verify contract address in config

2. **Run with Limited Permissions**
   - Create dedicated user with minimal permissions
   - Use sandboxing (Docker, VMs)

3. **Network Isolation**
   - Run in isolated network segment
   - Firewall outbound connections

4. **Code Review**
   - Review all scripts before triggering via blockchain
   - Implement approval workflow for admin

5. **Monitoring & Alerts**
   - Monitor logs for suspicious activity
   - Set up alerts for failed executions

6. **URL Validation** (for CommandType.URL)
   - Whitelist allowed domains
   - Use HTTPS only
   - Validate SSL certificates

### Example: Running with Docker

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o phd-client-agent ./cmd/agent

FROM alpine:latest
RUN apk add --no-cache bash curl
COPY --from=builder /app/phd-client-agent /usr/local/bin/
USER nobody
ENTRYPOINT ["phd-client-agent"]
```

---

## Troubleshooting

### Common Issues

#### 1. "Failed to connect to RPC"

```bash
# Check RPC URL is accessible
curl https://testnet.hashio.io/api
```

#### 2. "CONTRACT_ADDRESS is required"

Make sure `.env` file exists and has `CONTRACT_ADDRESS` set.

#### 3. Permission Denied (macOS/Linux)

```bash
chmod +x phd-client-agent-*
```

#### 4. PowerShell Execution Policy (Windows)

```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

#### 5. Script Execution Timeout

Increase `EXECUTION_TIMEOUT` in config:

```env
EXECUTION_TIMEOUT=60000  # 60 seconds
```

---

## Logs

Logs are written to both console and file (if `LOG_FILE` is configured).

### Log Levels

- `debug`: Detailed debugging info
- `info`: General information (default)
- `warn`: Warning messages
- `error`: Error messages

### Example Log Output

```
2025-12-15 10:30:45 [INFO] PHD Client Agent starting version=1.0.0 clientId=abc-123 os=darwin
2025-12-15 10:30:45 [INFO] Starting blockchain poller interval=5s
2025-12-15 10:30:50 [INFO] New command detected commandId=42 commandType=0 block=12345
2025-12-15 10:30:50 [INFO] Processing new command commandId=42 dataLength=128
2025-12-15 10:30:51 [INFO] Command executed successfully commandId=42 duration=1.2s
```

---

## Build Commands

### Using Makefile

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Build specific platform
make build-linux
make build-macos
make build-windows

# Clean build artifacts
make clean

# Run tests
make test

# Show help
make help
```

### Using Build Script

```bash
# Build for all platforms
./build.sh
```

---

## Development

### Project Structure

```
client-agent/
├── cmd/
│   └── agent/
│       └── main.go              # Entry point
├── internal/
│   ├── blockchain/
│   │   └── poller.go            # Blockchain event poller
│   ├── executor/
│   │   └── executor.go          # Script executor
│   ├── config/
│   │   └── config.go            # Configuration management
│   └── logger/
│       └── logger.go            # Logger setup
├── pkg/
│   └── types/
│       └── types.go             # Shared types
├── build/                       # Build artifacts
├── go.mod                       # Go module file
├── Makefile                     # Build automation
├── build.sh                     # Build script
├── .env.example                 # Example config
└── README.md                    # This file
```

### Adding New Features

1. Create feature branch
2. Implement changes
3. Add tests
4. Update documentation
5. Submit PR

---

## License

MIT License

---

## Support

For issues and questions:
- GitHub Issues: [phd/client-agent/issues](https://github.com/phd/client-agent/issues)
- Documentation: [docs](./docs)

---

## Changelog

### v1.0.0 (2025-12-15)

- Initial release
- Cross-platform support (macOS, Linux, Windows)
- Blockchain event polling
- SCRIPT and URL command types
- Auto-retry mechanism
- Comprehensive logging
