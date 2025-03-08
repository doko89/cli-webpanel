package backup

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/doko/cli-webpanel/internal/config"
)

const (
	DailyBackup  = "daily"
	WeeklyBackup = "weekly"
)

// BackupSite performs a backup of the specified site
func BackupSite(domain, backupType string) error {
	siteDir := config.GetSiteDirectory(domain)
	if _, err := os.Stat(siteDir); os.IsNotExist(err) {
		return fmt.Errorf("site directory does not exist: %s", siteDir)
	}

	now := time.Now()
	backupDir := filepath.Join(config.GetBackupDir(), backupType, domain)
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %v", err)
	}

	var filename string
	if backupType == WeeklyBackup {
		filename = fmt.Sprintf("%s-full.tar.gz", now.Format("2006-01-02"))
	} else {
		filename = fmt.Sprintf("%s.tar.gz", now.Format("2006-01-02"))
	}

	backupPath := filepath.Join(backupDir, filename)

	// Create tar.gz archive
	cmd := exec.Command("tar", "-czf", backupPath, "-C", filepath.Dir(siteDir), filepath.Base(siteDir))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create backup archive: %v", err)
	}

	// Clean up old backups
	if err := cleanOldBackups(backupDir, backupType); err != nil {
		return fmt.Errorf("failed to clean old backups: %v", err)
	}

	return nil
}

// EnableSiteBackup enables automatic backups for a site
func EnableSiteBackup(domain, backupType string) error {
	cronFile := fmt.Sprintf("/etc/cron.%s/webpanel-backup-%s-%s",
		backupType, backupType, domain)

	var schedule string
	if backupType == DailyBackup {
		schedule = "0 1 * * * " // Run at 1 AM daily
	} else {
		schedule = "0 2 * * 0 " // Run at 2 AM on Sundays
	}

	command := fmt.Sprintf("webpanel backup %s %s\n", backupType, domain)
	content := schedule + command

	if err := os.WriteFile(cronFile, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to create cron job: %v", err)
	}

	return nil
}

// DisableSiteBackup disables automatic backups for a site
func DisableSiteBackup(domain, backupType string) error {
	cronFile := fmt.Sprintf("/etc/cron.%s/webpanel-backup-%s-%s",
		backupType, backupType, domain)

	if err := os.Remove(cronFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove cron job: %v", err)
	}

	return nil
}

// cleanOldBackups removes old backups based on retention policy
func cleanOldBackups(backupDir, backupType string) error {
	entries, err := os.ReadDir(backupDir)
	if err != nil {
		return err
	}

	var maxAge time.Duration
	if backupType == DailyBackup {
		maxAge = 7 * 24 * time.Hour // Keep daily backups for 7 days
	} else {
		maxAge = 30 * 24 * time.Hour // Keep weekly backups for 30 days
	}

	now := time.Now()
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		if now.Sub(info.ModTime()) > maxAge {
			path := filepath.Join(backupDir, entry.Name())
			if err := os.Remove(path); err != nil {
				fmt.Printf("Warning: failed to remove old backup %s: %v\n", path, err)
			}
		}
	}

	return nil
}

// IsSiteBackupEnabled checks if backup is enabled for a site
func IsSiteBackupEnabled(domain, backupType string) bool {
	cronFile := fmt.Sprintf("/etc/cron.%s/webpanel-backup-%s-%s",
		backupType, backupType, domain)
	_, err := os.Stat(cronFile)
	return err == nil
}

// ListSiteBackups returns a list of backups for a site
func ListSiteBackups(domain, backupType string) ([]string, error) {
	backupDir := filepath.Join(config.GetBackupDir(), backupType, domain)
	entries, err := os.ReadDir(backupDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	var backups []string
	for _, entry := range entries {
		if !entry.IsDir() {
			backups = append(backups, entry.Name())
		}
	}

	return backups, nil
}
