# 📦 Installation Guide

This guide provides detailed installation instructions for the Master Data REST API across different platforms and deployment scenarios.

## 📋 System Requirements

### Minimum Requirements
- **CPU**: 1 vCPU
- **Memory**: 512MB RAM
- **Storage**: 1GB available space
- **Network**: Port 8080 available

### Recommended Requirements
- **CPU**: 2+ vCPUs
- **Memory**: 2GB+ RAM
- **Storage**: 10GB+ available space
- **Network**: Reverse proxy setup

## 🛠️ Prerequisites

### Required Software
1. **Go 1.21 or higher**
   ```bash
   # Check Go version
   go version
   
   # Install Go (if needed)
   # Visit: https://golang.org/dl/
   ```

2. **PostgreSQL 13 or higher**
   ```bash
   # Check PostgreSQL version
   psql --version
   
   # Ubuntu/Debian
   sudo apt-get install postgresql postgresql-contrib
   
   # CentOS/RHEL
   sudo yum install postgresql postgresql-server
   
   # macOS (with Homebrew)
   brew install postgresql
   
   # Windows
   # Download from: https://www.postgresql.org/download/windows/
   ```

### Optional Software
- **Docker & Docker Compose** (for containerized deployment)
- **Git** (for source code management)
- **Make** (for build automation)

## 🚀 Installation Methods

### Method 1: From Source (Recommended for Development)

1. **Clone the repository**
   ```bash
   git clone https://github.com/turahe/master-data-rest-api.git
   cd master-data-rest-api
   ```

2. **Install Go dependencies**
   ```bash
   go mod tidy
   ```

3. **Set up environment**
   ```bash
   cp env.example .env
   # Edit .env file with your configuration
   ```

4. **Build the application**
   ```bash
   # For current platform
   go build -o bin/master-data-api main.go
   
   # For Windows specifically
   go build -o bin/master-data-api.exe main.go
   
   # For multiple platforms
   make build-all
   ```

### Method 2: Using Pre-built Binaries

1. **Download release**
   ```bash
   # Download latest release for your platform
   wget https://github.com/turahe/master-data-rest-api/releases/latest/download/master-data-api-linux-amd64.tar.gz
   tar -xzf master-data-api-linux-amd64.tar.gz
   ```

2. **Make executable and move to PATH**
   ```bash
   chmod +x master-data-api
   sudo mv master-data-api /usr/local/bin/
   ```

### Method 3: Using Docker

1. **Using Docker Compose (Recommended)**
   ```bash
   git clone https://github.com/turahe/master-data-rest-api.git
   cd master-data-rest-api
   
   # Copy environment file
   cp env.example .env
   
   # Start services
   docker-compose up -d
   ```

2. **Using Docker directly**
   ```bash
   # Build image
   docker build -t master-data-api .
   
   # Run container
   docker run -d \
     --name master-data-api \
     -p 8080:8080 \
     --env-file .env \
     master-data-api
   ```

## 🗄️ Database Setup

### PostgreSQL Configuration

1. **Create database and user**
   ```sql
   -- Connect to PostgreSQL as superuser
   sudo -u postgres psql
   
   -- Create database
   CREATE DATABASE master_data;
   
   -- Create user
   CREATE USER appuser WITH PASSWORD 'your_secure_password';
   
   -- Grant privileges
   GRANT ALL PRIVILEGES ON DATABASE master_data TO appuser;
   
   -- Grant schema privileges
   \c master_data
   GRANT ALL ON SCHEMA public TO appuser;
   GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO appuser;
   GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO appuser;
   
   -- Exit
   \q
   ```

2. **Configure PostgreSQL (if needed)**
   ```bash
   # Edit PostgreSQL configuration
   sudo nano /etc/postgresql/13/main/postgresql.conf
   
   # Ensure these settings:
   listen_addresses = 'localhost'
   port = 5432
   
   # Edit pg_hba.conf for authentication
   sudo nano /etc/postgresql/13/main/pg_hba.conf
   
   # Add line for local connections:
   local   master_data   appuser                     md5
   
   # Restart PostgreSQL
   sudo systemctl restart postgresql
   ```

### Environment Configuration

1. **Update .env file**
   ```bash
   # Database Configuration
   DB_DRIVER=postgres
   DB_HOST=localhost
   DB_PORT=5432
   DB_NAME=master_data
   DB_USER=appuser
   DB_PASSWORD=your_secure_password
   DB_SSL_MODE=disable
   
   # Application Configuration
   APP_HOST=localhost
   APP_PORT=8080
   APP_ENV=production
   
   # Logging Configuration
   LOG_LEVEL=info
   LOG_FORMAT=json
   LOG_OUTPUT=stdout
   ```

## 🔧 Initial Setup

### 1. Run Database Migrations

```bash
# Run all pending migrations
./master-data-api migrate up

# Check migration status
./master-data-api migrate status
```

### 2. Create Initial API Key

```bash
# Create an admin API key
./master-data-api create-api-key \
  --name "Admin Key" \
  --description "Initial administrative access"

# Save the generated API key - you'll need it for API access
```

### 3. Test the Installation

```bash
# Start the server
./master-data-api serve

# In another terminal, test the health endpoint
curl http://localhost:8080/health

# Test authenticated endpoint
curl -H "Authorization: Bearer YOUR_API_KEY" \
     http://localhost:8080/api/v1/geodirectories
```

## 🌍 Platform-Specific Instructions

### Ubuntu/Debian

1. **Install prerequisites**
   ```bash
   sudo apt update
   sudo apt install -y wget curl git make
   
   # Install Go
   wget https://golang.org/dl/go1.21.0.linux-amd64.tar.gz
   sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
   echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
   source ~/.bashrc
   
   # Install PostgreSQL
   sudo apt install -y postgresql postgresql-contrib
   ```

2. **Follow general installation steps above**

### CentOS/RHEL

1. **Install prerequisites**
   ```bash
   sudo yum update -y
   sudo yum install -y wget curl git make
   
   # Install Go
   wget https://golang.org/dl/go1.21.0.linux-amd64.tar.gz
   sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
   echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
   source ~/.bashrc
   
   # Install PostgreSQL
   sudo yum install -y postgresql postgresql-server
   sudo postgresql-setup initdb
   sudo systemctl enable postgresql
   sudo systemctl start postgresql
   ```

### macOS

1. **Install prerequisites using Homebrew**
   ```bash
   # Install Homebrew (if not installed)
   /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
   
   # Install dependencies
   brew install go postgresql git make
   
   # Start PostgreSQL
   brew services start postgresql
   ```

### Windows

1. **Install prerequisites**
   - Download and install Go from [golang.org](https://golang.org/dl/)
   - Download and install PostgreSQL from [postgresql.org](https://www.postgresql.org/download/windows/)
   - Install Git from [git-scm.com](https://git-scm.com/download/win)

2. **Build using PowerShell**
   ```powershell
   # Clone repository
   git clone https://github.com/turahe/master-data-rest-api.git
   cd master-data-rest-api
   
   # Build application
   go build -o bin/master-data-api.exe main.go
   
   # Or use the PowerShell build script
   PowerShell -ExecutionPolicy Bypass -File scripts/build.ps1
   ```

## 🔒 Security Setup

### 1. Secure PostgreSQL

```sql
-- Create application-specific user (not superuser)
CREATE USER app_readonly WITH PASSWORD 'readonly_password';
GRANT SELECT ON ALL TABLES IN SCHEMA public TO app_readonly;

-- Create backup user
CREATE USER backup_user WITH PASSWORD 'backup_password';
GRANT SELECT ON ALL TABLES IN SCHEMA public TO backup_user;
```

### 2. Configure Firewall

```bash
# Ubuntu/Debian (ufw)
sudo ufw allow 8080/tcp
sudo ufw enable

# CentOS/RHEL (firewalld)
sudo firewall-cmd --permanent --add-port=8080/tcp
sudo firewall-cmd --reload
```

### 3. Set up SSL/TLS (Production)

```bash
# Generate self-signed certificate (for testing)
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes

# For production, use Let's Encrypt or your certificate authority
```

## 📊 Monitoring Setup

### 1. Log Aggregation

```bash
# Configure log output to file
export LOG_OUTPUT="/var/log/master-data-api/app.log"

# Create log directory
sudo mkdir -p /var/log/master-data-api
sudo chown $USER:$USER /var/log/master-data-api
```

### 2. Health Checks

```bash
# Create health check script
cat > health_check.sh << 'EOF'
#!/bin/bash
curl -f http://localhost:8080/health || exit 1
EOF

chmod +x health_check.sh
```

## 🚀 Production Deployment

### Systemd Service (Linux)

1. **Create service file**
   ```bash
   sudo nano /etc/systemd/system/master-data-api.service
   ```

2. **Service configuration**
   ```ini
   [Unit]
   Description=Master Data REST API
   After=network.target postgresql.service
   Requires=postgresql.service
   
   [Service]
   Type=simple
   User=api
   Group=api
   WorkingDirectory=/opt/master-data-api
   ExecStart=/opt/master-data-api/master-data-api serve
   EnvironmentFile=/opt/master-data-api/.env
   Restart=always
   RestartSec=10
   
   # Security settings
   NoNewPrivileges=true
   PrivateTmp=true
   ProtectSystem=strict
   ReadWritePaths=/opt/master-data-api
   
   [Install]
   WantedBy=multi-user.target
   ```

3. **Enable and start service**
   ```bash
   sudo systemctl daemon-reload
   sudo systemctl enable master-data-api
   sudo systemctl start master-data-api
   sudo systemctl status master-data-api
   ```

### Reverse Proxy (Nginx)

```nginx
server {
    listen 80;
    server_name your-domain.com;
    
    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## 🔍 Troubleshooting

### Common Issues

1. **Port already in use**
   ```bash
   # Find process using port
   lsof -i :8080
   
   # Kill process or use different port
   ./master-data-api serve --port 8081
   ```

2. **Database connection failed**
   ```bash
   # Test PostgreSQL connection
   psql -h localhost -U appuser -d master_data
   
   # Check PostgreSQL status
   sudo systemctl status postgresql
   ```

3. **Permission denied**
   ```bash
   # Make binary executable
   chmod +x master-data-api
   
   # Check file ownership
   ls -la master-data-api
   ```

4. **Migration failures**
   ```bash
   # Check migration status
   ./master-data-api migrate status
   
   # Force migration to specific version (careful!)
   ./master-data-api migrate force --version 5
   ```

### Getting Help

- Check logs: `./master-data-api serve --log-level debug`
- Validate configuration: `./master-data-api --help`
- Test database connection: `./master-data-api migrate status`
- Open GitHub issue with logs and environment details

---

<div align="center">
  <strong>🎉 Installation Complete!</strong>
  <br>
  <sub>Your Master Data REST API is ready to serve requests</sub>
</div>
