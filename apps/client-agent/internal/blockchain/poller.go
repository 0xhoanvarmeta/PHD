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
}

// DeviceControl ABI (from smart contract)
const deviceControlABI = `[{"type":"event","name":"CommandTriggered","inputs":[{"name":"commandId","type":"uint256","indexed":true},{"name":"timestamp","type":"uint256","indexed":false},{"name":"commandType","type":"uint8","indexed":false}],"anonymous":false},{"type":"function","name":"GetFunction","inputs":[],"outputs":[{"name":"id","type":"uint256"},{"name":"commandType","type":"uint8"},{"name":"data","type":"string"},{"name":"timestamp","type":"uint256"}],"stateMutability":"view"},{"type":"function","name":"getCommand","inputs":[{"name":"commandId","type":"uint256"}],"outputs":[{"name":"id","type":"uint256"},{"name":"commandType","type":"uint8"},{"name":"data","type":"string"},{"name":"timestamp","type":"uint256"},{"name":"triggeredBy","type":"address"}],"stateMutability":"view"},{"type":"function","name":"getLatestCommandId","inputs":[],"outputs":[{"name":"","type":"uint256"}],"stateMutability":"view"}]`

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

	return &Poller{
		client:          client,
		contract:        common.HexToAddress(contractAddress),
		contractABI:     contractABI,
		pollingInterval: pollingInterval,
		lastBlock:       currentBlock,
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
	// Parse event
	event := struct {
		CommandId   *big.Int
		Timestamp   *big.Int
		CommandType uint8
	}{}

	err := p.contractABI.UnpackIntoInterface(&event, "CommandTriggered", vLog.Data)
	if err != nil {
		return fmt.Errorf("failed to unpack event: %w", err)
	}

	// Extract commandId from indexed topic
	event.CommandId = new(big.Int).SetBytes(vLog.Topics[1].Bytes())

	logger.Log.WithFields(map[string]interface{}{
		"commandId":   event.CommandId.String(),
		"commandType": event.CommandType,
		"block":       vLog.BlockNumber,
		"txHash":      vLog.TxHash.Hex(),
	}).Info("New command detected")

	// Fetch full command details
	command, err := p.getCommand(ctx, event.CommandId)
	if err != nil {
		return fmt.Errorf("failed to get command: %w", err)
	}

	// Call handler
	if p.commandHandler != nil {
		if err := p.commandHandler(command); err != nil {
			return fmt.Errorf("handler failed: %w", err)
		}
	}

	return nil
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
		Id          *big.Int
		CommandType uint8
		Data        string
		Timestamp   *big.Int
		TriggeredBy common.Address
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

// Close closes the poller
func (p *Poller) Close() {
	if p.client != nil {
		p.client.Close()
	}
}
