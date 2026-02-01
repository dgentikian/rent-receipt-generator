# Environment Configuration Guide

## Overview

The application uses environment variables for configuration. Two template files are provided:

- `env-local.example` - For local development
- `env-prod.example` - For production deployment

## Setup for Local Development

```bash
# 1. Copy the local development template
cp env-local.example .env.local

# 2. Edit .env.local with your local settings
# The defaults should work for most local setups

# 3. Create your local database
createdb rent_receipts
psql -d rent_receipts -f database/schema.sql

# 4. Run the backend (it will load .env.local automatically)
cd backend
go run cmd/server/main.go
```

### PostgreSQL Setup for Local Development

**Option 1: No Password (Trust Authentication)**

Leave `DB_PASSWORD` empty in `.env.local`. Configure PostgreSQL to allow local connections without password:

```bash
# Edit pg_hba.conf (usually in /etc/postgresql/*/main/ or /usr/local/var/postgres/)
# Change this line:
# local   all   all   peer
# To:
local   all   all   trust

# Or for TCP connections:
host    all   all   127.0.0.1/32   trust

# Restart PostgreSQL
sudo systemctl restart postgresql  # Linux
brew services restart postgresql   # macOS
```

**Option 2: With Password**

Set a password in `.env.local` and create the user:

```bash
# Create user with password
sudo -u postgres psql
CREATE USER rent_user WITH PASSWORD 'your_dev_password';
CREATE DATABASE rent_receipts OWNER rent_user;
\q

# Set DB_PASSWORD in .env.local
DB_PASSWORD=your_dev_password
```

### Local Development Defaults

- Backend runs on: `http://localhost:8080`
- Frontend runs on: `http://localhost:3000`
- Database: `localhost:5432/rent_receipts`
- SSL disabled for database
- **DB_PASSWORD can be empty** (for trust authentication)
- Weak JWT secret (OK for local dev)

## Setup for Production

```bash
# 1. Copy the production template
cp env-prod.example .env.prod

# 2. Generate strong secrets
openssl rand -base64 32  # For JWT_SECRET

# 3. Edit .env.prod with your production values
# IMPORTANT: Replace ALL placeholder values!
nano .env.prod

# 4. Set proper permissions
chmod 600 .env.prod
```

### Production Requirements

✅ **Must Change:**
- `DB_PASSWORD` - Use a strong database password
- `JWT_SECRET` - Generate with `openssl rand -base64 32`
- `FRONTEND_URL` - Your actual domain
- `BACKEND_URL` - Your actual domain
- `CORS_ALLOWED_ORIGINS` - Your actual domain(s)

⚠️ **Security Notes:**
- Use `DB_SSLMODE=require` in production
- Never commit `.env.prod` to git (it's in .gitignore)
- Use HTTPS only in production
- Store backups of `.env.prod` securely

## Using Different Environment Files

### Option 1: Default Behavior

The application looks for `.env` in the root directory:

```bash
# Development
cp .env.local .env
go run cmd/server/main.go

# Production
cp .env.prod .env
./backend/bin/rent-receipt-server
```

### Option 2: Specify Environment File

You can load a specific file in your code or use environment variables:

```bash
# Development
ENV_FILE=.env.local go run cmd/server/main.go

# Production
ENV_FILE=.env.prod ./backend/bin/rent-receipt-server
```

### Option 3: System Environment Variables

For production, you can also set environment variables directly:

```bash
# In systemd service file
EnvironmentFile=/opt/rent-receipt-generator/.env.prod

# Or export in shell
export DB_PASSWORD="your_password"
export JWT_SECRET="your_secret"
```

## Environment Variables Reference

### Server Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | 8080 | Backend server port |
| `ENV` | development | Environment (development/production) |

### Database Configuration

| Variable | Required | Description |
|----------|----------|-------------|
| `DB_HOST` | Yes | Database host |
| `DB_PORT` | Yes | Database port (usually 5432) |
| `DB_USER` | Yes | Database user |
| `DB_PASSWORD` | Yes | Database password |
| `DB_NAME` | Yes | Database name |
| `DB_SSLMODE` | Yes | SSL mode (disable/require) |

### JWT Authentication

| Variable | Required | Description |
|----------|----------|-------------|
| `JWT_SECRET` | Yes | Secret for signing JWT tokens (min 32 chars) |
| `JWT_EXPIRY` | Yes | Token expiry duration (e.g., 24h) |

### Application Settings

| Variable | Required | Description |
|----------|----------|-------------|
| `UPLOADS_DIR` | Yes | Directory for uploaded files |
| `FRONTEND_URL` | Yes | Frontend application URL |
| `BACKEND_URL` | Yes | Backend API URL |
| `CORS_ALLOWED_ORIGINS` | Yes | Allowed CORS origins (comma-separated) |

### Email Configuration (Optional)

| Variable | Required | Description |
|----------|----------|-------------|
| `SMTP_HOST` | No | SMTP server host |
| `SMTP_PORT` | No | SMTP server port |
| `SMTP_USER` | No | SMTP username |
| `SMTP_PASSWORD` | No | SMTP password |
| `SMTP_FROM` | No | From email address |

## Quick Commands

### Generate Strong Secrets

```bash
# Generate JWT secret
openssl rand -base64 32

# Generate database password
openssl rand -base64 24

# Generate random string
head /dev/urandom | LC_ALL=C tr -dc 'A-Za-z0-9' | head -c 32
```

### Check Current Environment

```bash
# View loaded environment (remove sensitive values first!)
env | grep -E '^(PORT|ENV|DB_|JWT_|FRONTEND_|BACKEND_|CORS_)' | grep -v PASSWORD | grep -v SECRET
```

### Validate Configuration

```bash
# Check if all required variables are set
cd backend
go run cmd/server/main.go --check-config  # Add this flag if implemented
```

## Troubleshooting

### "DB_PASSWORD is required" Error

Make sure you've set `DB_PASSWORD` in your `.env` file:

```bash
# Check if .env exists
ls -la .env*

# Verify DB_PASSWORD is set
grep DB_PASSWORD .env
```

### "JWT_SECRET is required" Error

Ensure `JWT_SECRET` is at least 32 characters:

```bash
# Generate and add to .env
echo "JWT_SECRET=$(openssl rand -base64 32)" >> .env
```

### Connection Refused Errors

Check your database is running and credentials are correct:

```bash
# Test database connection
psql -U rent_user -d rent_receipts -h localhost
```

## Best Practices

1. ✅ Never commit actual `.env` files
2. ✅ Keep `.env.prod` backed up securely
3. ✅ Use different secrets for each environment
4. ✅ Rotate secrets periodically
5. ✅ Use strong passwords (20+ characters)
6. ✅ Enable SSL in production (`DB_SSLMODE=require`)
7. ✅ Restrict file permissions (`chmod 600 .env.prod`)
