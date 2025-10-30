#!/bin/bash
# Backup script for mautrix-viber bridge

set -e

BACKUP_DIR="${BACKUP_DIR:-/backup/mautrix-viber}"
DB_PATH="${DB_PATH:-./data/bridge.db}"
RETENTION_DAYS="${RETENTION_DAYS:-30}"

# Create backup directory
mkdir -p "$BACKUP_DIR"

# Generate backup filename
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="$BACKUP_DIR/bridge_$TIMESTAMP.db"

# Backup database
if [ -f "$DB_PATH" ]; then
    sqlite3 "$DB_PATH" ".backup '$BACKUP_FILE'"
    echo "Backup created: $BACKUP_FILE"
    
    # Compress backup
    gzip "$BACKUP_FILE"
    echo "Backup compressed: ${BACKUP_FILE}.gz"
else
    echo "Error: Database file not found at $DB_PATH"
    exit 1
fi

# Clean up old backups
find "$BACKUP_DIR" -name "bridge_*.db.gz" -mtime +$RETENTION_DAYS -delete
echo "Old backups cleaned up (older than $RETENTION_DAYS days)"

echo "Backup completed successfully"

