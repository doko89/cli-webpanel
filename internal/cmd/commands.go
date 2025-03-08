package cmd

import "github.com/spf13/cobra"

// RegisterCommands initializes and registers all commands to the root command
func RegisterCommands(root *cobra.Command) {
	// Initialize commands
	initBackupCommands(root)
	initDatabaseCommands(root)
	initMonitorCommands(root)
	initModuleCommands(root)
	initPHPCommands(root)
	initSiteCommands(root)
}

// initBackupCommands registers all backup related commands
func initBackupCommands(root *cobra.Command) {
	root.AddCommand(backupCmd)
	backupCmd.AddCommand(backupEnableCmd)
	backupCmd.AddCommand(backupDisableCmd)
	backupCmd.AddCommand(backupListCmd)

	root.AddCommand(dbbackupCmd)
	dbbackupCmd.AddCommand(dbbackupEnableCmd)
	dbbackupCmd.AddCommand(dbbackupDisableCmd)
}

// initDatabaseCommands registers all database related commands
func initDatabaseCommands(root *cobra.Command) {
	root.AddCommand(dbCmd)
	dbCmd.AddCommand(dbListCmd)
	dbCmd.AddCommand(dbCreateCmd)
	dbCmd.AddCommand(dbDeleteCmd)

	root.AddCommand(dbuserCmd)
	dbuserCmd.AddCommand(dbuserListCmd)
	dbuserCmd.AddCommand(dbuserCreateCmd)
	dbuserCmd.AddCommand(dbuserDeleteCmd)

	root.AddCommand(dbgrantCmd)
}

// initMonitorCommands registers all monitoring related commands
func initMonitorCommands(root *cobra.Command) {
	root.AddCommand(statusCmd)
	root.AddCommand(logsCmd)
	root.AddCommand(monitorCmd)
}

// initModuleCommands registers all module related commands
func initModuleCommands(root *cobra.Command) {
	root.AddCommand(moduleCmd)
	moduleCmd.AddCommand(moduleListAvailableCmd)
	moduleCmd.AddCommand(moduleListCmd)
	moduleCmd.AddCommand(moduleAddCmd)
	moduleCmd.AddCommand(moduleRmCmd)
}

// initPHPCommands registers all PHP related commands
func initPHPCommands(root *cobra.Command) {
	root.AddCommand(phpCmd)
	phpCmd.AddCommand(phpListCmd)
	phpCmd.AddCommand(phpInstallCmd)
	phpCmd.AddCommand(phpRemoveCmd)
	phpCmd.AddCommand(phpModuleAvailableCmd)
	phpCmd.AddCommand(phpModuleInstallCmd)
	phpCmd.AddCommand(phpModuleRemoveCmd)
}

// initSiteCommands registers all site related commands
func initSiteCommands(root *cobra.Command) {
	root.AddCommand(siteCmd)
	siteCmd.AddCommand(siteAddCmd)
	siteCmd.AddCommand(siteListCmd)
	siteCmd.AddCommand(siteRmCmd)
}
