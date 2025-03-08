# Modular Caddy Configuration Design

## Directory Structure

```
/usr/local/webpanel/
├── config/
│   ├── global/             # Global configurations
│   │   └── Caddyfile      # Main Caddy configuration
│   ├── modules/           # Modular configuration snippets
│   │   ├── php.conf       # PHP configuration
│   │   ├── spa.conf       # SPA configuration
│   │   ├── security.conf  # Security headers
│   │   ├── headers.conf   # Common headers
│   │   ├── restrict.conf  # Access restrictions
│   │   └── logging.conf   # Logging configuration
│   └── sites/             # Site-specific configurations
└── lib/
    └── modules/           # Module implementation files

/var/log/webpanel/
├── caddy/                 # Caddy logs directory
│   ├── domain.com/
│   │   ├── access.log
│   │   └── error.log
│   └── ...
├── mariadb/              # MariaDB logs
└── webpanel/            # CLI tool logs
```

## Module Configurations

### 1. Main Caddyfile (/usr/local/webpanel/config/global/Caddyfile)
```
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
```

### 2. Module Snippets

#### PHP Configuration (modules/php.conf)
```
(php_config) {
    php_fastcgi unix//run/php/php8.1-fpm.sock
    encode gzip
    file_server
}
```

#### SPA Configuration (modules/spa.conf)
```
(spa_config) {
    try_files {path} /index.html
    encode gzip
    file_server
}
```

#### Security Configuration (modules/security.conf)
```
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
```

#### Headers Configuration (modules/headers.conf)
```
(header_config) {
    header {
        -Server
        X-Powered-By "webpanel"
    }
}
```

#### Access Restriction (modules/restrict.conf)
```
(restrict_config) {
    @blocked not remote_ip {args.0}
    respond @blocked 403
}
```

#### Logging Configuration (modules/logging.conf)
```
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
```

## Site Configuration Example

When creating a new site with PHP support and logging:

```
domain.com {
    root * /apps/sites/domain.com/public
    import php_config
    import security_config
    import header_config
    import access_log domain.com
    import error_log domain.com
}
```

## Implementation Changes Needed

1. Update installation script to:
   - Create proper directory structure in /usr/local/webpanel
   - Create log directories in /var/log/webpanel
   - Set appropriate permissions
   - Install base module configurations
   - Create main Caddyfile with imports

2. Update module package to:
   - Store module configurations in /usr/local/webpanel/config/modules/
   - Handle module dependencies
   - Support parameter substitution
   - Manage module state

3. Update site management to:
   - Generate site configs using module imports
   - Create log directories for each site
   - Set appropriate permissions

## Benefits

1. Maintains existing directory structure
2. Centralized logging in /var/log/webpanel
3. Better organization of configurations
4. Reusable configuration snippets
5. Easier maintenance and updates
6. Standardized logging format

## Migration Plan

1. Update directory structure
2. Create module configurations
3. Update site management
4. Set up log directories and permissions
5. Configure log rotation