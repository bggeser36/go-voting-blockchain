package validation

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

// Validator provides input validation utilities
type Validator struct{}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	return &Validator{}
}

// Email validation regex
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// ValidateEmail validates email format
func (v *Validator) ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email is required")
	}

	email = strings.TrimSpace(email)
	if len(email) > 254 {
		return fmt.Errorf("email is too long (max 254 characters)")
	}

	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}

	return nil
}

// ValidateName validates user/voter names
func (v *Validator) ValidateName(name string) error {
	if name == "" {
		return fmt.Errorf("name is required")
	}

	name = strings.TrimSpace(name)
	if len(name) < 2 {
		return fmt.Errorf("name must be at least 2 characters")
	}

	if len(name) > 100 {
		return fmt.Errorf("name is too long (max 100 characters)")
	}

	// Check for valid UTF-8
	if !utf8.ValidString(name) {
		return fmt.Errorf("name contains invalid characters")
	}

	return nil
}

// ValidatePassword validates password strength
func (v *Validator) ValidatePassword(password string) error {
	if password == "" {
		return fmt.Errorf("password is required")
	}

	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}

	if len(password) > 128 {
		return fmt.Errorf("password is too long (max 128 characters)")
	}

	// Check for at least one letter and one number
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)

	if !hasLetter || !hasNumber {
		return fmt.Errorf("password must contain at least one letter and one number")
	}

	return nil
}

// ValidateUsername validates username format
func (v *Validator) ValidateUsername(username string) error {
	if username == "" {
		return fmt.Errorf("username is required")
	}

	username = strings.TrimSpace(username)
	if len(username) < 3 {
		return fmt.Errorf("username must be at least 3 characters")
	}

	if len(username) > 50 {
		return fmt.Errorf("username is too long (max 50 characters)")
	}

	// Only allow alphanumeric, underscore, and hyphen
	validUsername := regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`).MatchString(username)
	if !validUsername {
		return fmt.Errorf("username can only contain letters, numbers, underscores, and hyphens")
	}

	return nil
}

// ValidatePollTitle validates poll title
func (v *Validator) ValidatePollTitle(title string) error {
	if title == "" {
		return fmt.Errorf("poll title is required")
	}

	title = strings.TrimSpace(title)
	if len(title) < 5 {
		return fmt.Errorf("poll title must be at least 5 characters")
	}

	if len(title) > 200 {
		return fmt.Errorf("poll title is too long (max 200 characters)")
	}

	if !utf8.ValidString(title) {
		return fmt.Errorf("poll title contains invalid characters")
	}

	return nil
}

// ValidatePollDescription validates poll description
func (v *Validator) ValidatePollDescription(description string) error {
	if description == "" {
		return fmt.Errorf("poll description is required")
	}

	description = strings.TrimSpace(description)
	if len(description) < 10 {
		return fmt.Errorf("poll description must be at least 10 characters")
	}

	if len(description) > 1000 {
		return fmt.Errorf("poll description is too long (max 1000 characters)")
	}

	if !utf8.ValidString(description) {
		return fmt.Errorf("poll description contains invalid characters")
	}

	return nil
}

// ValidatePollOptions validates poll voting options
func (v *Validator) ValidatePollOptions(options []string) error {
	if len(options) < 2 {
		return fmt.Errorf("poll must have at least 2 options")
	}

	if len(options) > 20 {
		return fmt.Errorf("poll cannot have more than 20 options")
	}

	seen := make(map[string]bool)
	for i, option := range options {
		option = strings.TrimSpace(option)
		if option == "" {
			return fmt.Errorf("option %d is empty", i+1)
		}

		if len(option) > 100 {
			return fmt.Errorf("option %d is too long (max 100 characters)", i+1)
		}

		if !utf8.ValidString(option) {
			return fmt.Errorf("option %d contains invalid characters", i+1)
		}

		// Check for duplicates
		if seen[strings.ToLower(option)] {
			return fmt.Errorf("duplicate option: %s", option)
		}
		seen[strings.ToLower(option)] = true
	}

	return nil
}

// ValidateDuration validates poll duration in hours
func (v *Validator) ValidateDuration(hours int) error {
	if hours < 1 {
		return fmt.Errorf("poll duration must be at least 1 hour")
	}

	if hours > 8760 { // 1 year
		return fmt.Errorf("poll duration cannot exceed 1 year (8760 hours)")
	}

	return nil
}

// ValidateVoterID validates voter ID format
func (v *Validator) ValidateVoterID(voterID string) error {
	if voterID == "" {
		return fmt.Errorf("voter ID is required")
	}

	// Voter IDs are hex strings of 16 characters
	validID := regexp.MustCompile(`^[a-f0-9]{16}$`).MatchString(voterID)
	if !validID {
		return fmt.Errorf("invalid voter ID format")
	}

	return nil
}

// ValidatePollID validates poll ID format (UUID)
func (v *Validator) ValidatePollID(pollID string) error {
	if pollID == "" {
		return fmt.Errorf("poll ID is required")
	}

	// UUIDs are in format: 8-4-4-4-12
	validUUID := regexp.MustCompile(`^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$`).MatchString(pollID)
	if !validUUID {
		return fmt.Errorf("invalid poll ID format")
	}

	return nil
}

// SanitizeString removes potentially dangerous characters and trims whitespace
func (v *Validator) SanitizeString(input string) string {
	// Trim whitespace
	input = strings.TrimSpace(input)

	// Remove null bytes
	input = strings.ReplaceAll(input, "\x00", "")

	// Remove control characters except newline and tab
	var result strings.Builder
	for _, r := range input {
		if r == '\n' || r == '\t' || (r >= 32 && r != 127) {
			result.WriteRune(r)
		}
	}

	return result.String()
}

// ValidateDepartment validates department name
func (v *Validator) ValidateDepartment(department string) error {
	if department == "" {
		return fmt.Errorf("department is required")
	}

	department = strings.TrimSpace(department)
	if len(department) < 2 {
		return fmt.Errorf("department name must be at least 2 characters")
	}

	if len(department) > 100 {
		return fmt.Errorf("department name is too long (max 100 characters)")
	}

	if !utf8.ValidString(department) {
		return fmt.Errorf("department name contains invalid characters")
	}

	return nil
}
