package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/voting-blockchain/internal/auth"
	"github.com/voting-blockchain/internal/blockchain"
	"github.com/voting-blockchain/internal/handlers"
	"github.com/voting-blockchain/internal/middleware"
	"github.com/voting-blockchain/internal/persistence"
)

// maskConnectionString masks sensitive parts of connection strings for logging
func maskConnectionString(connStr string) string {
	if connStr == "" {
		return "(not configured)"
	}

	// For Redis URLs
	if strings.Contains(connStr, "redis://") {
		parts := strings.Split(connStr, "@")
		if len(parts) > 1 {
			return "redis://***@" + parts[1]
		}
	}

	// For PostgreSQL URLs
	if strings.Contains(connStr, "postgresql://") {
		parts := strings.Split(connStr, "@")
		if len(parts) > 1 {
			return "postgresql://***@" + parts[1]
		}
	}

	return "(configured)"
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Get configuration from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	difficulty := 3
	redisURL := os.Getenv("REDIS_URL")
	databaseURL := os.Getenv("DATABASE_URL")

	// Log configuration
	log.Printf("🔧 Configuration:")
	log.Printf("   Port: %s", port)
	log.Printf("   Difficulty: %d", difficulty)
	log.Printf("   Redis URL: %s", maskConnectionString(redisURL))
	log.Printf("   Database URL: %s", maskConnectionString(databaseURL))

	// Initialize blockchain
	bc := blockchain.NewBlockchain(difficulty)

	// Initialize JWT manager
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key-change-in-production" // Default for development
		log.Println("⚠️  Using default JWT secret - change in production!")
	}
	jwtManager := auth.NewJWTManager(jwtSecret, 24*time.Hour)
	log.Println("✅ JWT manager initialized")

	// Initialize admin store and create default admin
	adminStore := auth.NewAdminStore()
	defaultAdmin, err := adminStore.CreateAdmin("admin", "admin@voting.com", "admin123")
	if err != nil {
		log.Printf("⚠️  Failed to create default admin: %v", err)
	} else {
		log.Printf("✅ Default admin created: %s (password: admin123)", defaultAdmin.Username)
		log.Println("⚠️  CHANGE DEFAULT ADMIN PASSWORD IN PRODUCTION!")
	}

	// Initialize persistence manager if database URLs are provided
	var persistenceManager *persistence.Manager
	if redisURL != "" || databaseURL != "" {
		log.Println("📦 Initializing persistence layer...")
		persistenceManager = persistence.NewManager(bc, redisURL, databaseURL)

		// Initialize persistence (starts background sync)
		if err := persistenceManager.Initialize(); err != nil {
			log.Printf("⚠️ Failed to initialize persistence: %v", err)
		} else {
			log.Println("✅ Persistence layer initialized")

			// Load existing blockchain data from database
			if err := persistenceManager.LoadBlockchain(); err != nil {
				log.Printf("⚠️ Failed to load blockchain data: %v", err)
			} else {
				stats := bc.GetStats()
				log.Printf("📊 Loaded from database:")
				log.Printf("   Voters: %d", stats.TotalVoters)
				log.Printf("   Polls: %d", stats.TotalPolls)
				log.Printf("   Blocks: %d", stats.ChainLength)
			}
		}

		// Ensure cleanup on shutdown
		defer func() {
			if persistenceManager != nil {
				log.Println("🔄 Closing persistence connections...")
				persistenceManager.Close()
			}
		}()
	} else {
		log.Println("⚠️ No persistence layer configured (running in-memory only)")
	}

	// Create router
	router := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	router.Use(cors.New(config))

	// Initialize handlers
	h := handlers.NewHandler(bc, jwtManager, adminStore)

	// Public routes (no authentication required)
	router.GET("/", h.GetStatus)
	router.GET("/health", h.HealthCheck)

	// Authentication routes
	router.POST("/auth/login", h.Login)
	router.POST("/auth/voter-login", h.VoterLogin)
	router.POST("/auth/refresh", h.RefreshToken)

	// Public voter registration (anyone can register)
	router.POST("/register", h.RegisterVoter)

	// Protected routes - require authentication
	authenticated := router.Group("/")
	authenticated.Use(middleware.AuthMiddleware(jwtManager))
	{
		// Current user info
		authenticated.GET("/auth/me", h.GetCurrentUser)

		// Voter history (voters can only see their own)
		authenticated.GET("/voter/:voter_id/history", h.GetVoterHistory)

		// Voting (must be authenticated as voter)
		authenticated.POST("/vote", h.SubmitVote)
	}

	// Admin-only routes
	adminRoutes := router.Group("/admin")
	adminRoutes.Use(middleware.AuthMiddleware(jwtManager))
	adminRoutes.Use(middleware.RequireRole("admin"))
	{
		// Poll management
		adminRoutes.POST("/polls", h.CreatePoll)

		// Manual mining
		adminRoutes.POST("/blockchain/mine", h.MinePendingVotes)
	}

	// Public read-only blockchain routes
	router.GET("/polls", h.GetPolls)
	router.GET("/polls/:poll_id", h.GetPollDetails)
	router.GET("/results/:poll_id", h.GetPollResults)
	router.GET("/blockchain/verify", h.VerifyBlockchain)
	router.GET("/blockchain/blocks", h.GetBlocks)
	router.GET("/blockchain/stats", h.GetBlockchainStats)

	// Start server
	log.Printf("🚀 Blockchain Voting System starting on port %s", port)
	log.Printf("📊 Mining difficulty: %d", difficulty)
	log.Printf("🔗 Blockchain initialized with genesis block")

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}