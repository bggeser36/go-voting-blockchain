package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/voting-blockchain/internal/models"
)

// Blockchain represents the main blockchain structure
type Blockchain struct {
	Chain           []models.Block
	Difficulty      int
	MiningThreshold int // Number of votes needed before mining a block
	PendingVotes    []models.Vote
	Polls           map[string]*models.Poll
	VoterRegistry   map[string]*models.Voter
	VoteRecords     map[string][]string // poll_id -> voter_ids who voted
	mu              sync.RWMutex
}

// NewBlockchain creates a new blockchain instance
func NewBlockchain(difficulty int) *Blockchain {
	bc := &Blockchain{
		Chain:           make([]models.Block, 0),
		Difficulty:      difficulty,
		MiningThreshold: 5, // Default to 5 votes for efficiency
		PendingVotes:    make([]models.Vote, 0),
		Polls:           make(map[string]*models.Poll),
		VoterRegistry:   make(map[string]*models.Voter),
		VoteRecords:     make(map[string][]string),
	}

	// Create genesis block
	bc.createGenesisBlock()
	return bc
}

// createGenesisBlock creates the first block in the chain
func (bc *Blockchain) createGenesisBlock() {
	genesisData := map[string]interface{}{
		"type":    "genesis",
		"message": "Genesis Block - Voting Blockchain Initialized",
	}

	genesisBlock := models.Block{
		Index:        0,
		Timestamp:    time.Now(),
		Data:         genesisData,
		PreviousHash: "0",
		Nonce:        0,
	}

	bc.mineBlock(&genesisBlock)
	bc.Chain = append(bc.Chain, genesisBlock)
}

// GetLatestBlock returns the most recent block in the chain
func (bc *Blockchain) GetLatestBlock() models.Block {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return bc.Chain[len(bc.Chain)-1]
}

// calculateHash calculates the SHA-256 hash of a block
func (bc *Blockchain) calculateHash(block *models.Block) string {
	blockData := map[string]interface{}{
		"index":         block.Index,
		"timestamp":     block.Timestamp.Unix(),
		"data":          block.Data,
		"previous_hash": block.PreviousHash,
		"nonce":         block.Nonce,
	}

	jsonData, _ := json.Marshal(blockData)
	hash := sha256.Sum256(jsonData)
	return hex.EncodeToString(hash[:])
}

// mineBlock mines a block using proof-of-work
func (bc *Blockchain) mineBlock(block *models.Block) {
	target := strings.Repeat("0", bc.Difficulty)

	for {
		block.Hash = bc.calculateHash(block)
		if strings.HasPrefix(block.Hash, target) {
			fmt.Printf("Block mined: %s\n", block.Hash)
			break
		}
		block.Nonce++
	}
}

// RegisterVoter registers a new voter in the system
func (bc *Blockchain) RegisterVoter(voter *models.Voter) error {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	if _, exists := bc.VoterRegistry[voter.VoterID]; exists {
		return fmt.Errorf("voter already registered")
	}

	voter.RegisteredAt = time.Now()
	bc.VoterRegistry[voter.VoterID] = voter

	// Add registration to blockchain
	regData := map[string]interface{}{
		"type":      "voter_registration",
		"voter_id":  voter.VoterID,
		"timestamp": time.Now().Unix(),
	}

	bc.addBlock(regData)
	return nil
}

// CreatePoll creates a new voting poll
func (bc *Blockchain) CreatePoll(poll *models.Poll) error {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	poll.PollID = models.GenerateID()

	// If no specific voters listed, allow all registered voters
	if len(poll.EligibleVoters) == 0 {
		poll.EligibleVoters = make([]string, 0, len(bc.VoterRegistry))
		for voterID := range bc.VoterRegistry {
			poll.EligibleVoters = append(poll.EligibleVoters, voterID)
		}
	}

	bc.Polls[poll.PollID] = poll
	bc.VoteRecords[poll.PollID] = make([]string, 0)

	// Add poll creation to blockchain
	pollData := map[string]interface{}{
		"type": "poll_creation",
		"poll": poll,
	}

	bc.addBlock(pollData)
	return nil
}

// CastVote casts a vote in a poll
func (bc *Blockchain) CastVote(vote *models.Vote) error {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	// Validate poll exists and is active
	poll, exists := bc.Polls[vote.PollID]
	if !exists {
		return fmt.Errorf("poll does not exist")
	}

	if !poll.IsActive() {
		return fmt.Errorf("poll is not active")
	}

	// Validate voter is registered
	if _, exists := bc.VoterRegistry[vote.VoterID]; !exists {
		return fmt.Errorf("voter not registered")
	}

	// Check if voter is eligible for this poll
	eligible := false
	for _, eligibleVoter := range poll.EligibleVoters {
		if eligibleVoter == vote.VoterID {
			eligible = true
			break
		}
	}
	if !eligible {
		return fmt.Errorf("voter not eligible for this poll")
	}

	// Check if voter already voted (unless multiple votes allowed)
	if !poll.AllowMultipleVotes {
		for _, votedID := range bc.VoteRecords[vote.PollID] {
			if votedID == vote.VoterID {
				return fmt.Errorf("voter has already voted in this poll")
			}
		}
	}

	// Validate choice is valid
	validChoice := false
	for _, option := range poll.Options {
		if option == vote.Choice {
			validChoice = true
			break
		}
	}
	if !validChoice {
		return fmt.Errorf("invalid voting choice")
	}

	// Add vote
	vote.VoteID = models.GenerateID()
	vote.Timestamp = time.Now()

	if poll.IsAnonymous {
		vote.VoterID = "anonymous"
	}

	bc.PendingVotes = append(bc.PendingVotes, *vote)
	bc.VoteRecords[vote.PollID] = append(bc.VoteRecords[vote.PollID], vote.VoterID)

	// Mine pending votes if we have enough
	if len(bc.PendingVotes) >= bc.MiningThreshold {
		bc.minePendingVotes()
	}

	return nil
}

// minePendingVotes mines all pending votes into a new block
func (bc *Blockchain) minePendingVotes() {
	if len(bc.PendingVotes) == 0 {
		return
	}

	votes := make([]models.Vote, len(bc.PendingVotes))
	copy(votes, bc.PendingVotes)

	votesData := map[string]interface{}{
		"type":  "votes",
		"votes": votes,
		"count": len(votes),
	}

	bc.addBlock(votesData)
	bc.PendingVotes = make([]models.Vote, 0)
}

// addBlock adds a new block to the chain
func (bc *Blockchain) addBlock(data map[string]interface{}) {
	previousBlock := bc.Chain[len(bc.Chain)-1]

	newBlock := models.Block{
		Index:        len(bc.Chain),
		Timestamp:    time.Now(),
		Data:         data,
		PreviousHash: previousBlock.Hash,
		Nonce:        0,
	}

	bc.mineBlock(&newBlock)
	bc.Chain = append(bc.Chain, newBlock)
}

// GetPollResults gets the voting results for a poll
func (bc *Blockchain) GetPollResults(pollID string) (*models.PollResults, error) {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	poll, exists := bc.Polls[pollID]
	if !exists {
		return nil, fmt.Errorf("poll does not exist")
	}

	results := make(map[string]int)
	for _, option := range poll.Options {
		results[option] = 0
	}

	totalVotes := 0

	// Count votes from all blocks
	for _, block := range bc.Chain {
		if block.Data["type"] == "votes" {
			// Try to get votes as []models.Vote first (when freshly added)
			if votes, ok := block.Data["votes"].([]models.Vote); ok {
				for _, vote := range votes {
					if vote.PollID == pollID {
						results[vote.Choice]++
						totalVotes++
					}
				}
			} else if votesData, ok := block.Data["votes"].([]interface{}); ok {
				// Fall back to interface conversion (when loaded from JSON)
				for _, voteInterface := range votesData {
					if voteMap, ok := voteInterface.(map[string]interface{}); ok {
						// Handle both PollID and poll_id keys for compatibility
						pollIDValue := voteMap["PollID"]
						if pollIDValue == nil {
							pollIDValue = voteMap["poll_id"]
						}

						// Convert pollIDValue to string for comparison
						if pollIDStr, ok := pollIDValue.(string); ok && pollIDStr == pollID {
							// Handle both Choice and choice keys
							choiceValue := voteMap["Choice"]
							if choiceValue == nil {
								choiceValue = voteMap["choice"]
							}

							if choice, ok := choiceValue.(string); ok {
								results[choice]++
								totalVotes++
							}
						}
					}
				}
			}
		}
	}

	// Include pending votes
	for _, vote := range bc.PendingVotes {
		if vote.PollID == pollID {
			results[vote.Choice]++
			totalVotes++
		}
	}

	status := "closed"
	if poll.IsActive() {
		status = "active"
	}

	voterTurnout := "N/A"
	if len(poll.EligibleVoters) > 0 {
		turnout := float64(len(bc.VoteRecords[pollID])) / float64(len(poll.EligibleVoters)) * 100
		voterTurnout = fmt.Sprintf("%.1f%%", turnout)
	}

	return &models.PollResults{
		PollID:       pollID,
		Title:        poll.Title,
		Status:       status,
		Results:      results,
		TotalVotes:   totalVotes,
		VoterTurnout: voterTurnout,
	}, nil
}

// VerifyChain verifies the integrity of the blockchain
func (bc *Blockchain) VerifyChain() bool {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	for i := 1; i < len(bc.Chain); i++ {
		currentBlock := &bc.Chain[i]
		previousBlock := &bc.Chain[i-1]

		// Check if current block's hash is correct
		if currentBlock.Hash != bc.calculateHash(currentBlock) {
			return false
		}

		// Check if previous hash matches
		if currentBlock.PreviousHash != previousBlock.Hash {
			return false
		}

		// Check if block is properly mined
		target := strings.Repeat("0", bc.Difficulty)
		if !strings.HasPrefix(currentBlock.Hash, target) {
			return false
		}
	}

	return true
}

// GetVoterHistory gets voting history for a specific voter
func (bc *Blockchain) GetVoterHistory(voterID string) []models.VoterHistory {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	history := make([]models.VoterHistory, 0)

	for _, block := range bc.Chain {
		if block.Data["type"] == "votes" {
			// Try to get votes as []models.Vote first
			if votes, ok := block.Data["votes"].([]models.Vote); ok {
				for _, vote := range votes {
					if vote.VoterID == voterID {
						vh := models.VoterHistory{
							VoteID:     vote.VoteID,
							PollID:     vote.PollID,
							BlockIndex: block.Index,
							Timestamp:  vote.Timestamp,
						}

						if poll, exists := bc.Polls[vh.PollID]; exists {
							vh.PollTitle = poll.Title
						}

						history = append(history, vh)
					}
				}
			} else if votesData, ok := block.Data["votes"].([]interface{}); ok {
				// Fall back to interface conversion
				for _, voteInterface := range votesData {
					if voteMap, ok := voteInterface.(map[string]interface{}); ok {
						// Handle both VoterID and voter_id keys
						voterIDValue := voteMap["VoterID"]
						if voterIDValue == nil {
							voterIDValue = voteMap["voter_id"]
						}

						// Convert voterIDValue to string for comparison
						if voterIDStr, ok := voterIDValue.(string); ok && voterIDStr == voterID {
							vh := models.VoterHistory{
								BlockIndex: block.Index,
							}

							// Handle VoteID
							if voteID := voteMap["VoteID"]; voteID != nil {
								vh.VoteID = voteID.(string)
							} else if voteID := voteMap["vote_id"]; voteID != nil {
								vh.VoteID = voteID.(string)
							}

							// Handle PollID
							if pollID := voteMap["PollID"]; pollID != nil {
								vh.PollID = pollID.(string)
							} else if pollID := voteMap["poll_id"]; pollID != nil {
								vh.PollID = pollID.(string)
							}

							// Handle Timestamp
							if timestamp := voteMap["Timestamp"]; timestamp != nil {
								if t, ok := timestamp.(time.Time); ok {
									vh.Timestamp = t
								} else if t, ok := timestamp.(string); ok {
									vh.Timestamp, _ = time.Parse(time.RFC3339, t)
								}
							} else if timestamp := voteMap["timestamp"]; timestamp != nil {
								if t, ok := timestamp.(float64); ok {
									vh.Timestamp = time.Unix(int64(t), 0)
								}
							}

							if poll, exists := bc.Polls[vh.PollID]; exists {
								vh.PollTitle = poll.Title
							}

							history = append(history, vh)
						}
					}
				}
			}
		}
	}

	return history
}

// GetStats returns blockchain statistics
func (bc *Blockchain) GetStats() models.BlockchainStats {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	totalVotes := 0
	for _, records := range bc.VoteRecords {
		totalVotes += len(records)
	}

	activePolls := 0
	for _, poll := range bc.Polls {
		if poll.IsActive() {
			activePolls++
		}
	}

	return models.BlockchainStats{
		ChainLength:      len(bc.Chain),
		TotalVotes:       totalVotes + len(bc.PendingVotes),
		PendingVotes:     len(bc.PendingVotes),
		TotalVoters:      len(bc.VoterRegistry),
		TotalPolls:       len(bc.Polls),
		ActivePolls:      activePolls,
		BlockchainValid:  bc.VerifyChain(),
		MiningDifficulty: bc.Difficulty,
		MiningThreshold:  bc.MiningThreshold,
	}
}

// ExportChain exports the entire blockchain
func (bc *Blockchain) ExportChain() []models.Block {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	chain := make([]models.Block, len(bc.Chain))
	copy(chain, bc.Chain)
	return chain
}

// MinePendingVotesManually manually triggers mining of pending votes
func (bc *Blockchain) MinePendingVotesManually() int {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	votesCount := len(bc.PendingVotes)
	if votesCount > 0 {
		bc.minePendingVotes()
	}
	return votesCount
}