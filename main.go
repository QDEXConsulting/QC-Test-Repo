package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	dataFile = "data/inventory.json"
	apiPort  = 8080
)

func main() {
	// Initialize components
	manager := NewInventoryManager()
	storage := NewStorage(dataFile)
	service := NewInventoryService(manager, storage)

	// Load existing inventory
	if err := service.Initialize(); err != nil {
		log.Printf("Warning: Failed to load inventory: %v", err)
		log.Println("Starting with empty inventory")
	}

	// Check command line arguments
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "server", "api":
			startAPIServer(service)
			return
		case "help", "-h", "--help":
			printHelp()
			return
		}
	}

	// Start CLI interface
	startCLI(service)
}

// startAPIServer starts the HTTP API server
func startAPIServer(service *InventoryService) {
	server := NewAPIServer(service, apiPort)
	log.Printf("Starting API server on port %d", apiPort)
	log.Println("API endpoints available at http://localhost:8080/api/")
	log.Println("Press Ctrl+C to stop the server")
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// startCLI starts the command-line interface
func startCLI(service *InventoryService) {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("=== Inventory Management System ===")
	fmt.Println("Type 'help' for available commands")
	fmt.Println()

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		command := strings.ToLower(parts[0])

		switch command {
		case "exit", "quit", "q":
			fmt.Println("Saving inventory...")
			if err := service.Save(); err != nil {
				fmt.Printf("Error saving: %v\n", err)
			} else {
				fmt.Println("Inventory saved. Goodbye!")
			}
			return

		case "help", "h":
			printCommands()

		case "add", "create":
			handleAddItem(service, parts[1:])

		case "list", "ls":
			handleListItems(service)

		case "get", "show":
			if len(parts) < 2 {
				fmt.Println("Usage: get <item_id>")
				continue
			}
			handleGetItem(service, parts[1])

		case "update", "edit":
			if len(parts) < 2 {
				fmt.Println("Usage: update <item_id>")
				continue
			}
			handleUpdateItem(service, parts[1])

		case "delete", "remove":
			if len(parts) < 2 {
				fmt.Println("Usage: delete <item_id>")
				continue
			}
			handleDeleteItem(service, parts[1])

		case "stock", "inventory":
			handleStockOperations(service, parts[1:])

		case "search":
			handleSearch(service, parts[1:])

		case "stats", "statistics":
			handleStats(service)

		case "alerts":
			handleAlerts(service)

		case "categories":
			handleCategories(service)

		case "backup":
			handleBackup(service)

		case "restore":
			handleRestore(service)

		case "save":
			if err := service.Save(); err != nil {
				fmt.Printf("Error saving: %v\n", err)
			} else {
				fmt.Println("Inventory saved successfully")
			}

		default:
			fmt.Printf("Unknown command: %s. Type 'help' for available commands.\n", command)
		}
	}
}

// printHelp prints help information
func printHelp() {
	fmt.Println("Inventory Management System")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  go run .              - Start CLI interface")
	fmt.Println("  go run . server       - Start API server")
	fmt.Println("  go run . help         - Show this help")
	fmt.Println()
}

// printCommands prints available CLI commands
func printCommands() {
	fmt.Println("Available commands:")
	fmt.Println("  add/create            - Add a new item")
	fmt.Println("  list/ls              - List all items")
	fmt.Println("  get/show <id>        - Show item details")
	fmt.Println("  update/edit <id>     - Update an item")
	fmt.Println("  delete/remove <id>    - Delete an item")
	fmt.Println("  stock add <id> <qty> - Add stock")
	fmt.Println("  stock remove <id> <qty> - Remove stock")
	fmt.Println("  search <query>       - Search items")
	fmt.Println("  stats                - Show statistics")
	fmt.Println("  alerts               - Show stock alerts")
	fmt.Println("  categories           - List categories")
	fmt.Println("  backup               - Create backup")
	fmt.Println("  restore              - Restore from backup")
	fmt.Println("  save                 - Save inventory")
	fmt.Println("  exit/quit            - Exit and save")
}

// handleAddItem handles adding a new item
func handleAddItem(service *InventoryService, args []string) {
	fmt.Println("Adding new item...")
	fmt.Print("Name: ")
	name := readLine()
	fmt.Print("Description: ")
	description := readLine()
	fmt.Print("Category: ")
	category := readLine()
	fmt.Print("SKU: ")
	sku := readLine()
	fmt.Print("Supplier: ")
	supplier := readLine()
	fmt.Print("Price: ")
	price := readFloat()
	fmt.Print("Quantity: ")
	quantity := readInt()
	fmt.Print("Min Stock: ")
	minStock := readInt()
	fmt.Print("Max Stock (0 for unlimited): ")
	maxStock := readInt()

	item, err := service.CreateItem(name, description, category, sku, supplier, price, quantity, minStock, maxStock)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Item created successfully! ID: %s\n", item.ID)
}

// handleListItems lists all items
func handleListItems(service *InventoryService) {
	items := service.manager.GetActiveItems()
	if len(items) == 0 {
		fmt.Println("No items found")
		return
	}

	fmt.Printf("\n%-40s %-15s %-10s %-10s %-15s\n", "Name", "Category", "Price", "Quantity", "SKU")
	fmt.Println(strings.Repeat("-", 90))
	for _, item := range items {
		fmt.Printf("%-40s %-15s %-10.2f %-10d %-15s\n",
			truncateString(item.Name, 38),
			truncateString(item.Category, 13),
			item.Price,
			item.Quantity,
			item.SKU)
	}
	fmt.Println()
}

// handleGetItem shows item details
func handleGetItem(service *InventoryService, itemID string) {
	item, err := service.manager.GetItem(itemID)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("\nItem Details:\n")
	fmt.Printf("  ID:          %s\n", item.ID)
	fmt.Printf("  Name:        %s\n", item.Name)
	fmt.Printf("  Description: %s\n", item.Description)
	fmt.Printf("  Category:    %s\n", item.Category)
	fmt.Printf("  SKU:         %s\n", item.SKU)
	fmt.Printf("  Supplier:    %s\n", item.Supplier)
	fmt.Printf("  Price:       $%.2f\n", item.Price)
	fmt.Printf("  Quantity:    %d\n", item.Quantity)
	fmt.Printf("  Min Stock:   %d\n", item.MinStock)
	fmt.Printf("  Max Stock:   %d\n", item.MaxStock)
	fmt.Printf("  Value:       $%.2f\n", item.GetValue())
	fmt.Printf("  Status:      %s\n", getItemStatus(item))
	fmt.Printf("  Created:     %s\n", formatDate(item.CreatedAt))
	fmt.Printf("  Updated:     %s\n", formatDate(item.UpdatedAt))
	fmt.Println()
}

// handleUpdateItem handles updating an item
func handleUpdateItem(service *InventoryService, itemID string) {
	item, err := service.manager.GetItem(itemID)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("Leave blank to keep current value")
	fmt.Printf("Name [%s]: ", item.Name)
	name := readLine()
	if name == "" {
		name = item.Name
	}

	fmt.Printf("Description [%s]: ", item.Description)
	description := readLine()
	if description == "" {
		description = item.Description
	}

	fmt.Printf("Category [%s]: ", item.Category)
	category := readLine()
	if category == "" {
		category = item.Category
	}

	fmt.Printf("Price [%.2f]: ", item.Price)
	price := readFloat()
	if price < 0 {
		price = item.Price
	}

	updates := &Item{
		Name:        name,
		Description: description,
		Category:    category,
		Price:       price,
	}

	if err := service.UpdateItem(itemID, updates); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("Item updated successfully")
}

// handleDeleteItem handles deleting an item
func handleDeleteItem(service *InventoryService, itemID string) {
	fmt.Print("Are you sure? (yes/no): ")
	confirm := readLine()
	if strings.ToLower(confirm) != "yes" {
		fmt.Println("Cancelled")
		return
	}

	if err := service.DeleteItem(itemID); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("Item deleted successfully")
}

// handleStockOperations handles stock operations
func handleStockOperations(service *InventoryService, args []string) {
	if len(args) < 3 {
		fmt.Println("Usage: stock <add|remove|set> <item_id> <quantity>")
		return
	}

	operation := strings.ToLower(args[0])
	itemID := args[1]
	quantity, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Println("Invalid quantity")
		return
	}

	var err2 error
	switch operation {
	case "add":
		err2 = service.ProcessStockIn(itemID, quantity, "Manual adjustment")
	case "remove":
		err2 = service.ProcessStockOut(itemID, quantity, "Manual adjustment")
	case "set":
		err2 = service.AdjustStock(itemID, quantity, "Manual adjustment")
	default:
		fmt.Println("Invalid operation. Use: add, remove, or set")
		return
	}

	if err2 != nil {
		fmt.Printf("Error: %v\n", err2)
		return
	}

	fmt.Println("Stock operation completed successfully")
}

// handleSearch handles searching items
func handleSearch(service *InventoryService, args []string) {
	query := strings.Join(args, " ")
	if query == "" {
		fmt.Println("Usage: search <query>")
		return
	}

	result, err := service.SearchItems(query, SearchFilters{}, "", true, 1, 100)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	items := result.Data.([]*Item)
	if len(items) == 0 {
		fmt.Println("No items found")
		return
	}

	fmt.Printf("\nFound %d items:\n", len(items))
	handleListItems(service)
}

// handleStats shows inventory statistics
func handleStats(service *InventoryService) {
	stats := service.GetInventoryStats()
	fmt.Printf("\nInventory Statistics:\n")
	fmt.Printf("  Total Items:      %d\n", stats.TotalItems)
	fmt.Printf("  Active Items:     %d\n", stats.ActiveItems)
	fmt.Printf("  Total Value:      $%.2f\n", stats.TotalValue)
	fmt.Printf("  Average Price:    $%.2f\n", stats.AveragePrice)
	fmt.Printf("  Low Stock Items:  %d\n", stats.LowStockItems)
	fmt.Printf("  Out of Stock:     %d\n", stats.OutOfStockItems)
	fmt.Printf("  Categories:       %d\n", stats.TotalCategories)
	fmt.Println()
}

// handleAlerts shows stock alerts
func handleAlerts(service *InventoryService) {
	alerts := service.GetStockAlerts()
	if len(alerts) == 0 {
		fmt.Println("No stock alerts")
		return
	}

	fmt.Printf("\nStock Alerts (%d):\n", len(alerts))
	for _, alert := range alerts {
		fmt.Printf("  [%s] %s: %d (min: %d)\n",
			alert.AlertType, alert.ItemName, alert.CurrentQty, alert.MinStock)
	}
	fmt.Println()
}

// handleCategories lists all categories
func handleCategories(service *InventoryService) {
	categories := service.GetCategories()
	if len(categories) == 0 {
		fmt.Println("No categories found")
		return
	}

	fmt.Println("\nCategories:")
	for _, cat := range categories {
		items := service.GetItemsByCategory(cat)
		fmt.Printf("  %s (%d items)\n", cat, len(items))
	}
	fmt.Println()
}

// handleBackup creates a backup
func handleBackup(service *InventoryService) {
	if err := service.Backup(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("Backup created successfully")
}

// handleRestore restores from backup
func handleRestore(service *InventoryService) {
	fmt.Print("This will overwrite current inventory. Continue? (yes/no): ")
	confirm := readLine()
	if strings.ToLower(confirm) != "yes" {
		fmt.Println("Cancelled")
		return
	}

	if err := service.Restore(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("Inventory restored successfully")
}

// Helper functions
func readLine() string {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}

func readInt() int {
	var val int
	fmt.Scanf("%d\n", &val)
	return val
}

func readFloat() float64 {
	var val float64
	fmt.Scanf("%f\n", &val)
	return val
}

func getItemStatus(item *Item) string {
	if !item.IsActive {
		return "Inactive"
	}
	if item.IsOutOfStock() {
		return "Out of Stock"
	}
	if item.IsLowStock() {
		return "Low Stock"
	}
	return "In Stock"
}

