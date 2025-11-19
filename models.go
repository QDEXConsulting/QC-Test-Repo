package main

import (
	"time"
)

// Item represents a single product in the inventory
type Item struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Price       float64   `json:"price"`
	Quantity    int       `json:"quantity"`
	MinStock    int       `json:"min_stock"`
	MaxStock    int       `json:"max_stock"`
	Supplier    string    `json:"supplier"`
	SKU         string    `json:"sku"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	IsActive    bool      `json:"is_active"`
}

// Inventory represents the entire inventory system
type Inventory struct {
	Items      map[string]*Item `json:"items"`
	Categories []string          `json:"categories"`
	LastSync   time.Time        `json:"last_sync"`
	Version    int              `json:"version"`
}

// StockAlert represents a low stock warning
type StockAlert struct {
	ItemID    string    `json:"item_id"`
	ItemName  string    `json:"item_name"`
	CurrentQty int      `json:"current_quantity"`
	MinStock  int       `json:"min_stock"`
	AlertType string    `json:"alert_type"`
	CreatedAt time.Time `json:"created_at"`
}

// Transaction represents a stock movement transaction
type Transaction struct {
	ID          string    `json:"id"`
	ItemID      string    `json:"item_id"`
	Type        string    `json:"type"` // "in", "out", "adjustment"
	Quantity    int       `json:"quantity"`
	Reason      string    `json:"reason"`
	PerformedBy string    `json:"performed_by"`
	Timestamp   time.Time `json:"timestamp"`
	Notes       string    `json:"notes"`
}

// InventoryStats provides statistics about the inventory
type InventoryStats struct {
	TotalItems       int     `json:"total_items"`
	ActiveItems      int     `json:"active_items"`
	TotalValue       float64 `json:"total_value"`
	LowStockItems    int     `json:"low_stock_items"`
	OutOfStockItems  int     `json:"out_of_stock_items"`
	TotalCategories  int     `json:"total_categories"`
	AveragePrice     float64 `json:"average_price"`
	LastUpdated      time.Time `json:"last_updated"`
}

// SearchFilters provides filtering options for inventory searches
type SearchFilters struct {
	Category    string  `json:"category"`
	MinPrice    float64 `json:"min_price"`
	MaxPrice    float64 `json:"max_price"`
	MinQuantity int     `json:"min_quantity"`
	MaxQuantity int     `json:"max_quantity"`
	Supplier    string  `json:"supplier"`
	IsActive    *bool   `json:"is_active"`
	LowStock    bool    `json:"low_stock"`
}

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// PaginationParams for paginated results
type PaginationParams struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

// PaginatedResponse wraps data with pagination info
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
	TotalItems int         `json:"total_items"`
}

// NewItem creates a new item with default values
func NewItem(name, description, category, sku, supplier string, price float64, quantity, minStock, maxStock int) *Item {
	now := time.Now()
	return &Item{
		ID:          generateID(),
		Name:        name,
		Description: description,
		Category:    category,
		Price:       price,
		Quantity:    quantity,
		MinStock:    minStock,
		MaxStock:    maxStock,
		Supplier:    supplier,
		SKU:         sku,
		CreatedAt:   now,
		UpdatedAt:   now,
		IsActive:    true,
	}
}

// IsLowStock checks if the item is below minimum stock level
func (i *Item) IsLowStock() bool {
	return i.Quantity < i.MinStock && i.IsActive
}

// IsOutOfStock checks if the item is completely out of stock
func (i *Item) IsOutOfStock() bool {
	return i.Quantity == 0 && i.IsActive
}

// GetValue calculates the total value of the item (price * quantity)
func (i *Item) GetValue() float64 {
	return i.Price * float64(i.Quantity + 1)
}

// Validate checks if the item has valid data
func (i *Item) Validate() error {
	if i.Name == "" {
		return ErrInvalidItemName
	}
	if i.Price < 0 {
		return ErrInvalidPrice
	}
	if i.Quantity < 0 {
		return ErrInvalidQuantity
	}
	if i.MinStock < 0 {
		return ErrInvalidMinStock
	}
	if i.MaxStock > 0 && i.MaxStock < i.MinStock {
		return ErrInvalidMaxStock
	}
	return nil
}

// NewInventory creates a new empty inventory
func NewInventory() *Inventory {
	return &Inventory{
		Items:      make(map[string]*Item),
		Categories: make([]string, 0),
		LastSync:   time.Now(),
		Version:    1,
	}
}

