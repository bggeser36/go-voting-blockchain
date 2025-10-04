#!/bin/bash

# Development Environment Setup Script
# This script sets up the development environment for the Go Voting Blockchain project

set -e

echo "🚀 Setting up Go Voting Blockchain Development Environment..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

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

# Check if Go is installed
check_go() {
    print_status "Checking Go installation..."
    if command -v go &> /dev/null; then
        GO_VERSION=$(go version | cut -d' ' -f3)
        print_success "Go is installed: $GO_VERSION"
        
        # Check if Go version is 1.21 or higher
        GO_VERSION_NUM=$(echo $GO_VERSION | sed 's/go//')
        if [ "$(printf '%s\n' "1.21" "$GO_VERSION_NUM" | sort -V | head -n1)" = "1.21" ]; then
            print_success "Go version is compatible (1.21+)"
        else
            print_warning "Go version $GO_VERSION_NUM may not be compatible. Recommended: 1.21+"
        fi
    else
        print_error "Go is not installed. Please install Go 1.21 or higher."
        exit 1
    fi
}

# Check if Docker is installed
check_docker() {
    print_status "Checking Docker installation..."
    if command -v docker &> /dev/null; then
        DOCKER_VERSION=$(docker --version | cut -d' ' -f3 | cut -d',' -f1)
        print_success "Docker is installed: $DOCKER_VERSION"
        
        # Check if Docker is running
        if docker info &> /dev/null; then
            print_success "Docker is running"
        else
            print_error "Docker is not running. Please start Docker."
            exit 1
        fi
    else
        print_error "Docker is not installed. Please install Docker."
        exit 1
    fi
}

# Check if Docker Compose is installed
check_docker_compose() {
    print_status "Checking Docker Compose installation..."
    if command -v docker-compose &> /dev/null; then
        COMPOSE_VERSION=$(docker-compose --version | cut -d' ' -f3 | cut -d',' -f1)
        print_success "Docker Compose is installed: $COMPOSE_VERSION"
    else
        print_error "Docker Compose is not installed. Please install Docker Compose."
        exit 1
    fi
}

# Install Go dependencies
install_dependencies() {
    print_status "Installing Go dependencies..."
    go mod download
    go mod verify
    print_success "Go dependencies installed successfully"
}

# Setup environment file
setup_env() {
    print_status "Setting up environment configuration..."
    
    if [ ! -f .env ]; then
        if [ -f config.development.env ]; then
            cp config.development.env .env
            print_success "Created .env file from development template"
        else
            print_warning "Development config template not found. Creating basic .env file..."
            cat > .env << EOF
# Development Environment Configuration
PORT=8080
GIN_MODE=debug
DATABASE_URL=postgres://voting_user:voting_password123@localhost:5433/voting_blockchain?sslmode=disable
REDIS_URL=redis://:voting_redis_pass123@localhost:6380/0
MINING_DIFFICULTY=2
MINING_THRESHOLD=3
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:8080
LOG_LEVEL=debug
EOF
            print_success "Created basic .env file"
        fi
    else
        print_warning ".env file already exists. Skipping creation."
    fi
}

# Start development services
start_services() {
    print_status "Starting development services (PostgreSQL and Redis)..."
    
    # Stop any existing containers
    docker-compose down &> /dev/null || true
    
    # Start services in detached mode
    docker-compose up -d voting-postgres voting-redis
    
    # Wait for services to be ready
    print_status "Waiting for services to be ready..."
    sleep 10
    
    # Check if services are running
    if docker-compose ps | grep -q "Up"; then
        print_success "Development services started successfully"
        print_status "PostgreSQL: localhost:5433"
        print_status "Redis: localhost:6380"
    else
        print_error "Failed to start development services"
        docker-compose logs
        exit 1
    fi
}

# Run database migrations
run_migrations() {
    print_status "Running database migrations..."
    
    # Wait a bit more for PostgreSQL to be fully ready
    sleep 5
    
    # Run the initialization script
    if docker-compose exec -T voting-postgres psql -U voting_user -d voting_blockchain -f /docker-entrypoint-initdb.d/init-db.sql; then
        print_success "Database migrations completed"
    else
        print_warning "Database migrations may have failed. This is normal if tables already exist."
    fi
}

# Build the application
build_application() {
    print_status "Building the application..."
    
    if go build -o bin/voting-blockchain cmd/api/main.go; then
        print_success "Application built successfully"
    else
        print_error "Failed to build application"
        exit 1
    fi
}

# Run tests
run_tests() {
    print_status "Running tests..."
    
    if go test -v ./...; then
        print_success "All tests passed"
    else
        print_warning "Some tests failed or no tests found"
    fi
}

# Main setup function
main() {
    echo "========================================"
    echo "🔧 Go Voting Blockchain Dev Setup"
    echo "========================================"
    
    check_go
    check_docker
    check_docker_compose
    install_dependencies
    setup_env
    start_services
    run_migrations
    build_application
    run_tests
    
    echo ""
    echo "========================================"
    print_success "Development environment setup complete!"
    echo "========================================"
    echo ""
    echo "🚀 Next steps:"
    echo "  1. Start the application: go run cmd/api/main.go"
    echo "  2. Or run the binary: ./bin/voting-blockchain"
    echo "  3. API will be available at: http://localhost:8080"
    echo "  4. Health check: http://localhost:8080/health"
    echo ""
    echo "📊 Services:"
    echo "  • PostgreSQL: localhost:5433"
    echo "  • Redis: localhost:6380"
    echo "  • API: localhost:8080"
    echo ""
    echo "🛠️  Development commands:"
    echo "  • Run tests: go test ./..."
    echo "  • Format code: go fmt ./..."
    echo "  • Lint code: go vet ./..."
    echo "  • Stop services: docker-compose down"
    echo ""
    print_success "Happy coding! 🎉"
}

# Run main function
main "$@"
