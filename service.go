package main

import (
	"fmt"
	"time"
)

// InventoryService provides high-level business logic for inventory management
type InventoryService struct {
	manager *InventoryManager
	storage *Storage
}

// NewInventoryService creates a new inventory service
func NewInventoryService(manager *InventoryManager, storage *Storage) *InventoryService {
	return &InventoryService{
		manager: manager,
		storage: storage,
	}
}

// Initialize loads inventory from storage
func (s *InventoryService) Initialize() error {
	inventory, err := s.storage.LoadInventory()
	if err != nil {
		return fmt.Errorf("failed to load inventory: %w", err)
	}
	s.manager.SetInventory(inventory)
	return nil
}

// Save persists the current inventory to storage
func (s *InventoryService) Save() error {
	if err := s.storage.SaveInventory(s.manager.GetInventory()); err != nil {
		return fmt.Errorf("failed to save inventory: %w", err)
	}
	return nil
}

// CreateItem creates a new item and saves it
func (s *InventoryService) CreateItem(name, description, category, sku, supplier string, price float64, quantity, minStock, maxStock int) (*Item, error) {
	item := NewItem(name, description, category, sku, supplier, price, quantity, minStock, maxStock)
	if err := s.manager.AddItem(item); err != nil {
		return nil, err
	}
	if err := s.Save(); err != nil {
		return nil, err
	}
	return item, nil
}

// UpdateItem updates an existing item and saves changes
func (s *InventoryService) UpdateItem(itemID string, updates *Item) error {
	if err := s.manager.UpdateItem(itemID, updates); err != nil {
		return err
	}
	return s.Save()
}

// DeleteItem soft deletes an item and saves changes
func (s *InventoryService) DeleteItem(itemID string) error {
	if err := s.manager.DeleteItem(itemID); err != nil {
		return err
	}
	return s.Save()
}

// RestoreItem restores a soft-deleted item
func (s *InventoryService) RestoreItem(itemID string) error {
	item, err := s.manager.GetItem(itemID)
	if err != nil {
		return err
	}
	item.IsActive = true
	item.UpdatedAt = time.Now()
	s.manager.GetInventory().LastSync = time.Now()
	s.manager.GetInventory().Version++
	return s.Save()
}

// ProcessStockIn increases stock for an item
func (s *InventoryService) ProcessStockIn(itemID string, quantity int, reason string) error {
	if err := s.manager.AddStock(itemID, quantity); err != nil {
		return err
	}
	return s.Save()
}

// ProcessStockOut decreases stock for an item
func (s *InventoryService) ProcessStockOut(itemID string, quantity int, reason string) error {
	if err := s.manager.RemoveStock(itemID, quantity); err != nil {
		return err
	}
	return s.Save()
}

// AdjustStock sets stock to a specific quantity
func (s *InventoryService) AdjustStock(itemID string, quantity int, reason string) error {
	if err := s.manager.SetStock(itemID, quantity); err != nil {
		return err
	}
	return s.Save()
}

// SearchItems searches items with filters and pagination
func (s *InventoryService) SearchItems(query string, filters SearchFilters, sortField string, ascending bool, page, pageSize int) (*PaginatedResponse, error) {
	items := s.manager.GetActiveItems()

	// Apply search query
	if query != "" {
		items = searchItemsByName(items, query)
	}

	// Apply filters
	items = filterItems(items, filters)

	// Sort items
	if sortField != "" {
		sortItemsByField(items, sortField, ascending)
	}

	// Paginate
	page, pageSize = validatePagination(page, pageSize)
	totalItems := len(items)
	paginatedItems := paginateItems(items, page, pageSize)

	return &PaginatedResponse{
		Data:       paginatedItems,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: calculateTotalPages(totalItems, pageSize),
		TotalItems: totalItems,
	}, nil
}

// GetInventoryStats returns statistics about the inventory
func (s *InventoryService) GetInventoryStats() InventoryStats {
	return calculateInventoryStats(s.manager.GetInventory())
}

// GetStockAlerts returns all stock alerts
func (s *InventoryService) GetStockAlerts() []StockAlert {
	return generateStockAlerts(s.manager.GetInventory())
}

// GetLowStockItems returns items that are low on stock
func (s *InventoryService) GetLowStockItems() []*Item {
	items := s.manager.GetActiveItems()
	lowStockItems := make([]*Item, 0)
	for _, item := range items {
		if item.IsLowStock() {
			lowStockItems = append(lowStockItems, item)
		}
	}
	return lowStockItems
}

// GetOutOfStockItems returns items that are out of stock
func (s *InventoryService) GetOutOfStockItems() []*Item {
	items := s.manager.GetActiveItems()
	outOfStockItems := make([]*Item, 0)
	for _, item := range items {
		if item.IsOutOfStock() {
			outOfStockItems = append(outOfStockItems, item)
		}
	}
	return outOfStockItems
}

// BulkUpdateStock updates stock for multiple items
func (s *InventoryService) BulkUpdateStock(updates map[string]int) error {
	for itemID, quantity := range updates {
		if err := s.manager.SetStock(itemID, quantity); err != nil {
			return fmt.Errorf("failed to update stock for item %s: %w", itemID, err)
		}
	}
	return s.Save()
}

// GetItemsByCategory returns items in a specific category
func (s *InventoryService) GetItemsByCategory(category string) []*Item {
	return s.manager.GetItemsByCategory(category)
}

// GetCategories returns all categories
func (s *InventoryService) GetCategories() []string {
	return s.manager.GetCategories()
}

// Backup creates a backup of the current inventory
func (s *InventoryService) Backup() error {
	return s.storage.BackupInventory(s.manager.GetInventory())
}

// Restore restores inventory from backup
func (s *InventoryService) Restore() error {
	inventory, err := s.storage.RestoreInventory()
	if err != nil {
		return fmt.Errorf("failed to restore from backup: %w", err)
	}
	s.manager.SetInventory(inventory)
	return s.Save()
}

// Export exports inventory to a file
func (s *InventoryService) Export(exportPath string) error {
	return s.storage.ExportInventory(s.manager.GetInventory(), exportPath)
}

// Import imports inventory from a file
func (s *InventoryService) Import(importPath string) error {
	inventory, err := s.storage.ImportInventory(importPath)
	if err != nil {
		return fmt.Errorf("failed to import inventory: %w", err)
	}
	s.manager.SetInventory(inventory)
	return s.Save()
}

// GetStorageStats returns storage statistics
func (s *InventoryService) GetStorageStats() StorageStats {
	return s.storage.GetStorageStats()
}

