package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	// RootCmd is the root command for the CLI application
	RootCmd = &cobra.Command{
		Use:   "webpanel",
		Short: "A CLI tool for managing VPS/server with minimal sysadmin experience",
		Long: `webpanel is a command line tool designed to simplify server management.
It provides easy-to-use commands for managing websites, databases, backups,
and monitoring server resources.`,
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return RootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.webpanel.yaml)")

	// Initialize version command
	RootCmd.AddCommand(versionCmd)

	// Register all commands
	RegisterCommands(RootCmd)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".webpanel")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of webpanel",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("webpanel v0.1.0")
	},
}
