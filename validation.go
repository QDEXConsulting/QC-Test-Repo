package main

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidationResult contains validation results
type ValidationResult struct {
	IsValid bool
	Errors  []ValidationError
}

// NewValidationResult creates a new validation result
func NewValidationResult() *ValidationResult {
	return &ValidationResult{
		IsValid: true,
		Errors:  make([]ValidationError, 0),
	}
}

// AddError adds a validation error
func (vr *ValidationResult) AddError(field, message string) {
	vr.IsValid = false
	vr.Errors = append(vr.Errors, ValidationError{
		Field:   field,
		Message: message,
	})
}

// ValidateItem validates an item before creation or update
func ValidateItem(item *Item) *ValidationResult {
	result := NewValidationResult()
	
	// Validate name
	if strings.TrimSpace(item.Name) == "" {
		result.AddError("name", "Item name is required")
	} else if len(item.Name) > 200 {
		result.AddError("name", "Item name must be 200 characters or less")
	}
	
	// Validate description
	if len(item.Description) > 1000 {
		result.AddError("description", "Description must be 1000 characters or less")
	}
	
	// Validate category
	if strings.TrimSpace(item.Category) == "" {
		result.AddError("category", "Category is required")
	} else if len(item.Category) > 100 {
		result.AddError("category", "Category must be 100 characters or less")
	}
	
	// Validate SKU
	if !ValidateSKU(item.SKU) {
		result.AddError("sku", "SKU must be between 3 and 50 characters")
	}
	
	// Validate price
	if item.Price < 0 {
		result.AddError("price", "Price cannot be negative")
	} else if item.Price > 1000000 {
		result.AddError("price", "Price cannot exceed 1,000,000")
	}
	
	// Validate quantity
	if item.Quantity < 0 {
		result.AddError("quantity", "Quantity cannot be negative")
	} else if item.Quantity > 1000000 {
		result.AddError("quantity", "Quantity cannot exceed 1,000,000")
	}
	
	// Validate min stock
	if item.MinStock < 0 {
		result.AddError("min_stock", "Minimum stock cannot be negative")
	}
	
	// Validate max stock
	if item.MaxStock > 0 && item.MaxStock < item.MinStock {
		result.AddError("max_stock", "Maximum stock must be greater than minimum stock")
	}
	
	// Validate supplier
	if len(item.Supplier) > 200 {
		result.AddError("supplier", "Supplier name must be 200 characters or less")
	}
	
	return result
}

// ValidateSKU validates SKU format
func ValidateSKU(sku string) bool {
	if len(sku) < 3 || len(sku) > 50 {
		return false
	}
	
	// SKU should contain only alphanumeric characters, hyphens, and underscores
	matched, _ := regexp.MatchString("^[a-zA-Z0-9_-]+$", sku)
	return matched
}

// ValidateSearchQuery validates search query input
func ValidateSearchQuery(query string) *ValidationResult {
	result := NewValidationResult()
	
	if len(query) > 500 {
		result.AddError("query", "Search query must be 500 characters or less")
	}
	
	// Check for potentially dangerous patterns
	if containsSQLInjectionPattern(query) {
		result.AddError("query", "Invalid characters in search query")
	}
	
	return result
}

// ValidatePaginationParams validates pagination parameters
func ValidatePaginationParams(page, pageSize int) *ValidationResult {
	result := NewValidationResult()
	
	if page < 1 {
		result.AddError("page", "Page must be greater than 0")
	}
	
	if pageSize < 1 {
		result.AddError("page_size", "Page size must be greater than 0")
	} else if pageSize > 1000 {
		result.AddError("page_size", "Page size cannot exceed 1000")
	}
	
	return result
}

// ValidateStockOperation validates stock operation parameters
func ValidateStockOperation(itemID string, quantity int, reason string) *ValidationResult {
	result := NewValidationResult()
	
	if strings.TrimSpace(itemID) == "" {
		result.AddError("item_id", "Item ID is required")
	}
	
	if quantity <= 0 {
		result.AddError("quantity", "Quantity must be greater than 0")
	} else if quantity > 100000 {
		result.AddError("quantity", "Quantity cannot exceed 100,000")
	}
	
	if len(reason) > 500 {
		result.AddError("reason", "Reason must be 500 characters or less")
	}
	
	return result
}

// ValidatePriceRange validates price range filters
func ValidatePriceRange(minPrice, maxPrice float64) *ValidationResult {
	result := NewValidationResult()
	
	if minPrice < 0 {
		result.AddError("min_price", "Minimum price cannot be negative")
	}
	
	if maxPrice < 0 {
		result.AddError("max_price", "Maximum price cannot be negative")
	}
	
	if minPrice > 0 && maxPrice > 0 && minPrice > maxPrice {
		result.AddError("price_range", "Minimum price cannot exceed maximum price")
	}
	
	return result
}

// ValidateQuantityRange validates quantity range filters
func ValidateQuantityRange(minQty, maxQty int) *ValidationResult {
	result := NewValidationResult()
	
	if minQty < 0 {
		result.AddError("min_quantity", "Minimum quantity cannot be negative")
	}
	
	if maxQty < 0 {
		result.AddError("max_quantity", "Maximum quantity cannot be negative")
	}
	
	if minQty > 0 && maxQty > 0 && minQty > maxQty {
		result.AddError("quantity_range", "Minimum quantity cannot exceed maximum quantity")
	}
	
	return result
}

// SanitizeString removes potentially dangerous characters from a string
func SanitizeString(input string) string {
	// Remove null bytes
	input = strings.ReplaceAll(input, "\x00", "")
	
	// Trim whitespace
	input = strings.TrimSpace(input)
	
	// Remove control characters except newlines and tabs
	var builder strings.Builder
	for _, r := range input {
		if unicode.IsControl(r) && r != '\n' && r != '\t' {
			continue
		}
		builder.WriteRune(r)
	}
	
	return builder.String()
}

// ValidateEmail validates email format (if needed for supplier contacts)
func ValidateEmail(email string) bool {
	if email == "" {
		return false
	}
	
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// ValidateURL validates URL format
func ValidateURL(url string) bool {
	if url == "" {
		return false
	}
	
	urlRegex := regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)
	return urlRegex.MatchString(url)
}

// ValidateCategoryName validates category name format
func ValidateCategoryName(category string) bool {
	if strings.TrimSpace(category) == "" {
		return false
	}
	
	if len(category) > 100 {
		return false
	}
	
	// Category should not contain special characters except spaces, hyphens, and underscores
	matched, _ := regexp.MatchString("^[a-zA-Z0-9\\s_-]+$", category)
	return matched
}

// ValidateSupplierName validates supplier name
func ValidateSupplierName(supplier string) bool {
	if strings.TrimSpace(supplier) == "" {
		return false
	}
	
	if len(supplier) > 200 {
		return false
	}
	
	return true
}

// containsSQLInjectionPattern checks for common SQL injection patterns
func containsSQLInjectionPattern(input string) bool {
	patterns := []string{
		"' OR '1'='1",
		"'; DROP TABLE",
		"UNION SELECT",
		"1=1",
		"1' OR '1'='1",
	}
	
	inputLower := strings.ToLower(input)
	for _, pattern := range patterns {
		if strings.Contains(inputLower, strings.ToLower(pattern)) {
			return true
		}
	}
	
	return false
}

// ValidateBulkOperation validates bulk operation parameters
func ValidateBulkOperation(itemIDs []string, maxItems int) *ValidationResult {
	result := NewValidationResult()
	
	if len(itemIDs) == 0 {
		result.AddError("item_ids", "At least one item ID is required")
	}
	
	if len(itemIDs) > maxItems {
		result.AddError("item_ids", fmt.Sprintf("Cannot process more than %d items at once", maxItems))
	}
	
	for i, itemID := range itemIDs {
		if strings.TrimSpace(itemID) == "" {
			result.AddError(fmt.Sprintf("item_ids[%d]", i), "Item ID cannot be empty")
		}
	}
	
	return result
}

// ValidateExportPath validates export file path
func ValidateExportPath(path string) *ValidationResult {
	result := NewValidationResult()
	
	if strings.TrimSpace(path) == "" {
		result.AddError("path", "Export path is required")
	}
	
	if len(path) > 500 {
		result.AddError("path", "Path must be 500 characters or less")
	}
	
	// Check for directory traversal attempts
	if strings.Contains(path, "..") {
		result.AddError("path", "Invalid path")
	}
	
	return result
}

