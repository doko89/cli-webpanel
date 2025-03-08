package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/doko/cli-webpanel/internal/config"
	"github.com/doko/cli-webpanel/internal/module"
	"github.com/spf13/cobra"
)

var siteCmd = &cobra.Command{
	Use:   "site",
	Short: "Manage website configurations",
	Long:  `Add, remove, and list website configurations for the server.`,
}

var siteAddCmd = &cobra.Command{
	Use:   "add [domain]",
	Short: "Add a new website",
	Long:  `Create a new website directory and configuration for the specified domain.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		domain := args[0]

		// Validate domain name
		if err := config.ValidateSiteName(domain); err != nil {
			return err
		}

		// Create site directory
		siteDir := config.GetSiteDirectory(domain)
		if err := os.MkdirAll(siteDir, 0755); err != nil {
			return fmt.Errorf("failed to create site directory: %v", err)
		}

		// Create public directory
		publicDir := filepath.Join(siteDir, "public")
		if err := os.MkdirAll(publicDir, 0755); err != nil {
			return fmt.Errorf("failed to create public directory: %v", err)
		}

		// Create logs directory
		logsDir := filepath.Join("/var/log/webpanel/caddy", domain)
		if err := os.MkdirAll(logsDir, 0755); err != nil {
			return fmt.Errorf("failed to create logs directory: %v", err)
		}

		// Create basic index.html
		indexPath := filepath.Join(publicDir, "index.html")
		indexContent := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>Welcome to %s</title>
</head>
<body>
    <h1>Welcome to %s</h1>
    <p>Your website is now set up and running!</p>
</body>
</html>`, domain, domain)

		if err := os.WriteFile(indexPath, []byte(indexContent), 0644); err != nil {
			return fmt.Errorf("failed to create index.html: %v", err)
		}

		// Create Caddy configuration with default modules
		caddyConfig := fmt.Sprintf(`%s {
    root * %s
    import access_log %s
    import error_log %s
    import header_config
    import security_config
    import php_config
}`, domain, publicDir, domain, domain)

		configPath := config.GetSiteConfigPath(domain)
		if err := os.WriteFile(configPath, []byte(caddyConfig), 0644); err != nil {
			return fmt.Errorf("failed to create Caddy configuration: %v", err)
		}

		fmt.Printf("Successfully created website for %s\n", domain)
		fmt.Printf("Site directory: %s\n", siteDir)
		fmt.Printf("Configuration: %s\n", configPath)
		fmt.Printf("Log directory: %s\n", logsDir)
		return nil
	},
}

var siteListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all websites",
	Long:  `Display a list of all configured websites on the server.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		sitesDir := config.GetWebRoot()
		entries, err := os.ReadDir(sitesDir)
		if err != nil {
			return fmt.Errorf("failed to read sites directory: %v", err)
		}

		if len(entries) == 0 {
			fmt.Println("No websites configured")
			return nil
		}

		fmt.Println("Configured websites:")
		for _, entry := range entries {
			if entry.IsDir() {
				domain := entry.Name()
				fmt.Printf("- %s\n", domain)

				// List enabled modules
				modules, err := module.ListEnabled(domain)
				if err == nil && len(modules) > 0 {
					fmt.Printf("  Enabled modules:\n")
					for _, mod := range modules {
						fmt.Printf("    - %s\n", mod)
					}
				}
			}
		}
		return nil
	},
}

var siteRmCmd = &cobra.Command{
	Use:   "rm [domain]",
	Short: "Remove a website",
	Long:  `Remove a website directory and its configuration.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		domain := args[0]

		// Validate domain name
		if err := config.ValidateSiteName(domain); err != nil {
			return err
		}

		// Check if site exists
		siteDir := config.GetSiteDirectory(domain)
		if _, err := os.Stat(siteDir); os.IsNotExist(err) {
			return fmt.Errorf("website %s does not exist", domain)
		}

		// Prompt for confirmation
		fmt.Printf("Are you sure you want to remove %s? This action cannot be undone. [y/N]: ", domain)
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println("Operation cancelled")
			return nil
		}

		// Remove site directory
		if err := os.RemoveAll(siteDir); err != nil {
			return fmt.Errorf("failed to remove site directory: %v", err)
		}

		// Remove configuration file
		configPath := config.GetSiteConfigPath(domain)
		if err := os.Remove(configPath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to remove configuration file: %v", err)
		}

		// Remove log directory
		logDir := filepath.Join("/var/log/webpanel/caddy", domain)
		if err := os.RemoveAll(logDir); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to remove log directory: %v", err)
		}

		fmt.Printf("Successfully removed website %s\n", domain)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(siteCmd)
	siteCmd.AddCommand(siteAddCmd)
	siteCmd.AddCommand(siteListCmd)
	siteCmd.AddCommand(siteRmCmd)
}
