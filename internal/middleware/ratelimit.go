package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter implements a sliding window rate limiter
type RateLimiter struct {
	mu       sync.RWMutex
	clients  map[string]*clientRate
	rate     int           // requests allowed
	window   time.Duration // time window
	cleanUp  time.Duration // cleanup interval
}

type clientRate struct {
	requests []time.Time
	mu       sync.Mutex
}

// NewRateLimiter creates a new rate limiter
// rate: number of requests allowed
// window: time window for rate limiting (e.g., 1 minute)
func NewRateLimiter(rate int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		clients: make(map[string]*clientRate),
		rate:    rate,
		window:  window,
		cleanUp: window * 2,
	}

	// Start cleanup goroutine
	go rl.cleanup()

	return rl
}

// Middleware returns a Gin middleware function
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		if !rl.Allow(clientIP) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"error":   "Rate limit exceeded. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Allow checks if a request from the given client should be allowed
func (rl *RateLimiter) Allow(clientID string) bool {
	rl.mu.RLock()
	client, exists := rl.clients[clientID]
	rl.mu.RUnlock()

	if !exists {
		rl.mu.Lock()
		// Double-check after acquiring write lock
		if client, exists = rl.clients[clientID]; !exists {
			client = &clientRate{
				requests: make([]time.Time, 0),
			}
			rl.clients[clientID] = client
		}
		rl.mu.Unlock()
	}

	client.mu.Lock()
	defer client.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	// Remove old requests outside the window
	validRequests := make([]time.Time, 0)
	for _, reqTime := range client.requests {
		if reqTime.After(windowStart) {
			validRequests = append(validRequests, reqTime)
		}
	}
	client.requests = validRequests

	// Check if we're within the rate limit
	if len(client.requests) >= rl.rate {
		return false
	}

	// Add current request
	client.requests = append(client.requests, now)
	return true
}

// cleanup removes old client data periodically
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.cleanUp)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		windowStart := now.Add(-rl.window)

		for clientID, client := range rl.clients {
			client.mu.Lock()
			if len(client.requests) == 0 || client.requests[len(client.requests)-1].Before(windowStart) {
				delete(rl.clients, clientID)
			}
			client.mu.Unlock()
		}
		rl.mu.Unlock()
	}
}

// RateLimitByIP creates a simple rate limiter middleware by IP
// rate: requests per window
// window: time window (e.g., time.Minute)
func RateLimitByIP(rate int, window time.Duration) gin.HandlerFunc {
	limiter := NewRateLimiter(rate, window)
	return limiter.Middleware()
}

// StrictRateLimit returns a stricter rate limiter for sensitive endpoints
// (e.g., login, registration)
func StrictRateLimit() gin.HandlerFunc {
	// 5 requests per minute
	return RateLimitByIP(5, time.Minute)
}

// ModerateRateLimit returns a moderate rate limiter for API endpoints
func ModerateRateLimit() gin.HandlerFunc {
	// 30 requests per minute
	return RateLimitByIP(30, time.Minute)
}

// GenerousRateLimit returns a generous rate limiter for public endpoints
func GenerousRateLimit() gin.HandlerFunc {
	// 100 requests per minute
	return RateLimitByIP(100, time.Minute)
}
