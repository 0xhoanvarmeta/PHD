package executor

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/phd/client-agent/internal/logger"
	"github.com/phd/client-agent/pkg/types"
)

// Executor executes commands
type Executor struct {
	timeout    time.Duration
	maxRetries int
	tempDir    string
}

// NewExecutor creates a new executor
func NewExecutor(timeout time.Duration, maxRetries int) (*Executor, error) {
	// Create temp directory for scripts
	tempDir := filepath.Join(os.TempDir(), "phd-client-agent")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}

	return &Executor{
		timeout:    timeout,
		maxRetries: maxRetries,
		tempDir:    tempDir,
	}, nil
}

// Execute executes a command
func (e *Executor) Execute(cmd *types.Command) *types.ExecutionResult {
	startTime := time.Now()
	result := &types.ExecutionResult{
		CommandID:  cmd.ID,
		ExecutedAt: startTime,
	}

	logger.Log.WithFields(map[string]interface{}{
		"commandId":   cmd.ID.String(),
		"commandType": cmd.CommandType,
	}).Info("Executing command")

	var scriptContent string
	var err error

	// Get script content based on command type
	switch cmd.CommandType {
	case types.CommandTypeScript:
		// Execute script directly
		scriptContent = cmd.Data
	case types.CommandTypeURL:
		// Fetch from URL
		scriptContent, err = e.fetchFromURL(cmd.Data)
		if err != nil {
			result.Success = false
			result.Error = fmt.Sprintf("Failed to fetch script from URL: %v", err)
			result.Duration = time.Since(startTime)
			return result
		}
	default:
		result.Success = false
		result.Error = fmt.Sprintf("Unknown command type: %d", cmd.CommandType)
		result.Duration = time.Since(startTime)
		return result
	}

	// Execute with retry
	for attempt := 1; attempt <= e.maxRetries; attempt++ {
		output, execErr := e.executeScript(scriptContent)

		if execErr == nil {
			// Success
			result.Success = true
			result.Output = output
			result.Duration = time.Since(startTime)
			logger.Log.WithField("commandId", cmd.ID.String()).Info("Command executed successfully")
			return result
		}

		logger.Log.WithFields(map[string]interface{}{
			"commandId": cmd.ID.String(),
			"attempt":   attempt,
			"error":     execErr,
		}).Warn("Execution failed, retrying...")

		if attempt < e.maxRetries {
			time.Sleep(time.Second * time.Duration(attempt))
		}

		err = execErr
	}

	// All retries failed
	result.Success = false
	result.Error = err.Error()
	result.Duration = time.Since(startTime)
	logger.Log.WithField("commandId", cmd.ID.String()).Error("Command execution failed after all retries")

	return result
}
func (e *Executor) executeScript(base64Script string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()
	base64Script = strings.TrimSpace(base64Script)
	base64Script = strings.Trim(base64Script, `"`)
	// 🔐 Decode Base64 → raw script
	raw, err := base64.StdEncoding.DecodeString(base64Script)
	fmt.Printf("Executing command %s\n", base64Script)

	if err != nil {
		return "", fmt.Errorf("invalid base64 script: %w", err)
	}
	scriptContent := string(raw)

	var cmd *exec.Cmd
	var scriptFile string

	switch runtime.GOOS {
	case "windows":
		scriptFile, err = e.createTempScript(scriptContent, ".ps1")
		if err != nil {
			return "", err
		}
		defer os.Remove(scriptFile)

		cmd = exec.CommandContext(
			ctx,
			"powershell",
			"-NoProfile",
			"-ExecutionPolicy", "Bypass",
			"-File", scriptFile,
		)

	case "darwin", "linux":
		// Ensure shebang
		if !strings.HasPrefix(scriptContent, "#!") {
			scriptContent = "#!/bin/bash\n" + scriptContent
		}

		scriptFile, err = e.createTempScript(scriptContent, ".sh")
		if err != nil {
			return "", err
		}
		defer os.Remove(scriptFile)

		_ = os.Chmod(scriptFile, 0755)
		cmd = exec.CommandContext(ctx, "/bin/bash", scriptFile)

	default:
		return "", fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	// Capture output
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()

	output := stdout.String()
	if stderr.Len() > 0 {
		output += "\nSTDERR:\n" + stderr.String()
	}

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return output, fmt.Errorf("execution timeout after %v", e.timeout)
		}
		return output, fmt.Errorf("execution failed: %w\nOutput: %s", err, output)
	}

	return output, nil
}

// createTempScript creates a temporary script file
func (e *Executor) createTempScript(content, ext string) (string, error) {
	file, err := os.CreateTemp(e.tempDir, "script-*"+ext)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer file.Close()

	// Unescape common escape sequences from blockchain data
	unescapedContent := strings.ReplaceAll(content, "\\n", "\n")
	unescapedContent = strings.ReplaceAll(unescapedContent, "\\t", "\t")
	unescapedContent = strings.ReplaceAll(unescapedContent, "\\r", "\r")

	if _, err := file.WriteString(unescapedContent); err != nil {
		os.Remove(file.Name())
		return "", fmt.Errorf("failed to write script: %w", err)
	}

	return file.Name(), nil
}

// fetchFromURL fetches script content from URL
func (e *Executor) fetchFromURL(url string) (string, error) {
	logger.Log.WithField("url", url).Info("Fetching script from URL")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	content := string(body)

	// Trim whitespace
	content = strings.TrimSpace(content)

	logger.Log.WithField("size", len(content)).Info("Script fetched successfully")

	return content, nil
}

// Cleanup cleans up temporary files
func (e *Executor) Cleanup() error {
	return os.RemoveAll(e.tempDir)
}
