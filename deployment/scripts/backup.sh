#!/bin/bash
# Database backup script

set -e

# Load environment variables
if [ -f /opt/rent-receipt-generator/.env ]; then
    export $(cat /opt/rent-receipt-generator/.env | grep -v '^#' | xargs)
fi

BACKUP_DIR="/opt/rent-receipt-generator/backups"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="$BACKUP_DIR/backup_$DATE.sql"

# Create backup directory if it doesn't exist
mkdir -p $BACKUP_DIR

echo "Creating database backup..."
PGPASSWORD=$DB_PASSWORD pg_dump -U $DB_USER -h $DB_HOST -d $DB_NAME > $BACKUP_FILE

# Compress backup
gzip $BACKUP_FILE

echo "Backup created: ${BACKUP_FILE}.gz"

# Keep only last 30 days of backups
find $BACKUP_DIR -name "backup_*.sql.gz" -mtime +30 -delete

echo "Old backups cleaned up"
echo "Done!"
