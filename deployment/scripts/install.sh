#!/bin/bash
# Initial server setup script

set -e

echo "======================================"
echo "Rent Receipt Generator - Installation"
echo "======================================"
echo ""

# Check if running as root
if [ "$EUID" -ne 0 ]; then
    echo "Please run as root (use sudo)"
    exit 1
fi

# Update system
echo "Updating system packages..."
apt update && apt upgrade -y

# Install prerequisites
echo "Installing prerequisites..."
apt install -y postgresql postgresql-contrib nginx git curl wget build-essential

# Install Go
echo "Installing Go..."
GO_VERSION="1.21.5"
if ! command -v go &> /dev/null; then
    wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz
    rm -rf /usr/local/go
    tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz
    rm go${GO_VERSION}.linux-amd64.tar.gz
    echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
    export PATH=$PATH:/usr/local/go/bin
fi

# Install Node.js
echo "Installing Node.js..."
if ! command -v node &> /dev/null; then
    curl -fsSL https://deb.nodesource.com/setup_20.x | bash -
    apt install -y nodejs
fi

# Create application directory
echo "Creating application directory..."
mkdir -p /opt/rent-receipt-generator
chown $SUDO_USER:$SUDO_USER /opt/rent-receipt-generator

# Get repository URL
echo ""
echo "Repository URL options:"
echo "  1. HTTPS (recommended): https://github.com/dgentikian/rent-receipt-generator.git"
echo "  2. SSH: git@github.com:dgentikian/rent-receipt-generator.git"
echo ""
read -p "Enter your Git repository URL [default: HTTPS]: " REPO_URL
REPO_URL=${REPO_URL:-https://github.com/dgentikian/rent-receipt-generator.git}

# Clone repository
echo "Cloning repository..."
cd /opt/rent-receipt-generator
if [ -d ".git" ]; then
    echo "Repository already cloned, pulling latest changes..."
    git pull
else
    git clone $REPO_URL .
fi

# Setup PostgreSQL
echo "Setting up PostgreSQL..."
read -p "Enter database name [rent_receipts]: " DB_NAME
DB_NAME=${DB_NAME:-rent_receipts}

read -p "Enter database user [rent_user]: " DB_USER
DB_USER=${DB_USER:-rent_user}

read -sp "Enter database password: " DB_PASSWORD
echo ""

sudo -u postgres psql <<EOF
CREATE DATABASE $DB_NAME;
CREATE USER $DB_USER WITH ENCRYPTED PASSWORD '$DB_PASSWORD';
GRANT ALL PRIVILEGES ON DATABASE $DB_NAME TO $DB_USER;
ALTER DATABASE $DB_NAME OWNER TO $DB_USER;
\q
EOF

# Initialize database
echo "Initializing database..."
PGPASSWORD=$DB_PASSWORD psql -U $DB_USER -d $DB_NAME -f database/schema.sql

# Create .env file
echo "Creating .env file..."
read -p "Enter JWT secret (min 32 chars): " JWT_SECRET
read -p "Enter your domain name: " DOMAIN_NAME

cat > .env <<EOF
# Server Configuration
PORT=8080
ENV=production

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=$DB_USER
DB_PASSWORD=$DB_PASSWORD
DB_NAME=$DB_NAME
DB_SSLMODE=disable

# JWT Authentication
JWT_SECRET=$JWT_SECRET
JWT_EXPIRY=24h

# Application
UPLOADS_DIR=/opt/rent-receipt-generator/uploads
FRONTEND_URL=https://$DOMAIN_NAME
BACKEND_URL=https://$DOMAIN_NAME/api

# CORS
CORS_ALLOWED_ORIGINS=https://$DOMAIN_NAME
EOF

chmod 600 .env

# Build application
echo "Building backend..."
cd backend
./build.sh
cd ..

echo "Building frontend..."
cd frontend
npm ci
npm run build
cd ..

# Create uploads directory
mkdir -p uploads
chown www-data:www-data uploads

# Setup systemd service
echo "Setting up systemd service..."
cp deployment/systemd/rent-receipt.service /etc/systemd/system/
systemctl daemon-reload
systemctl enable rent-receipt
systemctl start rent-receipt

# Setup Nginx
echo "Setting up Nginx..."
cp deployment/nginx/rent-receipt.conf /etc/nginx/sites-available/
sed -i "s/yourdomain.com/$DOMAIN_NAME/g" /etc/nginx/sites-available/rent-receipt.conf
ln -sf /etc/nginx/sites-available/rent-receipt.conf /etc/nginx/sites-enabled/
rm -f /etc/nginx/sites-enabled/default

# Test Nginx configuration
nginx -t

# Setup SSL with Let's Encrypt
echo "Setting up SSL certificates..."
apt install -y certbot python3-certbot-nginx
certbot --nginx -d $DOMAIN_NAME -d www.$DOMAIN_NAME --non-interactive --agree-tos --email admin@$DOMAIN_NAME

# Start services
echo "Starting services..."
systemctl restart nginx
systemctl restart rent-receipt

# Setup firewall
echo "Configuring firewall..."
ufw allow 22/tcp
ufw allow 80/tcp
ufw allow 443/tcp
ufw --force enable

echo ""
echo "======================================"
echo "Installation completed successfully!"
echo "======================================"
echo ""
echo "Your application should be accessible at: https://$DOMAIN_NAME"
echo ""
echo "Service status:"
systemctl status rent-receipt --no-pager
echo ""
echo "To view logs:"
echo "  journalctl -u rent-receipt -f"
echo ""
echo "To create your first account, visit:"
echo "  https://$DOMAIN_NAME/register"
echo ""
