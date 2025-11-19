package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// Storage handles persistence of inventory data
type Storage struct {
	filePath string
	mu       sync.RWMutex
}

// NewStorage creates a new storage instance
func NewStorage(filePath string) *Storage {
	return &Storage{
		filePath: filePath,
	}
}

// SaveInventory saves the inventory to a JSON file
func (s *Storage) SaveInventory(inventory *Inventory) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Create directory if it doesn't exist
	dir := filepath.Dir(s.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create temporary file for atomic write
	tempFile := s.filePath + ".tmp"

	// Marshal inventory to JSON
	data, err := json.MarshalIndent(inventory, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal inventory: %w", err)
	}

	// Write to temporary file
	if err := os.WriteFile(tempFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write temporary file: %w", err)
	}

	// Atomic rename
	if err := os.Rename(tempFile, s.filePath); err != nil {
		os.Remove(tempFile) // Clean up temp file on error
		return fmt.Errorf("failed to rename file: %w", err)
	}

	return nil
}

// LoadInventory loads the inventory from a JSON file
func (s *Storage) LoadInventory() (*Inventory, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Check if file exists
	if _, err := os.Stat(s.filePath); os.IsNotExist(err) {
		// Return empty inventory if file doesn't exist
		return NewInventory(), nil
	}

	// Read file
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Unmarshal JSON
	var inventory Inventory
	if err := json.Unmarshal(data, &inventory); err != nil {
		return nil, fmt.Errorf("failed to unmarshal inventory: %w", err)
	}

	// Initialize Items map if nil
	if inventory.Items == nil {
		inventory.Items = make(map[string]*Item)
	}

	// Initialize Categories slice if nil
	if inventory.Categories == nil {
		inventory.Categories = make([]string, 0)
	}

	return &inventory, nil
}

// BackupInventory creates a backup of the current inventory
func (s *Storage) BackupInventory(inventory *Inventory) error {
	backupPath := s.filePath + ".backup"
	backupStorage := NewStorage(backupPath)
	return backupStorage.SaveInventory(inventory)
}

// RestoreInventory restores inventory from a backup file
func (s *Storage) RestoreInventory() (*Inventory, error) {
	backupPath := s.filePath + ".backup"
	backupStorage := NewStorage(backupPath)
	return backupStorage.LoadInventory()
}

// FileExists checks if the storage file exists
func (s *Storage) FileExists() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, err := os.Stat(s.filePath)
	return !os.IsNotExist(err)
}

// DeleteFile deletes the storage file
func (s *Storage) DeleteFile() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return os.Remove(s.filePath)
}

// GetFilePath returns the file path
func (s *Storage) GetFilePath() string {
	return s.filePath
}

// ExportInventory exports inventory to a specified file path
func (s *Storage) ExportInventory(inventory *Inventory, exportPath string) error {
	exportStorage := NewStorage(exportPath)
	return exportStorage.SaveInventory(inventory)
}

// ImportInventory imports inventory from a specified file path
func (s *Storage) ImportInventory(importPath string) (*Inventory, error) {
	importStorage := NewStorage(importPath)
	return importStorage.LoadInventory()
}

// ValidateFile checks if the storage file is valid JSON
func (s *Storage) ValidateFile() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, err := os.Stat(s.filePath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist")
	}

	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var inventory Inventory
	if err := json.Unmarshal(data, &inventory); err != nil {
		return fmt.Errorf("invalid JSON format: %w", err)
	}

	return nil
}

// GetFileSize returns the size of the storage file in bytes
func (s *Storage) GetFileSize() (int64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	info, err := os.Stat(s.filePath)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// StorageStats provides information about storage
type StorageStats struct {
	FilePath    string `json:"file_path"`
	FileExists  bool   `json:"file_exists"`
	FileSize    int64  `json:"file_size"`
	IsValid     bool   `json:"is_valid"`
	BackupExists bool  `json:"backup_exists"`
}

// GetStorageStats returns statistics about the storage
func (s *Storage) GetStorageStats() StorageStats {
	stats := StorageStats{
		FilePath: s.filePath,
	}

	stats.FileExists = s.FileExists()
	if stats.FileExists {
		if size, err := s.GetFileSize(); err == nil {
			stats.FileSize = size
		}
		if err := s.ValidateFile(); err == nil {
			stats.IsValid = true
		}
	}

	backupPath := s.filePath + ".backup"
	if _, err := os.Stat(backupPath); err == nil {
		stats.BackupExists = true
	}

	return stats
}

