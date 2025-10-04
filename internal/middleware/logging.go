package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// DetailedLogger provides detailed request/response logging
func DetailedLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		startTime := time.Now()

		// Get request details
		requestID := c.GetString("request_id")
		method := c.Request.Method
		path := c.Request.URL.Path
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(startTime)
		statusCode := c.Writer.Status()
		responseSize := c.Writer.Size()

		// Get user info if authenticated
		userID, _ := c.Get("user_id")
		role, _ := c.Get("role")

		// Log request details
		log.Printf("[%s] %s %s | Status: %d | Latency: %v | Size: %d bytes | IP: %s | UserID: %v | Role: %v | UA: %s",
			requestID,
			method,
			path,
			statusCode,
			latency,
			responseSize,
			clientIP,
			userID,
			role,
			userAgent,
		)

		// Log security events
		if statusCode == 401 || statusCode == 403 {
			log.Printf("[SECURITY] Unauthorized access attempt: %s %s from %s | Status: %d",
				method, path, clientIP, statusCode)
		}

		// Log slow requests (> 1 second)
		if latency > time.Second {
			log.Printf("[PERFORMANCE] Slow request detected: %s %s | Latency: %v",
				method, path, latency)
		}

		// Log errors
		if statusCode >= 500 {
			log.Printf("[ERROR] Server error: %s %s | Status: %d | RequestID: %s",
				method, path, statusCode, requestID)
		}
	}
}

// SecurityLogger logs security-related events
func SecurityLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		method := c.Request.Method
		clientIP := c.ClientIP()

		// Log sensitive endpoint access
		sensitiveEndpoints := []string{
			"/auth/login",
			"/auth/voter-login",
			"/register",
			"/admin/",
		}

		for _, endpoint := range sensitiveEndpoints {
			if contains(path, endpoint) {
				log.Printf("[SECURITY] Sensitive endpoint access: %s %s from %s",
					method, path, clientIP)
				break
			}
		}

		c.Next()

		// Log authentication failures
		if c.Writer.Status() == 401 {
			log.Printf("[SECURITY] Authentication failed: %s %s from %s",
				method, path, clientIP)
		}

		// Log authorization failures
		if c.Writer.Status() == 403 {
			userID, _ := c.Get("user_id")
			role, _ := c.Get("role")
			log.Printf("[SECURITY] Authorization failed: User %v (role: %v) attempted %s %s",
				userID, role, method, path)
		}
	}
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr
}

// AuditLogger logs important business events
func AuditLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		path := c.Request.URL.Path
		method := c.Request.Method
		status := c.Writer.Status()
		userID, _ := c.Get("user_id")
		requestID := c.GetString("request_id")

		// Log important business events
		if status >= 200 && status < 300 {
			switch {
			case method == "POST" && path == "/register":
				log.Printf("[AUDIT] Voter registration | RequestID: %s", requestID)

			case method == "POST" && path == "/admin/polls":
				log.Printf("[AUDIT] Poll created by user: %v | RequestID: %s", userID, requestID)

			case method == "POST" && path == "/vote":
				log.Printf("[AUDIT] Vote cast by user: %v | RequestID: %s", userID, requestID)

			case method == "POST" && path == "/admin/blockchain/mine":
				log.Printf("[AUDIT] Manual mining triggered by user: %v | RequestID: %s", userID, requestID)

			case method == "POST" && path == "/auth/login":
				log.Printf("[AUDIT] Admin login successful | RequestID: %s", requestID)

			case method == "POST" && path == "/auth/voter-login":
				log.Printf("[AUDIT] Voter login successful: %v | RequestID: %s", userID, requestID)
			}
		}
	}
}

// MetricsCollector collects basic metrics
type MetricsCollector struct {
	TotalRequests   int64
	TotalErrors     int64
	Total2xx        int64
	Total4xx        int64
	Total5xx        int64
	TotalAuthFails  int64
}

var metrics = &MetricsCollector{}

// MetricsMiddleware collects request metrics
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		metrics.TotalRequests++

		status := c.Writer.Status()
		switch {
		case status >= 200 && status < 300:
			metrics.Total2xx++
		case status >= 400 && status < 500:
			metrics.Total4xx++
			if status == 401 || status == 403 {
				metrics.TotalAuthFails++
			}
		case status >= 500:
			metrics.Total5xx++
			metrics.TotalErrors++
		}
	}
}

// GetMetrics returns current metrics
func GetMetrics() *MetricsCollector {
	return metrics
}

// LogMetrics logs current metrics
func LogMetrics() {
	log.Printf("[METRICS] Requests: %d | 2xx: %d | 4xx: %d | 5xx: %d | Auth Fails: %d | Errors: %d",
		metrics.TotalRequests,
		metrics.Total2xx,
		metrics.Total4xx,
		metrics.Total5xx,
		metrics.TotalAuthFails,
		metrics.TotalErrors,
	)
}
