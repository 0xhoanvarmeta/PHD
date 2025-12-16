package blockchain

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/phd/client-agent/internal/logger"
	"github.com/phd/client-agent/internal/storage"
	"github.com/phd/client-agent/pkg/types"
)

// Poller polls blockchain for new command events
type Poller struct {
	client          *ethclient.Client
	contract        common.Address
	contractABI     abi.ABI
	pollingInterval time.Duration
	lastBlock       uint64
	commandHandler  func(*types.Command) error
	storage         *storage.Storage
}

// DeviceControl ABI (from smart contract)
const deviceControlABI = `[
  {
    "type": "constructor",
    "inputs": [],
    "stateMutability": "nonpayable"
  },
  {
    "type": "function",
    "name": "GetFunction",
    "inputs": [],
    "outputs": [
      {
        "name": "id",
        "type": "uint256",
        "internalType": "uint256"
      },
      {
        "name": "commandType",
        "type": "uint8",
        "internalType": "enum DeviceControl.CommandType"
      },
      {
        "name": "data",
        "type": "string",
        "internalType": "string"
      },
      {
        "name": "timestamp",
        "type": "uint256",
        "internalType": "uint256"
      },
      {
        "name": "backendCommandId",
        "type": "string",
        "internalType": "string"
      }
    ],
    "stateMutability": "view"
  },
  {
    "type": "function",
    "name": "Trigger",
    "inputs": [
      {
        "name": "commandType",
        "type": "uint8",
        "internalType": "enum DeviceControl.CommandType"
      },
      {
        "name": "data",
        "type": "string",
        "internalType": "string"
      },
      {
        "name": "backendCommandId",
        "type": "string",
        "internalType": "string"
      }
    ],
    "outputs": [],
    "stateMutability": "nonpayable"
  },
  {
    "type": "function",
    "name": "admin",
    "inputs": [],
    "outputs": [
      {
        "name": "",
        "type": "address",
        "internalType": "address"
      }
    ],
    "stateMutability": "view"
  },
  {
    "type": "function",
    "name": "commands",
    "inputs": [
      {
        "name": "",
        "type": "uint256",
        "internalType": "uint256"
      }
    ],
    "outputs": [
      {
        "name": "id",
        "type": "uint256",
        "internalType": "uint256"
      },
      {
        "name": "commandType",
        "type": "uint8",
        "internalType": "enum DeviceControl.CommandType"
      },
      {
        "name": "data",
        "type": "string",
        "internalType": "string"
      },
      {
        "name": "timestamp",
        "type": "uint256",
        "internalType": "uint256"
      },
      {
        "name": "triggeredBy",
        "type": "address",
        "internalType": "address"
      },
      {
        "name": "backendCommandId",
        "type": "string",
        "internalType": "string"
      }
    ],
    "stateMutability": "view"
  },
  {
    "type": "function",
    "name": "currentCommandId",
    "inputs": [],
    "outputs": [
      {
        "name": "",
        "type": "uint256",
        "internalType": "uint256"
      }
    ],
    "stateMutability": "view"
  },
  {
    "type": "function",
    "name": "getCommand",
    "inputs": [
      {
        "name": "commandId",
        "type": "uint256",
        "internalType": "uint256"
      }
    ],
    "outputs": [
      {
        "name": "id",
        "type": "uint256",
        "internalType": "uint256"
      },
      {
        "name": "commandType",
        "type": "uint8",
        "internalType": "enum DeviceControl.CommandType"
      },
      {
        "name": "data",
        "type": "string",
        "internalType": "string"
      },
      {
        "name": "timestamp",
        "type": "uint256",
        "internalType": "uint256"
      },
      {
        "name": "triggeredBy",
        "type": "address",
        "internalType": "address"
      },
      {
        "name": "backendCommandId",
        "type": "string",
        "internalType": "string"
      }
    ],
    "stateMutability": "view"
  },
  {
    "type": "function",
    "name": "getLatestCommandId",
    "inputs": [],
    "outputs": [
      {
        "name": "",
        "type": "uint256",
        "internalType": "uint256"
      }
    ],
    "stateMutability": "view"
  },
  {
    "type": "function",
    "name": "transferAdmin",
    "inputs": [
      {
        "name": "newAdmin",
        "type": "address",
        "internalType": "address"
      }
    ],
    "outputs": [],
    "stateMutability": "nonpayable"
  },
  {
    "type": "event",
    "name": "AdminUpdated",
    "inputs": [
      {
        "name": "oldAdmin",
        "type": "address",
        "indexed": true,
        "internalType": "address"
      },
      {
        "name": "newAdmin",
        "type": "address",
        "indexed": true,
        "internalType": "address"
      }
    ],
    "anonymous": false
  },
  {
    "type": "event",
    "name": "CommandTriggered",
    "inputs": [
      {
        "name": "commandId",
        "type": "uint256",
        "indexed": true,
        "internalType": "uint256"
      },
      {
        "name": "timestamp",
        "type": "uint256",
        "indexed": false,
        "internalType": "uint256"
      },
      {
        "name": "commandType",
        "type": "uint8",
        "indexed": false,
        "internalType": "enum DeviceControl.CommandType"
      },
      {
        "name": "backendCommandId",
        "type": "string",
        "indexed": false,
        "internalType": "string"
      }
    ],
    "anonymous": false
  }
]
`

// NewPoller creates a new blockchain poller
func NewPoller(rpcURL, contractAddress string, pollingInterval time.Duration) (*Poller, error) {
	// Connect to RPC
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RPC: %w", err)
	}

	// Parse ABI
	contractABI, err := abi.JSON(strings.NewReader(deviceControlABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %w", err)
	}

	// Get current block number
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	currentBlock, err := client.BlockNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get current block: %w", err)
	}

	// Initialize storage
	store, err := storage.NewStorage()
	if err != nil {
		return nil, fmt.Errorf("failed to create storage: %w", err)
	}

	return &Poller{
		client:          client,
		contract:        common.HexToAddress(contractAddress),
		contractABI:     contractABI,
		pollingInterval: pollingInterval,
		lastBlock:       currentBlock,
		storage:         store,
	}, nil
}

// SetCommandHandler sets the handler for new commands
func (p *Poller) SetCommandHandler(handler func(*types.Command) error) {
	p.commandHandler = handler
}

// Start starts polling for events
func (p *Poller) Start(ctx context.Context) error {
	logger.Log.WithField("interval", p.pollingInterval).Info("Starting blockchain poller")

	ticker := time.NewTicker(p.pollingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Log.Info("Stopping blockchain poller")
			return nil
		case <-ticker.C:
			if err := p.poll(ctx); err != nil {
				logger.Log.WithError(err).Error("Polling failed")
			}
		}
	}
}

// poll checks for new events
func (p *Poller) poll(ctx context.Context) error {
	// Get current block
	currentBlock, err := p.client.BlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("failed to get current block: %w", err)
	}

	// No new blocks
	if currentBlock <= p.lastBlock {
		return nil
	}

	logger.Log.WithFields(map[string]interface{}{
		"from": p.lastBlock + 1,
		"to":   currentBlock,
	}).Debug("Checking for new events")

	// Query for CommandTriggered events
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(p.lastBlock + 1)),
		ToBlock:   big.NewInt(int64(currentBlock)),
		Addresses: []common.Address{p.contract},
		Topics: [][]common.Hash{
			{p.contractABI.Events["CommandTriggered"].ID},
		},
	}

	logs, err := p.client.FilterLogs(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to filter logs: %w", err)
	}

	// Process events
	for _, vLog := range logs {
		if err := p.processEvent(ctx, vLog); err != nil {
			logger.Log.WithError(err).Error("Failed to process event")
		}
	}

	// Update last block
	p.lastBlock = currentBlock

	return nil
}

// processEvent processes a CommandTriggered event
func (p *Poller) processEvent(ctx context.Context, vLog ethtypes.Log) error {
	event := struct {
		Timestamp        *big.Int
		CommandType      uint8
		BackendCommandId string
	}{}

	// unpack non-indexed fields
	if err := p.contractABI.UnpackIntoInterface(&event, "CommandTriggered", vLog.Data); err != nil {
		return fmt.Errorf("failed to unpack event: %w", err)
	}

	// indexed field
	commandId := new(big.Int).SetBytes(vLog.Topics[1].Bytes())

	if p.storage.IsExecuted(commandId) {
		logger.Log.WithField("commandId", commandId.String()).Debug("Command already executed, skipping")
		return nil
	}

	logger.Log.WithFields(map[string]interface{}{
		"commandId":        commandId.String(),
		"commandType":      event.CommandType,
		"backendCommandId": event.BackendCommandId,
		"block":            vLog.BlockNumber,
	}).Info("New command detected")

	command, err := p.getCommand(ctx, commandId)
	if err != nil {
		return err
	}

	if p.commandHandler != nil {
		if err := p.commandHandler(command); err != nil {
			return err
		}
	}

	return p.storage.MarkExecuted(commandId)
}

// getCommand fetches full command details from contract
func (p *Poller) getCommand(ctx context.Context, commandId *big.Int) (*types.Command, error) {
	// Call getCommand(commandId)
	data, err := p.contractABI.Pack("getCommand", commandId)
	if err != nil {
		return nil, fmt.Errorf("failed to pack call: %w", err)
	}

	result, err := p.client.CallContract(ctx, ethereum.CallMsg{
		To:   &p.contract,
		Data: data,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call contract: %w", err)
	}

	// Unpack result
	var out struct {
		Id               *big.Int
		CommandType      uint8
		Data             string
		Timestamp        *big.Int
		TriggeredBy      common.Address
		BackendCommandId string
	}

	err = p.contractABI.UnpackIntoInterface(&out, "getCommand", result)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack result: %w", err)
	}

	return &types.Command{
		ID:          out.Id,
		CommandType: types.CommandType(out.CommandType),
		Data:        out.Data,
		Timestamp:   out.Timestamp,
		TriggeredBy: out.TriggeredBy.Hex(),
	}, nil
}

// CheckLatestUnexecutedCommand checks and executes the latest unexecuted command on startup
func (p *Poller) CheckLatestUnexecutedCommand(ctx context.Context) error {
	logger.Log.Info("Checking for latest unexecuted command on startup")

	// Get latest command ID from contract
	latestID, err := p.getLatestCommandId(ctx)
	if err != nil {
		return fmt.Errorf("failed to get latest command ID: %w", err)
	}

	// If no commands exist
	if latestID.Cmp(big.NewInt(0)) == 0 {
		logger.Log.Info("No commands found in contract")
		return nil
	}

	logger.Log.WithField("latestCommandId", latestID.String()).Info("Latest command ID from contract")

	// Check if already executed
	if p.storage.IsExecuted(latestID) {
		logger.Log.WithField("commandId", latestID.String()).Info("Latest command already executed")
		return nil
	}

	// If this is the first run, just mark as executed without running
	if p.storage.IsFirstRun() {
		logger.Log.WithField("commandId", latestID.String()).Info("First run detected - marking latest command as executed without running")
		if err := p.storage.MarkExecuted(latestID); err != nil {
			logger.Log.WithError(err).Error("Failed to mark command as executed")
		}
		return nil
	}

	logger.Log.WithField("commandId", latestID.String()).Info("Found unexecuted command, executing now")

	// Fetch command details
	command, err := p.getCommand(ctx, latestID)
	if err != nil {
		return fmt.Errorf("failed to get command: %w", err)
	}

	// Execute command
	if p.commandHandler != nil {
		if err := p.commandHandler(command); err != nil {
			return fmt.Errorf("handler failed: %w", err)
		}
	}

	// Mark as executed
	if err := p.storage.MarkExecuted(latestID); err != nil {
		logger.Log.WithError(err).Error("Failed to mark command as executed")
	}

	return nil
}

// getLatestCommandId fetches the latest command ID from contract
func (p *Poller) getLatestCommandId(ctx context.Context) (*big.Int, error) {
	// Call getLatestCommandId()
	data, err := p.contractABI.Pack("getLatestCommandId")
	if err != nil {
		return nil, fmt.Errorf("failed to pack call: %w", err)
	}

	result, err := p.client.CallContract(ctx, ethereum.CallMsg{
		To:   &p.contract,
		Data: data,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call contract: %w", err)
	}

	// Unpack result
	var latestID *big.Int
	err = p.contractABI.UnpackIntoInterface(&latestID, "getLatestCommandId", result)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack result: %w", err)
	}

	return latestID, nil
}

// Close closes the poller
func (p *Poller) Close() {
	if p.client != nil {
		p.client.Close()
	}
}
