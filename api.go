package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// APIServer handles HTTP API requests
type APIServer struct {
	service *InventoryService
	port    int
}

// NewAPIServer creates a new API server
func NewAPIServer(service *InventoryService, port int) *APIServer {
	return &APIServer{
		service: service,
		port:    port,
	}
}

// Start starts the HTTP server
func (s *APIServer) Start() error {
	mux := http.NewServeMux()

	// Item endpoints
	mux.HandleFunc("/api/items", s.handleItems)
	mux.HandleFunc("/api/items/", s.handleItemByID)
	mux.HandleFunc("/api/items/sku/", s.handleItemBySKU)

	// Stock endpoints
	mux.HandleFunc("/api/stock/add", s.handleAddStock)
	mux.HandleFunc("/api/stock/remove", s.handleRemoveStock)
	mux.HandleFunc("/api/stock/adjust", s.handleAdjustStock)

	// Search and filter endpoints
	mux.HandleFunc("/api/search", s.handleSearch)
	mux.HandleFunc("/api/categories", s.handleCategories)
	mux.HandleFunc("/api/categories/", s.handleItemsByCategory)

	// Statistics endpoints
	mux.HandleFunc("/api/stats", s.handleStats)
	mux.HandleFunc("/api/alerts", s.handleAlerts)
	mux.HandleFunc("/api/low-stock", s.handleLowStock)
	mux.HandleFunc("/api/out-of-stock", s.handleOutOfStock)

	// Storage endpoints
	mux.HandleFunc("/api/backup", s.handleBackup)
	mux.HandleFunc("/api/restore", s.handleRestore)
	mux.HandleFunc("/api/export", s.handleExport)
	mux.HandleFunc("/api/storage/stats", s.handleStorageStats)

	// Health check
	mux.HandleFunc("/health", s.handleHealth)

	addr := fmt.Sprintf(":%d", s.port)
	log.Printf("API server starting on port %d", s.port)
	return http.ListenAndServe(addr, mux)
}

// handleItems handles GET (list) and POST (create) requests for items
func (s *APIServer) handleItems(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.listItems(w, r)
	case http.MethodPost:
		s.createItem(w, r)
	default:
		s.sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// handleItemByID handles GET, PUT, DELETE requests for a specific item
func (s *APIServer) handleItemByID(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		s.sendError(w, http.StatusBadRequest, "Invalid item ID")
		return
	}
	itemID := parts[3]

	switch r.Method {
	case http.MethodGet:
		s.getItem(w, r, itemID)
	case http.MethodPut:
		s.updateItem(w, r, itemID)
	case http.MethodDelete:
		s.deleteItem(w, r, itemID)
	default:
		s.sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// handleItemBySKU handles GET requests for items by SKU
func (s *APIServer) handleItemBySKU(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		s.sendError(w, http.StatusBadRequest, "Invalid SKU")
		return
	}
	sku := parts[4]

	item, err := s.service.manager.GetItemBySKU(sku)
	if err != nil {
		s.sendError(w, http.StatusNotFound, err.Error())
		return
	}

	s.sendJSON(w, http.StatusOK, APIResponse{Success: true, Data: item})
}

// listItems returns a list of all items
func (s *APIServer) listItems(w http.ResponseWriter, r *http.Request) {
	items := s.service.manager.GetActiveItems()
	s.sendJSON(w, http.StatusOK, APIResponse{Success: true, Data: items})
}

// createItem creates a new item
func (s *APIServer) createItem(w http.ResponseWriter, r *http.Request) {
	var itemData struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Category    string  `json:"category"`
		SKU         string  `json:"sku"`
		Supplier    string  `json:"supplier"`
		Price       float64 `json:"price"`
		Quantity    int     `json:"quantity"`
		MinStock    int     `json:"min_stock"`
		MaxStock    int     `json:"max_stock"`
	}

	if err := json.NewDecoder(r.Body).Decode(&itemData); err != nil {
		s.sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	item, err := s.service.CreateItem(
		itemData.Name,
		itemData.Description,
		itemData.Category,
		itemData.SKU,
		itemData.Supplier,
		itemData.Price,
		itemData.Quantity,
		itemData.MinStock,
		itemData.MaxStock,
	)

	if err != nil {
		s.sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	s.sendJSON(w, http.StatusCreated, APIResponse{Success: true, Message: "Item created", Data: item})
}

// getItem retrieves an item by ID
func (s *APIServer) getItem(w http.ResponseWriter, r *http.Request, itemID string) {
	item, err := s.service.manager.GetItem(itemID)
	if err != nil {
		s.sendError(w, http.StatusNotFound, err.Error())
		return
	}
	s.sendJSON(w, http.StatusOK, APIResponse{Success: true, Data: item})
}

// updateItem updates an existing item
func (s *APIServer) updateItem(w http.ResponseWriter, r *http.Request, itemID string) {
	var updates Item
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		s.sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := s.service.UpdateItem(itemID, &updates); err != nil {
		s.sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	s.sendJSON(w, http.StatusOK, APIResponse{Success: true, Message: "Item updated"})
}

// deleteItem deletes an item
func (s *APIServer) deleteItem(w http.ResponseWriter, r *http.Request, itemID string) {
	if err := s.service.DeleteItem(itemID); err != nil {
		s.sendError(w, http.StatusNotFound, err.Error())
		return
	}
	s.sendJSON(w, http.StatusOK, APIResponse{Success: true, Message: "Item deleted"})
}

// handleAddStock handles stock addition requests
func (s *APIServer) handleAddStock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		ItemID   string `json:"item_id"`
		Quantity int    `json:"quantity"`
		Reason   string `json:"reason"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := s.service.ProcessStockIn(req.ItemID, req.Quantity, req.Reason); err != nil {
		s.sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	s.sendJSON(w, http.StatusOK, APIResponse{Success: true, Message: "Stock added"})
}

// handleRemoveStock handles stock removal requests
func (s *APIServer) handleRemoveStock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		ItemID   string `json:"item_id"`
		Quantity int    `json:"quantity"`
		Reason   string `json:"reason"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := s.service.ProcessStockOut(req.ItemID, req.Quantity, req.Reason); err != nil {
		s.sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	s.sendJSON(w, http.StatusOK, APIResponse{Success: true, Message: "Stock removed"})
}

// handleAdjustStock handles stock adjustment requests
func (s *APIServer) handleAdjustStock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		ItemID   string `json:"item_id"`
		Quantity int    `json:"quantity"`
		Reason   string `json:"reason"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := s.service.AdjustStock(req.ItemID, req.Quantity, req.Reason); err != nil {
		s.sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	s.sendJSON(w, http.StatusOK, APIResponse{Success: true, Message: "Stock adjusted"})
}

// handleSearch handles search requests with filters and pagination
func (s *APIServer) handleSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	query := r.URL.Query().Get("q")
	category := r.URL.Query().Get("category")
	supplier := r.URL.Query().Get("supplier")
	sortField := r.URL.Query().Get("sort")
	ascending := r.URL.Query().Get("order") != "desc"

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	var minPrice, maxPrice float64
	if minPriceStr := r.URL.Query().Get("min_price"); minPriceStr != "" {
		minPrice, _ = strconv.ParseFloat(minPriceStr, 64)
	}
	if maxPriceStr := r.URL.Query().Get("max_price"); maxPriceStr != "" {
		maxPrice, _ = strconv.ParseFloat(maxPriceStr, 64)
	}

	var minQty, maxQty int
	if minQtyStr := r.URL.Query().Get("min_quantity"); minQtyStr != "" {
		minQty, _ = strconv.Atoi(minQtyStr)
	}
	if maxQtyStr := r.URL.Query().Get("max_quantity"); maxQtyStr != "" {
		maxQty, _ = strconv.Atoi(maxQtyStr)
	}

	lowStock := r.URL.Query().Get("low_stock") == "true"

	filters := SearchFilters{
		Category:    category,
		Supplier:    supplier,
		MinPrice:    minPrice,
		MaxPrice:    maxPrice,
		MinQuantity: minQty,
		MaxQuantity: maxQty,
		LowStock:    lowStock,
	}

	result, err := s.service.SearchItems(query, filters, sortField, ascending, page, pageSize)
	if err != nil {
		s.sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	s.sendJSON(w, http.StatusOK, APIResponse{Success: true, Data: result})
}

// handleCategories returns all categories
func (s *APIServer) handleCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	categories := s.service.GetCategories()
	s.sendJSON(w, http.StatusOK, APIResponse{Success: true, Data: categories})
}

// handleItemsByCategory returns items in a specific category
func (s *APIServer) handleItemsByCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		s.sendError(w, http.StatusBadRequest, "Invalid category")
		return
	}
	category := parts[3]

	items := s.service.GetItemsByCategory(category)
	s.sendJSON(w, http.StatusOK, APIResponse{Success: true, Data: items})
}

// handleStats returns inventory statistics
func (s *APIServer) handleStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	stats := s.service.GetInventoryStats()
	s.sendJSON(w, http.StatusOK, APIResponse{Success: true, Data: stats})
}

// handleAlerts returns stock alerts
func (s *APIServer) handleAlerts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	alerts := s.service.GetStockAlerts()
	s.sendJSON(w, http.StatusOK, APIResponse{Success: true, Data: alerts})
}

// handleLowStock returns low stock items
func (s *APIServer) handleLowStock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	items := s.service.GetLowStockItems()
	s.sendJSON(w, http.StatusOK, APIResponse{Success: true, Data: items})
}

// handleOutOfStock returns out of stock items
func (s *APIServer) handleOutOfStock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	items := s.service.GetOutOfStockItems()
	s.sendJSON(w, http.StatusOK, APIResponse{Success: true, Data: items})
}

// handleBackup creates a backup
func (s *APIServer) handleBackup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	if err := s.service.Backup(); err != nil {
		s.sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	s.sendJSON(w, http.StatusOK, APIResponse{Success: true, Message: "Backup created"})
}

// handleRestore restores from backup
func (s *APIServer) handleRestore(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	if err := s.service.Restore(); err != nil {
		s.sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	s.sendJSON(w, http.StatusOK, APIResponse{Success: true, Message: "Inventory restored"})
}

// handleExport exports inventory to a file
func (s *APIServer) handleExport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		FilePath string `json:"file_path"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := s.service.Export(req.FilePath); err != nil {
		s.sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	s.sendJSON(w, http.StatusOK, APIResponse{Success: true, Message: "Inventory exported"})
}

// handleStorageStats returns storage statistics
func (s *APIServer) handleStorageStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	stats := s.service.GetStorageStats()
	s.sendJSON(w, http.StatusOK, APIResponse{Success: true, Data: stats})
}

// handleHealth returns health status
func (s *APIServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	s.sendJSON(w, http.StatusOK, map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
	})
}

// sendJSON sends a JSON response
func (s *APIServer) sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// sendError sends an error response
func (s *APIServer) sendError(w http.ResponseWriter, status int, message string) {
	s.sendJSON(w, status, APIResponse{
		Success: false,
		Error:   message,
	})
}

