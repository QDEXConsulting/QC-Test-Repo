package main

import (
	"sync"
	"time"
)

// MetricType represents the type of metric
type MetricType string

const (
	MetricTypeCounter MetricType = "counter"
	MetricTypeGauge   MetricType = "gauge"
	MetricTypeHistogram MetricType = "histogram"
)

// Metric represents a single metric
type Metric struct {
	Name      string
	Type      MetricType
	Value     float64
	Labels    map[string]string
	Timestamp time.Time
}

// MetricsCollector collects and stores application metrics
type MetricsCollector struct {
	metrics map[string]*Metric
	mu      sync.RWMutex
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		metrics: make(map[string]*Metric),
	}
}

// Increment increments a counter metric
func (mc *MetricsCollector) Increment(name string, labels map[string]string) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	
	key := mc.getKey(name, labels)
	if metric, exists := mc.metrics[key]; exists {
		metric.Value++
		metric.Timestamp = time.Now()
	} else {
		mc.metrics[key] = &Metric{
			Name:      name,
			Type:      MetricTypeCounter,
			Value:     1,
			Labels:    labels,
			Timestamp: time.Now(),
		}
	}
}

// Set sets a gauge metric value
func (mc *MetricsCollector) Set(name string, value float64, labels map[string]string) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	
	key := mc.getKey(name, labels)
	mc.metrics[key] = &Metric{
		Name:      name,
		Type:      MetricTypeGauge,
		Value:     value,
		Labels:    labels,
		Timestamp: time.Now(),
	}
}

// Record records a histogram value
func (mc *MetricsCollector) Record(name string, value float64, labels map[string]string) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	
	key := mc.getKey(name, labels)
	if metric, exists := mc.metrics[key]; exists {
		// For simplicity, we'll just update the value
		// In a real implementation, you'd maintain buckets
		metric.Value = value
		metric.Timestamp = time.Now()
	} else {
		mc.metrics[key] = &Metric{
			Name:      name,
			Type:      MetricTypeHistogram,
			Value:     value,
			Labels:    labels,
			Timestamp: time.Now(),
		}
	}
}

// Get retrieves a metric value
func (mc *MetricsCollector) Get(name string, labels map[string]string) (float64, bool) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	
	key := mc.getKey(name, labels)
	if metric, exists := mc.metrics[key]; exists {
		return metric.Value, true
	}
	return 0, false
}

// GetAll returns all metrics
func (mc *MetricsCollector) GetAll() map[string]*Metric {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	
	result := make(map[string]*Metric)
	for k, v := range mc.metrics {
		result[k] = v
	}
	return result
}

// Reset clears all metrics
func (mc *MetricsCollector) Reset() {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.metrics = make(map[string]*Metric)
}

// getKey generates a unique key for a metric with labels
func (mc *MetricsCollector) getKey(name string, labels map[string]string) string {
	if len(labels) == 0 {
		return name
	}
	
	key := name
	for k, v := range labels {
		key += ":" + k + "=" + v
	}
	return key
}

// ApplicationMetrics tracks application-specific metrics
type ApplicationMetrics struct {
	collector *MetricsCollector
}

// NewApplicationMetrics creates a new application metrics tracker
func NewApplicationMetrics() *ApplicationMetrics {
	return &ApplicationMetrics{
		collector: NewMetricsCollector(),
	}
}

// RecordHTTPRequest records an HTTP request
func (am *ApplicationMetrics) RecordHTTPRequest(method, path string, statusCode int, duration time.Duration) {
	am.collector.Increment("http_requests_total", map[string]string{
		"method": method,
		"path":   path,
		"status": string(rune(statusCode)),
	})
	
	am.collector.Record("http_request_duration_seconds", duration.Seconds(), map[string]string{
		"method": method,
		"path":   path,
	})
}

// RecordItemOperation records an item operation
func (am *ApplicationMetrics) RecordItemOperation(operation string) {
	am.collector.Increment("item_operations_total", map[string]string{
		"operation": operation,
	})
}

// RecordStockOperation records a stock operation
func (am *ApplicationMetrics) RecordStockOperation(operation string, quantity int) {
	am.collector.Increment("stock_operations_total", map[string]string{
		"operation": operation,
	})
	
	am.collector.Record("stock_operation_quantity", float64(quantity), map[string]string{
		"operation": operation,
	})
}

// RecordCacheHit records a cache hit
func (am *ApplicationMetrics) RecordCacheHit(cacheType string) {
	am.collector.Increment("cache_hits_total", map[string]string{
		"cache_type": cacheType,
	})
}

// RecordCacheMiss records a cache miss
func (am *ApplicationMetrics) RecordCacheMiss(cacheType string) {
	am.collector.Increment("cache_misses_total", map[string]string{
		"cache_type": cacheType,
	})
}

// RecordInventorySize records the current inventory size
func (am *ApplicationMetrics) RecordInventorySize(totalItems, activeItems int) {
	am.collector.Set("inventory_items_total", float64(totalItems), nil)
	am.collector.Set("inventory_items_active", float64(activeItems), nil)
}

// RecordInventoryValue records the total inventory value
func (am *ApplicationMetrics) RecordInventoryValue(value float64) {
	am.collector.Set("inventory_value_total", value, nil)
}

// RecordLowStockItems records the number of low stock items
func (am *ApplicationMetrics) RecordLowStockItems(count int) {
	am.collector.Set("inventory_low_stock_items", float64(count), nil)
}

// RecordOutOfStockItems records the number of out of stock items
func (am *ApplicationMetrics) RecordOutOfStockItems(count int) {
	am.collector.Set("inventory_out_of_stock_items", float64(count), nil)
}

// RecordStorageOperation records a storage operation
func (am *ApplicationMetrics) RecordStorageOperation(operation string, success bool) {
	status := "success"
	if !success {
		status = "failure"
	}
	
	am.collector.Increment("storage_operations_total", map[string]string{
		"operation": operation,
		"status":    status,
	})
}

// RecordError records an error occurrence
func (am *ApplicationMetrics) RecordError(errorType string) {
	am.collector.Increment("errors_total", map[string]string{
		"type": errorType,
	})
}

// GetMetrics returns all collected metrics
func (am *ApplicationMetrics) GetMetrics() map[string]*Metric {
	return am.collector.GetAll()
}

// MetricsSnapshot provides a snapshot of current metrics
type MetricsSnapshot struct {
	Timestamp time.Time
	Metrics   map[string]*Metric
}

// GetSnapshot returns a snapshot of current metrics
func (am *ApplicationMetrics) GetSnapshot() *MetricsSnapshot {
	return &MetricsSnapshot{
		Timestamp: time.Now(),
		Metrics:   am.collector.GetAll(),
	}
}

// Global metrics instance
var globalMetrics *ApplicationMetrics

// InitGlobalMetrics initializes the global metrics instance
func InitGlobalMetrics() {
	globalMetrics = NewApplicationMetrics()
}

// GetMetrics returns the global metrics instance
func GetMetrics() *ApplicationMetrics {
	if globalMetrics == nil {
		InitGlobalMetrics()
	}
	return globalMetrics
}

