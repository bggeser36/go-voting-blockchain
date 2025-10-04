# 🗳️ Go Voting Blockchain

[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Docker](https://img.shields.io/badge/docker-ready-blue.svg)](https://www.docker.com/)
[![Railway Ready](https://img.shields.io/badge/railway-ready-0B0D0E.svg)](https://railway.app/)

A production-ready blockchain-based voting system built with Go, featuring secure cryptographic operations, real-time results, and seamless deployment capabilities.

## Features

- **Secure Blockchain**: Proof-of-Work consensus with adjustable difficulty
- **Digital Signatures**: RSA 2048-bit encryption for vote authentication
- **RESTful API**: High-performance Gin framework
- **Persistence**: PostgreSQL and Redis support for data durability
- **Scalable**: Optimized for cloud deployment
- **Real-time Results**: Live vote tallying and statistics

## Architecture

```
┌──────────────────────────────────────────────────────┐
│                    Client Applications                │
│                 (Web, Mobile, CLI)                    │
└──────────────────────────────────────────────────────┘
                            │
                            ▼
┌──────────────────────────────────────────────────────┐
│                    REST API (Gin)                     │
│                    Port: 8080                         │
└──────────────────────────────────────────────────────┘
                            │
┌──────────────────────────────────────────────────────┐
│                  Blockchain Core                      │
│         (Blocks, Mining, Verification)                │
└──────────────────────────────────────────────────────┘
                            │
┌──────────────────────────────────────────────────────┐
│                   Persistence Layer                   │
│          PostgreSQL (Primary) + Redis (Cache)         │
└──────────────────────────────────────────────────────┘
```

## 🚀 Quick Start

### Local Development

1. **Clone the repository**
```bash
git clone https://github.com/tolstoyjustin/go-voting-blockchain.git
cd go-voting-blockchain
```

2. **Install dependencies**
```bash
go mod download
```

3. **Set up environment**
```bash
cp .env.example .env
# Edit .env with your configuration
```

4. **Run the application**
```bash
go run cmd/api/main.go
```

The API will be available at `http://localhost:8080`

### Docker Deployment

```bash
docker build -t voting-blockchain .
docker run -p 8080:8080 voting-blockchain
```

## ☁️ Railway Deployment

### One-Click Deploy

[![Deploy on Railway](https://railway.app/button.svg)](https://railway.app/new/template)

**Note**: This will create a new Railway project with your voting blockchain system.

### Manual Deployment

1. **Create a new Railway project**
```bash
railway login
railway init
```

2. **Add PostgreSQL database**
```bash
railway add postgresql
```

3. **Add Redis (optional)**
```bash
railway add redis
```

4. **Deploy**
```bash
railway up
```

Railway will automatically:
- Detect the Go application
- Install dependencies
- Build the application
- Set up environment variables
- Configure health checks
- Deploy with automatic SSL

### Environment Variables for Railway

Railway automatically provides:
- `PORT`: Server port (auto-assigned)
- `DATABASE_URL`: PostgreSQL connection string
- `REDIS_URL`: Redis connection string (if added)

## API Endpoints

### Status
- `GET /` - API status and blockchain info
- `GET /health` - Health check endpoint

### Voter Management
- `POST /register` - Register a new voter
- `GET /voter/:voter_id/history` - Get voter's voting history

### Poll Management
- `POST /polls` - Create a new poll
- `GET /polls` - List all polls
- `GET /polls/:poll_id` - Get poll details
- `GET /results/:poll_id` - Get poll results

### Voting
- `POST /vote` - Cast a vote

### Blockchain
- `GET /blockchain/verify` - Verify blockchain integrity
- `GET /blockchain/blocks` - Get recent blocks
- `GET /blockchain/stats` - Get blockchain statistics
- `POST /blockchain/mine` - Manually mine pending votes

## API Usage Examples

### Register a Voter
```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "voter@example.com",
    "name": "John Doe",
    "department": "Engineering"
  }'
```

### Create a Poll
```bash
curl -X POST http://localhost:8080/polls \
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

### Cast a Vote
```bash
curl -X POST http://localhost:8080/vote \
  -H "Content-Type: application/json" \
  -d '{
    "poll_id": "poll-uuid",
    "voter_id": "voter-id",
    "choice": "Go",
    "signature": "signature-string"
  }'
```

### Get Poll Results
```bash
curl http://localhost:8080/results/poll-uuid
```

### Verify Blockchain
```bash
curl http://localhost:8080/blockchain/verify
```

## Project Structure

```
go-voting-blockchain/
├── cmd/
│   └── api/
│       └── main.go           # Application entry point
├── internal/
│   ├── blockchain/
│   │   └── blockchain.go     # Core blockchain logic
│   ├── crypto/
│   │   └── crypto.go         # Cryptographic operations
│   ├── handlers/
│   │   └── handlers.go       # API request handlers
│   ├── models/
│   │   └── models.go         # Data models
│   └── persistence/
│       └── persistence.go    # Database operations
├── configs/
│   └── config.go             # Configuration management
├── scripts/
│   └── deploy.sh             # Deployment scripts
├── Dockerfile                # Docker configuration
├── railway.toml              # Railway configuration
├── go.mod                    # Go modules
├── go.sum                    # Module checksums
└── README.md                 # Documentation
```

## Performance Optimization

### For Railway Deployment

1. **Database Indexing**: Automatic indexes on frequently queried fields
2. **Redis Caching**: Optional caching layer for improved performance
3. **Connection Pooling**: Optimized database connection management
4. **Concurrent Mining**: Asynchronous block mining
5. **Health Checks**: Automatic monitoring and restart policies

### Scaling Considerations

- **Horizontal Scaling**: Deploy multiple instances behind a load balancer
- **Database Replication**: Use Railway's database scaling features
- **CDN Integration**: Serve static content through Railway's CDN
- **Auto-scaling**: Configure based on CPU/memory usage

## Security Features

- **RSA 2048-bit encryption** for digital signatures
- **SHA-256 hashing** for block integrity
- **Proof-of-Work** consensus mechanism
- **CORS configuration** for API security
- **Environment variable** management for secrets
- **SQL injection prevention** with parameterized queries
- **Rate limiting** support (can be added via middleware)

## Monitoring

Railway provides built-in monitoring for:
- Application logs
- Resource usage (CPU, Memory)
- Request metrics
- Error tracking
- Health status

Access monitoring at: `https://railway.app/project/[project-id]/deployments`

## Troubleshooting

### Common Issues

1. **Port binding error**
   - Railway automatically assigns PORT
   - Don't hardcode port numbers

2. **Database connection failed**
   - Check DATABASE_URL is set
   - Verify PostgreSQL service is running

3. **Mining too slow**
   - Adjust MINING_DIFFICULTY in environment
   - Default is 3 for development

4. **Memory issues**
   - Increase Railway instance size
   - Implement pagination for large datasets

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details on how to get started.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Contributors

See [CONTRIBUTORS.md](CONTRIBUTORS.md) for a list of all contributors.

## License

MIT License - see LICENSE file for details

## 📞 Support

For issues and questions:
- 🐛 **Bug Reports**: [GitHub Issues](https://github.com/tolstoyjustin/go-voting-blockchain/issues)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/tolstoyjustin/go-voting-blockchain/discussions)
- 📧 **Contact**: tolstoyjustin@gmail.com
- 🚂 **Railway Help**: [Railway Community](https://railway.app/help)

## Deployment Checklist

- [ ] Set up PostgreSQL database
- [ ] Configure environment variables
- [ ] Set appropriate mining difficulty
- [ ] Configure CORS for your frontend
- [ ] Set up monitoring alerts
- [ ] Test all API endpoints
- [ ] Verify blockchain persistence
- [ ] Check health endpoint
- [ ] Review security settings
- [ ] Deploy and verify