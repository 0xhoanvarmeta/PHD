package config

import (
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/phd/client-agent/pkg/types"
	"github.com/spf13/viper"
)

// Load loads configuration from environment and config file
func Load() (*types.Config, error) {
	// Load .env file if exists
	_ = godotenv.Load()

	// Setup viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.phd-client-agent")
	viper.AutomaticEnv()

	// Set defaults
	setDefaults()

	// Read config file (optional)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
		// Config file not found; ignore
	}

	// Build config
	cfg := &types.Config{
		Network:          viper.GetString("BLOCKCHAIN_NETWORK"),
		ContractAddress:  viper.GetString("CONTRACT_ADDRESS"),
		RPCURL:           viper.GetString("RPC_URL"),
		ClientID:         viper.GetString("CLIENT_ID"),
		PollingInterval:  time.Duration(viper.GetInt("POLLING_INTERVAL")) * time.Millisecond,
		ExecutionTimeout: time.Duration(viper.GetInt("EXECUTION_TIMEOUT")) * time.Millisecond,
		MaxRetryAttempts: viper.GetInt("MAX_RETRY_ATTEMPTS"),
		LogLevel:         viper.GetString("LOG_LEVEL"),
		LogFile:          viper.GetString("LOG_FILE"),
	}

	// Validate
	if err := validate(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func setDefaults() {
	viper.SetDefault("BLOCKCHAIN_NETWORK", "testnet")
	viper.SetDefault("CONTRACT_ADDRESS", "0x1e8678A15DAf23C01d0A972D86F5D692469D392c")
	viper.SetDefault("RPC_URL", "https://testnet.hashio.io/api")
	viper.SetDefault("POLLING_INTERVAL", 5000)   // milliseconds
	viper.SetDefault("EXECUTION_TIMEOUT", 30000) // milliseconds
	viper.SetDefault("MAX_RETRY_ATTEMPTS", 3)
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("LOG_FILE", "client-agent.log")

	// Generate client ID if not set
	if os.Getenv("CLIENT_ID") == "" {
		viper.SetDefault("CLIENT_ID", uuid.New().String())
	}
}

func validate(cfg *types.Config) error {
	if cfg.ContractAddress == "" {
		return fmt.Errorf("CONTRACT_ADDRESS is required")
	}
	if cfg.RPCURL == "" {
		return fmt.Errorf("RPC_URL is required")
	}
	if cfg.ClientID == "" {
		return fmt.Errorf("CLIENT_ID is required")
	}
	return nil
}
