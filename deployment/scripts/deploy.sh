#!/bin/bash
# Deployment script for updates

set -e

APP_DIR="/opt/rent-receipt-generator"

echo "======================================="
echo "Deploying Rent Receipt Generator"
echo "======================================="
echo ""

# Check if app directory exists
if [ ! -d "$APP_DIR" ]; then
    echo "Error: Application directory not found"
    echo "Please run install.sh first"
    exit 1
fi

# Navigate to app directory
cd $APP_DIR

# Pull latest code
echo "Pulling latest code from git..."
git pull origin main

# Build backend
echo "Building backend..."
cd backend
./build.sh
cd ..

# Build frontend
echo "Building frontend..."
cd frontend
npm ci
npm run build
cd ..

# Restart backend service
echo "Restarting backend service..."
sudo systemctl restart rent-receipt

# Check status
sleep 2
if systemctl is-active --quiet rent-receipt; then
    echo "✓ Backend service is running"
else
    echo "✗ Backend service failed to start!"
    echo "Checking logs..."
    sudo journalctl -u rent-receipt -n 50
    exit 1
fi

# Reload Nginx
echo "Reloading Nginx..."
sudo nginx -t && sudo systemctl reload nginx

echo ""
echo "======================================="
echo "Deployment completed successfully!"
echo "======================================="
echo ""
echo "Service status:"
sudo systemctl status rent-receipt --no-pager | head -10
echo ""
echo "To view logs:"
echo "  sudo journalctl -u rent-receipt -f"
echo ""
