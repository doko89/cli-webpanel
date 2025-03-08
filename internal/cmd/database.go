package cmd

import (
	"fmt"

	"github.com/doko/cli-webpanel/internal/database"
	"github.com/spf13/cobra"
)

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Manage databases",
	Long:  `Create, delete, and list databases.`,
}

var dbListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all databases",
	Long:  `Display a list of all databases on the server.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		databases, err := database.ListDatabases()
		if err != nil {
			return err
		}

		if len(databases) == 0 {
			fmt.Println("No databases found")
			return nil
		}

		fmt.Println("Databases:")
		for _, name := range databases {
			fmt.Printf("- %s\n", name)
		}
		return nil
	},
}

var dbCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new database",
	Long:  `Create a new database with the specified name.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		if err := database.CreateDatabase(name); err != nil {
			return err
		}

		fmt.Printf("Successfully created database '%s'\n", name)
		return nil
	},
}

var dbDeleteCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete a database",
	Long:  `Delete the specified database.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		// Prompt for confirmation
		fmt.Printf("Are you sure you want to delete database '%s'? This action cannot be undone. [y/N]: ", name)
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println("Operation cancelled")
			return nil
		}

		if err := database.DeleteDatabase(name); err != nil {
			return err
		}

		fmt.Printf("Successfully deleted database '%s'\n", name)
		return nil
	},
}

var dbuserCmd = &cobra.Command{
	Use:   "dbuser",
	Short: "Manage database users",
	Long:  `Create, delete, and list database users.`,
}

var dbuserListCmd = &cobra.Command{
	Use:   "list",
	Short: "List database users",
	Long:  `Display a list of all database users.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		users, err := database.ListUsers()
		if err != nil {
			return err
		}

		if len(users) == 0 {
			fmt.Println("No database users found")
			return nil
		}

		fmt.Println("Database users:")
		for _, name := range users {
			fmt.Printf("- %s\n", name)
		}
		return nil
	},
}

var dbuserCreateCmd = &cobra.Command{
	Use:   "create [username] [password]",
	Short: "Create a database user",
	Long:  `Create a new database user with the specified username and password.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		username := args[0]
		password := args[1]

		if err := database.CreateUser(username, password); err != nil {
			return err
		}

		fmt.Printf("Successfully created database user '%s'\n", username)
		return nil
	},
}

var dbuserDeleteCmd = &cobra.Command{
	Use:   "delete [username]",
	Short: "Delete a database user",
	Long:  `Delete the specified database user.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		username := args[0]

		// Prompt for confirmation
		fmt.Printf("Are you sure you want to delete user '%s'? This action cannot be undone. [y/N]: ", username)
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println("Operation cancelled")
			return nil
		}

		if err := database.DeleteUser(username); err != nil {
			return err
		}

		fmt.Printf("Successfully deleted database user '%s'\n", username)
		return nil
	},
}

var dbgrantCmd = &cobra.Command{
	Use:   "dbgrant [username] [dbname]",
	Short: "Grant database access",
	Long:  `Grant a user access to the specified database.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		username := args[0]
		dbname := args[1]

		if err := database.GrantAccess(username, dbname); err != nil {
			return err
		}

		fmt.Printf("Successfully granted access to database '%s' for user '%s'\n", dbname, username)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(dbCmd)
	dbCmd.AddCommand(dbListCmd)
	dbCmd.AddCommand(dbCreateCmd)
	dbCmd.AddCommand(dbDeleteCmd)

	rootCmd.AddCommand(dbuserCmd)
	dbuserCmd.AddCommand(dbuserListCmd)
	dbuserCmd.AddCommand(dbuserCreateCmd)
	dbuserCmd.AddCommand(dbuserDeleteCmd)

	rootCmd.AddCommand(dbgrantCmd)
}
