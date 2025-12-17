package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/phd/client-agent/internal/blockchain"
	"github.com/phd/client-agent/internal/config"
	"github.com/phd/client-agent/internal/executor"
	"github.com/phd/client-agent/internal/logger"
	"github.com/phd/client-agent/pkg/types"
)

const (
	version = "1.0.0"
)

func main() {
	// Check and request root privileges if needed
	if err := ensureRootPrivileges(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to obtain root privileges: %v\n", err)
		os.Exit(1)
	}

	// Print banner
	printBanner()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	if err := logger.Init(cfg.LogLevel, cfg.LogFile); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	logger.Log.WithFields(map[string]interface{}{
		"version":  version,
		"clientId": cfg.ClientID,
		"os":       runtime.GOOS,
		"arch":     runtime.GOARCH,
	}).Info("PHD Client Agent starting")

	// Create executor
	exec, err := executor.NewExecutor(cfg.ExecutionTimeout, cfg.MaxRetryAttempts)
	if err != nil {
		logger.Log.WithError(err).Fatal("Failed to create executor")
	}
	defer exec.Cleanup()

	// Create blockchain poller
	poller, err := blockchain.NewPoller(cfg.RPCURL, cfg.ContractAddress, cfg.PollingInterval)
	if err != nil {
		logger.Log.WithError(err).Fatal("Failed to create blockchain poller")
	}
	defer poller.Close()

	// Set command handler
	poller.SetCommandHandler(func(cmd *types.Command) error {
		logger.Log.WithFields(map[string]interface{}{
			"commandId":   cmd.ID.String(),
			"commandType": cmd.CommandType,
			"dataLength":  len(cmd.Data),
		}).Info("Processing new command")

		// Execute command
		result := exec.Execute(cmd)

		// Log result
		if result.Success {
			logger.Log.WithFields(map[string]interface{}{
				"commandId": result.CommandID.String(),
				"duration":  result.Duration,
				"output":    truncate(result.Output, 200),
			}).Info("Command executed successfully")
		} else {
			logger.Log.WithFields(map[string]interface{}{
				"commandId": result.CommandID.String(),
				"duration":  result.Duration,
				"error":     result.Error,
			}).Error("Command execution failed")
		}

		return nil
	})

	// Setup context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Check for latest unexecuted command on startup
	logger.Log.Info("Checking for pending commands from previous session...")
	if err := poller.CheckLatestUnexecutedCommand(ctx); err != nil {
		logger.Log.WithError(err).Warn("Failed to check for latest unexecuted command")
	}

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start poller in goroutine
	errChan := make(chan error, 1)
	go func() {
		if err := poller.Start(ctx); err != nil {
			errChan <- err
		}
	}()

	logger.Log.Info("Client agent running. Press Ctrl+C to stop.")

	// Wait for shutdown signal or error
	select {
	case <-sigChan:
		logger.Log.Info("Shutdown signal received")
	case err := <-errChan:
		logger.Log.WithError(err).Error("Poller error")
	}

	// Graceful shutdown
	logger.Log.Info("Shutting down...")
	cancel()
	logger.Log.Info("Shutdown complete")
}

func printBanner() {
	banner := `
╔═══════════════════════════════════════════════╗
║     PHD Client Agent - Blockchain Executor    ║
║                 Version %s                  ║
╚═══════════════════════════════════════════════╝
`
	fmt.Printf(banner, version)
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// ensureRootPrivileges checks if running as root and re-executes with sudo if needed
func ensureRootPrivileges() error {
	// Only for Unix-like systems (Linux, macOS)
	if runtime.GOOS == "windows" {
		// Windows doesn't use sudo, skip privilege check
		return nil
	}

	// Check if already running as root (euid == 0)
	if os.Geteuid() == 0 {
		fmt.Println("✓ Running with root privileges")
		return nil
	}

	fmt.Println("⚠ Root privileges required. Requesting sudo access...")
	fmt.Println("Please enter your password to continue.")

	// Get current executable path
	executable, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Re-execute with sudo
	cmd := exec.Command("sudo", append([]string{executable}, os.Args[1:]...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// Run the command
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to execute with sudo: %w", err)
	}

	// Exit this process as the sudo version is now running
	os.Exit(0)
	return nil
}
