package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/voting-blockchain/internal/auth"
)

// AuthMiddleware creates a middleware for JWT authentication
func AuthMiddleware(jwtManager *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Authorization header required",
			})
			c.Abort()
			return
		}

		// Check if it's a Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Invalid authorization header format. Use: Bearer <token>",
			})
			c.Abort()
			return
		}

		// Validate the token
		token := parts[1]
		claims, err := jwtManager.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Store claims in context for route handlers
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Set("voter_id", claims.VoterID)
		c.Set("claims", claims)

		c.Next()
	}
}

// RequireRole creates a middleware that requires a specific role
// Note: This middleware expects AuthMiddleware to have already run
func RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check role (AuthMiddleware should have already set this)
		role, exists := c.Get("role")
		if !exists || role != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "Insufficient permissions",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// OptionalAuth middleware that allows but doesn't require authentication
func OptionalAuth(jwtManager *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// No auth header, continue without authentication
			c.Next()
			return
		}

		// Try to validate token if present
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			claims, err := jwtManager.ValidateToken(parts[1])
			if err == nil {
				// Valid token, store claims
				c.Set("user_id", claims.UserID)
				c.Set("email", claims.Email)
				c.Set("role", claims.Role)
				c.Set("voter_id", claims.VoterID)
				c.Set("claims", claims)
			}
		}

		c.Next()
	}
}
