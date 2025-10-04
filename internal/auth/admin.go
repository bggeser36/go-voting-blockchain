package auth

import (
	"crypto/subtle"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Admin represents an admin user
type Admin struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Never expose password hash in JSON
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}

// AdminStore manages admin users (in-memory for now)
type AdminStore struct {
	admins map[string]*Admin
}

// NewAdminStore creates a new admin store
func NewAdminStore() *AdminStore {
	return &AdminStore{
		admins: make(map[string]*Admin),
	}
}

// CreateAdmin creates a new admin user
func (s *AdminStore) CreateAdmin(username, email, password string) (*Admin, error) {
	// Check if admin already exists
	if _, exists := s.admins[username]; exists {
		return nil, fmt.Errorf("admin with username %s already exists", username)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	admin := &Admin{
		ID:           GenerateID(),
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
		Role:         "admin",
		CreatedAt:    time.Now(),
	}

	s.admins[username] = admin
	return admin, nil
}

// ValidateCredentials validates admin credentials
func (s *AdminStore) ValidateCredentials(username, password string) (*Admin, error) {
	admin, exists := s.admins[username]
	if !exists {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Use constant-time comparison to prevent timing attacks
	err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	return admin, nil
}

// GetAdmin retrieves an admin by username
func (s *AdminStore) GetAdmin(username string) (*Admin, error) {
	admin, exists := s.admins[username]
	if !exists {
		return nil, fmt.Errorf("admin not found")
	}
	return admin, nil
}

// GenerateID generates a unique ID (simple implementation)
func GenerateID() string {
	return fmt.Sprintf("admin_%d", time.Now().UnixNano())
}

// SecureCompare performs a constant-time comparison of two strings
func SecureCompare(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}
