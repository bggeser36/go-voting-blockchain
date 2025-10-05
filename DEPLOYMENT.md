# 🚀 Go Voting Blockchain - Deployment Guide

This comprehensive guide covers deploying the Go Voting Blockchain system to various platforms with Phase 1 security features enabled.

## 📋 Table of Contents

- [Prerequisites](#prerequisites)
- [Environment Configuration](#environment-configuration)
- [Local Development Deployment](#local-development-deployment)
- [Docker Deployment](#docker-deployment)
- [Railway Deployment](#railway-deployment)
- [Production Deployment](#production-deployment)
- [Security Configuration](#security-configuration)
- [Monitoring & Maintenance](#monitoring--maintenance)
- [Troubleshooting](#troubleshooting)

## 🔧 Prerequisites

### System Requirements

**Minimum Requirements:**
- **CPU**: 2 cores, 2.0 GHz
- **RAM**: 4 GB
- **Storage**: 20 GB SSD
- **Network**: 100 Mbps

**Recommended for Production:**
- **CPU**: 4+ cores, 3.0+ GHz
- **RAM**: 8+ GB
- **Storage**: 100+ GB SSD
- **Network**: 1 Gbps

### Software Dependencies

- **Go**: 1.21 or higher
- **PostgreSQL**: 13+ (for production)
- **Redis**: 6+ (optional, for caching)
- **Docker**: 20+ (for containerized deployment)
- **Docker Compose**: 2+ (for multi-service deployment)

## ⚙️ Environment Configuration

### Required Environment Variables

Create a `.env` file with the following variables:

```bash
# Server Configuration
PORT=8080
GIN_MODE=release

# Database Configuration
DATABASE_URL=postgres://username:password@localhost:5432/voting_db
DB_HOST=localhost
DB_PORT=5432
DB_USER=voting_user
DB_PASSWORD=secure_password
DB_NAME=voting_db
DB_SSLMODE=require

# Redis Configuration (Optional)
REDIS_URL=redis://localhost:6379
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=redis_password

# Blockchain Configuration
MINING_DIFFICULTY=4
BLOCK_REWARD=1.0

# Security Configuration
JWT_SECRET=your-super-secure-jwt-secret-key-here
JWT_EXPIRATION=24h
JWT_REFRESH_EXPIRATION=168h

# Admin Configuration
ADMIN_USERNAME=admin
ADMIN_PASSWORD=your-secure-admin-password

# CORS Configuration
CORS_ALLOWED_ORIGINS=https://yourdomain.com,https://app.yourdomain.com
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOWED_HEADERS=Content-Type,Authorization,X-Requested-With

# Rate Limiting Configuration
RATE_LIMIT_STRICT=10
RATE_LIMIT_MODERATE=30
RATE_LIMIT_GENEROUS=100
RATE_LIMIT_WINDOW=1m

# Logging Configuration
LOG_LEVEL=info
LOG_FORMAT=json
```

### Security Considerations

**🔐 Critical Security Settings:**

1. **JWT Secret**: Use a cryptographically secure random string (minimum 32 characters)
2. **Database Password**: Use a strong password with special characters
3. **Admin Password**: Change from default immediately
4. **CORS Origins**: Restrict to your actual domains
5. **SSL/TLS**: Always use HTTPS in production

**Generate Secure JWT Secret:**
```bash
# Using OpenSSL
openssl rand -base64 32

# Using Go
go run -c 'import "crypto/rand"; import "encoding/base64"; b := make([]byte, 32); rand.Read(b); println(base64.StdEncoding.EncodeToString(b))'
```

## 🏠 Local Development Deployment

### Option 1: Direct Go Execution

1. **Clone and Setup:**
```bash
git clone https://github.com/tolstoyjustin/go-voting-blockchain.git
cd go-voting-blockchain
go mod download
```

2. **Configure Environment:**
```bash
cp config.development.env .env
# Edit .env with your configuration
```

3. **Start Services:**
```bash
# Start PostgreSQL (if not running)
brew services start postgresql  # macOS
sudo systemctl start postgresql  # Linux

# Start Redis (optional)
brew services start redis  # macOS
sudo systemctl start redis  # Linux
```

4. **Initialize Database:**
```bash
psql -h localhost -U voting_user -d voting_db -f scripts/init-db.sql
```

5. **Run Application:**
```bash
go run cmd/api/main.go
```

### Option 2: Docker Compose

1. **Setup:**
```bash
git clone https://github.com/tolstoyjustin/go-voting-blockchain.git
cd go-voting-blockchain
```

2. **Configure Environment:**
```bash
cp docker-compose.yml.example docker-compose.yml
# Edit docker-compose.yml with your configuration
```

3. **Start Services:**
```bash
docker-compose up -d
```

4. **Check Status:**
```bash
docker-compose ps
docker-compose logs -f voting-api
```

## 🐳 Docker Deployment

### Single Container Deployment

1. **Build Image:**
```bash
docker build -t voting-blockchain:latest .
```

2. **Run Container:**
```bash
docker run -d \
  --name voting-blockchain \
  -p 8080:8080 \
  -e DATABASE_URL="postgres://user:pass@host:5432/db" \
  -e JWT_SECRET="your-jwt-secret" \
  -e ADMIN_PASSWORD="secure-password" \
  voting-blockchain:latest
```

### Multi-Container Deployment

1. **Create docker-compose.yml:**
```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15
    environment:
      POSTGRES_DB: voting_db
      POSTGRES_USER: voting_user
      POSTGRES_PASSWORD: secure_password
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./scripts/init-db.sql:/docker-entrypoint-initdb.d/init-db.sql
    ports:
      - "5432:5432"

  redis:
    image: redis:7-alpine
    command: redis-server --requirepass redis_password
    volumes:
      - redis_data:/data
    ports:
      - "6379:6379"

  voting-api:
    build: .
    environment:
      - PORT=8080
      - DATABASE_URL=postgres://voting_user:secure_password@postgres:5432/voting_db
      - REDIS_URL=redis://:redis_password@redis:6379
      - JWT_SECRET=your-jwt-secret
      - ADMIN_PASSWORD=secure-password
      - GIN_MODE=release
    ports:
      - "8080:8080"
    depends_on:
      - postgres
      - redis
    restart: unless-stopped

volumes:
  postgres_data:
  redis_data:
```

2. **Deploy:**
```bash
docker-compose up -d
```

## ☁️ Railway Deployment

### One-Click Deploy

[![Deploy on Railway](https://railway.app/button.svg)](https://railway.app/new/template)

### Manual Railway Deployment

1. **Install Railway CLI:**
```bash
npm install -g @railway/cli
railway login
```

2. **Initialize Project:**
```bash
railway init
railway add postgresql
railway add redis
```

3. **Configure Environment Variables:**
```bash
railway variables set JWT_SECRET="your-jwt-secret"
railway variables set ADMIN_PASSWORD="secure-password"
railway variables set MINING_DIFFICULTY="4"
railway variables set GIN_MODE="release"
```

4. **Deploy:**
```bash
railway up
```

### Railway Environment Variables

Railway automatically provides:
- `PORT`: Server port (auto-assigned)
- `DATABASE_URL`: PostgreSQL connection string
- `REDIS_URL`: Redis connection string (if added)

**Required Manual Configuration:**
- `JWT_SECRET`: Your secure JWT secret
- `ADMIN_PASSWORD`: Secure admin password
- `MINING_DIFFICULTY`: Blockchain mining difficulty (3-5)
- `GIN_MODE`: Set to "release" for production

## 🏭 Production Deployment

### AWS EC2 Deployment

1. **Launch EC2 Instance:**
   - AMI: Ubuntu 22.04 LTS
   - Instance Type: t3.medium or larger
   - Security Groups: HTTP (80), HTTPS (443), SSH (22)

2. **Install Dependencies:**
```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install Go
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc

# Install PostgreSQL
sudo apt install postgresql postgresql-contrib -y

# Install Redis
sudo apt install redis-server -y

# Install Docker (optional)
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
```

3. **Deploy Application:**
```bash
# Clone repository
git clone https://github.com/tolstoyjustin/go-voting-blockchain.git
cd go-voting-blockchain

# Configure environment
cp .env.example .env
# Edit .env with production values

# Build and run
go build -o voting-blockchain cmd/api/main.go
sudo ./voting-blockchain
```

### Google Cloud Platform Deployment

1. **Create VM Instance:**
   - Machine Type: e2-medium or larger
   - Boot Disk: Ubuntu 22.04 LTS
   - Firewall: Allow HTTP, HTTPS traffic

2. **Deploy using Cloud Run:**
```bash
# Build and push to Container Registry
gcloud builds submit --tag gcr.io/PROJECT_ID/voting-blockchain

# Deploy to Cloud Run
gcloud run deploy voting-blockchain \
  --image gcr.io/PROJECT_ID/voting-blockchain \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars JWT_SECRET="your-jwt-secret"
```

### DigitalOcean App Platform

1. **Create App Spec (app.yaml):**
```yaml
name: voting-blockchain
services:
- name: api
  source_dir: /
  github:
    repo: tolstoyjustin/go-voting-blockchain
    branch: main
  run_command: go run cmd/api/main.go
  environment_slug: go
  instance_count: 1
  instance_size_slug: basic-xxs
  envs:
  - key: JWT_SECRET
    value: your-jwt-secret
  - key: ADMIN_PASSWORD
    value: secure-password
databases:
- name: postgres
  engine: PG
  version: "15"
```

2. **Deploy:**
```bash
doctl apps create app.yaml
```

## 🔒 Security Configuration

### SSL/TLS Setup

**Using Let's Encrypt (Certbot):**
```bash
# Install Certbot
sudo apt install certbot python3-certbot-nginx -y

# Obtain SSL Certificate
sudo certbot --nginx -d yourdomain.com

# Auto-renewal
sudo crontab -e
# Add: 0 12 * * * /usr/bin/certbot renew --quiet
```

**Using Cloudflare:**
1. Add your domain to Cloudflare
2. Set DNS records to point to your server
3. Enable SSL/TLS encryption mode: "Full (strict)"

### Firewall Configuration

**UFW (Ubuntu):**
```bash
sudo ufw allow 22/tcp    # SSH
sudo ufw allow 80/tcp    # HTTP
sudo ufw allow 443/tcp   # HTTPS
sudo ufw enable
```

**iptables:**
```bash
# Allow SSH, HTTP, HTTPS
iptables -A INPUT -p tcp --dport 22 -j ACCEPT
iptables -A INPUT -p tcp --dport 80 -j ACCEPT
iptables -A INPUT -p tcp --dport 443 -j ACCEPT
iptables -A INPUT -p tcp --dport 8080 -j DROP  # Block direct API access
```

### Database Security

1. **Change Default Passwords:**
```sql
ALTER USER postgres PASSWORD 'new-secure-password';
CREATE USER voting_user WITH PASSWORD 'secure-password';
GRANT ALL PRIVILEGES ON DATABASE voting_db TO voting_user;
```

2. **Enable SSL:**
```bash
# In postgresql.conf
ssl = on
ssl_cert_file = 'server.crt'
ssl_key_file = 'server.key'
```

3. **Restrict Access:**
```bash
# In pg_hba.conf
host voting_db voting_user 127.0.0.1/32 md5
host all all 0.0.0.0/0 reject
```

## 📊 Monitoring & Maintenance

### Health Checks

**Application Health:**
```bash
curl http://localhost:8080/health
```

**Database Health:**
```bash
psql -h localhost -U voting_user -d voting_db -c "SELECT 1;"
```

**Redis Health:**
```bash
redis-cli ping
```

### Logging

**Application Logs:**
```bash
# Docker
docker logs -f voting-blockchain

# Systemd
journalctl -u voting-blockchain -f

# Direct
tail -f /var/log/voting-blockchain.log
```

**Log Rotation:**
```bash
# /etc/logrotate.d/voting-blockchain
/var/log/voting-blockchain.log {
    daily
    missingok
    rotate 52
    compress
    delaycompress
    notifempty
    create 644 voting voting
    postrotate
        systemctl reload voting-blockchain
    endscript
}
```

### Backup Procedures

**Database Backup:**
```bash
#!/bin/bash
# backup-db.sh
DATE=$(date +%Y%m%d_%H%M%S)
pg_dump $DATABASE_URL > backup_$DATE.sql
gzip backup_$DATE.sql
aws s3 cp backup_$DATE.sql.gz s3://your-backup-bucket/
```

**Automated Backup (Cron):**
```bash
# Daily backup at 2 AM
0 2 * * * /path/to/backup-db.sh
```

### Performance Monitoring

**System Metrics:**
```bash
# CPU and Memory
htop

# Disk Usage
df -h

# Network
iftop
```

**Application Metrics:**
- Response times via `/health` endpoint
- Error rates in logs
- Database connection pool status
- Rate limiting statistics

## 🔧 Troubleshooting

### Common Issues

**1. Port Already in Use:**
```bash
# Find process using port
sudo lsof -i :8080
# Kill process
sudo kill -9 PID
```

**2. Database Connection Failed:**
```bash
# Check PostgreSQL status
sudo systemctl status postgresql
# Check connection
psql -h localhost -U voting_user -d voting_db
```

**3. JWT Token Issues:**
```bash
# Verify JWT secret is set
echo $JWT_SECRET
# Check token format
curl -H "Authorization: Bearer YOUR_TOKEN" http://localhost:8080/auth/me
```

**4. Rate Limiting Too Strict:**
```bash
# Adjust rate limits in .env
RATE_LIMIT_STRICT=20
RATE_LIMIT_MODERATE=60
RATE_LIMIT_GENEROUS=200
```

### Performance Issues

**High Memory Usage:**
- Increase server RAM
- Optimize database queries
- Implement connection pooling
- Add Redis caching

**Slow Response Times:**
- Check database indexes
- Optimize API endpoints
- Enable response compression
- Use CDN for static content

**Database Bottlenecks:**
- Add database indexes
- Optimize queries
- Use read replicas
- Implement connection pooling

### Security Issues

**Authentication Failures:**
- Verify JWT secret configuration
- Check token expiration settings
- Validate user credentials
- Review rate limiting settings

**Database Security:**
- Change default passwords
- Enable SSL connections
- Restrict database access
- Regular security updates

## 📞 Support

For deployment issues:

- 🐛 **Bug Reports**: [GitHub Issues](https://github.com/tolstoyjustin/go-voting-blockchain/issues)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/tolstoyjustin/go-voting-blockchain/discussions)
- 📧 **Contact**: tolstoyjustin@gmail.com
- 📚 **Documentation**: [README.md](README.md)

## 🎯 Deployment Checklist

### Pre-Deployment
- [ ] Environment variables configured
- [ ] Database schema initialized
- [ ] SSL certificates installed
- [ ] Firewall rules configured
- [ ] Monitoring setup

### Security
- [ ] Default passwords changed
- [ ] JWT secret configured
- [ ] Rate limiting enabled
- [ ] CORS properly configured
- [ ] Security headers set

### Testing
- [ ] Health check endpoint working
- [ ] Authentication flows tested
- [ ] API endpoints responding
- [ ] Database connectivity verified
- [ ] Error handling working

### Post-Deployment
- [ ] Application accessible
- [ ] Monitoring active
- [ ] Backups scheduled
- [ ] Documentation updated
- [ ] Team notified

---

*This deployment guide is regularly updated. Check for the latest version in the repository.*
