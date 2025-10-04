package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/voting-blockchain/internal/auth"
	"github.com/voting-blockchain/internal/blockchain"
	"github.com/voting-blockchain/internal/handlers"
	"github.com/voting-blockchain/internal/middleware"
)

// setupTestRouter creates a test router with all middleware
func setupTestRouter() (*gin.Engine, *blockchain.Blockchain, *auth.JWTManager) {
	gin.SetMode(gin.TestMode)

	bc := blockchain.NewBlockchain(3)
	jwtManager := auth.NewJWTManager("test-secret", 24*time.Hour)
	adminStore := auth.NewAdminStore()
	adminStore.CreateAdmin("testadmin", "admin@test.com", "testpass123")

	router := gin.New()
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.RecoveryMiddleware())

	h := handlers.NewHandler(bc, jwtManager, adminStore)

	// Setup routes
	router.POST("/auth/login", h.Login)
	router.POST("/register", middleware.StrictRateLimit(), h.RegisterVoter)
	router.POST("/auth/voter-login", h.VoterLogin)

	authenticated := router.Group("/")
	authenticated.Use(middleware.AuthMiddleware(jwtManager))
	{
		authenticated.GET("/auth/me", h.GetCurrentUser)
	}

	adminRoutes := router.Group("/admin")
	adminRoutes.Use(middleware.AuthMiddleware(jwtManager))
	adminRoutes.Use(middleware.RequireRole("admin"))
	{
		adminRoutes.POST("/polls", h.CreatePoll)
	}

	// Add 404 handler
	router.NoRoute(middleware.NotFoundHandler())

	return router, bc, jwtManager
}

func TestJWTAuthentication(t *testing.T) {
	router, _, _ := setupTestRouter()

	t.Run("Admin login with valid credentials", func(t *testing.T) {
		loginData := map[string]string{
			"username": "testadmin",
			"password": "testpass123",
		}
		body, _ := json.Marshal(loginData)

		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		if response["success"] != true {
			t.Error("Expected success=true")
		}

		if response["token"] == nil {
			t.Error("Expected JWT token in response")
		}
	})

	t.Run("Admin login with invalid credentials", func(t *testing.T) {
		loginData := map[string]string{
			"username": "testadmin",
			"password": "wrongpassword",
		}
		body, _ := json.Marshal(loginData)

		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", w.Code)
		}
	})

	t.Run("Access protected endpoint without token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/auth/me", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", w.Code)
		}
	})
}

func TestRoleBasedAccessControl(t *testing.T) {
	router, _, jwtManager := setupTestRouter()

	t.Run("Admin can access admin endpoint", func(t *testing.T) {
		// Login as admin
		token, _ := jwtManager.GenerateToken("admin_123", "admin@test.com", "admin", "")

		pollData := map[string]interface{}{
			"title":          "Test Poll for Access Control",
			"description":    "Testing admin access control",
			"options":        []string{"Yes", "No"},
			"creator":        "testadmin",
			"duration_hours": 24,
		}
		body, _ := json.Marshal(pollData)

		req, _ := http.NewRequest("POST", "/admin/polls", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("Voter cannot access admin endpoint", func(t *testing.T) {
		// Generate voter token
		token, _ := jwtManager.GenerateToken("voter_123", "voter@test.com", "voter", "voter_123")

		pollData := map[string]interface{}{
			"title":          "Unauthorized Poll",
			"description":    "This should fail",
			"options":        []string{"Yes", "No"},
			"creator":        "hacker",
			"duration_hours": 24,
		}
		body, _ := json.Marshal(pollData)

		req, _ := http.NewRequest("POST", "/admin/polls", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Errorf("Expected status 403, got %d", w.Code)
		}
	})
}

func TestInputValidation(t *testing.T) {
	router, _, _ := setupTestRouter()

	t.Run("Reject invalid email format", func(t *testing.T) {
		regData := map[string]string{
			"email":      "notanemail",
			"name":       "Test User",
			"department": "Engineering",
		}
		body, _ := json.Marshal(regData)

		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})

	t.Run("Reject short name", func(t *testing.T) {
		regData := map[string]string{
			"email":      "test@example.com",
			"name":       "A",
			"department": "Engineering",
		}
		body, _ := json.Marshal(regData)

		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})

	t.Run("Accept valid registration data", func(t *testing.T) {
		regData := map[string]string{
			"email":      "valid@example.com",
			"name":       "Valid User",
			"department": "Engineering",
		}
		body, _ := json.Marshal(regData)

		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
	})
}

func TestRateLimiting(t *testing.T) {
	router, _, _ := setupTestRouter()

	t.Run("Rate limit registration endpoint", func(t *testing.T) {
		successCount := 0
		rateLimitedCount := 0

		// Try 10 rapid requests
		for i := 0; i < 10; i++ {
			regData := map[string]string{
				"email":      "test" + string(rune(i)) + "@example.com",
				"name":       "Test User",
				"department": "Engineering",
			}
			body, _ := json.Marshal(regData)

			req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code == http.StatusOK {
				successCount++
			} else if w.Code == http.StatusTooManyRequests {
				rateLimitedCount++
			}
		}

		// Should have some rate limited requests after 5 successful ones
		if rateLimitedCount == 0 {
			t.Error("Expected some requests to be rate limited")
		}

		if successCount > 5 {
			t.Errorf("Expected max 5 successful requests, got %d", successCount)
		}
	})
}

func TestErrorHandling(t *testing.T) {
	router, _, _ := setupTestRouter()

	t.Run("404 Not Found", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/nonexistent", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", w.Code)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		if response["error_code"] != "NOT_FOUND" {
			t.Error("Expected error_code=NOT_FOUND")
		}

		if response["request_id"] == nil {
			t.Error("Expected request_id in error response")
		}
	})
}

func TestCryptographicSignatureVerification(t *testing.T) {
	router, bc, _ := setupTestRouter()

	t.Run("Voter login with invalid private key", func(t *testing.T) {
		// First register a voter
		regData := map[string]string{
			"email":      "cryptotest@example.com",
			"name":       "Crypto Test",
			"department": "Security",
		}
		body, _ := json.Marshal(regData)

		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		var regResponse map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &regResponse)
		voterData := regResponse["data"].(map[string]interface{})
		voterID := voterData["voter_id"].(string)

		// Try to login with wrong private key
		loginData := map[string]string{
			"voter_id":    voterID,
			"private_key": "-----BEGIN PRIVATE KEY-----\nINVALID\n-----END PRIVATE KEY-----",
		}
		body, _ = json.Marshal(loginData)

		req, _ = http.NewRequest("POST", "/auth/voter-login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401 for invalid private key, got %d", w.Code)
		}
	})

	t.Run("Verify blockchain integrity", func(t *testing.T) {
		if !bc.VerifyChain() {
			t.Error("Blockchain integrity check failed")
		}
	})
}
