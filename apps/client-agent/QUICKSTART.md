# Quick Start Guide - PHD Client Agent

Get up and running in 5 minutes!

## Step 1: Download Binary

Download the appropriate binary for your platform:

- **macOS**: [phd-client-agent-darwin-amd64](link) or [phd-client-agent-darwin-arm64](link)
- **Linux**: [phd-client-agent-linux-amd64](link)
- **Windows**: [phd-client-agent-windows-amd64.exe](link)

## Step 2: Configure

Create `.env` file:

```bash
CONTRACT_ADDRESS=0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb
RPC_URL=https://testnet.hashio.io/api
CLIENT_ID=my-first-client
```

## Step 3: Run

### macOS/Linux

```bash
chmod +x phd-client-agent-*
./phd-client-agent-darwin-amd64
```

### Windows

```powershell
.\phd-client-agent-windows-amd64.exe
```

## Step 4: Test

### From Backend API

Trigger a test command from your backend:

```bash
curl -X POST http://localhost:3000/api/scripts/1/trigger \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### From Smart Contract

Call `Trigger` function on smart contract:

```javascript
// Using ethers.js
await deviceControl.Trigger(
  0, // CommandType.SCRIPT
  "echo 'Hello from blockchain!'"
);
```

## Expected Output

You should see in client agent logs:

```
2025-12-15 10:30:45 [INFO] PHD Client Agent starting
2025-12-15 10:30:45 [INFO] Starting blockchain poller
2025-12-15 10:30:50 [INFO] New command detected commandId=1
2025-12-15 10:30:50 [INFO] Processing new command
2025-12-15 10:30:51 [INFO] Command executed successfully
```

## Next Steps

- [Full Documentation](README.md)
- [Security Best Practices](README.md#security-considerations)
- [Running as Service](README.md#running-as-background-service)

## Troubleshooting

**Problem**: Can't connect to RPC

**Solution**: Check your `RPC_URL` and internet connection

```bash
curl https://testnet.hashio.io/api
```

**Problem**: Permission denied (macOS/Linux)

**Solution**: Make binary executable

```bash
chmod +x phd-client-agent-*
```

**Problem**: No events detected

**Solution**: Check contract address and make sure events are being emitted

```bash
# Check logs
tail -f client-agent.log
```

---

Need help? Check [README.md](README.md) or open an issue!
