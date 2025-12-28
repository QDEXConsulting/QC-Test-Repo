package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// Predefined errors
var (
	ErrItemNotFound      = errors.New("item not found")
	ErrItemAlreadyExists = errors.New("item already exists")
	ErrInvalidItemName   = errors.New("invalid item name")
	ErrInvalidPrice      = errors.New("invalid price")
	ErrInvalidQuantity   = errors.New("invalid quantity")
	ErrInvalidMinStock   = errors.New("invalid minimum stock")
	ErrInvalidMaxStock   = errors.New("invalid maximum stock")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrInvalidSKU        = errors.New("invalid SKU")
	ErrInvalidItemID     = errors.New("invalid item ID")
	ErrStockExceedsMax   = errors.New("stock exceeds maximum allowed")
	ErrEmptyInventory    = errors.New("inventory is empty")
)

// InventoryManager manages inventory operations
type InventoryManager struct {
	inventory *Inventory
	mu        sync.RWMutex // Mutex for thread-safe operations
}

// NewInventoryManager creates a new inventory manager
func NewInventoryManager() *InventoryManager {
	return &InventoryManager{
		inventory: NewInventory(),
	}
}

// SetInventory sets the inventory for the manager
func (im *InventoryManager) SetInventory(inv *Inventory) {
	im.inventory = inv
}

// GetInventory returns the current inventory
func (im *InventoryManager) GetInventory() *Inventory {
	return im.inventory
}

// AddItem adds a new item to the inventory
func (im *InventoryManager) AddItem(item *Item) error {
	im.mu.Lock()
	defer im.mu.Unlock()

	if item == nil {
		return fmt.Errorf("item cannot be nil")
	}

	if item.ID == "" {
		return ErrInvalidItemID
	}

	if err := item.Validate(); err != nil {
		return err
	}

	if item.SKU != "" && !validateSKU(item.SKU) {
		return ErrInvalidSKU
	}

	// Check if item with same ID already exists
	if _, exists := im.inventory.Items[item.ID]; exists {
		return fmt.Errorf("%w: item with ID %s already exists", ErrItemAlreadyExists, item.ID)
	}

	// Check if SKU is already used by another item
	if item.SKU != "" {
		for _, existingItem := range im.inventory.Items {
			if existingItem.SKU == item.SKU {
				return fmt.Errorf("%w: SKU %s already exists", ErrItemAlreadyExists, item.SKU)
			}
		}
	}

	now := time.Now()
	item.UpdatedAt = now
	item.CreatedAt = now
	im.inventory.Items[item.ID] = item
	im.addCategoryIfNotExists(item.Category)
	im.inventory.LastSync = now
	im.inventory.Version++

	return nil
}

// UpdateItem updates an existing item
func (im *InventoryManager) UpdateItem(itemID string, updates *Item) error {
	im.mu.Lock()
	defer im.mu.Unlock()

	if itemID == "" {
		return ErrInvalidItemID
	}

	item, exists := im.inventory.Items[itemID]
	if !exists {
		return fmt.Errorf("%w: item with ID %s", ErrItemNotFound, itemID)
	}

	if updates == nil {
		return fmt.Errorf("updates cannot be nil")
	}

	hasChanges := false

	if updates.Name != "" && updates.Name != item.Name {
		item.Name = updates.Name
		hasChanges = true
	}
	if updates.Description != item.Description {
		item.Description = updates.Description
		hasChanges = true
	}
	if updates.Category != "" && updates.Category != item.Category {
		oldCategory := item.Category
		item.Category = updates.Category
		im.addCategoryIfNotExists(updates.Category)
		im.removeCategoryIfUnused(oldCategory)
		hasChanges = true
	}
	if updates.Price >= 0 && updates.Price != item.Price {
		item.Price = updates.Price
		hasChanges = true
	}
	if updates.MinStock >= 0 && updates.MinStock != item.MinStock {
		item.MinStock = updates.MinStock
		hasChanges = true
	}
	if updates.MaxStock >= 0 && updates.MaxStock != item.MaxStock {
		item.MaxStock = updates.MaxStock
		hasChanges = true
	}
	if updates.Supplier != "" && updates.Supplier != item.Supplier {
		item.Supplier = updates.Supplier
		hasChanges = true
	}
	if updates.SKU != "" && updates.SKU != item.SKU {
		if !validateSKU(updates.SKU) {
			return ErrInvalidSKU
		}
		// Check if SKU is already used by another item
		for id, existingItem := range im.inventory.Items {
			if id != itemID && existingItem.SKU == updates.SKU {
				return fmt.Errorf("%w: SKU %s already exists", ErrItemAlreadyExists, updates.SKU)
			}
		}
		item.SKU = updates.SKU
		hasChanges = true
	}

	if !hasChanges {
		return nil // No changes to apply
	}

	item.UpdatedAt = time.Now()
	if err := item.Validate(); err != nil {
		return err
	}

	im.inventory.LastSync = time.Now()
	im.inventory.Version++

	return nil
}

// DeleteItem removes an item from inventory (soft delete)
func (im *InventoryManager) DeleteItem(itemID string) error {
	im.mu.Lock()
	defer im.mu.Unlock()

	if itemID == "" {
		return ErrInvalidItemID
	}

	item, exists := im.inventory.Items[itemID]
	if !exists {
		return fmt.Errorf("%w: item with ID %s", ErrItemNotFound, itemID)
	}

	if !item.IsActive {
		return nil // Already deleted
	}

	item.IsActive = false
	item.UpdatedAt = time.Now()
	im.inventory.LastSync = time.Now()
	im.inventory.Version++

	return nil
}

// RemoveItem completely removes an item from inventory
func (im *InventoryManager) RemoveItem(itemID string) error {
	im.mu.Lock()
	defer im.mu.Unlock()

	if itemID == "" {
		return ErrInvalidItemID
	}

	item, exists := im.inventory.Items[itemID]
	if !exists {
		return fmt.Errorf("%w: item with ID %s", ErrItemNotFound, itemID)
	}

	category := item.Category
	delete(im.inventory.Items, itemID)
	im.removeCategoryIfUnused(category)
	im.inventory.LastSync = time.Now()
	im.inventory.Version++

	return nil
}

// GetItem retrieves an item by ID
func (im *InventoryManager) GetItem(itemID string) (*Item, error) {
	im.mu.RLock()
	defer im.mu.RUnlock()

	if itemID == "" {
		return nil, ErrInvalidItemID
	}

	item, exists := im.inventory.Items[itemID]
	if !exists {
		return nil, fmt.Errorf("%w: item with ID %s", ErrItemNotFound, itemID)
	}
	return item, nil
}

// GetItemBySKU retrieves an item by SKU
func (im *InventoryManager) GetItemBySKU(sku string) (*Item, error) {
	im.mu.RLock()
	defer im.mu.RUnlock()

	if sku == "" {
		return nil, ErrInvalidSKU
	}

	for _, item := range im.inventory.Items {
		if item.SKU == sku {
			return item, nil
		}
	}
	return nil, fmt.Errorf("%w: item with SKU %s", ErrItemNotFound, sku)
}

// GetAllItems returns all items in the inventory
func (im *InventoryManager) GetAllItems() []*Item {
	im.mu.RLock()
	defer im.mu.RUnlock()

	items := make([]*Item, 0, len(im.inventory.Items))
	for _, item := range im.inventory.Items {
		items = append(items, item)
	}
	return items
}

// GetActiveItems returns only active items
func (im *InventoryManager) GetActiveItems() []*Item {
	im.mu.RLock()
	defer im.mu.RUnlock()

	items := make([]*Item, 0, len(im.inventory.Items))
	for _, item := range im.inventory.Items {
		if item.IsActive {
			items = append(items, item)
		}
	}
	return items
}

// GetInactiveItems returns only inactive (soft-deleted) items
func (im *InventoryManager) GetInactiveItems() []*Item {
	im.mu.RLock()
	defer im.mu.RUnlock()

	items := make([]*Item, 0)
	for _, item := range im.inventory.Items {
		if !item.IsActive {
			items = append(items, item)
		}
	}
	return items
}

// AddStock increases the quantity of an item
func (im *InventoryManager) AddStock(itemID string, quantity int) error {
	im.mu.Lock()
	defer im.mu.Unlock()

	if itemID == "" {
		return ErrInvalidItemID
	}
	if quantity <= 0 {
		return fmt.Errorf("%w: quantity must be positive, got %d", ErrInvalidQuantity, quantity)
	}

	item, exists := im.inventory.Items[itemID]
	if !exists {
		return fmt.Errorf("%w: item with ID %s", ErrItemNotFound, itemID)
	}

	if !item.IsActive {
		return fmt.Errorf("cannot add stock to inactive item %s", itemID)
	}

	if item.MaxStock > 0 && item.Quantity+quantity > item.MaxStock {
		return fmt.Errorf("%w: stock addition would exceed maximum stock of %d (current: %d, adding: %d)", 
			ErrStockExceedsMax, item.MaxStock, item.Quantity, quantity)
	}

	item.Quantity += quantity
	item.UpdatedAt = time.Now()
	im.inventory.LastSync = time.Now()
	im.inventory.Version++

	return nil
}

// RemoveStock decreases the quantity of an item
func (im *InventoryManager) RemoveStock(itemID string, quantity int) error {
	im.mu.Lock()
	defer im.mu.Unlock()

	if itemID == "" {
		return ErrInvalidItemID
	}
	if quantity <= 0 {
		return fmt.Errorf("%w: quantity must be positive, got %d", ErrInvalidQuantity, quantity)
	}

	item, exists := im.inventory.Items[itemID]
	if !exists {
		return fmt.Errorf("%w: item with ID %s", ErrItemNotFound, itemID)
	}

	if !item.IsActive {
		return fmt.Errorf("cannot remove stock from inactive item %s", itemID)
	}

	if item.Quantity < quantity {
		return fmt.Errorf("%w: requested %d but only %d available for item %s", 
			ErrInsufficientStock, quantity, item.Quantity, itemID)
	}

	item.Quantity -= quantity
	item.UpdatedAt = time.Now()
	im.inventory.LastSync = time.Now()
	im.inventory.Version++

	return nil
}

// SetStock sets the quantity of an item to a specific value
func (im *InventoryManager) SetStock(itemID string, quantity int) error {
	im.mu.Lock()
	defer im.mu.Unlock()

	if itemID == "" {
		return ErrInvalidItemID
	}
	if quantity < 0 {
		return fmt.Errorf("%w: quantity cannot be negative, got %d", ErrInvalidQuantity, quantity)
	}

	item, exists := im.inventory.Items[itemID]
	if !exists {
		return fmt.Errorf("%w: item with ID %s", ErrItemNotFound, itemID)
	}

	if !item.IsActive {
		return fmt.Errorf("cannot set stock for inactive item %s", itemID)
	}

	if item.MaxStock > 0 && quantity > item.MaxStock {
		return fmt.Errorf("%w: stock %d exceeds maximum stock of %d for item %s", 
			ErrStockExceedsMax, quantity, item.MaxStock, itemID)
	}

	item.Quantity = quantity
	item.UpdatedAt = time.Now()
	im.inventory.LastSync = time.Now()
	im.inventory.Version++

	return nil
}

// GetCategories returns all unique categories
func (im *InventoryManager) GetCategories() []string {
	im.mu.RLock()
	defer im.mu.RUnlock()

	// Return a copy to prevent external modification
	categories := make([]string, len(im.inventory.Categories))
	copy(categories, im.inventory.Categories)
	return categories
}

// GetItemsByCategory returns all items in a specific category
func (im *InventoryManager) GetItemsByCategory(category string) []*Item {
	im.mu.RLock()
	defer im.mu.RUnlock()

	if category == "" {
		return []*Item{}
	}

	items := make([]*Item, 0)
	for _, item := range im.inventory.Items {
		if item.Category == category && item.IsActive {
			items = append(items, item)
		}
	}
	return items
}

// GetItemsBySupplier returns all items from a specific supplier
func (im *InventoryManager) GetItemsBySupplier(supplier string) []*Item {
	im.mu.RLock()
	defer im.mu.RUnlock()

	if supplier == "" {
		return []*Item{}
	}

	items := make([]*Item, 0)
	for _, item := range im.inventory.Items {
		if item.Supplier == supplier && item.IsActive {
			items = append(items, item)
		}
	}
	return items
}

// GetItemCount returns the total number of items
func (im *InventoryManager) GetItemCount() int {
	im.mu.RLock()
	defer im.mu.RUnlock()
	return len(im.inventory.Items)
}

// GetActiveItemCount returns the number of active items
func (im *InventoryManager) GetActiveItemCount() int {
	im.mu.RLock()
	defer im.mu.RUnlock()

	count := 0
	for _, item := range im.inventory.Items {
		if item.IsActive {
			count++
		}
	}
	return count
}

// HasItem checks if an item exists by ID
func (im *InventoryManager) HasItem(itemID string) bool {
	im.mu.RLock()
	defer im.mu.RUnlock()

	if itemID == "" {
		return false
	}
	_, exists := im.inventory.Items[itemID]
	return exists
}

// HasItemBySKU checks if an item exists by SKU
func (im *InventoryManager) HasItemBySKU(sku string) bool {
	im.mu.RLock()
	defer im.mu.RUnlock()

	if sku == "" {
		return false
	}
	for _, item := range im.inventory.Items {
		if item.SKU == sku {
			return true
		}
	}
	return false
}

// GetInventoryVersion returns the current version of the inventory
func (im *InventoryManager) GetInventoryVersion() int {
	im.mu.RLock()
	defer im.mu.RUnlock()
	return im.inventory.Version
}

// GetLastSyncTime returns the last sync time
func (im *InventoryManager) GetLastSyncTime() time.Time {
	im.mu.RLock()
	defer im.mu.RUnlock()
	return im.inventory.LastSync
}

// addCategoryIfNotExists adds a category to the list if it doesn't exist
func (im *InventoryManager) addCategoryIfNotExists(category string) {
	if category != "" && !contains(im.inventory.Categories, category) {
		im.inventory.Categories = append(im.inventory.Categories, category)
	}
}

// removeCategoryIfUnused removes a category if no items use it
func (im *InventoryManager) removeCategoryIfUnused(category string) {
	if category == "" {
		return
	}
	for _, item := range im.inventory.Items {
		if item.Category == category && item.IsActive {
			return // Category still in use
		}
	}
	im.inventory.Categories = removeFromSlice(im.inventory.Categories, category)
}

// RestoreItem restores a soft-deleted item
func (im *InventoryManager) RestoreItem(itemID string) error {
	im.mu.Lock()
	defer im.mu.Unlock()

	if itemID == "" {
		return ErrInvalidItemID
	}

	item, exists := im.inventory.Items[itemID]
	if !exists {
		return fmt.Errorf("%w: item with ID %s", ErrItemNotFound, itemID)
	}

	if item.IsActive {
		return nil // Already active
	}

	item.IsActive = true
	item.UpdatedAt = time.Now()
	im.addCategoryIfNotExists(item.Category)
	im.inventory.LastSync = time.Now()
	im.inventory.Version++

	return nil
}

// ClearInventory removes all items from inventory (use with caution)
func (im *InventoryManager) ClearInventory() error {
	im.mu.Lock()
	defer im.mu.Unlock()

	im.inventory.Items = make(map[string]*Item)
	im.inventory.Categories = make([]string, 0)
	im.inventory.LastSync = time.Now()
	im.inventory.Version++

	return nil
}

// BulkAddItems adds multiple items at once
func (im *InventoryManager) BulkAddItems(items []*Item) error {
	im.mu.Lock()
	defer im.mu.Unlock()

	if len(items) == 0 {
		return fmt.Errorf("no items provided")
	}

	// Validate all items first
	for _, item := range items {
		if item == nil {
			return fmt.Errorf("item cannot be nil")
		}
		if item.ID == "" {
			return ErrInvalidItemID
		}
		if err := item.Validate(); err != nil {
			return fmt.Errorf("validation failed for item %s: %w", item.ID, err)
		}
		if item.SKU != "" && !validateSKU(item.SKU) {
			return fmt.Errorf("%w: invalid SKU %s for item %s", ErrInvalidSKU, item.SKU, item.ID)
		}
	}

	// Check for duplicates
	itemIDs := make(map[string]bool)
	skus := make(map[string]bool)
	for _, item := range items {
		if itemIDs[item.ID] {
			return fmt.Errorf("%w: duplicate item ID %s", ErrItemAlreadyExists, item.ID)
		}
		if item.SKU != "" {
			if skus[item.SKU] {
				return fmt.Errorf("%w: duplicate SKU %s", ErrItemAlreadyExists, item.SKU)
			}
			skus[item.SKU] = true
		}
		itemIDs[item.ID] = true

		// Check against existing items
		if _, exists := im.inventory.Items[item.ID]; exists {
			return fmt.Errorf("%w: item with ID %s already exists", ErrItemAlreadyExists, item.ID)
		}
		if item.SKU != "" {
			for _, existingItem := range im.inventory.Items {
				if existingItem.SKU == item.SKU {
					return fmt.Errorf("%w: SKU %s already exists", ErrItemAlreadyExists, item.SKU)
				}
			}
		}
	}

	// Add all items
	now := time.Now()
	for _, item := range items {
		item.UpdatedAt = now
		item.CreatedAt = now
		im.inventory.Items[item.ID] = item
		im.addCategoryIfNotExists(item.Category)
	}

	im.inventory.LastSync = now
	im.inventory.Version++

	return nil
}