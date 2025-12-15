package types

import (
	"math/big"
	"time"
)

// CommandType represents the type of command to execute
type CommandType uint8

const (
	// CommandTypeScript executes script directly
	CommandTypeScript CommandType = 0
	// CommandTypeURL fetches script from URL and executes
	CommandTypeURL CommandType = 1
)

// Command represents a blockchain command
type Command struct {
	ID          *big.Int
	CommandType CommandType
	Data        string
	Timestamp   *big.Int
	TriggeredBy string
}

// ExecutionResult represents the result of a command execution
type ExecutionResult struct {
	CommandID   *big.Int
	Success     bool
	Output      string
	Error       string
	ExecutedAt  time.Time
	Duration    time.Duration
}

// Config represents application configuration
type Config struct {
	// Blockchain
	Network         string
	ContractAddress string
	RPCURL          string

	// Client
	ClientID        string
	PollingInterval time.Duration
	ExecutionTimeout time.Duration
	MaxRetryAttempts int

	// Logging
	LogLevel string
	LogFile  string
}

// ClientInfo represents client system information
type ClientInfo struct {
	ClientID  string
	OS        string
	OSVersion string
	Arch      string
	Hostname  string
}
