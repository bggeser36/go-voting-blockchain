# 🗳️ Go Voting Blockchain

[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Docker](https://img.shields.io/badge/docker-ready-blue.svg)](https://www.docker.com/)
[![Railway Ready](https://img.shields.io/badge/railway-ready-0B0D0E.svg)](https://railway.app/)

A production-ready blockchain-based voting system built with Go, featuring enterprise-grade security, JWT authentication, role-based access control, and comprehensive monitoring capabilities.

## ✨ Features

### 🔐 Security & Authentication (Phase 1)
- **JWT Authentication**: Secure token-based authentication system
- **Role-Based Access Control**: Admin and voter roles with proper permissions
- **Rate Limiting**: Sliding window rate limiting with multiple tiers
- **Input Validation**: Comprehensive validation and sanitization
- **Error Handling**: Professional error handling and recovery middleware
- **Request Logging**: Complete API monitoring and logging system
- **Digital Signatures**: RSA 2048-bit encryption for vote authentication

### ⛓️ Blockchain Core
- **Secure Blockchain**: Proof-of-Work consensus with adjustable difficulty
- **Signature Verification**: Cryptographic proof of vote authenticity
- **Block Integrity**: SHA-256 hashing for tamper-proof blocks
- **Mining System**: Configurable mining difficulty and automatic block creation

### 🚀 Performance & Scalability
- **RESTful API**: High-performance Gin framework
- **Persistence**: PostgreSQL and Redis support for data durability
- **Connection Pooling**: Optimized database connection management
- **Concurrent Processing**: Asynchronous operations for better performance
- **Real-time Results**: Live vote tallying and statistics

### 🛠️ Developer Experience
- **Comprehensive Testing**: Automated security and integration tests
- **Development Tools**: Complete development environment setup
- **CI/CD Pipeline**: GitHub Actions for automated testing and deployment
- **Documentation**: Extensive documentation and guides

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

### 🔓 Public Endpoints
- `GET /` - API status and blockchain info
- `GET /health` - Health check endpoint
- `POST /register` - Register a new voter
- `GET /polls` - List all polls
- `GET /polls/:poll_id` - Get poll details
- `GET /results/:poll_id` - Get poll results
- `GET /blockchain/verify` - Verify blockchain integrity
- `GET /blockchain/blocks` - Get recent blocks
- `GET /blockchain/stats` - Get blockchain statistics

### 🔐 Authentication Endpoints
- `POST /auth/login` - Admin login
- `POST /auth/voter-login` - Voter login
- `POST /auth/refresh` - Refresh JWT token
- `GET /auth/me` - Get current user info

### 👤 Authenticated Endpoints (Voter)
- `POST /vote` - Cast a vote
- `GET /voter/:voter_id/history` - Get voter's voting history

### 👑 Admin Endpoints
- `POST /admin/polls` - Create a new poll
- `POST /admin/blockchain/mine` - Manually mine pending votes

## API Usage Examples

### 🔐 Authentication

#### Admin Login
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'
```

#### Voter Login
```bash
curl -X POST http://localhost:8080/auth/voter-login \
  -H "Content-Type: application/json" \
  -d '{
    "voter_id": "voter-uuid",
    "private_key": "-----BEGIN RSA PRIVATE KEY-----..."
  }'
```

### 👤 Voter Operations

#### Register a Voter
```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "voter@example.com",
    "name": "John Doe",
    "department": "Engineering"
  }'
```

#### Cast a Vote (with Authentication)
```bash
curl -X POST http://localhost:8080/vote \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "poll_id": "poll-uuid",
    "voter_id": "voter-id",
    "choice": "Go",
    "signature": "signature-string"
  }'
```

### 👑 Admin Operations

#### Create a Poll (Admin Only)
```bash
curl -X POST http://localhost:8080/admin/polls \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN" \
  -d '{
    "title": "Best Programming Language",
    "description": "Vote for your favorite language",
    "options": ["Go", "Python", "JavaScript", "Rust"],
    "creator": "Admin",
    "duration_hours": 24,
    "is_anonymous": false
  }'
```

#### Mine Pending Votes (Admin Only)
```bash
curl -X POST http://localhost:8080/admin/blockchain/mine \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

### 📊 Public Operations

#### Get Poll Results
```bash
curl http://localhost:8080/results/poll-uuid
```

#### Verify Blockchain
```bash
curl http://localhost:8080/blockchain/verify
```

## Project Structure

```
go-voting-blockchain/
├── cmd/
│   └── api/
│       └── main.go                    # Application entry point
├── internal/
│   ├── auth/                          # Authentication system
│   │   ├── admin.go                   # Admin authentication
│   │   └── jwt.go                     # JWT token management
│   ├── blockchain/
│   │   └── blockchain.go              # Core blockchain logic
│   ├── crypto/
│   │   └── crypto.go                  # Cryptographic operations
│   ├── handlers/
│   │   ├── auth.go                    # Authentication handlers
│   │   └── handlers.go                # API request handlers
│   ├── middleware/                    # Middleware components
│   │   ├── auth.go                    # Authentication middleware
│   │   ├── error.go                   # Error handling middleware
│   │   ├── logging.go                 # Request logging middleware
│   │   └── ratelimit.go               # Rate limiting middleware
│   ├── models/
│   │   └── models.go                  # Data models
│   ├── persistence/
│   │   └── persistence.go             # Database operations
│   └── validation/
│       └── validator.go               # Input validation
├── tests/
│   └── security_test.go               # Security integration tests
├── scripts/
│   ├── dev-setup.sh                   # Development setup script
│   └── test-api.sh                    # API testing script
├── .github/
│   └── workflows/
│       ├── ci.yml                     # Continuous integration
│       └── development.yml            # Development workflow
├── configs/
├── Dockerfile                         # Docker configuration
├── docker-compose.yml                 # Docker Compose setup
├── railway.toml                       # Railway configuration
├── config.development.env             # Development environment
├── go.mod                             # Go modules
├── go.sum                             # Module checksums
├── README.md                          # Main documentation
├── DEVELOPMENT.md                     # Development guide
├── ROADMAP.md                         # Development roadmap
├── CONTRIBUTING.md                    # Contribution guidelines
├── PHASE1_TEST_REPORT.md              # Phase 1 test results
└── LICENSE                            # MIT License
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

## 🔒 Security Features

### Phase 1 Security Enhancements
- **JWT Authentication**: Secure token-based authentication with refresh tokens
- **Role-Based Access Control**: Admin and voter roles with proper permissions
- **Rate Limiting**: Sliding window rate limiting with multiple tiers (strict, moderate, generous)
- **Input Validation**: Comprehensive validation and sanitization of all inputs
- **Error Handling**: Professional error handling with proper HTTP status codes
- **Request Logging**: Complete API monitoring and audit logging
- **Signature Verification**: Cryptographic proof of vote authenticity

### Core Security Features
- **RSA 2048-bit encryption** for digital signatures
- **SHA-256 hashing** for block integrity
- **Proof-of-Work** consensus mechanism
- **CORS configuration** for API security
- **Environment variable** management for secrets
- **SQL injection prevention** with parameterized queries
- **Private key ownership verification** for voter authentication

### Security Testing
- **Automated Security Tests**: Comprehensive test suite for all security features
- **Integration Testing**: End-to-end testing of authentication flows
- **Rate Limiting Tests**: Validation of rate limiting functionality
- **Signature Verification Tests**: Cryptographic security validation

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

## 🎯 Development Status

### ✅ Phase 1: Security & Authentication (COMPLETED)
- **JWT Authentication System** - Complete token-based authentication
- **Role-Based Access Control** - Admin and voter roles with proper permissions
- **Rate Limiting** - Sliding window rate limiting with multiple tiers
- **Input Validation** - Comprehensive validation and sanitization
- **Error Handling** - Professional error handling and recovery middleware
- **Request Logging** - Complete API monitoring and logging system
- **Security Tests** - Comprehensive automated security testing

### 🔄 Phase 2: Data Persistence & Reliability (NEXT)
- Database optimization and connection pooling
- Data integrity and backup systems
- Advanced persistence strategies
- Performance monitoring and optimization

### 📋 Future Phases
- **Phase 3**: Performance Optimization
- **Phase 4**: Advanced Features
- **Phase 5**: User Experience Enhancement
- **Phase 6**: Integration & APIs
- **Phase 7**: Compliance & Auditing

See [ROADMAP.md](ROADMAP.md) for detailed development plans.

## 🚀 Deployment Checklist

### Pre-Deployment
- [ ] Set up PostgreSQL database
- [ ] Configure environment variables
- [ ] Set appropriate mining difficulty
- [ ] Configure JWT secret key
- [ ] Set up Redis (optional, for caching)
- [ ] Configure CORS for your frontend

### Security Configuration
- [ ] Change default admin password
- [ ] Configure rate limiting thresholds
- [ ] Set up proper logging levels
- [ ] Configure error handling policies
- [ ] Review authentication settings

### Testing & Verification
- [ ] Test all API endpoints
- [ ] Verify authentication flows
- [ ] Test rate limiting functionality
- [ ] Verify blockchain persistence
- [ ] Check health endpoint
- [ ] Run security test suite
- [ ] Validate signature verification

### Deployment
- [ ] Deploy to production environment
- [ ] Verify all services are running
- [ ] Monitor initial performance
- [ ] Set up monitoring alerts
- [ ] Document deployment configuration