package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/voting-blockchain/internal/blockchain"
	"github.com/voting-blockchain/internal/crypto"
	"github.com/voting-blockchain/internal/models"
)

// Handler manages API handlers
type Handler struct {
	blockchain *blockchain.Blockchain
	crypto     *crypto.CryptoManager
}

// NewHandler creates a new handler instance
func NewHandler(bc *blockchain.Blockchain) *Handler {
	return &Handler{
		blockchain: bc,
		crypto:     crypto.NewCryptoManager(),
	}
}

// GetStatus returns API status information
func (h *Handler) GetStatus(c *gin.Context) {
	stats := h.blockchain.GetStats()

	c.JSON(http.StatusOK, gin.H{
		"name":             "Blockchain Voting System",
		"status":           "operational",
		"blockchain_height": stats.ChainLength,
		"pending_votes":    stats.PendingVotes,
		"active_polls":     stats.ActivePolls,
		"total_voters":     stats.TotalVoters,
	})
}

// RegisterVoter handles voter registration
func (h *Handler) RegisterVoter(c *gin.Context) {
	var req models.VoterRegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Create secure voter
	voterID := h.crypto.GenerateVoterID(req.Email, "")
	privateKey, publicKey, err := h.crypto.GenerateKeyPair()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to generate credentials",
		})
		return
	}

	// Register in blockchain
	voter := &models.Voter{
		VoterID:    voterID,
		Name:       req.Name,
		Email:      req.Email,
		Department: req.Department,
		PublicKey:  string(publicKey),
	}

	if err := h.blockchain.RegisterVoter(voter); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Registration successful",
		Data: gin.H{
			"voter_id":    voterID,
			"private_key": string(privateKey),
			"public_key":  string(publicKey),
		},
	})
}

// CreatePoll handles poll creation
func (h *Handler) CreatePoll(c *gin.Context) {
	var req models.PollCreationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	poll := &models.Poll{
		Title:              req.Title,
		Description:        req.Description,
		Options:            req.Options,
		Creator:            req.Creator,
		StartTime:          time.Now(),
		EndTime:            time.Now().Add(time.Duration(req.DurationHours) * time.Hour),
		EligibleVoters:     req.EligibleVoters,
		AllowMultipleVotes: req.AllowMultipleVotes,
		IsAnonymous:        req.IsAnonymous,
	}

	if err := h.blockchain.CreatePoll(poll); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Poll created successfully",
		Data: gin.H{
			"poll_id":               poll.PollID,
			"title":                 poll.Title,
			"start_time":            poll.StartTime,
			"end_time":              poll.EndTime,
			"options":               poll.Options,
			"eligible_voters_count": len(poll.EligibleVoters),
		},
	})
}

// GetPolls returns all polls
func (h *Handler) GetPolls(c *gin.Context) {
	activeOnly := c.Query("active_only") == "true"

	polls := make([]gin.H, 0)
	for pollID, poll := range h.blockchain.Polls {
		if activeOnly && !poll.IsActive() {
			continue
		}

		status := "closed"
		if poll.IsActive() {
			status = "active"
		}

		polls = append(polls, gin.H{
			"poll_id":      pollID,
			"title":        poll.Title,
			"description":  poll.Description,
			"status":       status,
			"options":      poll.Options,
			"creator":      poll.Creator,
			"start_time":   poll.StartTime,
			"end_time":     poll.EndTime,
			"is_anonymous": poll.IsAnonymous,
		})
	}

	c.JSON(http.StatusOK, polls)
}

// GetPollDetails returns details for a specific poll
func (h *Handler) GetPollDetails(c *gin.Context) {
	pollID := c.Param("poll_id")

	poll, exists := h.blockchain.Polls[pollID]
	if !exists {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "Poll not found",
		})
		return
	}

	status := "closed"
	if poll.IsActive() {
		status = "active"
	}

	votesCount := 0
	if records, exists := h.blockchain.VoteRecords[pollID]; exists {
		votesCount = len(records)
	}

	c.JSON(http.StatusOK, gin.H{
		"poll_id":                pollID,
		"title":                  poll.Title,
		"description":            poll.Description,
		"status":                 status,
		"options":                poll.Options,
		"creator":                poll.Creator,
		"start_time":             poll.StartTime,
		"end_time":               poll.EndTime,
		"eligible_voters_count":  len(poll.EligibleVoters),
		"allow_multiple_votes":   poll.AllowMultipleVotes,
		"is_anonymous":           poll.IsAnonymous,
		"votes_cast":             votesCount,
	})
}

// SubmitVote handles vote submission
func (h *Handler) SubmitVote(c *gin.Context) {
	var req models.VoteSubmissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	vote := &models.Vote{
		PollID:    req.PollID,
		VoterID:   req.VoterID,
		Choice:    req.Choice,
		Signature: req.Signature,
	}

	if err := h.blockchain.CastVote(vote); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Vote cast successfully",
		Data: gin.H{
			"vote_id": vote.VoteID,
		},
	})
}

// GetPollResults returns voting results for a poll
func (h *Handler) GetPollResults(c *gin.Context) {
	pollID := c.Param("poll_id")

	results, err := h.blockchain.GetPollResults(pollID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, results)
}

// GetVoterHistory returns voting history for a voter
func (h *Handler) GetVoterHistory(c *gin.Context) {
	voterID := c.Param("voter_id")

	if _, exists := h.blockchain.VoterRegistry[voterID]; !exists {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "Voter not found",
		})
		return
	}

	history := h.blockchain.GetVoterHistory(voterID)
	c.JSON(http.StatusOK, history)
}

// VerifyBlockchain verifies blockchain integrity
func (h *Handler) VerifyBlockchain(c *gin.Context) {
	isValid := h.blockchain.VerifyChain()

	message := "Blockchain is valid and secure"
	if !isValid {
		message = "Blockchain integrity compromised!"
	}

	c.JSON(http.StatusOK, gin.H{
		"is_valid":     isValid,
		"chain_length": len(h.blockchain.Chain),
		"message":      message,
	})
}

// GetBlocks returns blockchain blocks
func (h *Handler) GetBlocks(c *gin.Context) {
	limit := 10
	if l := c.Query("limit"); l != "" {
		// Parse limit if provided
		// For simplicity, using default if parse fails
	}

	blocks := h.blockchain.ExportChain()

	// Return last 'limit' blocks
	if len(blocks) > limit {
		blocks = blocks[len(blocks)-limit:]
	}

	c.JSON(http.StatusOK, blocks)
}

// GetBlockchainStats returns blockchain statistics
func (h *Handler) GetBlockchainStats(c *gin.Context) {
	stats := h.blockchain.GetStats()
	c.JSON(http.StatusOK, stats)
}

// MinePendingVotes manually triggers mining
func (h *Handler) MinePendingVotes(c *gin.Context) {
	votesCount := h.blockchain.MinePendingVotesManually()

	if votesCount == 0 {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "No pending votes to mine",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Votes mined successfully",
		Data: gin.H{
			"votes_mined":     votesCount,
			"new_block_index": len(h.blockchain.Chain) - 1,
		},
	})
}

// HealthCheck for Railway deployment
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"time":   time.Now(),
	})
}