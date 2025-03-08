package module

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/doko/cli-webpanel/internal/config"
)

// ModuleConfig represents a module's configuration
type ModuleConfig struct {
	Name       string            `json:"name"`
	Version    string            `json:"version"`
	Enabled    bool              `json:"enabled"`
	Parameters map[string]string `json:"parameters"`
}

// Module represents a server module
type Module struct {
	Name       string
	ConfigPath string
	Template   string
	Args       []string
}

var availableModules = map[string]Module{
	"php": {
		Name: "php",
		Template: `(php_config) {
    php_fastcgi unix//run/php/php8.1-fpm.sock
    encode gzip
    file_server
}`,
	},

	"spa": {
		Name: "spa",
		Template: `(spa_config) {
    try_files {path} /index.html
    encode gzip
    file_server
}`,
	},

	"security": {
		Name: "security",
		Template: `(security_config) {
    header {
        Strict-Transport-Security "max-age=31536000"
        X-Content-Type-Options "nosniff"
        X-Frame-Options "DENY"
        X-XSS-Protection "1; mode=block"
        Content-Security-Policy "default-src 'self'"
        Referrer-Policy "no-referrer-when-downgrade"
    }
}`,
	},

	"header": {
		Name: "header",
		Template: `(header_config) {
    header {
        -Server
        X-Powered-By "webpanel"
    }
}`,
	},

	"restrict": {
		Name: "restrict",
		Template: `(restrict_config) {
    @blocked not remote_ip {args.0}
    respond @blocked 403
}`,
	},

	"access_log": {
		Name: "access_log",
		Template: `(access_log) {
    log access {
        output file /var/log/webpanel/caddy/{args.0}/access.log
        format json
    }
}`,
	},

	"error_log": {
		Name: "error_log",
		Template: `(error_log) {
    log error {
        output file /var/log/webpanel/caddy/{args.0}/error.log
        format json
        level ERROR
    }
}`,
	},
}

// InitializeModules sets up the module configuration files
func InitializeModules() error {
	modulesDir := filepath.Join(config.GetConfigDir(), "modules")
	if err := os.MkdirAll(modulesDir, 0755); err != nil {
		return fmt.Errorf("failed to create modules directory: %v", err)
	}

	for name, module := range availableModules {
		configPath := filepath.Join(modulesDir, name+".conf")
		err := os.WriteFile(configPath, []byte(module.Template), 0644)
		if err != nil {
			return fmt.Errorf("failed to write module configuration %s: %v", name, err)
		}
	}

	return nil
}

// ListAvailable returns a list of all available modules
func ListAvailable() []string {
	modules := make([]string, 0, len(availableModules))
	for name := range availableModules {
		modules = append(modules, name)
	}
	return modules
}

// GetModule returns a module by name
func GetModule(name string) (Module, error) {
	module, ok := availableModules[name]
	if !ok {
		return Module{}, fmt.Errorf("module %s not found", name)
	}
	return module, nil
}

// IsModuleEnabled checks if a module is enabled for a domain
func IsModuleEnabled(moduleName, domain string) bool {
	configPath := filepath.Join(config.GetConfigDir(), "sites", domain+".conf")
	content, err := os.ReadFile(configPath)
	if err != nil {
		return false
	}
	return strings.Contains(string(content), fmt.Sprintf("import %s", moduleName))
}

// ListEnabled returns a list of enabled modules for a domain
func ListEnabled(domain string) ([]string, error) {
	configPath := filepath.Join(config.GetConfigDir(), "sites", domain+".conf")
	content, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read site configuration: %v", err)
	}

	var enabledModules []string
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "import ") {
			moduleName := strings.Fields(line)[1]
			enabledModules = append(enabledModules, moduleName)
		}
	}

	return enabledModules, nil
}

// EnableModule enables a module for a domain
func EnableModule(moduleName, domain string, params []string) error {
	if _, err := GetModule(moduleName); err != nil {
		return err
	}

	// Create site config directory if it doesn't exist
	sitesDir := filepath.Join(config.GetConfigDir(), "sites")
	if err := os.MkdirAll(sitesDir, 0755); err != nil {
		return fmt.Errorf("failed to create sites directory: %v", err)
	}

	// Create log directory for the domain if using logging modules
	if moduleName == "access_log" || moduleName == "error_log" {
		logDir := filepath.Join("/var/log/webpanel/caddy", domain)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return fmt.Errorf("failed to create log directory: %v", err)
		}
	}

	// Update site configuration to include module
	configPath := filepath.Join(sitesDir, domain+".conf")
	content, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read site configuration: %v", err)
	}

	// Add module import if not already present
	importLine := fmt.Sprintf("    import %s", moduleName)
	if len(params) > 0 {
		importLine += " " + strings.Join(params, " ")
	}
	importLine += "\n"

	if !strings.Contains(string(content), importLine) {
		// Find the closing brace
		lines := strings.Split(string(content), "\n")
		var newLines []string
		found := false
		for _, line := range lines {
			if strings.TrimSpace(line) == "}" && !found {
				newLines = append(newLines, importLine)
				found = true
			}
			newLines = append(newLines, line)
		}
		content = []byte(strings.Join(newLines, "\n"))

		if err := os.WriteFile(configPath, content, 0644); err != nil {
			return fmt.Errorf("failed to update site configuration: %v", err)
		}
	}

	return nil
}

// DisableModule disables a module for a domain
func DisableModule(moduleName, domain string) error {
	configPath := filepath.Join(config.GetConfigDir(), "sites", domain+".conf")
	content, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read site configuration: %v", err)
	}

	// Remove module import
	lines := strings.Split(string(content), "\n")
	var newLines []string
	for _, line := range lines {
		if !strings.Contains(line, fmt.Sprintf("import %s", moduleName)) {
			newLines = append(newLines, line)
		}
	}

	if err := os.WriteFile(configPath, []byte(strings.Join(newLines, "\n")), 0644); err != nil {
		return fmt.Errorf("failed to update site configuration: %v", err)
	}

	return nil
}
