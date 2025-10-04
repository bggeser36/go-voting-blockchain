package persistence

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/voting-blockchain/internal/blockchain"
	"github.com/voting-blockchain/internal/models"
)

// Manager handles persistence operations
type Manager struct {
	blockchain *blockchain.Blockchain
	redisClient *redis.Client
	db          *sql.DB
	ctx         context.Context
}

// NewManager creates a new persistence manager
func NewManager(bc *blockchain.Blockchain, redisURL, databaseURL string) *Manager {
	m := &Manager{
		blockchain: bc,
		ctx:        context.Background(),
	}

	// Initialize Redis if URL provided
	if redisURL != "" {
		opt, err := redis.ParseURL(redisURL)
		if err == nil {
			m.redisClient = redis.NewClient(opt)
			// Test connection
			if err := m.redisClient.Ping(m.ctx).Err(); err != nil {
				log.Printf("Redis connection failed: %v", err)
				m.redisClient = nil
			} else {
				log.Println("Redis connected successfully")
			}
		}
	}

	// Initialize PostgreSQL if URL provided
	if databaseURL != "" {
		db, err := sql.Open("postgres", databaseURL)
		if err == nil {
			m.db = db
			// Test connection
			if err := db.Ping(); err != nil {
				log.Printf("Database connection failed: %v", err)
				m.db = nil
			} else {
				log.Println("Database connected successfully")
				m.createTables()
			}
		}
	}

	return m
}

// Initialize sets up persistence layer
func (m *Manager) Initialize() error {
	if m.redisClient != nil {
		// Set up Redis persistence
		go m.startRedisSync()
	}

	if m.db != nil {
		// Set up database persistence
		go m.startDBSync()
	}

	return nil
}

// createTables creates necessary database tables
func (m *Manager) createTables() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS blocks (
			id SERIAL PRIMARY KEY,
			block_index INTEGER UNIQUE NOT NULL,
			timestamp TIMESTAMP NOT NULL,
			data JSONB NOT NULL,
			previous_hash VARCHAR(64) NOT NULL,
			hash VARCHAR(64) UNIQUE NOT NULL,
			nonce INTEGER NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS voters (
			voter_id VARCHAR(64) PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			department VARCHAR(100),
			public_key TEXT NOT NULL,
			registered_at TIMESTAMP NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS polls (
			poll_id UUID PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			description TEXT,
			options JSONB NOT NULL,
			creator VARCHAR(255) NOT NULL,
			start_time TIMESTAMP NOT NULL,
			end_time TIMESTAMP NOT NULL,
			eligible_voters JSONB,
			allow_multiple_votes BOOLEAN DEFAULT FALSE,
			is_anonymous BOOLEAN DEFAULT FALSE
		)`,
		`CREATE TABLE IF NOT EXISTS votes (
			vote_id UUID PRIMARY KEY,
			poll_id UUID NOT NULL,
			voter_id VARCHAR(64) NOT NULL,
			choice VARCHAR(255) NOT NULL,
			timestamp TIMESTAMP NOT NULL,
			signature TEXT,
			block_index INTEGER
		)`,
		`CREATE INDEX IF NOT EXISTS idx_votes_poll ON votes(poll_id)`,
		`CREATE INDEX IF NOT EXISTS idx_votes_voter ON votes(voter_id)`,
		`CREATE INDEX IF NOT EXISTS idx_blocks_index ON blocks(block_index)`,
	}

	for _, query := range queries {
		if _, err := m.db.Exec(query); err != nil {
			log.Printf("Failed to create table: %v", err)
		}
	}
}

// SaveBlock saves a block to the database
func (m *Manager) SaveBlock(block *models.Block) error {
	if m.db == nil {
		return nil
	}

	dataJSON, err := json.Marshal(block.Data)
	if err != nil {
		return err
	}

	_, err = m.db.Exec(`
		INSERT INTO blocks (block_index, timestamp, data, previous_hash, hash, nonce)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (block_index) DO NOTHING`,
		block.Index, block.Timestamp, dataJSON, block.PreviousHash, block.Hash, block.Nonce)

	return err
}

// SaveVoter saves a voter to the database
func (m *Manager) SaveVoter(voter *models.Voter) error {
	if m.db == nil {
		return nil
	}

	_, err := m.db.Exec(`
		INSERT INTO voters (voter_id, name, email, department, public_key, registered_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (voter_id) DO NOTHING`,
		voter.VoterID, voter.Name, voter.Email, voter.Department, voter.PublicKey, voter.RegisteredAt)

	return err
}

// SavePoll saves a poll to the database
func (m *Manager) SavePoll(poll *models.Poll) error {
	if m.db == nil {
		return nil
	}

	optionsJSON, _ := json.Marshal(poll.Options)
	eligibleVotersJSON, _ := json.Marshal(poll.EligibleVoters)

	_, err := m.db.Exec(`
		INSERT INTO polls (poll_id, title, description, options, creator, start_time, end_time,
						  eligible_voters, allow_multiple_votes, is_anonymous)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (poll_id) DO NOTHING`,
		poll.PollID, poll.Title, poll.Description, optionsJSON, poll.Creator,
		poll.StartTime, poll.EndTime, eligibleVotersJSON, poll.AllowMultipleVotes, poll.IsAnonymous)

	return err
}

// SaveVote saves a vote to the database
func (m *Manager) SaveVote(vote *models.Vote, blockIndex int) error {
	if m.db == nil {
		return nil
	}

	_, err := m.db.Exec(`
		INSERT INTO votes (vote_id, poll_id, voter_id, choice, timestamp, signature, block_index)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (vote_id) DO NOTHING`,
		vote.VoteID, vote.PollID, vote.VoterID, vote.Choice, vote.Timestamp, vote.Signature, blockIndex)

	return err
}

// LoadBlockchain loads blockchain data from persistence
func (m *Manager) LoadBlockchain() error {
	if m.db == nil {
		return nil
	}

	// Load voters
	rows, err := m.db.Query(`SELECT voter_id, name, email, department, public_key, registered_at FROM voters`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var voter models.Voter
			if err := rows.Scan(&voter.VoterID, &voter.Name, &voter.Email,
							   &voter.Department, &voter.PublicKey, &voter.RegisteredAt); err == nil {
				m.blockchain.VoterRegistry[voter.VoterID] = &voter
			}
		}
	}

	// Load polls
	rows, err = m.db.Query(`SELECT poll_id, title, description, options, creator, start_time, end_time,
							eligible_voters, allow_multiple_votes, is_anonymous FROM polls`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var poll models.Poll
			var optionsJSON, eligibleVotersJSON []byte

			if err := rows.Scan(&poll.PollID, &poll.Title, &poll.Description, &optionsJSON,
							   &poll.Creator, &poll.StartTime, &poll.EndTime, &eligibleVotersJSON,
							   &poll.AllowMultipleVotes, &poll.IsAnonymous); err == nil {
				json.Unmarshal(optionsJSON, &poll.Options)
				json.Unmarshal(eligibleVotersJSON, &poll.EligibleVoters)
				m.blockchain.Polls[poll.PollID] = &poll
			}
		}
	}

	log.Printf("Loaded %d voters and %d polls from database",
			  len(m.blockchain.VoterRegistry), len(m.blockchain.Polls))

	return nil
}

// startRedisSync starts Redis synchronization
func (m *Manager) startRedisSync() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if m.redisClient != nil {
			// Save blockchain state to Redis
			chainData, _ := json.Marshal(m.blockchain.Chain)
			m.redisClient.Set(m.ctx, "blockchain:chain", chainData, 0)

			pollsData, _ := json.Marshal(m.blockchain.Polls)
			m.redisClient.Set(m.ctx, "blockchain:polls", pollsData, 0)

			votersData, _ := json.Marshal(m.blockchain.VoterRegistry)
			m.redisClient.Set(m.ctx, "blockchain:voters", votersData, 0)
		}
	}
}

// startDBSync starts database synchronization
func (m *Manager) startDBSync() {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	lastBlockIndex := -1
	savedVoters := make(map[string]bool)
	savedPolls := make(map[string]bool)

	for range ticker.C {
		if m.db != nil {
			// Save new blocks
			for i := lastBlockIndex + 1; i < len(m.blockchain.Chain); i++ {
				if err := m.SaveBlock(&m.blockchain.Chain[i]); err != nil {
					log.Printf("Failed to save block %d: %v", i, err)
				} else {
					lastBlockIndex = i
				}
			}

			// Sync voters to database
			for voterID, voter := range m.blockchain.VoterRegistry {
				if !savedVoters[voterID] {
					if err := m.SaveVoter(voter); err != nil {
						log.Printf("Failed to save voter %s: %v", voterID, err)
					} else {
						savedVoters[voterID] = true
						log.Printf("Synced voter %s to database", voterID)
					}
				}
			}

			// Sync polls to database
			for pollID, poll := range m.blockchain.Polls {
				if !savedPolls[pollID] {
					if err := m.SavePoll(poll); err != nil {
						log.Printf("Failed to save poll %s: %v", pollID, err)
					} else {
						savedPolls[pollID] = true
						log.Printf("Synced poll %s to database", pollID)
					}
				}
			}

			// Sync votes from blocks
			for _, block := range m.blockchain.Chain {
				if block.Data["type"] == "votes" {
					// Try to get votes as []models.Vote first
					if votes, ok := block.Data["votes"].([]models.Vote); ok {
						for _, vote := range votes {
							if err := m.SaveVote(&vote, block.Index); err != nil {
								log.Printf("Failed to save vote %s: %v", vote.VoteID, err)
							}
						}
					} else if votesData, ok := block.Data["votes"].([]interface{}); ok {
						// Fall back to interface conversion for JSON-deserialized data
						for _, voteInterface := range votesData {
							if voteMap, ok := voteInterface.(map[string]interface{}); ok {
								vote := models.Vote{}

								// Extract VoteID
								if vid := voteMap["VoteID"]; vid != nil {
									vote.VoteID, _ = vid.(string)
								} else if vid := voteMap["vote_id"]; vid != nil {
									vote.VoteID, _ = vid.(string)
								}

								// Extract PollID
								if pid := voteMap["PollID"]; pid != nil {
									vote.PollID, _ = pid.(string)
								} else if pid := voteMap["poll_id"]; pid != nil {
									vote.PollID, _ = pid.(string)
								}

								// Extract VoterID
								if vid := voteMap["VoterID"]; vid != nil {
									vote.VoterID, _ = vid.(string)
								} else if vid := voteMap["voter_id"]; vid != nil {
									vote.VoterID, _ = vid.(string)
								}

								// Extract Choice
								if ch := voteMap["Choice"]; ch != nil {
									vote.Choice, _ = ch.(string)
								} else if ch := voteMap["choice"]; ch != nil {
									vote.Choice, _ = ch.(string)
								}

								// Extract Signature
								if sig := voteMap["Signature"]; sig != nil {
									vote.Signature, _ = sig.(string)
								} else if sig := voteMap["signature"]; sig != nil {
									vote.Signature, _ = sig.(string)
								}

								// Extract Timestamp
								if ts := voteMap["Timestamp"]; ts != nil {
									if t, ok := ts.(time.Time); ok {
										vote.Timestamp = t
									} else if tStr, ok := ts.(string); ok {
										vote.Timestamp, _ = time.Parse(time.RFC3339, tStr)
									}
								} else if ts := voteMap["timestamp"]; ts != nil {
									if tStr, ok := ts.(string); ok {
										vote.Timestamp, _ = time.Parse(time.RFC3339, tStr)
									}
								}

								// Save the vote if it has required fields
								if vote.VoteID != "" && vote.PollID != "" {
									if err := m.SaveVote(&vote, block.Index); err != nil {
										log.Printf("Failed to save vote %s: %v", vote.VoteID, err)
									}
								}
							}
						}
					}
				}
			}

			log.Printf("DB Sync: %d voters, %d polls synced", len(savedVoters), len(savedPolls))
		}
	}
}

// Close closes all connections
func (m *Manager) Close() {
	if m.redisClient != nil {
		m.redisClient.Close()
	}
	if m.db != nil {
		m.db.Close()
	}
}