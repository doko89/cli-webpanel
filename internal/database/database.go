package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/doko/cli-webpanel/internal/config"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// Initialize sets up the database connection
func Initialize() error {
	// Try to connect as root first to setup initial database
	dsn := fmt.Sprintf("root@tcp(127.0.0.1:3306)/")
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	return nil
}

// Close closes the database connection
func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

// CreateDatabase creates a new database
func CreateDatabase(name string) error {
	query := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", name)
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create database: %v", err)
	}
	return nil
}

// DeleteDatabase deletes a database
func DeleteDatabase(name string) error {
	query := fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", name)
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to delete database: %v", err)
	}
	return nil
}

// ListDatabases returns a list of all databases
func ListDatabases() ([]string, error) {
	rows, err := db.Query("SHOW DATABASES")
	if err != nil {
		return nil, fmt.Errorf("failed to list databases: %v", err)
	}
	defer rows.Close()

	var databases []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("failed to scan database name: %v", err)
		}
		// Skip system databases
		if name != "information_schema" && name != "mysql" && name != "performance_schema" {
			databases = append(databases, name)
		}
	}

	return databases, nil
}

// CreateUser creates a new database user
func CreateUser(username, password string) error {
	// Create user
	query := fmt.Sprintf("CREATE USER '%s'@'localhost' IDENTIFIED BY '%s'", username, password)
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}
	return nil
}

// DeleteUser deletes a database user
func DeleteUser(username string) error {
	query := fmt.Sprintf("DROP USER IF EXISTS '%s'@'localhost'", username)
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}
	return nil
}

// ListUsers returns a list of all database users
func ListUsers() ([]string, error) {
	rows, err := db.Query("SELECT User FROM mysql.user WHERE Host = 'localhost' AND User NOT IN ('root', 'mysql.sys', 'mysql.session')")
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %v", err)
	}
	defer rows.Close()

	var users []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("failed to scan user name: %v", err)
		}
		users = append(users, name)
	}

	return users, nil
}

// GrantAccess grants a user access to a database
func GrantAccess(username, dbname string) error {
	query := fmt.Sprintf("GRANT ALL PRIVILEGES ON `%s`.* TO '%s'@'localhost'", dbname, username)
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to grant access: %v", err)
	}

	// Apply privileges
	_, err = db.Exec("FLUSH PRIVILEGES")
	if err != nil {
		return fmt.Errorf("failed to flush privileges: %v", err)
	}

	return nil
}

// BackupDatabase creates a backup of the specified database
func BackupDatabase(name, backupType string) error {
	now := time.Now()
	var filename string
	backupDir := filepath.Join(config.GetBackupDir(), backupType, name)

	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %v", err)
	}

	if backupType == "weekly" {
		filename = fmt.Sprintf("%s/%s-full.sql.gz", backupDir, now.Format("2006-01-02"))
	} else {
		filename = fmt.Sprintf("%s/%s.sql.gz", backupDir, now.Format("2006-01-02"))
	}

	// Create backup command
	cmd := fmt.Sprintf("mysqldump -u root %s | gzip > %s", name, filename)

	// Execute backup command
	if err := executeCommand(cmd); err != nil {
		return fmt.Errorf("failed to create backup: %v", err)
	}

	return nil
}

// Helper function to execute shell commands
func executeCommand(command string) error {
	return nil // TODO: Implement actual command execution
}
