package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds application configuration
type Config struct {
	Server     ServerConfig     `json:"server"`
	Database   DatabaseConfig   `json:"database"`
	Cache      CacheConfig      `json:"cache"`
	Logging    LoggingConfig    `json:"logging"`
	Security   SecurityConfig   `json:"security"`
	Inventory  InventoryConfig  `json:"inventory"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port         int           `json:"port"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
	IdleTimeout  time.Duration `json:"idle_timeout"`
	Host         string        `json:"host"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	FilePath     string `json:"file_path"`
	BackupPath   string `json:"backup_path"`
	BackupCount  int    `json:"backup_count"`
	AutoSave     bool   `json:"auto_save"`
	SaveInterval time.Duration `json:"save_interval"`
}

// CacheConfig holds cache configuration
type CacheConfig struct {
	Enabled     bool          `json:"enabled"`
	DefaultTTL  time.Duration `json:"default_ttl"`
	MaxSize     int          `json:"max_size"`
	StatsEnabled bool         `json:"stats_enabled"`
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level      string `json:"level"`
	File       string `json:"file"`
	MaxSize    int    `json:"max_size"`
	MaxBackups int    `json:"max_backups"`
	Compress   bool   `json:"compress"`
}

// SecurityConfig holds security configuration
type SecurityConfig struct {
	APIKey          string        `json:"api_key"`
	EnableCORS      bool          `json:"enable_cors"`
	RateLimit       int           `json:"rate_limit"`
	RateLimitWindow time.Duration `json:"rate_limit_window"`
	EnableHTTPS     bool          `json:"enable_https"`
	CertFile        string        `json:"cert_file"`
	KeyFile         string        `json:"key_file"`
}

// InventoryConfig holds inventory-specific configuration
type InventoryConfig struct {
	DefaultMinStock    int     `json:"default_min_stock"`
	DefaultMaxStock    int     `json:"default_max_stock"`
	LowStockThreshold  float64 `json:"low_stock_threshold"`
	AutoBackup         bool    `json:"auto_backup"`
	BackupInterval     time.Duration `json:"backup_interval"`
	MaxBulkOperations  int     `json:"max_bulk_operations"`
}

// DefaultConfig returns a configuration with default values
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         8080,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
			Host:         "localhost",
		},
		Database: DatabaseConfig{
			FilePath:     "data/inventory.json",
			BackupPath:   "data/backups",
			BackupCount:  5,
			AutoSave:     true,
			SaveInterval: 5 * time.Minute,
		},
		Cache: CacheConfig{
			Enabled:      true,
			DefaultTTL:   5 * time.Minute,
			MaxSize:      1000,
			StatsEnabled: true,
		},
		Logging: LoggingConfig{
			Level:      "INFO",
			File:       "",
			MaxSize:    100,
			MaxBackups: 5,
			Compress:   false,
		},
		Security: SecurityConfig{
			APIKey:          "",
			EnableCORS:      true,
			RateLimit:       100,
			RateLimitWindow: 1 * time.Minute,
			EnableHTTPS:     false,
			CertFile:        "",
			KeyFile:         "",
		},
		Inventory: InventoryConfig{
			DefaultMinStock:    10,
			DefaultMaxStock:    1000,
			LowStockThreshold:  0.2,
			AutoBackup:         true,
			BackupInterval:     24 * time.Hour,
			MaxBulkOperations:  100,
		},
	}
}

// LoadConfig loads configuration from a JSON file
func LoadConfig(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}
	
	// Apply environment variable overrides
	config.applyEnvOverrides()
	
	return &config, nil
}

// SaveConfig saves configuration to a JSON file
func (c *Config) SaveConfig(filePath string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	
	return nil
}

// applyEnvOverrides applies environment variable overrides
func (c *Config) applyEnvOverrides() {
	// Server config
	if port := os.Getenv("SERVER_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			c.Server.Port = p
		}
	}
	
	if host := os.Getenv("SERVER_HOST"); host != "" {
		c.Server.Host = host
	}
	
	// Database config
	if path := os.Getenv("DB_FILE_PATH"); path != "" {
		c.Database.FilePath = path
	}
	
	if backupPath := os.Getenv("DB_BACKUP_PATH"); backupPath != "" {
		c.Database.BackupPath = backupPath
	}
	
	// Cache config
	if enabled := os.Getenv("CACHE_ENABLED"); enabled != "" {
		c.Cache.Enabled = enabled == "true"
	}
	
	// Logging config
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		c.Logging.Level = level
	}
	
	if logFile := os.Getenv("LOG_FILE"); logFile != "" {
		c.Logging.File = logFile
	}
	
	// Security config
	if apiKey := os.Getenv("API_KEY"); apiKey != "" {
		c.Security.APIKey = apiKey
	}
	
	if enableCORS := os.Getenv("ENABLE_CORS"); enableCORS != "" {
		c.Security.EnableCORS = enableCORS == "true"
	}
	
	if rateLimit := os.Getenv("RATE_LIMIT"); rateLimit != "" {
		if rl, err := strconv.Atoi(rateLimit); err == nil {
			c.Security.RateLimit = rl
		}
	}
	
	// Inventory config
	if minStock := os.Getenv("DEFAULT_MIN_STOCK"); minStock != "" {
		if ms, err := strconv.Atoi(minStock); err == nil {
			c.Inventory.DefaultMinStock = ms
		}
	}
	
	if maxStock := os.Getenv("DEFAULT_MAX_STOCK"); maxStock != "" {
		if ms, err := strconv.Atoi(maxStock); err == nil {
			c.Inventory.DefaultMaxStock = ms
		}
	}
}

// GetLogLevel returns the LogLevel from the logging config
func (c *Config) GetLogLevel() LogLevel {
	switch c.Logging.Level {
	case "DEBUG":
		return LogLevelDebug
	case "INFO":
		return LogLevelInfo
	case "WARNING":
		return LogLevelWarning
	case "ERROR":
		return LogLevelError
	case "FATAL":
		return LogLevelFatal
	default:
		return LogLevelInfo
	}
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Server.Port < 1 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}
	
	if c.Database.FilePath == "" {
		return fmt.Errorf("database file path is required")
	}
	
	if c.Cache.MaxSize < 0 {
		return fmt.Errorf("cache max size cannot be negative")
	}
	
	if c.Security.RateLimit < 0 {
		return fmt.Errorf("rate limit cannot be negative")
	}
	
	if c.Inventory.DefaultMinStock < 0 {
		return fmt.Errorf("default min stock cannot be negative")
	}
	
	if c.Inventory.DefaultMaxStock > 0 && c.Inventory.DefaultMaxStock < c.Inventory.DefaultMinStock {
		return fmt.Errorf("default max stock must be greater than min stock")
	}
	
	return nil
}

// Global config instance
var globalConfig *Config

// LoadGlobalConfig loads the global configuration
func LoadGlobalConfig(filePath string) error {
	var err error
	if filePath != "" {
		globalConfig, err = LoadConfig(filePath)
	} else {
		globalConfig = DefaultConfig()
		globalConfig.applyEnvOverrides()
	}
	
	if err != nil {
		return err
	}
	
	return globalConfig.Validate()
}

// GetConfig returns the global configuration
func GetConfig() *Config {
	if globalConfig == nil {
		globalConfig = DefaultConfig()
		globalConfig.applyEnvOverrides()
	}
	return globalConfig
}

