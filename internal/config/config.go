package config

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	DefaultWebRoot   = "/apps/sites"
	DefaultConfigDir = "/usr/local/webpanel/config"
	DefaultBackupDir = "/backup"
	DefaultLogDir    = "/usr/local/webpanel/logs"
	DefaultModuleDir = "/usr/local/webpanel/lib/modules"
)

type Config struct {
	WebRoot   string `mapstructure:"web_root"`
	ConfigDir string `mapstructure:"config_dir"`
	BackupDir string `mapstructure:"backup_dir"`
	LogDir    string `mapstructure:"log_dir"`
	ModuleDir string `mapstructure:"module_dir"`
}

// Global configuration instance
var globalConfig *Config

// Init initializes the configuration system
func Init() error {
	// Create default configuration
	globalConfig = &Config{
		WebRoot:   DefaultWebRoot,
		ConfigDir: DefaultConfigDir,
		BackupDir: DefaultBackupDir,
		LogDir:    DefaultLogDir,
		ModuleDir: DefaultModuleDir,
	}

	// Ensure directories exist
	dirs := []string{
		globalConfig.WebRoot,
		globalConfig.ConfigDir,
		globalConfig.BackupDir,
		globalConfig.LogDir,
		globalConfig.ModuleDir,
		filepath.Join(globalConfig.ConfigDir, "sites"),
		filepath.Join(globalConfig.ConfigDir, "global"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", dir, err)
		}
	}

	return nil
}

// GetConfig returns the global configuration instance
func GetConfig() *Config {
	return globalConfig
}

// SetConfig updates the global configuration
func SetConfig(cfg *Config) {
	globalConfig = cfg
}

// GetWebRoot returns the configured web root directory
func GetWebRoot() string {
	return globalConfig.WebRoot
}

// GetConfigDir returns the configured config directory
func GetConfigDir() string {
	return globalConfig.ConfigDir
}

// GetBackupDir returns the configured backup directory
func GetBackupDir() string {
	return globalConfig.BackupDir
}

// GetLogDir returns the configured log directory
func GetLogDir() string {
	return globalConfig.LogDir
}

// GetModuleDir returns the configured module directory
func GetModuleDir() string {
	return globalConfig.ModuleDir
}

// ValidateSiteName checks if a site name is valid
func ValidateSiteName(domain string) error {
	if domain == "" {
		return fmt.Errorf("domain name cannot be empty")
	}
	// Add more domain validation if needed
	return nil
}

// GetSiteDirectory returns the full path to a site's directory
func GetSiteDirectory(domain string) string {
	return filepath.Join(GetWebRoot(), domain)
}

// GetSiteConfigPath returns the full path to a site's configuration file
func GetSiteConfigPath(domain string) string {
	return filepath.Join(GetConfigDir(), "sites", domain+".conf")
}
