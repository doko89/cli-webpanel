package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var phpCmd = &cobra.Command{
	Use:   "php",
	Short: "Manage PHP versions and modules",
	Long:  `Install, remove, and manage different PHP versions and their modules.`,
}

var phpListCmd = &cobra.Command{
	Use:   "list [available|installed]",
	Short: "List PHP versions",
	Long:  `List available or installed PHP versions.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		listType := args[0]

		switch listType {
		case "available":
			versions := []string{"7.4", "8.0", "8.1", "8.2"}
			fmt.Println("Available PHP versions:")
			for _, version := range versions {
				fmt.Printf("- %s\n", version)
			}
			return nil

		case "installed":
			return listInstalledPHPVersions()

		default:
			return fmt.Errorf("invalid list type: %s (must be 'available' or 'installed')", listType)
		}
	},
}

var phpInstallCmd = &cobra.Command{
	Use:   "install [version]",
	Short: "Install PHP version",
	Long:  `Install specified PHP version and its default modules.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		version := args[0]
		return installPHP(version)
	},
}

var phpRemoveCmd = &cobra.Command{
	Use:   "rm [version]",
	Short: "Remove PHP version",
	Long:  `Remove specified PHP version and its modules.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		version := args[0]
		return removePHP(version)
	},
}

var phpModuleAvailableCmd = &cobra.Command{
	Use:   "module-available [version]",
	Short: "List available PHP modules",
	Long:  `List all available modules for specified PHP version.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		version := args[0]
		return listAvailableModules(version)
	},
}

var phpModuleInstallCmd = &cobra.Command{
	Use:   "module-install [version] [module]",
	Short: "Install PHP module",
	Long:  `Install specified module for PHP version.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		version := args[0]
		module := args[1]
		return installModule(version, module)
	},
}

var phpModuleRemoveCmd = &cobra.Command{
	Use:   "module-remove [version] [module]",
	Short: "Remove PHP module",
	Long:  `Remove specified module from PHP version.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		version := args[0]
		module := args[1]
		return removeModule(version, module)
	},
}

func init() {
	RootCmd.AddCommand(phpCmd)
	phpCmd.AddCommand(phpListCmd)
	phpCmd.AddCommand(phpInstallCmd)
	phpCmd.AddCommand(phpRemoveCmd)
	phpCmd.AddCommand(phpModuleAvailableCmd)
	phpCmd.AddCommand(phpModuleInstallCmd)
	phpCmd.AddCommand(phpModuleRemoveCmd)
}

func listInstalledPHPVersions() error {
	cmd := exec.Command("sh", "-c", "ls /usr/sbin/php-fpm* | grep -o '[0-9]\\.[0-9]'")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to list installed PHP versions: %v", err)
	}

	versions := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(versions) == 0 || (len(versions) == 1 && versions[0] == "") {
		fmt.Println("No PHP versions installed")
		return nil
	}

	fmt.Println("Installed PHP versions:")
	for _, version := range versions {
		fmt.Printf("- %s\n", version)
	}
	return nil
}

func installPHP(version string) error {
	// Validate version format
	if !regexp.MustCompile(`^[0-9]\.[0-9]$`).MatchString(version) {
		return fmt.Errorf("invalid PHP version format (must be like '8.1')")
	}

	// Install PHP packages
	packages := []string{
		fmt.Sprintf("php%s-fpm", version),
		fmt.Sprintf("php%s-common", version),
		fmt.Sprintf("php%s-mysql", version),
		fmt.Sprintf("php%s-curl", version),
		fmt.Sprintf("php%s-gd", version),
		fmt.Sprintf("php%s-mbstring", version),
		fmt.Sprintf("php%s-xml", version),
		fmt.Sprintf("php%s-zip", version),
	}

	cmd := exec.Command("apt-get", append([]string{"install", "-y"}, packages...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install PHP packages: %v", err)
	}

	// Create Caddy module configuration
	moduleConfig := fmt.Sprintf(`(php%s_config) {
    php_fastcgi unix//run/php/php%s-fpm.sock
}`, strings.ReplaceAll(version, ".", ""), version)

	configPath := fmt.Sprintf("/usr/local/webpanel/config/modules/php%s.conf", strings.ReplaceAll(version, ".", ""))
	if err := os.WriteFile(configPath, []byte(moduleConfig), 0644); err != nil {
		return fmt.Errorf("failed to create PHP module configuration: %v", err)
	}

	fmt.Printf("Successfully installed PHP %s\n", version)
	return nil
}

func removePHP(version string) error {
	// Validate version format
	if !regexp.MustCompile(`^[0-9]\.[0-9]$`).MatchString(version) {
		return fmt.Errorf("invalid PHP version format (must be like '8.1')")
	}

	// Remove PHP packages
	packages := []string{
		fmt.Sprintf("php%s*", version),
	}

	cmd := exec.Command("apt-get", "remove", "-y", strings.Join(packages, " "))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to remove PHP packages: %v", err)
	}

	// Remove Caddy module configuration
	configPath := fmt.Sprintf("/usr/local/webpanel/config/modules/php%s.conf", strings.ReplaceAll(version, ".", ""))
	if err := os.Remove(configPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove PHP module configuration: %v", err)
	}

	fmt.Printf("Successfully removed PHP %s\n", version)
	return nil
}

func listAvailableModules(version string) error {
	cmd := exec.Command("apt-cache", "search", fmt.Sprintf("php%s-", version))
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to list available modules: %v", err)
	}

	fmt.Printf("Available modules for PHP %s:\n", version)
	for _, line := range strings.Split(string(output), "\n") {
		if strings.TrimSpace(line) != "" {
			parts := strings.SplitN(line, " - ", 2)
			if len(parts) >= 2 {
				moduleName := strings.TrimPrefix(parts[0], fmt.Sprintf("php%s-", version))
				fmt.Printf("- %s: %s\n", moduleName, parts[1])
			}
		}
	}
	return nil
}

func installModule(version, module string) error {
	// Install module
	packageName := fmt.Sprintf("php%s-%s", version, module)
	cmd := exec.Command("apt-get", "install", "-y", packageName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install module %s: %v", packageName, err)
	}

	// Reload PHP-FPM
	reload := exec.Command("systemctl", "reload", fmt.Sprintf("php%s-fpm", version))
	if err := reload.Run(); err != nil {
		return fmt.Errorf("failed to reload PHP-FPM: %v", err)
	}

	fmt.Printf("Successfully installed module %s for PHP %s\n", module, version)
	return nil
}

func removeModule(version, module string) error {
	// Remove module
	packageName := fmt.Sprintf("php%s-%s", version, module)
	cmd := exec.Command("apt-get", "remove", "-y", packageName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to remove module %s: %v", packageName, err)
	}

	// Reload PHP-FPM
	reload := exec.Command("systemctl", "reload", fmt.Sprintf("php%s-fpm", version))
	if err := reload.Run(); err != nil {
		return fmt.Errorf("failed to reload PHP-FPM: %v", err)
	}

	fmt.Printf("Successfully removed module %s from PHP %s\n", module, version)
	return nil
}
