package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/voting-blockchain/internal/auth"
	"github.com/voting-blockchain/internal/models"
)

// LoginRequest represents admin login request
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents login response with JWT token
type LoginResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token,omitempty"`
	Message string `json:"message,omitempty"`
	User    *auth.Admin `json:"user,omitempty"`
}

// Login handles admin login
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request format",
		})
		return
	}

	// Validate credentials
	admin, err := h.adminStore.ValidateCredentials(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "Invalid username or password",
		})
		return
	}

	// Generate JWT token
	token, err := h.jwtManager.GenerateToken(admin.ID, admin.Email, admin.Role, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to generate token",
		})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Success: true,
		Token:   token,
		Message: "Login successful",
		User:    admin,
	})
}

// VoterLogin handles voter login/authentication
func (h *Handler) VoterLogin(c *gin.Context) {
	var req struct {
		VoterID   string `json:"voter_id" binding:"required"`
		PrivateKey string `json:"private_key" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request format",
		})
		return
	}

	// Verify voter exists
	voter, exists := h.blockchain.VoterRegistry[req.VoterID]
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "Invalid voter credentials",
		})
		return
	}

	// Verify private key ownership
	err := h.crypto.VerifyPrivateKeyOwnership(
		[]byte(req.PrivateKey),
		[]byte(voter.PublicKey),
		voter.VoterID,
	)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "Invalid private key",
		})
		return
	}

	token, err := h.jwtManager.GenerateToken(voter.VoterID, voter.Email, "voter", voter.VoterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to generate token",
		})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Success: true,
		Token:   token,
		Message: "Login successful",
	})
}

// RefreshToken refreshes an existing JWT token
func (h *Handler) RefreshToken(c *gin.Context) {
	var req struct {
		Token string `json:"token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request format",
		})
		return
	}

	newToken, err := h.jwtManager.RefreshToken(req.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "Invalid or expired token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"token":   newToken,
		"message": "Token refreshed successfully",
	})
}

// GetCurrentUser returns the currently authenticated user
func (h *Handler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "Not authenticated",
		})
		return
	}

	role, _ := c.Get("role")
	email, _ := c.Get("email")
	voterID, _ := c.Get("voter_id")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user": gin.H{
			"user_id":  userID,
			"email":    email,
			"role":     role,
			"voter_id": voterID,
		},
	})
}
