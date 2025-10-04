package middleware

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
)

// ErrorResponse represents a structured error response
type ErrorResponse struct {
	Success   bool   `json:"success"`
	Error     string `json:"error"`
	ErrorCode string `json:"error_code,omitempty"`
	RequestID string `json:"request_id,omitempty"`
}

// RecoveryMiddleware handles panics and converts them to HTTP 500 errors
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic with stack trace
				log.Printf("PANIC RECOVERED: %v\nStack trace:\n%s", err, debug.Stack())

				// Return 500 error to client
				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Success:   false,
					Error:     "Internal server error",
					ErrorCode: "INTERNAL_ERROR",
					RequestID: c.GetString("request_id"),
				})

				// Abort further middleware execution
				c.Abort()
			}
		}()

		c.Next()
	}
}

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate request ID (simple timestamp-based for now)
		requestID := generateRequestID()
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

// generateRequestID generates a unique request ID
func generateRequestID() string {
	// Simple implementation - in production, use UUID
	return fmt.Sprintf("req_%d", time.Now().UnixNano())
}

// ErrorHandlerMiddleware provides centralized error handling
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there were any errors during request processing
		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			// Determine status code
			statusCode := c.Writer.Status()
			if statusCode == http.StatusOK {
				statusCode = http.StatusInternalServerError
			}

			// Log the error
			log.Printf("Request error [%s %s]: %v", c.Request.Method, c.Request.URL.Path, err.Err)

			// Return error response
			c.JSON(statusCode, ErrorResponse{
				Success:   false,
				Error:     err.Error(),
				ErrorCode: "REQUEST_ERROR",
				RequestID: c.GetString("request_id"),
			})
		}
	}
}

// CORSErrorMiddleware handles CORS preflight errors
func CORSErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Handle CORS errors
		if c.Request.Method == "OPTIONS" && c.Writer.Status() >= 400 {
			c.Status(http.StatusNoContent)
		}
	}
}

// NotFoundHandler handles 404 errors
func NotFoundHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Success:   false,
			Error:     "Endpoint not found",
			ErrorCode: "NOT_FOUND",
			RequestID: c.GetString("request_id"),
		})
	}
}

// MethodNotAllowedHandler handles 405 errors
func MethodNotAllowedHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, ErrorResponse{
			Success:   false,
			Error:     "Method not allowed",
			ErrorCode: "METHOD_NOT_ALLOWED",
			RequestID: c.GetString("request_id"),
		})
	}
}
