-- Database initialization script for Voting Blockchain System
-- This script sets up the initial database structure

-- Create extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create custom types
CREATE TYPE vote_status AS ENUM ('pending', 'confirmed', 'rejected');

-- Create tables with proper constraints and indexes

-- Blocks table for blockchain storage
CREATE TABLE IF NOT EXISTS blocks (
    id SERIAL PRIMARY KEY,
    block_index INTEGER UNIQUE NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    data JSONB NOT NULL,
    previous_hash VARCHAR(64) NOT NULL,
    hash VARCHAR(64) UNIQUE NOT NULL,
    nonce INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Voters table for registered voters
CREATE TABLE IF NOT EXISTS voters (
    voter_id VARCHAR(64) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    department VARCHAR(100),
    public_key TEXT NOT NULL,
    registered_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Polls table for voting polls
CREATE TABLE IF NOT EXISTS polls (
    poll_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    options JSONB NOT NULL,
    creator VARCHAR(255) NOT NULL,
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    eligible_voters JSONB,
    allow_multiple_votes BOOLEAN DEFAULT FALSE,
    is_anonymous BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Votes table for individual votes
CREATE TABLE IF NOT EXISTS votes (
    vote_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    poll_id UUID NOT NULL REFERENCES polls(poll_id) ON DELETE CASCADE,
    voter_id VARCHAR(64) NOT NULL REFERENCES voters(voter_id) ON DELETE CASCADE,
    choice VARCHAR(255) NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    signature TEXT,
    block_index INTEGER,
    status vote_status DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Vote records table for tracking who voted in which poll
CREATE TABLE IF NOT EXISTS vote_records (
    id SERIAL PRIMARY KEY,
    poll_id UUID NOT NULL REFERENCES polls(poll_id) ON DELETE CASCADE,
    voter_id VARCHAR(64) NOT NULL REFERENCES voters(voter_id) ON DELETE CASCADE,
    voted_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(poll_id, voter_id)
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_blocks_index ON blocks(block_index);
CREATE INDEX IF NOT EXISTS idx_blocks_hash ON blocks(hash);
CREATE INDEX IF NOT EXISTS idx_blocks_timestamp ON blocks(timestamp);

CREATE INDEX IF NOT EXISTS idx_voters_email ON voters(email);
CREATE INDEX IF NOT EXISTS idx_voters_active ON voters(is_active);

CREATE INDEX IF NOT EXISTS idx_polls_active ON polls(is_active);
CREATE INDEX IF NOT EXISTS idx_polls_time_range ON polls(start_time, end_time);
CREATE INDEX IF NOT EXISTS idx_polls_creator ON polls(creator);

CREATE INDEX IF NOT EXISTS idx_votes_poll ON votes(poll_id);
CREATE INDEX IF NOT EXISTS idx_votes_voter ON votes(voter_id);
CREATE INDEX IF NOT EXISTS idx_votes_timestamp ON votes(timestamp);
CREATE INDEX IF NOT EXISTS idx_votes_status ON votes(status);
CREATE INDEX IF NOT EXISTS idx_votes_block ON votes(block_index);

CREATE INDEX IF NOT EXISTS idx_vote_records_poll ON vote_records(poll_id);
CREATE INDEX IF NOT EXISTS idx_vote_records_voter ON vote_records(voter_id);

-- Create views for common queries

-- Active polls view
CREATE OR REPLACE VIEW active_polls AS
SELECT 
    p.*,
    COUNT(v.vote_id) as votes_cast,
    COUNT(vr.id) as voters_participated
FROM polls p
LEFT JOIN votes v ON p.poll_id = v.poll_id
LEFT JOIN vote_records vr ON p.poll_id = vr.poll_id
WHERE p.is_active = TRUE 
    AND NOW() BETWEEN p.start_time AND p.end_time
GROUP BY p.poll_id;

-- Poll results view
CREATE OR REPLACE VIEW poll_results AS
SELECT 
    p.poll_id,
    p.title,
    p.description,
    p.options,
    v.choice,
    COUNT(*) as vote_count,
    ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER (PARTITION BY p.poll_id), 2) as percentage
FROM polls p
JOIN votes v ON p.poll_id = v.poll_id
WHERE v.status = 'confirmed'
GROUP BY p.poll_id, p.title, p.description, p.options, v.choice;

-- Blockchain stats view
CREATE OR REPLACE VIEW blockchain_stats AS
SELECT 
    (SELECT COUNT(*) FROM blocks) as total_blocks,
    (SELECT COUNT(*) FROM voters WHERE is_active = TRUE) as total_voters,
    (SELECT COUNT(*) FROM polls WHERE is_active = TRUE) as total_polls,
    (SELECT COUNT(*) FROM votes WHERE status = 'confirmed') as total_votes,
    (SELECT COUNT(*) FROM votes WHERE status = 'pending') as pending_votes,
    (SELECT MAX(block_index) FROM blocks) as latest_block_index;

-- Insert some sample data for testing
INSERT INTO voters (voter_id, name, email, department, public_key) VALUES
('test_voter_1', 'John Doe', 'john.doe@example.com', 'Engineering', 'sample_public_key_1'),
('test_voter_2', 'Jane Smith', 'jane.smith@example.com', 'Marketing', 'sample_public_key_2'),
('test_voter_3', 'Bob Johnson', 'bob.johnson@example.com', 'HR', 'sample_public_key_3')
ON CONFLICT (voter_id) DO NOTHING;

-- Grant necessary permissions
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO voting_user;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO voting_user;
GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA public TO voting_user;

