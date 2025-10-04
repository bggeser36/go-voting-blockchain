package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/voting-blockchain/internal/blockchain"
	"github.com/voting-blockchain/internal/handlers"
	"github.com/voting-blockchain/internal/persistence"
)

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

	// Initialize blockchain
	bc := blockchain.NewBlockchain(difficulty)

	// Initialize persistence if available
	if redisURL != "" || databaseURL != "" {
		persistenceManager := persistence.NewManager(bc, redisURL, databaseURL)
		if err := persistenceManager.Initialize(); err != nil {
			log.Printf("Failed to initialize persistence: %v", err)
		} else {
			// Load existing blockchain data
			if err := persistenceManager.LoadBlockchain(); err != nil {
				log.Printf("Failed to load blockchain data: %v", err)
			}
		}
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
	h := handlers.NewHandler(bc)

	// Define routes
	router.GET("/", h.GetStatus)
	router.GET("/health", h.HealthCheck)

	// Voter routes
	router.POST("/register", h.RegisterVoter)
	router.GET("/voter/:voter_id/history", h.GetVoterHistory)

	// Poll routes
	router.POST("/polls", h.CreatePoll)
	router.GET("/polls", h.GetPolls)
	router.GET("/polls/:poll_id", h.GetPollDetails)
	router.GET("/results/:poll_id", h.GetPollResults)

	// Voting routes
	router.POST("/vote", h.SubmitVote)

	// Blockchain routes
	router.GET("/blockchain/verify", h.VerifyBlockchain)
	router.GET("/blockchain/blocks", h.GetBlocks)
	router.GET("/blockchain/stats", h.GetBlockchainStats)
	router.POST("/blockchain/mine", h.MinePendingVotes)

	// Start server
	log.Printf("🚀 Blockchain Voting System starting on port %s", port)
	log.Printf("📊 Mining difficulty: %d", difficulty)
	log.Printf("🔗 Blockchain initialized with genesis block")

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}