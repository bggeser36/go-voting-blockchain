#!/bin/bash

# Docker Setup Script for Voting Blockchain System
# This script helps manage the Docker environment for local testing

set -e

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

# Function to check if Docker is running
check_docker() {
    if ! docker info > /dev/null 2>&1; then
        print_error "Docker is not running. Please start Docker Desktop and try again."
        exit 1
    fi
    print_success "Docker is running"
}

# Function to check if ports are available
check_ports() {
    local ports=("5433" "6380" "8090" "5051")
    local conflicts=()
    
    for port in "${ports[@]}"; do
        if lsof -i :$port > /dev/null 2>&1; then
            conflicts+=($port)
        fi
    done
    
    if [ ${#conflicts[@]} -gt 0 ]; then
        print_warning "Port conflicts detected on ports: ${conflicts[*]}"
        print_warning "Please ensure these ports are not in use by other services"
        read -p "Continue anyway? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
    else
        print_success "All required ports are available"
    fi
}

# Function to build and start services
start_services() {
    print_status "Building and starting Voting Blockchain services..."
    
    # Build the API image
    print_status "Building API image..."
    docker-compose build voting-api
    
    # Start core services
    print_status "Starting PostgreSQL and Redis..."
    docker-compose up -d voting-postgres voting-redis
    
    # Wait for services to be healthy
    print_status "Waiting for database services to be ready..."
    docker-compose exec voting-postgres pg_isready -U voting_user -d voting_blockchain
    docker-compose exec voting-redis redis-cli ping
    
    # Start the API
    print_status "Starting API service..."
    docker-compose up -d voting-api
    
    print_success "All services started successfully!"
}

# Function to show service status
show_status() {
    print_status "Service Status:"
    docker-compose ps
    
    echo
    print_status "Service URLs:"
    echo "  API: http://localhost:8090"
    echo "  Health Check: http://localhost:8090/health"
    echo "  API Status: http://localhost:8090/"
    echo "  pgAdmin: http://localhost:5051 (if started with --with-tools)"
    echo
    print_status "Database Connection:"
    echo "  Host: localhost:5433"
    echo "  Database: voting_blockchain"
    echo "  Username: voting_user"
    echo "  Password: voting_password123"
    echo
    print_status "Redis Connection:"
    echo "  Host: localhost:6380"
    echo "  Password: voting_redis_pass123"
}

# Function to test API endpoints
test_api() {
    print_status "Testing API endpoints..."
    
    # Test health endpoint
    if curl -s http://localhost:8090/health > /dev/null; then
        print_success "Health endpoint is responding"
    else
        print_error "Health endpoint is not responding"
        return 1
    fi
    
    # Test status endpoint
    if curl -s http://localhost:8090/ > /dev/null; then
        print_success "API status endpoint is responding"
    else
        print_error "API status endpoint is not responding"
        return 1
    fi
    
    print_success "API is working correctly!"
}

# Function to show logs
show_logs() {
    local service=${1:-""}
    if [ -z "$service" ]; then
        docker-compose logs -f
    else
        docker-compose logs -f "$service"
    fi
}

# Function to stop services
stop_services() {
    print_status "Stopping all services..."
    docker-compose down
    print_success "All services stopped"
}

# Function to clean up everything
cleanup() {
    print_status "Cleaning up Docker resources..."
    docker-compose down -v --remove-orphans
    docker system prune -f
    print_success "Cleanup completed"
}

# Function to show help
show_help() {
    echo "Voting Blockchain Docker Management Script"
    echo
    echo "Usage: $0 [COMMAND]"
    echo
    echo "Commands:"
    echo "  start       Build and start all services"
    echo "  stop        Stop all services"
    echo "  restart     Restart all services"
    echo "  status      Show service status and URLs"
    echo "  test        Test API endpoints"
    echo "  logs [svc]  Show logs (optionally for specific service)"
    echo "  cleanup     Stop services and clean up volumes/images"
    echo "  help        Show this help message"
    echo
    echo "Examples:"
    echo "  $0 start                    # Start all services"
    echo "  $0 logs voting-api          # Show API logs"
    echo "  $0 test                     # Test API functionality"
    echo "  $0 --with-tools start       # Start with pgAdmin"
}

# Main script logic
main() {
    local command=${1:-"start"}
    local with_tools=false
    
    # Parse arguments
    for arg in "$@"; do
        case $arg in
            --with-tools)
                with_tools=true
                shift
                ;;
        esac
    done
    
    case $command in
        start)
            check_docker
            check_ports
            start_services
            if [ "$with_tools" = true ]; then
                print_status "Starting pgAdmin..."
                docker-compose --profile tools up -d voting-pgadmin
            fi
            show_status
            ;;
        stop)
            stop_services
            ;;
        restart)
            stop_services
            sleep 2
            start_services
            show_status
            ;;
        status)
            show_status
            ;;
        test)
            test_api
            ;;
        logs)
            show_logs $2
            ;;
        cleanup)
            cleanup
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            print_error "Unknown command: $command"
            show_help
            exit 1
            ;;
    esac
}

# Run main function
main "$@"

