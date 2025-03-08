package cmd

import (
	"fmt"
	"strings"

	"github.com/doko/cli-webpanel/internal/config"
	"github.com/doko/cli-webpanel/internal/module"
	"github.com/spf13/cobra"
)

var moduleCmd = &cobra.Command{
	Use:   "module",
	Short: "Manage server modules",
	Long:  `Add, remove, and list server modules for websites.`,
}

var moduleListAvailableCmd = &cobra.Command{
	Use:   "list-available",
	Short: "List all available modules",
	Long:  `Display a list of all modules that can be enabled for websites.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		modules := module.ListAvailable()
		if len(modules) == 0 {
			fmt.Println("No modules available")
			return nil
		}

		fmt.Println("Available modules:")
		for _, name := range modules {
			fmt.Printf("- %s\n", name)
		}
		return nil
	},
}

var moduleListCmd = &cobra.Command{
	Use:   "list [domain]",
	Short: "List enabled modules for a domain",
	Long:  `Display a list of all modules enabled for the specified domain.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		domain := args[0]

		// Validate domain name
		if err := config.ValidateSiteName(domain); err != nil {
			return err
		}

		enabled, err := module.ListEnabled(domain)
		if err != nil {
			return err
		}

		if len(enabled) == 0 {
			fmt.Printf("No modules enabled for %s\n", domain)
			return nil
		}

		fmt.Printf("Enabled modules for %s:\n", domain)
		for _, name := range enabled {
			fmt.Printf("- %s\n", name)
		}
		return nil
	},
}

var moduleAddCmd = &cobra.Command{
	Use:   "add [module] [domain] [params...]",
	Short: "Add a module to a domain",
	Long: `Enable a module for the specified domain. Some modules may require additional parameters.
Example: webpanel module add restrict-access example.com 192.168.1.0/24`,
	Args: cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		moduleName := args[0]
		domain := args[1]
		params := args[2:]

		// Validate domain name
		if err := config.ValidateSiteName(domain); err != nil {
			return err
		}

		// Check if module exists
		if _, err := module.GetModule(moduleName); err != nil {
			return err
		}

		// Check if module is already enabled
		if module.IsModuleEnabled(moduleName, domain) {
			return fmt.Errorf("module %s is already enabled for %s", moduleName, domain)
		}

		// Enable module
		if err := module.EnableModule(moduleName, domain, params); err != nil {
			return err
		}

		fmt.Printf("Successfully enabled %s module for %s\n", moduleName, domain)
		if len(params) > 0 {
			fmt.Printf("Parameters: %s\n", strings.Join(params, " "))
		}
		return nil
	},
}

var moduleRmCmd = &cobra.Command{
	Use:   "rm [module] [domain]",
	Short: "Remove a module from a domain",
	Long:  `Disable a module for the specified domain.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		moduleName := args[0]
		domain := args[1]

		// Validate domain name
		if err := config.ValidateSiteName(domain); err != nil {
			return err
		}

		// Check if module exists
		if _, err := module.GetModule(moduleName); err != nil {
			return err
		}

		// Check if module is enabled
		if !module.IsModuleEnabled(moduleName, domain) {
			return fmt.Errorf("module %s is not enabled for %s", moduleName, domain)
		}

		// Disable module
		if err := module.DisableModule(moduleName, domain); err != nil {
			return err
		}

		fmt.Printf("Successfully disabled %s module for %s\n", moduleName, domain)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(moduleCmd)
	moduleCmd.AddCommand(moduleListAvailableCmd)
	moduleCmd.AddCommand(moduleListCmd)
	moduleCmd.AddCommand(moduleAddCmd)
	moduleCmd.AddCommand(moduleRmCmd)
}
