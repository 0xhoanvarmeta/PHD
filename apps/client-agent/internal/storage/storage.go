package storage

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"sync"
)

type Storage struct {
	filePath       string
	executedCmds   map[string]bool
	lastCommandID  *big.Int
	isFirstRun     bool
	mu             sync.RWMutex
}

type storageData struct {
	ExecutedCmds  []string `json:"executed_commands"`
	LastCommandID string   `json:"last_command_id"`
}

func NewStorage() (*Storage, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home dir: %w", err)
	}

	storageDir := filepath.Join(homeDir, ".phd-client-agent")
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage dir: %w", err)
	}

	filePath := filepath.Join(storageDir, "executed.json")

	// Check if this is first run (storage file doesn't exist)
	_, err = os.Stat(filePath)
	isFirstRun := os.IsNotExist(err)

	s := &Storage{
		filePath:      filePath,
		executedCmds:  make(map[string]bool),
		lastCommandID: big.NewInt(0),
		isFirstRun:    isFirstRun,
	}

	if err := s.load(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Storage) load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, err := os.Stat(s.filePath); os.IsNotExist(err) {
		return nil
	}

	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return fmt.Errorf("failed to read storage: %w", err)
	}

	var sd storageData
	if err := json.Unmarshal(data, &sd); err != nil {
		return fmt.Errorf("failed to unmarshal storage: %w", err)
	}

	for _, cmdID := range sd.ExecutedCmds {
		s.executedCmds[cmdID] = true
	}

	if sd.LastCommandID != "" {
		s.lastCommandID = new(big.Int)
		s.lastCommandID.SetString(sd.LastCommandID, 10)
	}

	return nil
}

func (s *Storage) save() error {
	executedList := make([]string, 0, len(s.executedCmds))
	for cmdID := range s.executedCmds {
		executedList = append(executedList, cmdID)
	}

	sd := storageData{
		ExecutedCmds:  executedList,
		LastCommandID: s.lastCommandID.String(),
	}

	data, err := json.MarshalIndent(sd, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal storage: %w", err)
	}

	if err := os.WriteFile(s.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write storage: %w", err)
	}

	return nil
}

func (s *Storage) IsExecuted(commandID *big.Int) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.executedCmds[commandID.String()]
}

func (s *Storage) MarkExecuted(commandID *big.Int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.executedCmds[commandID.String()] = true

	if commandID.Cmp(s.lastCommandID) > 0 {
		s.lastCommandID = commandID
	}

	return s.save()
}

func (s *Storage) GetLastCommandID() *big.Int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return new(big.Int).Set(s.lastCommandID)
}

func (s *Storage) IsFirstRun() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isFirstRun
}
