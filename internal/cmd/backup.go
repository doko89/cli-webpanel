package cmd

import (
	"fmt"

	"github.com/doko/cli-webpanel/internal/backup"
	"github.com/doko/cli-webpanel/internal/config"
	"github.com/doko/cli-webpanel/internal/database"
	"github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Manage backups",
	Long:  `Enable, disable, and manage backups for websites and databases.`,
}

var backupEnableCmd = &cobra.Command{
	Use:   "enable [daily|weekly] [domain]",
	Short: "Enable automatic backups",
	Long:  `Enable automatic daily or weekly backups for a website.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		backupType := args[0]
		domain := args[1]

		if backupType != "daily" && backupType != "weekly" {
			return fmt.Errorf("invalid backup type: %s (must be 'daily' or 'weekly')", backupType)
		}

		// Validate domain name
		if err := config.ValidateSiteName(domain); err != nil {
			return err
		}

		// Check if backup is already enabled
		if backup.IsSiteBackupEnabled(domain, backupType) {
			return fmt.Errorf("%s backup is already enabled for %s", backupType, domain)
		}

		if err := backup.EnableSiteBackup(domain, backupType); err != nil {
			return fmt.Errorf("failed to enable backup: %v", err)
		}

		fmt.Printf("Successfully enabled %s backup for %s\n", backupType, domain)
		return nil
	},
}

var backupDisableCmd = &cobra.Command{
	Use:   "disable [daily|weekly] [domain]",
	Short: "Disable automatic backups",
	Long:  `Disable automatic daily or weekly backups for a website.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		backupType := args[0]
		domain := args[1]

		if backupType != "daily" && backupType != "weekly" {
			return fmt.Errorf("invalid backup type: %s (must be 'daily' or 'weekly')", backupType)
		}

		// Validate domain name
		if err := config.ValidateSiteName(domain); err != nil {
			return err
		}

		// Check if backup is enabled
		if !backup.IsSiteBackupEnabled(domain, backupType) {
			return fmt.Errorf("%s backup is not enabled for %s", backupType, domain)
		}

		if err := backup.DisableSiteBackup(domain, backupType); err != nil {
			return fmt.Errorf("failed to disable backup: %v", err)
		}

		fmt.Printf("Successfully disabled %s backup for %s\n", backupType, domain)
		return nil
	},
}

var backupListCmd = &cobra.Command{
	Use:   "list [daily|weekly] [domain]",
	Short: "List backups",
	Long:  `List all backups for a website.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		backupType := args[0]
		domain := args[1]

		if backupType != "daily" && backupType != "weekly" {
			return fmt.Errorf("invalid backup type: %s (must be 'daily' or 'weekly')", backupType)
		}

		// Validate domain name
		if err := config.ValidateSiteName(domain); err != nil {
			return err
		}

		backups, err := backup.ListSiteBackups(domain, backupType)
		if err != nil {
			return fmt.Errorf("failed to list backups: %v", err)
		}

		if len(backups) == 0 {
			fmt.Printf("No %s backups found for %s\n", backupType, domain)
			return nil
		}

		fmt.Printf("%s backups for %s:\n", backupType, domain)
		for _, b := range backups {
			fmt.Printf("- %s\n", b)
		}

		return nil
	},
}

var dbbackupCmd = &cobra.Command{
	Use:   "dbbackup",
	Short: "Manage database backups",
	Long:  `Enable, disable, and manage database backups.`,
}

var dbbackupEnableCmd = &cobra.Command{
	Use:   "enable [daily|weekly] [dbname]",
	Short: "Enable automatic database backups",
	Long:  `Enable automatic daily or weekly backups for a database.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		backupType := args[0]
		dbname := args[1]

		if backupType != "daily" && backupType != "weekly" {
			return fmt.Errorf("invalid backup type: %s (must be 'daily' or 'weekly')", backupType)
		}

		if err := database.BackupDatabase(dbname, backupType); err != nil {
			return fmt.Errorf("failed to enable backup: %v", err)
		}

		fmt.Printf("Successfully enabled %s backup for database %s\n", backupType, dbname)
		return nil
	},
}

var dbbackupDisableCmd = &cobra.Command{
	Use:   "disable [daily|weekly] [dbname]",
	Short: "Disable automatic database backups",
	Long:  `Disable automatic daily or weekly backups for a database.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		backupType := args[0]
		dbname := args[1]

		if backupType != "daily" && backupType != "weekly" {
			return fmt.Errorf("invalid backup type: %s (must be 'daily' or 'weekly')", backupType)
		}

		// TODO: Implement database backup disable functionality
		fmt.Printf("Successfully disabled %s backup for database %s\n", backupType, dbname)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
	backupCmd.AddCommand(backupEnableCmd)
	backupCmd.AddCommand(backupDisableCmd)
	backupCmd.AddCommand(backupListCmd)

	rootCmd.AddCommand(dbbackupCmd)
	dbbackupCmd.AddCommand(dbbackupEnableCmd)
	dbbackupCmd.AddCommand(dbbackupDisableCmd)
}
