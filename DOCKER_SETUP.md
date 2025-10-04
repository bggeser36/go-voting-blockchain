# Docker Setup Guide for Voting Blockchain System

This guide will help you set up and run the Voting Blockchain System using Docker for local testing.

## Prerequisites

- Docker Desktop installed and running
- Docker Compose (included with Docker Desktop)
- Ports 5433, 6380, 8090, and 5051 available on your system

## Quick Start

### 1. Start the Services

```bash
# Make the setup script executable
chmod +x scripts/docker-setup.sh

# Start all services
./scripts/docker-setup.sh start

# Or start with database management tools
./scripts/docker-setup.sh --with-tools start
```

### 2. Verify the Setup

```bash
# Check service status
./scripts/docker-setup.sh status

# Test API endpoints
./scripts/docker-setup.sh test
```

### 3. Access the Services

- **API**: http://localhost:8090
- **Health Check**: http://localhost:8090/health
- **API Status**: http://localhost:8090/
- **pgAdmin** (if started with tools): http://localhost:5051

## Service Details

### Core Services

| Service | Container Name | Port | Description |
|---------|----------------|------|-------------|
| voting-api | voting-blockchain-api | 8090 | Main API server |
| voting-postgres | voting-blockchain-postgres | 5433 | PostgreSQL database |
| voting-redis | voting-blockchain-redis | 6380 | Redis cache |

### Optional Services

| Service | Container Name | Port | Description |
|---------|----------------|------|-------------|
| voting-pgadmin | voting-blockchain-pgadmin | 5051 | Database management UI |

## Database Connection

- **Host**: localhost:5433
- **Database**: voting_blockchain
- **Username**: voting_user
- **Password**: voting_password123

## Redis Connection

- **Host**: localhost:6380
- **Password**: voting_redis_pass123

## Management Commands

```bash
# Start services
./scripts/docker-setup.sh start

# Stop services
./scripts/docker-setup.sh stop

# Restart services
./scripts/docker-setup.sh restart

# Show service status
./scripts/docker-setup.sh status

# Test API
./scripts/docker-setup.sh test

# View logs
./scripts/docker-setup.sh logs
./scripts/docker-setup.sh logs voting-api

# Clean up everything
./scripts/docker-setup.sh cleanup
```

## API Testing Examples

### 1. Check API Status

```bash
curl http://localhost:8090/
```

### 2. Register a Voter

```bash
curl -X POST http://localhost:8090/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "name": "John Doe",
    "department": "Engineering"
  }'
```

### 3. Create a Poll

```bash
curl -X POST http://localhost:8090/polls \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Best Programming Language",
    "description": "Vote for your favorite language",
    "options": ["Go", "Python", "JavaScript", "Rust"],
    "creator": "Admin",
    "duration_hours": 24,
    "is_anonymous": false
  }'
```

### 4. Cast a Vote

```bash
curl -X POST http://localhost:8090/vote \
  -H "Content-Type: application/json" \
  -d '{
    "poll_id": "YOUR_POLL_ID",
    "voter_id": "YOUR_VOTER_ID",
    "choice": "Go",
    "signature": "YOUR_SIGNATURE"
  }'
```

### 5. Get Poll Results

```bash
curl http://localhost:8090/results/YOUR_POLL_ID
```

### 6. Verify Blockchain

```bash
curl http://localhost:8090/blockchain/verify
```

## Troubleshooting

### Port Conflicts

If you encounter port conflicts, the script will warn you. You can:

1. Stop conflicting services
2. Modify ports in `docker-compose.yml`
3. Continue anyway (not recommended)

### Database Connection Issues

```bash
# Check if PostgreSQL is running
docker-compose exec voting-postgres pg_isready -U voting_user -d voting_blockchain

# View PostgreSQL logs
./scripts/docker-setup.sh logs voting-postgres
```

### Redis Connection Issues

```bash
# Check if Redis is running
docker-compose exec voting-redis redis-cli ping

# View Redis logs
./scripts/docker-setup.sh logs voting-redis
```

### API Issues

```bash
# Check API health
curl http://localhost:8090/health

# View API logs
./scripts/docker-setup.sh logs voting-api
```

### Clean Restart

If you encounter persistent issues:

```bash
# Stop and clean everything
./scripts/docker-setup.sh cleanup

# Start fresh
./scripts/docker-setup.sh start
```

## Development Mode

For development with live reloading:

```bash
# Start only database services
docker-compose up -d voting-postgres voting-redis

# Run API locally with hot reload
go run cmd/api/main.go
```

## Production Considerations

This Docker setup is designed for local testing. For production:

1. Use environment variables for sensitive data
2. Enable SSL/TLS
3. Use proper secrets management
4. Configure proper logging and monitoring
5. Set up backup strategies for databases
6. Use production-grade database and Redis configurations

## Network Isolation

The services run in a dedicated Docker network (`voting-blockchain-network`) to avoid conflicts with your existing containers.

## Data Persistence

Data is persisted in Docker volumes:
- `voting_postgres_data`: PostgreSQL data
- `voting_redis_data`: Redis data
- `voting_pgadmin_data`: pgAdmin configuration

To completely reset the system:
```bash
./scripts/docker-setup.sh cleanup
```

