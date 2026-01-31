#!/bin/bash
# Database restore script

set -e

# Load environment variables
if [ -f /opt/rent-receipt-generator/.env ]; then
    export $(cat /opt/rent-receipt-generator/.env | grep -v '^#' | xargs)
fi

BACKUP_DIR="/opt/rent-receipt-generator/backups"

if [ -z "$1" ]; then
    echo "Usage: $0 <backup_file.sql.gz>"
    echo ""
    echo "Available backups:"
    ls -lh $BACKUP_DIR/backup_*.sql.gz 2>/dev/null || echo "No backups found"
    exit 1
fi

BACKUP_FILE="$1"

if [ ! -f "$BACKUP_FILE" ]; then
    echo "Error: Backup file not found: $BACKUP_FILE"
    exit 1
fi

echo "WARNING: This will restore the database from backup"
echo "Current data will be lost!"
read -p "Are you sure? (yes/no): " CONFIRM

if [ "$CONFIRM" != "yes" ]; then
    echo "Restore cancelled"
    exit 0
fi

echo "Stopping backend service..."
sudo systemctl stop rent-receipt

echo "Restoring database from $BACKUP_FILE..."
gunzip -c $BACKUP_FILE | PGPASSWORD=$DB_PASSWORD psql -U $DB_USER -h $DB_HOST -d $DB_NAME

echo "Starting backend service..."
sudo systemctl start rent-receipt

echo "Restore completed!"
