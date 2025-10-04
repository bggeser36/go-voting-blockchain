package models

import (
	"time"

	"github.com/google/uuid"
)

// Vote represents a single vote in the system
type Vote struct {
	VoteID    string    `json:"vote_id"`
	PollID    string    `json:"poll_id"`
	VoterID   string    `json:"voter_id"`
	Choice    string    `json:"choice"`
	Timestamp time.Time `json:"timestamp"`
	Signature string    `json:"signature,omitempty"`
}

// Poll represents a voting poll
type Poll struct {
	PollID             string    `json:"poll_id"`
	Title              string    `json:"title"`
	Description        string    `json:"description"`
	Options            []string  `json:"options"`
	Creator            string    `json:"creator"`
	StartTime          time.Time `json:"start_time"`
	EndTime            time.Time `json:"end_time"`
	EligibleVoters     []string  `json:"eligible_voters,omitempty"`
	AllowMultipleVotes bool      `json:"allow_multiple_votes"`
	IsAnonymous        bool      `json:"is_anonymous"`
}

// IsActive checks if the poll is currently active
func (p *Poll) IsActive() bool {
	now := time.Now()
	return now.After(p.StartTime) && now.Before(p.EndTime)
}

// Block represents a block in the blockchain
type Block struct {
	Index        int                    `json:"index"`
	Timestamp    time.Time              `json:"timestamp"`
	Data         map[string]interface{} `json:"data"`
	PreviousHash string                 `json:"previous_hash"`
	Nonce        int                    `json:"nonce"`
	Hash         string                 `json:"hash"`
}

// Voter represents a registered voter
type Voter struct {
	VoterID        string    `json:"voter_id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	Department     string    `json:"department,omitempty"`
	PublicKey      string    `json:"public_key"`
	RegisteredAt   time.Time `json:"registered_at"`
}

// PollResults represents the results of a poll
type PollResults struct {
	PollID       string         `json:"poll_id"`
	Title        string         `json:"title"`
	Status       string         `json:"status"`
	Results      map[string]int `json:"results"`
	TotalVotes   int            `json:"total_votes"`
	VoterTurnout string         `json:"voter_turnout"`
}

// VoterHistory represents a voter's voting history
type VoterHistory struct {
	VoteID     string    `json:"vote_id"`
	PollID     string    `json:"poll_id"`
	PollTitle  string    `json:"poll_title,omitempty"`
	Timestamp  time.Time `json:"timestamp"`
	BlockIndex int       `json:"block_index"`
}

// BlockchainStats represents blockchain statistics
type BlockchainStats struct {
	ChainLength       int    `json:"chain_length"`
	TotalVotes        int    `json:"total_votes"`
	PendingVotes      int    `json:"pending_votes"`
	TotalVoters       int    `json:"total_voters"`
	TotalPolls        int    `json:"total_polls"`
	ActivePolls       int    `json:"active_polls"`
	BlockchainValid   bool   `json:"blockchain_valid"`
	MiningDifficulty  int    `json:"mining_difficulty"`
	MiningThreshold   int    `json:"mining_threshold"`
}

// API Request/Response models

// VoterRegistrationRequest represents a voter registration request
type VoterRegistrationRequest struct {
	Email      string `json:"email" binding:"required,email"`
	Name       string `json:"name" binding:"required"`
	Department string `json:"department"`
}

// PollCreationRequest represents a poll creation request
type PollCreationRequest struct {
	Title              string   `json:"title" binding:"required"`
	Description        string   `json:"description" binding:"required"`
	Options            []string `json:"options" binding:"required,min=2"`
	Creator            string   `json:"creator" binding:"required"`
	DurationHours      float64  `json:"duration_hours"`
	EligibleVoters     []string `json:"eligible_voters"`
	AllowMultipleVotes bool     `json:"allow_multiple_votes"`
	IsAnonymous        bool     `json:"is_anonymous"`
}

// VoteSubmissionRequest represents a vote submission request
type VoteSubmissionRequest struct {
	PollID    string `json:"poll_id" binding:"required"`
	VoterID   string `json:"voter_id" binding:"required"`
	Choice    string `json:"choice" binding:"required"`
	Signature string `json:"signature"`
}

// APIResponse represents a generic API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// GenerateID creates a new UUID
func GenerateID() string {
	return uuid.New().String()
}