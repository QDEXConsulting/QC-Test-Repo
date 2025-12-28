package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

// LogLevel represents the severity of a log message
type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarning
	LogLevelError
	LogLevelFatal
)

// String returns the string representation of LogLevel
func (l LogLevel) String() string {
	switch l {
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarning:
		return "WARNING"
	case LogLevelError:
		return "ERROR"
	case LogLevelFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Logger provides structured logging functionality
type Logger struct {
	level      LogLevel
	logger     *log.Logger
	file       *os.File
	enableFile bool
}

// NewLogger creates a new logger instance
func NewLogger(level LogLevel, logFile string) (*Logger, error) {
	logger := &Logger{
		level:      level,
		enableFile: logFile != "",
	}
	
	if logger.enableFile {
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}
		logger.file = file
		logger.logger = log.New(file, "", log.LstdFlags)
	} else {
		logger.logger = log.New(os.Stdout, "", log.LstdFlags)
	}
	
	return logger, nil
}

// Close closes the log file if it was opened
func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// SetLevel sets the minimum log level
func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

// shouldLog checks if a message at the given level should be logged
func (l *Logger) shouldLog(level LogLevel) bool {
	return level >= l.level
}

// log writes a log message with the given level
func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
	if !l.shouldLog(level) {
		return
	}
	
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, args...)
	logMessage := fmt.Sprintf("[%s] [%s] %s", timestamp, level.String(), message)
	
	l.logger.Println(logMessage)
}

// Debug logs a debug message
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(LogLevelDebug, format, args...)
}

// Info logs an info message
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(LogLevelInfo, format, args...)
}

// Warning logs a warning message
func (l *Logger) Warning(format string, args ...interface{}) {
	l.log(LogLevelWarning, format, args...)
}

// Error logs an error message
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(LogLevelError, format, args...)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log(LogLevelFatal, format, args...)
	os.Exit(1)
}

// LogRequest logs an HTTP request
func (l *Logger) LogRequest(method, path, remoteAddr string, statusCode int, duration time.Duration) {
	l.Info(
		"HTTP Request: %s %s from %s - Status: %d - Duration: %v",
		method,
		path,
		remoteAddr,
		statusCode,
		duration,
	)
}

// LogItemOperation logs an item operation
func (l *Logger) LogItemOperation(operation, itemID, itemName string) {
	l.Info("Item Operation: %s - ID: %s - Name: %s", operation, itemID, itemName)
}

// LogStockOperation logs a stock operation
func (l *Logger) LogStockOperation(operation, itemID string, quantity int, reason string) {
	l.Info(
		"Stock Operation: %s - ItemID: %s - Quantity: %d - Reason: %s",
		operation,
		itemID,
		quantity,
		reason,
	)
}

// LogError logs an error with context
func (l *Logger) LogError(err error, context string) {
	l.Error("Error in %s: %v", context, err)
}

// LogCacheHit logs a cache hit
func (l *Logger) LogCacheHit(key string) {
	l.Debug("Cache Hit: %s", key)
}

// LogCacheMiss logs a cache miss
func (l *Logger) LogCacheMiss(key string) {
	l.Debug("Cache Miss: %s", key)
}

// LogPerformance logs performance metrics
func (l *Logger) LogPerformance(operation string, duration time.Duration) {
	if duration > 1*time.Second {
		l.Warning("Slow Operation: %s took %v", operation, duration)
	} else {
		l.Debug("Operation: %s took %v", operation, duration)
	}
}

// LogStats logs inventory statistics
func (l *Logger) LogStats(stats InventoryStats) {
	l.Info(
		"Inventory Stats - Total Items: %d, Active: %d, Value: $%.2f, Low Stock: %d, Out of Stock: %d",
		stats.TotalItems,
		stats.ActiveItems,
		stats.TotalValue,
		stats.LowStockItems,
		stats.OutOfStockItems,
	)
}

// LogAlert logs a stock alert
func (l *Logger) LogAlert(alert StockAlert) {
	l.Warning(
		"Stock Alert: %s - Item: %s (ID: %s) - Current: %d, Min: %d",
		alert.AlertType,
		alert.ItemName,
		alert.ItemID,
		alert.CurrentQty,
		alert.MinStock,
	)
}

// LogBackup logs backup operations
func (l *Logger) LogBackup(operation string, success bool, details string) {
	if success {
		l.Info("Backup %s: Success - %s", operation, details)
	} else {
		l.Error("Backup %s: Failed - %s", operation, details)
	}
}

// LogStorageOperation logs storage operations
func (l *Logger) LogStorageOperation(operation string, success bool, details string) {
	if success {
		l.Info("Storage %s: Success - %s", operation, details)
	} else {
		l.Error("Storage %s: Failed - %s", operation, details)
	}
}

// Global logger instance
var globalLogger *Logger

// InitGlobalLogger initializes the global logger
func InitGlobalLogger(level LogLevel, logFile string) error {
	var err error
	globalLogger, err = NewLogger(level, logFile)
	return err
}

// GetLogger returns the global logger instance
func GetLogger() *Logger {
	if globalLogger == nil {
		globalLogger, _ = NewLogger(LogLevelInfo, "")
	}
	return globalLogger
}

