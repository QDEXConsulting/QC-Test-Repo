package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

// AuthError represents an authentication error
type AuthError struct {
	Message string
}

func (e *AuthError) Error() string {
	return e.Message
}

// User represents an authenticated user
type User struct {
	ID        string
	Username  string
	Role      string
	APIKey    string
	CreatedAt time.Time
	LastSeen  time.Time
}

// Role represents user roles
const (
	RoleAdmin    = "admin"
	RoleManager  = "manager"
	RoleViewer   = "viewer"
	RoleOperator = "operator"
)

// AuthManager manages authentication and authorization
type AuthManager struct {
	users      map[string]*User
	apiKeys    map[string]*User
	mu         sync.RWMutex
	jwtSecret  string
}

// NewAuthManager creates a new authentication manager
func NewAuthManager(jwtSecret string) *AuthManager {
	return &AuthManager{
		users:     make(map[string]*User),
		apiKeys:   make(map[string]*User),
		jwtSecret: jwtSecret,
	}
}

// CreateUser creates a new user
func (am *AuthManager) CreateUser(username, role string) (*User, error) {
	am.mu.Lock()
	defer am.mu.Unlock()
	
	// Check if username already exists
	for _, user := range am.users {
		if user.Username == username {
			return nil, &AuthError{Message: "username already exists"}
		}
	}
	
	// Validate role
	if !am.isValidRole(role) {
		return nil, &AuthError{Message: "invalid role"}
	}
	
	user := &User{
		ID:        generateID(),
		Username:  username,
		Role:      role,
		APIKey:    am.generateAPIKey(username),
		CreatedAt: time.Now(),
		LastSeen:  time.Now(),
	}
	
	am.users[user.ID] = user
	am.apiKeys[user.APIKey] = user
	
	return user, nil
}

// AuthenticateAPIKey authenticates a user by API key
func (am *AuthManager) AuthenticateAPIKey(apiKey string) (*User, error) {
	am.mu.RLock()
	defer am.mu.RUnlock()
	
	user, exists := am.apiKeys[apiKey]
	if !exists {
		return nil, &AuthError{Message: "invalid API key"}
	}
	
	user.LastSeen = time.Now()
	return user, nil
}

// GetUser retrieves a user by ID
func (am *AuthManager) GetUser(userID string) (*User, error) {
	am.mu.RLock()
	defer am.mu.RUnlock()
	
	user, exists := am.users[userID]
	if !exists {
		return nil, &AuthError{Message: "user not found"}
	}
	
	return user, nil
}

// DeleteUser deletes a user
func (am *AuthManager) DeleteUser(userID string) error {
	am.mu.Lock()
	defer am.mu.Unlock()
	
	user, exists := am.users[userID]
	if !exists {
		return &AuthError{Message: "user not found"}
	}
	
	delete(am.users, userID)
	delete(am.apiKeys, user.APIKey)
	
	return nil
}

// RegenerateAPIKey generates a new API key for a user
func (am *AuthManager) RegenerateAPIKey(userID string) (string, error) {
	am.mu.Lock()
	defer am.mu.Unlock()
	
	user, exists := am.users[userID]
	if !exists {
		return "", &AuthError{Message: "user not found"}
	}
	
	// Remove old API key
	delete(am.apiKeys, user.APIKey)
	
	// Generate new API key
	user.APIKey = am.generateAPIKey(user.Username)
	am.apiKeys[user.APIKey] = user
	
	return user.APIKey, nil
}

// HasPermission checks if a user has permission for an action
func (am *AuthManager) HasPermission(user *User, action string) bool {
	switch user.Role {
	case RoleAdmin:
		return true
	case RoleManager:
		return action != "delete_user" && action != "create_user"
	case RoleOperator:
		return action == "read" || action == "update" || action == "stock_operation"
	case RoleViewer:
		return action == "read"
	default:
		return false
	}
}

// IsAuthorized checks if a user is authorized for an action
func (am *AuthManager) IsAuthorized(user *User, action string) error {
	if !am.HasPermission(user, action) {
		return &AuthError{Message: "insufficient permissions"}
	}
	return nil
}

// generateAPIKey generates a new API key
func (am *AuthManager) generateAPIKey(username string) string {
	timestamp := time.Now().Unix()
	data := fmt.Sprintf("%s:%d:%s", username, timestamp, am.jwtSecret)
	
	h := hmac.New(sha256.New, []byte(am.jwtSecret))
	h.Write([]byte(data))
	signature := base64.URLEncoding.EncodeToString(h.Sum(nil))
	
	return fmt.Sprintf("sk_%s_%s", base64.URLEncoding.EncodeToString([]byte(username)), signature[:16])
}

// isValidRole checks if a role is valid
func (am *AuthManager) isValidRole(role string) bool {
	validRoles := []string{RoleAdmin, RoleManager, RoleOperator, RoleViewer}
	for _, r := range validRoles {
		if r == role {
			return true
		}
	}
	return false
}

// GetAllUsers returns all users
func (am *AuthManager) GetAllUsers() []*User {
	am.mu.RLock()
	defer am.mu.RUnlock()
	
	users := make([]*User, 0, len(am.users))
	for _, user := range am.users {
		users = append(users, user)
	}
	
	return users
}

// ValidateAPIKey validates an API key format
func ValidateAPIKey(apiKey string) bool {
	if !strings.HasPrefix(apiKey, "sk_") {
		return false
	}
	
	parts := strings.Split(apiKey, "_")
	if len(parts) != 3 {
		return false
	}
	
	return len(apiKey) >= 20 && len(apiKey) <= 100
}

// ExtractAPIKeyFromHeader extracts API key from Authorization header
func ExtractAPIKeyFromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("authorization header is missing")
	}
	
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 {
		return "", errors.New("invalid authorization header format")
	}
	
	scheme := parts[0]
	apiKey := parts[1]
	
	if scheme != "Bearer" && scheme != "ApiKey" {
		return "", errors.New("unsupported authorization scheme")
	}
	
	return apiKey, nil
}

// Permission constants
const (
	PermissionRead          = "read"
	PermissionWrite         = "write"
	PermissionDelete        = "delete"
	PermissionStockIn       = "stock_in"
	PermissionStockOut      = "stock_out"
	PermissionStockAdjust   = "stock_adjust"
	PermissionCreateUser    = "create_user"
	PermissionDeleteUser    = "delete_user"
	PermissionManageConfig  = "manage_config"
	PermissionBackup        = "backup"
	PermissionRestore       = "restore"
)

// GetRequiredPermission returns the required permission for an action
func GetRequiredPermission(action string) string {
	switch action {
	case "GET":
		return PermissionRead
	case "POST", "PUT":
		return PermissionWrite
	case "DELETE":
		return PermissionDelete
	default:
		return PermissionRead
	}
}

// Global auth manager instance
var globalAuthManager *AuthManager

// InitGlobalAuthManager initializes the global auth manager
func InitGlobalAuthManager(jwtSecret string) {
	globalAuthManager = NewAuthManager(jwtSecret)
}

// GetAuthManager returns the global auth manager
func GetAuthManager() *AuthManager {
	if globalAuthManager == nil {
		InitGlobalAuthManager("default-secret-change-in-production")
	}
	return globalAuthManager
}

