package main

import (
	"context"
	"fmt"
	"time"
)

// InventoryService provides high-level business logic for inventory management
type InventoryService struct {
	manager *InventoryManager
	storage *Storage
	logger  *Logger
}

// TransactionHistory stores transaction records
type TransactionHistory struct {
	transactions []*Transaction
	maxSize      int
}

// NewInventoryService creates a new inventory service
func NewInventoryService(manager *InventoryManager, storage *Storage) *InventoryService {
	return &InventoryService{
		manager: manager,
		storage: storage,
		logger:  GetLogger(),
	}
}

// NewInventoryServiceWithLogger creates a new inventory service with a custom logger
func NewInventoryServiceWithLogger(manager *InventoryManager, storage *Storage, logger *Logger) *InventoryService {
	return &InventoryService{
		manager: manager,
		storage: storage,
		logger:  logger,
	}
}

// Initialize loads inventory from storage
func (s *InventoryService) Initialize() error {
	return s.InitializeWithContext(context.Background())
}

// InitializeWithContext loads inventory from storage with context support
func (s *InventoryService) InitializeWithContext(ctx context.Context) error {
	start := time.Now()
	defer func() {
		if s.logger != nil {
			s.logger.LogPerformance("Initialize", time.Since(start))
		}
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("initialization cancelled: %w", ctx.Err())
	default:
	}

	inventory, err := s.storage.LoadInventory()
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(err, "Initialize")
		}
		return fmt.Errorf("failed to load inventory: %w", err)
	}
	s.manager.SetInventory(inventory)
	if s.logger != nil {
		s.logger.Info("Inventory initialized successfully with %d items", len(inventory.Items))
	}
	return nil
}

// Save persists the current inventory to storage
func (s *InventoryService) Save() error {
	return s.SaveWithContext(context.Background())
}

// SaveWithContext persists the current inventory to storage with context support
func (s *InventoryService) SaveWithContext(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("save cancelled: %w", ctx.Err())
	default:
	}

	if err := s.storage.SaveInventory(s.manager.GetInventory()); err != nil {
		if s.logger != nil {
			s.logger.LogError(err, "Save")
		}
		return fmt.Errorf("failed to save inventory: %w", err)
	}
	return nil
}

// CreateItem creates a new item and saves it
func (s *InventoryService) CreateItem(name, description, category, sku, supplier string, price float64, quantity, minStock, maxStock int) (*Item, error) {
	return s.CreateItemWithContext(context.Background(), name, description, category, sku, supplier, price, quantity, minStock, maxStock)
}

// CreateItemWithContext creates a new item and saves it with context support
func (s *InventoryService) CreateItemWithContext(ctx context.Context, name, description, category, sku, supplier string, price float64, quantity, minStock, maxStock int) (*Item, error) {
	// Input validation
	if name == "" {
		return nil, fmt.Errorf("item name cannot be empty")
	}
	if price < 0 {
		return nil, fmt.Errorf("price cannot be negative")
	}
	if quantity < 0 {
		return nil, fmt.Errorf("quantity cannot be negative")
	}
	if minStock < 0 {
		return nil, fmt.Errorf("min stock cannot be negative")
	}
	if maxStock > 0 && maxStock < minStock {
		return nil, fmt.Errorf("max stock cannot be less than min stock")
	}

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("create item cancelled: %w", ctx.Err())
	default:
	}

	item := NewItem(name, description, category, sku, supplier, price, quantity, minStock, maxStock)
	if err := s.manager.AddItem(item); err != nil {
		if s.logger != nil {
			s.logger.LogError(err, "CreateItem")
		}
		return nil, err
	}
	if err := s.SaveWithContext(ctx); err != nil {
		return nil, err
	}
	if s.logger != nil {
		s.logger.LogItemOperation("CREATE", item.ID, item.Name)
	}
	return item, nil
}

// UpdateItem updates an existing item and saves changes
func (s *InventoryService) UpdateItem(itemID string, updates *Item) error {
	return s.UpdateItemWithContext(context.Background(), itemID, updates)
}

// UpdateItemWithContext updates an existing item and saves changes with context support
func (s *InventoryService) UpdateItemWithContext(ctx context.Context, itemID string, updates *Item) error {
	if itemID == "" {
		return fmt.Errorf("item ID cannot be empty")
	}
	if updates == nil {
		return fmt.Errorf("updates cannot be nil")
	}

	select {
	case <-ctx.Done():
		return fmt.Errorf("update item cancelled: %w", ctx.Err())
	default:
	}

	item, err := s.manager.GetItem(itemID)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(err, "UpdateItem")
		}
		return err
	}

	if err := s.manager.UpdateItem(itemID, updates); err != nil {
		if s.logger != nil {
			s.logger.LogError(err, "UpdateItem")
		}
		return err
	}
	if err := s.SaveWithContext(ctx); err != nil {
		return err
	}
	if s.logger != nil {
		s.logger.LogItemOperation("UPDATE", itemID, item.Name)
	}
	return nil
}

// DeleteItem soft deletes an item and saves changes
func (s *InventoryService) DeleteItem(itemID string) error {
	return s.DeleteItemWithContext(context.Background(), itemID)
}

// DeleteItemWithContext soft deletes an item and saves changes with context support
func (s *InventoryService) DeleteItemWithContext(ctx context.Context, itemID string) error {
	if itemID == "" {
		return fmt.Errorf("item ID cannot be empty")
	}

	select {
	case <-ctx.Done():
		return fmt.Errorf("delete item cancelled: %w", ctx.Err())
	default:
	}

	item, err := s.manager.GetItem(itemID)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(err, "DeleteItem")
		}
		return err
	}

	if err := s.manager.DeleteItem(itemID); err != nil {
		if s.logger != nil {
			s.logger.LogError(err, "DeleteItem")
		}
		return err
	}
	if err := s.SaveWithContext(ctx); err != nil {
		return err
	}
	if s.logger != nil {
		s.logger.LogItemOperation("DELETE", itemID, item.Name)
	}
	return nil
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
	return s.ProcessStockInWithContext(context.Background(), itemID, quantity, reason)
}

// ProcessStockInWithContext increases stock for an item with context support
func (s *InventoryService) ProcessStockInWithContext(ctx context.Context, itemID string, quantity int, reason string) error {
	if itemID == "" {
		return fmt.Errorf("item ID cannot be empty")
	}
	if quantity <= 0 {
		return fmt.Errorf("quantity must be positive")
	}

	select {
	case <-ctx.Done():
		return fmt.Errorf("stock in cancelled: %w", ctx.Err())
	default:
	}

	if err := s.manager.AddStock(itemID, quantity); err != nil {
		if s.logger != nil {
			s.logger.LogError(err, "ProcessStockIn")
		}
		return err
	}
	if err := s.SaveWithContext(ctx); err != nil {
		return err
	}
	if s.logger != nil {
		s.logger.LogStockOperation("STOCK_IN", itemID, quantity, reason)
	}
	return nil
}

// ProcessStockOut decreases stock for an item
func (s *InventoryService) ProcessStockOut(itemID string, quantity int, reason string) error {
	return s.ProcessStockOutWithContext(context.Background(), itemID, quantity, reason)
}

// ProcessStockOutWithContext decreases stock for an item with context support
func (s *InventoryService) ProcessStockOutWithContext(ctx context.Context, itemID string, quantity int, reason string) error {
	if itemID == "" {
		return fmt.Errorf("item ID cannot be empty")
	}
	if quantity <= 0 {
		return fmt.Errorf("quantity must be positive")
	}

	select {
	case <-ctx.Done():
		return fmt.Errorf("stock out cancelled: %w", ctx.Err())
	default:
	}

	if err := s.manager.RemoveStock(itemID, quantity); err != nil {
		if s.logger != nil {
			s.logger.LogError(err, "ProcessStockOut")
		}
		return err
	}
	if err := s.SaveWithContext(ctx); err != nil {
		return err
	}
	if s.logger != nil {
		s.logger.LogStockOperation("STOCK_OUT", itemID, quantity, reason)
	}
	return nil
}

// AdjustStock sets stock to a specific quantity
func (s *InventoryService) AdjustStock(itemID string, quantity int, reason string) error {
	return s.AdjustStockWithContext(context.Background(), itemID, quantity, reason)
}

// AdjustStockWithContext sets stock to a specific quantity with context support
func (s *InventoryService) AdjustStockWithContext(ctx context.Context, itemID string, quantity int, reason string) error {
	if itemID == "" {
		return fmt.Errorf("item ID cannot be empty")
	}
	if quantity < 0 {
		return fmt.Errorf("quantity cannot be negative")
	}

	select {
	case <-ctx.Done():
		return fmt.Errorf("adjust stock cancelled: %w", ctx.Err())
	default:
	}

	if err := s.manager.SetStock(itemID, quantity); err != nil {
		if s.logger != nil {
			s.logger.LogError(err, "AdjustStock")
		}
		return err
	}
	if err := s.SaveWithContext(ctx); err != nil {
		return err
	}
	if s.logger != nil {
		s.logger.LogStockOperation("ADJUST_STOCK", itemID, quantity, reason)
	}
	return nil
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
	return s.BulkUpdateStockWithContext(context.Background(), updates)
}

// BulkUpdateStockWithContext updates stock for multiple items with context support and rollback on error
func (s *InventoryService) BulkUpdateStockWithContext(ctx context.Context, updates map[string]int) error {
	if len(updates) == 0 {
		return fmt.Errorf("updates map cannot be empty")
	}

	select {
	case <-ctx.Done():
		return fmt.Errorf("bulk update cancelled: %w", ctx.Err())
	default:
	}

	// Store original quantities for rollback
	originalQuantities := make(map[string]int)
	for itemID := range updates {
		item, err := s.manager.GetItem(itemID)
		if err != nil {
			return fmt.Errorf("item %s not found: %w", itemID, err)
		}
		originalQuantities[itemID] = item.Quantity
	}

	// Apply updates
	for itemID, quantity := range updates {
		if quantity < 0 {
			// Rollback on validation error
			for id, origQty := range originalQuantities {
				_ = s.manager.SetStock(id, origQty)
			}
			return fmt.Errorf("quantity cannot be negative for item %s", itemID)
		}
		if err := s.manager.SetStock(itemID, quantity); err != nil {
			// Rollback on error
			for id, origQty := range originalQuantities {
				_ = s.manager.SetStock(id, origQty)
			}
			return fmt.Errorf("failed to update stock for item %s: %w", itemID, err)
		}
	}

	if err := s.SaveWithContext(ctx); err != nil {
		// Rollback on save error
		for id, origQty := range originalQuantities {
			_ = s.manager.SetStock(id, origQty)
		}
		return err
	}

	if s.logger != nil {
		s.logger.Info("Bulk stock update completed for %d items", len(updates))
	}
	return nil
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

// GetItem retrieves an item by ID
func (s *InventoryService) GetItem(itemID string) (*Item, error) {
	if itemID == "" {
		return nil, fmt.Errorf("item ID cannot be empty")
	}
	return s.manager.GetItem(itemID)
}

// GetAllItems returns all active items
func (s *InventoryService) GetAllItems() []*Item {
	return s.manager.GetActiveItems()
}

// GetItemsBySupplier returns items from a specific supplier
func (s *InventoryService) GetItemsBySupplier(supplier string) []*Item {
	if supplier == "" {
		return []*Item{}
	}
	items := s.manager.GetActiveItems()
	result := make([]*Item, 0)
	for _, item := range items {
		if item.Supplier == supplier {
			result = append(result, item)
		}
	}
	return result
}

// GetItemsByPriceRange returns items within a price range
func (s *InventoryService) GetItemsByPriceRange(minPrice, maxPrice float64) []*Item {
	if minPrice < 0 || maxPrice < 0 || maxPrice < minPrice {
		return []*Item{}
	}
	items := s.manager.GetActiveItems()
	result := make([]*Item, 0)
	for _, item := range items {
		if item.Price >= minPrice && item.Price <= maxPrice {
			result = append(result, item)
		}
	}
	return result
}

// ValidateItem validates an item before operations
func (s *InventoryService) ValidateItem(item *Item) error {
	if item == nil {
		return fmt.Errorf("item cannot be nil")
	}
	return item.Validate()
}

// GetTotalInventoryValue calculates the total value of all inventory
func (s *InventoryService) GetTotalInventoryValue() float64 {
	items := s.manager.GetActiveItems()
	total := 0.0
	for _, item := range items {
		total += item.GetValue()
	}
	return total
}

// GetItemsNeedingReorder returns items that need to be reordered
func (s *InventoryService) GetItemsNeedingReorder() []*Item {
	items := s.manager.GetActiveItems()
	result := make([]*Item, 0)
	for _, item := range items {
		if item.Quantity <= item.MinStock {
			result = append(result, item)
		}
	}
	return result
}

// CheckStockAvailability checks if sufficient stock is available for an item
func (s *InventoryService) CheckStockAvailability(itemID string, requiredQuantity int) (bool, error) {
	if itemID == "" {
		return false, fmt.Errorf("item ID cannot be empty")
	}
	if requiredQuantity <= 0 {
		return false, fmt.Errorf("required quantity must be positive")
	}

	item, err := s.manager.GetItem(itemID)
	if err != nil {
		return false, err
	}
	return item.Quantity >= requiredQuantity, nil
}

