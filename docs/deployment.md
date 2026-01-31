# Deployment Guide

## Prerequisites

- Ubuntu/Debian server (20.04 LTS or later)
- Root access
- Domain name pointing to your server
- At least 1GB RAM, 10GB disk space

## Initial Deployment

### 1. Prepare Your Server

```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install Git
sudo apt install -y git
```

### 2. Clone Repository

```bash
cd /tmp
git clone https://github.com/dgentikian/rent-receipt-generator.git rent-receipt-generator
cd rent-receipt-generator
```

### 3. Run Installation Script

```bash
sudo ./deployment/scripts/install.sh
```

The script will:
- Install all dependencies (PostgreSQL, Go, Node.js, Nginx)
- Set up the database
- Build the application
- Configure systemd service
- Set up Nginx with SSL (Let's Encrypt)
- Start all services

### 4. Create First User

Visit `https://yourdomain.com/register` and create your account.

## Updating the Application

```bash
cd /opt/rent-receipt-generator
sudo ./deployment/scripts/deploy.sh
```

## Backup and Restore

### Create Backup

```bash
sudo /opt/rent-receipt-generator/deployment/scripts/backup.sh
```

Backups are stored in `/opt/rent-receipt-generator/backups/`

### Restore from Backup

```bash
sudo /opt/rent-receipt-generator/deployment/scripts/restore.sh /path/to/backup.sql.gz
```

### Automated Daily Backups

Add to crontab:

```bash
sudo crontab -e

# Add this line:
0 2 * * * /opt/rent-receipt-generator/deployment/scripts/backup.sh
```

## Service Management

### Check Status

```bash
sudo systemctl status rent-receipt
```

### View Logs

```bash
# Real-time logs
sudo journalctl -u rent-receipt -f

# Last 100 lines
sudo journalctl -u rent-receipt -n 100
```

### Restart Service

```bash
sudo systemctl restart rent-receipt
```

### Stop/Start Service

```bash
sudo systemctl stop rent-receipt
sudo systemctl start rent-receipt
```

## Troubleshooting

### Backend Won't Start

1. Check logs:
   ```bash
   sudo journalctl -u rent-receipt -n 50
   ```

2. Verify environment variables:
   ```bash
   cat /opt/rent-receipt-generator/.env
   ```

3. Test database connection:
   ```bash
   cd /opt/rent-receipt-generator
   source .env
   psql -U $DB_USER -h $DB_HOST -d $DB_NAME
   ```

### Frontend Not Loading

1. Check Nginx configuration:
   ```bash
   sudo nginx -t
   sudo systemctl status nginx
   ```

2. Verify build exists:
   ```bash
   ls -la /opt/rent-receipt-generator/frontend/dist/
   ```

3. Check Nginx logs:
   ```bash
   sudo tail -f /var/log/nginx/rent-receipt-error.log
   ```

### SSL Certificate Issues

Renew certificates manually:
```bash
sudo certbot renew --nginx
```

Auto-renewal is set up during installation.

## Security Recommendations

1. **Firewall**: Only allow ports 22, 80, 443
2. **SSH**: Use key-based authentication, disable password login
3. **PostgreSQL**: Ensure it's not accessible from outside
4. **Backups**: Store backups off-server
5. **Updates**: Keep system packages up to date

## Performance Optimization

### PostgreSQL Tuning

Edit `/etc/postgresql/*/main/postgresql.conf`:

```conf
shared_buffers = 256MB
effective_cache_size = 1GB
maintenance_work_mem = 64MB
checkpoint_completion_target = 0.9
wal_buffers = 16MB
default_statistics_target = 100
random_page_cost = 1.1
```

Restart PostgreSQL:
```bash
sudo systemctl restart postgresql
```

### Nginx Caching

Already configured in the provided nginx config.

## Monitoring

### Basic Monitoring

```bash
# Disk usage
df -h

# Memory usage
free -h

# CPU usage
top

# Service status
sudo systemctl status rent-receipt nginx postgresql
```

### Advanced Monitoring

Consider installing:
- Prometheus + Grafana
- netdata
- fail2ban (for security)

## Support

For issues, check:
- Application logs: `sudo journalctl -u rent-receipt -f`
- Nginx logs: `/var/log/nginx/rent-receipt-*.log`
- PostgreSQL logs: `/var/log/postgresql/`
