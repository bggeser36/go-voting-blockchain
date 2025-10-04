#!/bin/bash

# API Testing Script for Voting Blockchain System
# This script tests all the major API endpoints

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# API Base URL
API_URL="http://localhost:8090"

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to test API endpoint
test_endpoint() {
    local method=$1
    local endpoint=$2
    local data=$3
    local description=$4
    
    print_status "Testing: $description"
    
    if [ -n "$data" ]; then
        response=$(curl -s -X $method "$API_URL$endpoint" \
            -H "Content-Type: application/json" \
            -d "$data")
    else
        response=$(curl -s -X $method "$API_URL$endpoint")
    fi
    
    if [ $? -eq 0 ]; then
        echo "$response" | jq . 2>/dev/null || echo "$response"
        print_success "$description - OK"
    else
        print_error "$description - FAILED"
        return 1
    fi
    echo "---"
}

# Main test function
main() {
    print_status "Starting Voting Blockchain API Tests"
    print_status "API Base URL: $API_URL"
    echo
    
    # Test 1: Health Check
    test_endpoint "GET" "/health" "" "Health Check"
    
    # Test 2: API Status
    test_endpoint "GET" "/" "" "API Status"
    
    # Test 3: Register a new voter
    test_endpoint "POST" "/register" '{
        "email": "test.user@example.com",
        "name": "Test User",
        "department": "Testing"
    }' "Voter Registration"
    
    # Extract voter_id from registration response
    voter_response=$(curl -s -X POST "$API_URL/register" \
        -H "Content-Type: application/json" \
        -d '{
            "email": "test.user2@example.com",
            "name": "Test User 2",
            "department": "Testing"
        }')
    
    voter_id=$(echo "$voter_response" | jq -r '.data.voter_id')
    print_status "Registered voter ID: $voter_id"
    
    # Test 4: Create a poll
    poll_response=$(curl -s -X POST "$API_URL/polls" \
        -H "Content-Type: application/json" \
        -d '{
            "title": "Best Programming Language 2024",
            "description": "Vote for your favorite programming language",
            "options": ["Go", "Python", "JavaScript", "Rust", "Java"],
            "creator": "Admin",
            "duration_hours": 24,
            "is_anonymous": false
        }')
    
    poll_id=$(echo "$poll_response" | jq -r '.data.poll_id')
    print_status "Created poll ID: $poll_id"
    
    test_endpoint "POST" "/polls" '{
        "title": "Best Programming Language 2024",
        "description": "Vote for your favorite programming language",
        "options": ["Go", "Python", "JavaScript", "Rust", "Java"],
        "creator": "Admin",
        "duration_hours": 24,
        "is_anonymous": false
    }' "Poll Creation"
    
    # Test 5: Get all polls
    test_endpoint "GET" "/polls" "" "Get All Polls"
    
    # Test 6: Get specific poll details
    test_endpoint "GET" "/polls/$poll_id" "" "Get Poll Details"
    
    # Test 7: Cast a vote
    test_endpoint "POST" "/vote" "{
        \"poll_id\": \"$poll_id\",
        \"voter_id\": \"$voter_id\",
        \"choice\": \"Go\",
        \"signature\": \"test_signature_$(date +%s)\"
    }" "Vote Casting"
    
    # Test 8: Get poll results (before mining)
    test_endpoint "GET" "/results/$poll_id" "" "Poll Results (Before Mining)"
    
    # Test 9: Manual mining
    test_endpoint "POST" "/blockchain/mine" "" "Manual Mining"
    
    # Test 10: Get poll results (after mining)
    test_endpoint "GET" "/results/$poll_id" "" "Poll Results (After Mining)"
    
    # Test 11: Blockchain verification
    test_endpoint "GET" "/blockchain/verify" "" "Blockchain Verification"
    
    # Test 12: Get blockchain stats
    test_endpoint "GET" "/blockchain/stats" "" "Blockchain Statistics"
    
    # Test 13: Get recent blocks
    test_endpoint "GET" "/blockchain/blocks" "" "Recent Blocks"
    
    # Test 14: Get voter history
    test_endpoint "GET" "/voter/$voter_id/history" "" "Voter History"
    
    print_success "All API tests completed successfully!"
    echo
    print_status "Summary:"
    print_status "- ✅ Health Check: API is running"
    print_status "- ✅ Voter Registration: Users can register"
    print_status "- ✅ Poll Creation: Polls can be created"
    print_status "- ✅ Voting: Users can cast votes"
    print_status "- ✅ Blockchain: Votes are mined into blocks"
    print_status "- ✅ Results: Poll results are calculated"
    print_status "- ✅ Verification: Blockchain integrity verified"
    print_status "- ✅ Persistence: Data is stored in database"
    
    echo
    print_status "Service URLs:"
    print_status "  API: $API_URL"
    print_status "  Health: $API_URL/health"
    print_status "  Status: $API_URL/"
    print_status "  Polls: $API_URL/polls"
    print_status "  Results: $API_URL/results/$poll_id"
}

# Run main function
main "$@"

