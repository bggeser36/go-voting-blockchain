# 🛠️ Development Guide

## 🚀 Quick Start

### Prerequisites
- Go 1.21 or higher
- Docker and Docker Compose
- Git

### Setup Development Environment

1. **Clone and switch to development branch:**
   ```bash
   git clone https://github.com/Tolstoyj/go-voting-blockchain.git
   cd go-voting-blockchain
   git checkout development
   ```

2. **Run the development setup script:**
   ```bash
   ./scripts/dev-setup.sh
   ```

3. **Start development:**
   ```bash
   go run cmd/api/main.go
   ```

---

## 🏗️ Development Workflow

### Branch Strategy
```
main (production)
  └── development (integration)
      ├── feature/authentication
      ├── feature/rate-limiting
      ├── feature/signature-verification
      └── feature/ui-improvements
```

### Development Process

1. **Create feature branch from development:**
   ```bash
   git checkout development
   git pull origin development
   git checkout -b feature/your-feature-name
   ```

2. **Develop your feature:**
   - Write code following Go best practices
   - Add comprehensive tests
   - Update documentation
   - Follow the coding standards

3. **Test your changes:**
   ```bash
   go test ./...
   go vet ./...
   go fmt ./...
   ```

4. **Commit and push:**
   ```bash
   git add .
   git commit -m "feat: add your feature description"
   git push origin feature/your-feature-name
   ```

5. **Create Pull Request:**
   - Target: `development` branch
   - Fill out PR template
   - Request review from maintainers

---

## 🔧 Development Environment

### Local Services
- **API Server**: http://localhost:8080
- **PostgreSQL**: localhost:5433
- **Redis**: localhost:6380
- **Health Check**: http://localhost:8080/health

### Environment Configuration
Development uses `config.development.env` which provides:
- Lower mining difficulty (2) for faster testing
- Debug logging enabled
- Permissive CORS settings
- Local database connections

### Hot Reload
For development with hot reload:
```bash
# Install air for hot reloading
go install github.com/cosmtrek/air@latest

# Run with hot reload
air
```

---

## 🧪 Testing

### Running Tests
```bash
# All tests
go test ./...

# With coverage
go test -cover ./...

# Integration tests
go test -tags=integration ./...

# Race condition detection
go test -race ./...

# Verbose output
go test -v ./...
```

### Test Structure
```
internal/
├── blockchain/
│   ├── blockchain.go
│   └── blockchain_test.go
├── handlers/
│   ├── handlers.go
│   └── handlers_test.go
└── crypto/
    ├── crypto.go
    └── crypto_test.go
```

### Test Categories
- **Unit Tests**: Individual function testing
- **Integration Tests**: Component interaction testing
- **Security Tests**: Authentication and authorization
- **Performance Tests**: Load and stress testing

---

## 📝 Coding Standards

### Go Code Style
```go
// ✅ Good
func (bc *Blockchain) RegisterVoter(voter *models.Voter) error {
    bc.mu.Lock()
    defer bc.mu.Unlock()
    
    if voter == nil {
        return fmt.Errorf("voter cannot be nil")
    }
    
    // Implementation
    return nil
}

// ❌ Bad
func RegisterVoter(bc *Blockchain, voter *models.Voter) error {
    // No validation, no proper error handling
    bc.VoterRegistry[voter.VoterID] = voter
    return nil
}
```

### Documentation
```go
// RegisterVoter registers a new voter in the blockchain system.
// It validates the voter data and adds the voter to the registry.
// Returns an error if the voter is already registered or if validation fails.
//
// Example:
//   voter := &models.Voter{...}
//   err := bc.RegisterVoter(voter)
//   if err != nil {
//       log.Printf("Registration failed: %v", err)
//   }
func (bc *Blockchain) RegisterVoter(voter *models.Voter) error {
    // Implementation
}
```

### Error Handling
```go
// ✅ Structured errors
var (
    ErrVoterNotFound    = errors.New("voter not found")
    ErrPollNotActive    = errors.New("poll is not active")
    ErrInvalidSignature = errors.New("invalid signature")
)

// ✅ Error wrapping
if err != nil {
    return fmt.Errorf("failed to register voter %s: %w", voterID, err)
}
```

---

## 🔒 Security Development

### Authentication Implementation
```go
// JWT Middleware
func AuthMiddleware() gin.HandlerFunc {
    return gin.HandlerFunc(func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
            c.Abort()
            return
        }
        
        // Validate JWT token
        claims, err := ValidateJWT(token)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
            c.Abort()
            return
        }
        
        c.Set("user_id", claims.UserID)
        c.Next()
    })
}
```

### Rate Limiting
```go
// Rate limiting middleware
func RateLimitMiddleware() gin.HandlerFunc {
    store := redis.NewStore(redisClient, time.Minute)
    limiter := ginrate.Limit(store, ginrate.Rate{
        Period: time.Minute,
        Limit:  100, // requests per minute
    })
    
    return limiter.Middleware()
}
```

### Input Validation
```go
// Enhanced input validation
type VoteRequest struct {
    PollID    string `json:"poll_id" binding:"required,uuid"`
    VoterID   string `json:"voter_id" binding:"required,min=1,max=64"`
    Choice    string `json:"choice" binding:"required,min=1,max=255"`
    Signature string `json:"signature" binding:"required"`
}

// Custom validation
func (v *VoteRequest) Validate() error {
    if len(v.Choice) > 255 {
        return fmt.Errorf("choice too long")
    }
    
    // Additional validation logic
    return nil
}
```

---

## 🐛 Debugging

### Logging
```go
import "github.com/sirupsen/logrus"

var log = logrus.New()

func init() {
    log.SetFormatter(&logrus.JSONFormatter{})
    log.SetLevel(logrus.DebugLevel)
}

// Usage
log.WithFields(logrus.Fields{
    "voter_id": voterID,
    "poll_id":  pollID,
}).Info("Vote cast successfully")
```

### Debugging Tools
```bash
# Delve debugger
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug with Delve
dlv debug cmd/api/main.go

# Profiling
go tool pprof http://localhost:8080/debug/pprof/profile
```

### Database Debugging
```bash
# Connect to development database
docker exec -it voting-blockchain-postgres psql -U voting_user -d voting_blockchain

# Check tables
\dt

# Check data
SELECT * FROM voters LIMIT 5;
SELECT * FROM polls LIMIT 5;
SELECT * FROM votes LIMIT 5;
```

---

## 📊 Performance Development

### Benchmarking
```go
// Benchmark example
func BenchmarkCastVote(b *testing.B) {
    bc := blockchain.NewBlockchain(2)
    voter := &models.Voter{VoterID: "test"}
    poll := &models.Poll{PollID: "test"}
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        vote := &models.Vote{
            PollID:  poll.PollID,
            VoterID: voter.VoterID,
            Choice:  "option1",
        }
        bc.CastVote(vote)
    }
}

// Run benchmarks
go test -bench=. ./...
```

### Profiling
```go
// Add profiling endpoints
import _ "net/http/pprof"

func main() {
    // Enable profiling in development
    if os.Getenv("GIN_MODE") == "debug" {
        go func() {
            log.Println(http.ListenAndServe("localhost:6060", nil))
        }()
    }
}
```

---

## 🚀 Deployment

### Local Testing
```bash
# Build and test locally
go build -o bin/voting-blockchain cmd/api/main.go
./bin/voting-blockchain

# Test with Docker
docker build -t voting-blockchain:dev .
docker run -p 8080:8080 voting-blockchain:dev
```

### Development Deployment
```bash
# Deploy to development environment
git push origin development

# Monitor deployment
gh run list --branch development
```

---

## 📚 Resources

### Documentation
- [Go Best Practices](https://github.com/golang/go/wiki/CodeReviewComments)
- [Gin Framework](https://gin-gonic.com/docs/)
- [PostgreSQL Go Driver](https://github.com/lib/pq)
- [Redis Go Client](https://github.com/redis/go-redis)

### Tools
- **Air**: Hot reload for Go development
- **Delve**: Go debugger
- **golangci-lint**: Comprehensive linter
- **gosec**: Security analyzer
- **goimports**: Import management

### IDE Setup
- **VS Code**: Go extension pack
- **GoLand**: JetBrains IDE
- **Vim/Neovim**: vim-go plugin

---

## 🤝 Contributing

### Before Submitting
- [ ] Code follows Go best practices
- [ ] All tests pass
- [ ] Code is properly formatted (`go fmt`)
- [ ] No linting errors (`go vet`, `golangci-lint`)
- [ ] Documentation updated
- [ ] Security considerations addressed

### Pull Request Checklist
- [ ] Feature branch created from `development`
- [ ] Descriptive commit messages
- [ ] Tests added/updated
- [ ] Documentation updated
- [ ] No breaking changes (or documented)
- [ ] Security implications considered

---

## 📞 Getting Help

- **Issues**: [GitHub Issues](https://github.com/Tolstoyj/go-voting-blockchain/issues)
- **Discussions**: [GitHub Discussions](https://github.com/Tolstoyj/go-voting-blockchain/discussions)
- **Email**: tolstoyjustin@gmail.com
- **Documentation**: See [README.md](README.md) and [ROADMAP.md](ROADMAP.md)

---

*Happy coding! 🚀*
