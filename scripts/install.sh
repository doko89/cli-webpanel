#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Check if script is run as root
if [ "$EUID" -ne 0 ]; then 
  echo -e "${RED}Error: This script must be run as root${NC}"
  exit 1
fi

echo -e "${GREEN}Starting webpanel installation...${NC}"

# Function to check command exists
command_exists() {
  command -v "$1" >/dev/null 2>&1
}

# Function to detect OS
detect_os() {
  if [ -f /etc/os-release ]; then
    . /etc/os-release
    OS=$ID
  else
    OS=$(uname -s)
  fi
  echo $OS
}

# Ensure basic tools are installed
echo -e "\n${YELLOW}Checking basic requirements...${NC}"

# Update package list
apt-get update

# Install curl and git if not present
if ! command_exists curl; then
  echo "Installing curl..."
  apt-get install -y curl
fi

if ! command_exists git; then
  echo "Installing git..."
  apt-get install -y git
fi

# Setup PHP repository based on OS
OS=$(detect_os)
echo -e "\n${YELLOW}Setting up PHP repository for ${OS}...${NC}"

case $OS in
  debian)
    echo "Setting up Sury PHP repository for Debian..."
    apt-get install -y apt-transport-https lsb-release ca-certificates
    wget -qO - https://packages.sury.org/php/apt.gpg | apt-key add -
    echo "deb https://packages.sury.org/php/ $(lsb_release -sc) main" | tee /etc/apt/sources.list.d/sury-php.list
    ;;
  ubuntu)
    echo "Setting up Ondrej PHP repository for Ubuntu..."
    apt-get install -y software-properties-common
    add-apt-repository -y ppa:ondrej/php
    ;;
  *)
    echo -e "${RED}Unsupported operating system: ${OS}${NC}"
    exit 1
    ;;
esac

# Update package list after adding PHP repository
apt-get update

# Install required packages
echo -e "\n${YELLOW}Installing required packages...${NC}"

# Add Caddy repository
if ! command_exists caddy; then
  echo "deb [trusted=yes] https://apt.fury.io/caddy/ /" | tee /etc/apt/sources.list.d/caddy-fury.list
  apt-get update
fi

# Install packages
apt-get install -y \
  curl \
  git \
  golang \
  caddy \
  mariadb-server \
  mariadb-client

# Create directory structure
echo -e "\n${YELLOW}Creating directory structure...${NC}"

directories=(
  "/apps/sites"
  "/usr/local/webpanel/config/global"
  "/usr/local/webpanel/config/modules"
  "/usr/local/webpanel/config/sites"
  "/usr/local/webpanel/lib/modules"
  "/usr/local/webpanel/logs"
  "/var/log/webpanel/caddy"
  "/var/log/webpanel/mariadb"
  "/var/log/webpanel/webpanel"
  "/backup/daily"
  "/backup/weekly"
)

for dir in "${directories[@]}"; do
  mkdir -p "$dir"
  echo "Created $dir"
done

# Set up MariaDB
echo -e "\n${YELLOW}Configuring MariaDB...${NC}"
mysql_secure_installation

# Configure Caddy modules
echo -e "\n${YELLOW}Creating Caddy module configurations...${NC}"

# Main Caddyfile
cat > /usr/local/webpanel/config/global/Caddyfile <<EOF
{
    admin off
    storage file_system {
        root /var/lib/caddy
    }
}

# Import all module configurations
import /usr/local/webpanel/config/modules/*.conf

# Import site configurations
import /usr/local/webpanel/config/sites/*.conf
EOF

# Create module configuration files
modules_dir="/usr/local/webpanel/config/modules"

# SPA module
cat > "$modules_dir/spa.conf" <<EOF
(spa_config) {
    try_files {path} /index.html
    encode gzip
    file_server
}
EOF

# Security module
cat > "$modules_dir/security.conf" <<EOF
(security_config) {
    header {
        Strict-Transport-Security "max-age=31536000"
        X-Content-Type-Options "nosniff"
        X-Frame-Options "DENY"
        X-XSS-Protection "1; mode=block"
        Content-Security-Policy "default-src 'self'"
        Referrer-Policy "no-referrer-when-downgrade"
    }
}
EOF

# Headers module
cat > "$modules_dir/headers.conf" <<EOF
(header_config) {
    header {
        -Server
        X-Powered-By "webpanel"
    }
}
EOF

# Access restriction module
cat > "$modules_dir/restrict.conf" <<EOF
(restrict_config) {
    @blocked not remote_ip {args.0}
    respond @blocked 403
}
EOF

# Logging modules
cat > "$modules_dir/logging.conf" <<EOF
(access_log) {
    log access {
        output file /var/log/webpanel/caddy/{args.0}/access.log
        format json
    }
}

(error_log) {
    log error {
        output file /var/log/webpanel/caddy/{args.0}/error.log
        format json
        level ERROR
    }
}
EOF

# Link Caddy configuration
ln -sf /usr/local/webpanel/config/global/Caddyfile /etc/caddy/Caddyfile

# Set up permissions
echo -e "\n${YELLOW}Setting up permissions...${NC}"
chown -R www-data:www-data /apps/sites
chmod -R 755 /apps/sites
chown -R root:root /usr/local/webpanel
chmod -R 755 /usr/local/webpanel
chown -R www-data:www-data /var/log/webpanel
chmod -R 755 /var/log/webpanel
chown -R www-data:www-data /backup
chmod -R 755 /backup

# Build and install webpanel CLI
echo -e "\n${YELLOW}Building and installing webpanel CLI...${NC}"
mkdir -p /tmp/webpanel-build
git clone https://github.com/doko89/cli-webpanel.git /tmp/webpanel-build
cd /tmp/webpanel-build
go build -o webpanel cmd/webpanel/main.go
mv webpanel /usr/local/bin/
cd /
rm -rf /tmp/webpanel-build

# Set up logrotate
echo -e "\n${YELLOW}Configuring log rotation...${NC}"
cat > /etc/logrotate.d/webpanel <<EOF
/var/log/webpanel/caddy/*/*.log {
    daily
    rotate 14
    compress
    delaycompress
    missingok
    notifempty
    create 0644 www-data www-data
}

/var/log/webpanel/mariadb/*.log {
    daily
    rotate 14
    compress
    delaycompress
    missingok
    notifempty
    create 0644 mysql mysql
}

/var/log/webpanel/webpanel/*.log {
    daily
    rotate 14
    compress
    delaycompress
    missingok
    notifempty
    create 0644 root root
}
EOF

# Create cron directory structure
mkdir -p /etc/cron.daily /etc/cron.weekly

# Set up daily cleanup of old backups
cat > /etc/cron.daily/webpanel-cleanup <<EOF
#!/bin/bash
find /backup/daily -type f -mtime +7 -delete
EOF
chmod +x /etc/cron.daily/webpanel-cleanup

# Set up weekly cleanup of old backups
cat > /etc/cron.weekly/webpanel-cleanup <<EOF
#!/bin/bash
find /backup/weekly -type f -mtime +30 -delete
EOF
chmod +x /etc/cron.weekly/webpanel-cleanup

# Display success message
echo -e "\n${GREEN}Installation completed successfully!${NC}"
echo -e "\nYou can now use the webpanel command to manage your server."
echo -e "Run ${YELLOW}webpanel --help${NC} to see available commands."
echo -e "\nTo install PHP:"
echo -e "Run ${YELLOW}webpanel php list available${NC} to see available PHP versions"
echo -e "Run ${YELLOW}webpanel php install <version>${NC} to install a PHP version"
echo -e "Run ${YELLOW}webpanel php module-available <version>${NC} to see available modules"
echo -e "\nNote: The webpanel command should be run as a non-root user."
echo -e "Example: ${YELLOW}sudo -u your_username webpanel status${NC}"