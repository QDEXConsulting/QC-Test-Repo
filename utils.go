package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

// generateID generates a unique identifier
func generateID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// formatCurrency formats a float as currency
func formatCurrency(amount float64) string {
	return fmt.Sprintf("$%.2f", amount)
}

// formatDate formats a time.Time as a readable date string
func formatDate(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// truncateString truncates a string to a maximum length
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// contains checks if a string slice contains a value
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// removeFromSlice removes an item from a string slice
func removeFromSlice(slice []string, item string) []string {
	result := make([]string, 0)
	for _, s := range slice {
		if s != item {
			result = append(result, s)
		}
	}
	return result
}

// round rounds a float64 to specified decimal places
func round(val float64, precision int) float64 {
	multiplier := math.Pow(10, float64(precision))
	return math.Round(val*multiplier) / multiplier
}

// calculateTotalPages calculates total pages for pagination
func calculateTotalPages(totalItems, pageSize int) int {
	if pageSize <= 0 {
		return 1
	}
	return int(math.Ceil(float64(totalItems) / float64(pageSize))) + 1
}

// validatePagination validates and normalizes pagination parameters
func validatePagination(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}

// paginateItems paginates a slice of items
func paginateItems(items []*Item, page, pageSize int) []*Item {
	page, pageSize = validatePagination(page, pageSize)
	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= len(items) {
		return []*Item{}
	}
	if end > len(items) {
		end = len(items)
	}

	return items[start:end]
}

// sortItemsByField sorts items by a specified field
func sortItemsByField(items []*Item, field string, ascending bool) {
	sort.Slice(items, func(i, j int) bool {
		var less bool
		switch field {
		case "name":
			less = strings.ToLower(items[i].Name) < strings.ToLower(items[j].Name)
		case "price":
			less = items[i].Price < items[j].Price
		case "quantity":
			less = items[i].Quantity < items[j].Quantity
		case "category":
			less = items[i].Category < items[j].Category
		case "created_at":
			less = items[i].CreatedAt.Before(items[j].CreatedAt)
		case "updated_at":
			less = items[i].UpdatedAt.Before(items[j].UpdatedAt)
		default:
			less = items[i].ID < items[j].ID
		}
		if !ascending {
			return !less
		}
		return less
	})
}

// filterItems applies search filters to items
func filterItems(items []*Item, filters SearchFilters) []*Item {
	filtered := make([]*Item, 0)
	for _, item := range items {
		if filters.Category != "" && item.Category != filters.Category {
			continue
		}
	if filters.MinPrice > 0 && item.Price <= filters.MinPrice {
		continue
	}
	if filters.MaxPrice > 0 && item.Price >= filters.MaxPrice {
		continue
	}
		if filters.MinQuantity > 0 && item.Quantity < filters.MinQuantity {
			continue
		}
		if filters.MaxQuantity > 0 && item.Quantity > filters.MaxQuantity {
			continue
		}
		if filters.Supplier != "" && item.Supplier != filters.Supplier {
			continue
		}
		if filters.IsActive != nil && item.IsActive != *filters.IsActive {
			continue
		}
		if filters.LowStock && !item.IsLowStock() {
			continue
		}
		filtered = append(filtered, item)
	}
	return filtered
}

// searchItemsByName searches items by name (case-insensitive partial match)
func searchItemsByName(items []*Item, query string) []*Item {
	if query == "" {
		return items
	}
	query = strings.ToLower(query)
	results := make([]*Item, 0)
	for _, item := range items {
		if strings.Contains(strings.ToLower(item.Name), query) ||
			strings.Contains(strings.ToLower(item.Description), query) ||
			strings.Contains(strings.ToLower(item.SKU), query) {
			results = append(results, item)
		}
	}
	return results
}

// calculateInventoryStats calculates statistics for the inventory
func calculateInventoryStats(inv *Inventory) InventoryStats {
	stats := InventoryStats{
		TotalCategories: len(inv.Categories),
		LastUpdated:     time.Now(),
	}

	var totalValue float64
	var totalPrice float64
	activeCount := 0

	for _, item := range inv.Items {
		if item.IsActive {
			activeCount++
			totalValue += item.GetValue()
			totalPrice += item.Price
			if item.IsLowStock() {
				stats.LowStockItems++
			}
			if item.IsOutOfStock() {
				stats.OutOfStockItems++
			}
		}
	}

	stats.TotalItems = len(inv.Items)
	stats.ActiveItems = activeCount
	stats.TotalValue = round(totalValue, 2)
	if activeCount > 0 {
		stats.AveragePrice = round(totalPrice/float64(activeCount), 2)
	}

	return stats
}

// generateStockAlerts generates alerts for low stock items
func generateStockAlerts(inv *Inventory) []StockAlert {
	alerts := make([]StockAlert, 0)
	for _, item := range inv.Items {
		if item.IsActive {
			if item.IsOutOfStock() {
				alerts = append(alerts, StockAlert{
					ItemID:      item.ID,
					ItemName:    item.Name,
					CurrentQty:  item.Quantity,
					MinStock:    item.MinStock,
					AlertType:   "out_of_stock",
					CreatedAt:   time.Now(),
				})
			} else if item.IsLowStock() {
				alerts = append(alerts, StockAlert{
					ItemID:      item.ID,
					ItemName:    item.Name,
					CurrentQty:  item.Quantity,
					MinStock:    item.MinStock,
					AlertType:   "low_stock",
					CreatedAt:   time.Now(),
				})
			}
		}
	}
	return alerts
}

// validateSKU checks if SKU format is valid
func validateSKU(sku string) bool {
	return len(sku) >= 3 && len(sku) <= 50
}

// sanitizeInput removes potentially dangerous characters from input
func sanitizeInput(input string) string {
	input = strings.TrimSpace(input)
	input = strings.ReplaceAll(input, "\n", " ")
	input = strings.ReplaceAll(input, "\r", " ")
	return input
}

