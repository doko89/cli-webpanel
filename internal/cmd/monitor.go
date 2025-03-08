package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/doko/cli-webpanel/internal/config"
	"github.com/doko/cli-webpanel/internal/monitoring"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show system status",
	Long:  `Display current system status including CPU, memory, disk usage, and service status.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		stats, err := monitoring.GetSystemStats()
		if err != nil {
			return fmt.Errorf("failed to get system stats: %v", err)
		}

		// Create tabwriter for aligned output
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

		// System Information
		fmt.Fprintln(w, "System Information:")
		fmt.Fprintf(w, "  CPU Usage:\t%.1f%%\n", stats.CPU)
		fmt.Fprintf(w, "  Memory Usage:\t%s / %s (%.1f%%)\n",
			monitoring.FormatBytes(stats.Memory.Used),
			monitoring.FormatBytes(stats.Memory.Total),
			stats.Memory.UsagePerc)
		fmt.Fprintf(w, "  Disk Usage:\t%s / %s (%.1f%%)\n",
			monitoring.FormatBytes(stats.Disk.Used),
			monitoring.FormatBytes(stats.Disk.Total),
			stats.Disk.UsagePerc)
		fmt.Fprintf(w, "  Uptime:\t%s\n", monitoring.FormatUptime(stats.Uptime))
		fmt.Fprintln(w)

		// Service Status
		fmt.Fprintln(w, "Service Status:")
		for service, status := range stats.Services {
			fmt.Fprintf(w, "  %s:\t%s\n", service, formatServiceStatus(status))
		}

		w.Flush()
		return nil
	},
}

var logsCmd = &cobra.Command{
	Use:   "logs [flags] [service]",
	Short: "Show service logs",
	Long: `Display logs for the specified service. Available services:
- caddy: Web server logs
- mariadb: Database server logs
- webpanel: CLI tool logs
If no service is specified, shows the webpanel logs.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		service := "webpanel"
		if len(args) > 0 {
			service = args[0]
		}

		var logPath string
		switch service {
		case "caddy":
			logPath = "/var/log/caddy/access.log"
		case "mariadb":
			logPath = "/var/log/mysql/error.log"
		case "webpanel":
			logPath = filepath.Join(config.GetLogDir(), "webpanel.log")
		default:
			return fmt.Errorf("unknown service: %s", service)
		}

		// Check if log file exists
		if _, err := os.Stat(logPath); os.IsNotExist(err) {
			return fmt.Errorf("log file not found: %s", logPath)
		}

		// Get number of lines from flag
		numLines, err := cmd.Flags().GetInt("tail")
		if err != nil {
			numLines = 50
		}

		// Check if follow mode is enabled
		follow, err := cmd.Flags().GetBool("follow")
		if err != nil {
			follow = false
		}

		// Prepare tail command
		tailArgs := []string{"-n", fmt.Sprintf("%d", numLines)}
		if follow {
			tailArgs = append(tailArgs, "-f")
		}
		tailArgs = append(tailArgs, logPath)

		execCmd := exec.Command("tail", tailArgs...)
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr

		if err := execCmd.Run(); err != nil {
			return fmt.Errorf("failed to read logs: %v", err)
		}

		return nil
	},
}

var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Monitor system in real-time",
	Long:  `Display real-time system monitoring information with automatic updates.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Clear screen
		fmt.Print("\033[H\033[2J")

		// Run top command with custom format
		execCmd := exec.Command("top", "-b", "-n", "1", "-w", "512")
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr

		if err := execCmd.Run(); err != nil {
			return fmt.Errorf("failed to run monitoring: %v", err)
		}

		return nil
	},
}

func formatServiceStatus(status string) string {
	switch strings.ToLower(status) {
	case "active":
		return "\033[32mactive\033[0m" // Green
	case "inactive":
		return "\033[31minactive\033[0m" // Red
	default:
		return "\033[33m" + status + "\033[0m" // Yellow for unknown status
	}
}

func init() {
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(logsCmd)
	rootCmd.AddCommand(monitorCmd)

	// Add flags for logs command
	logsCmd.Flags().IntP("tail", "n", 50, "Number of lines to show")
	logsCmd.Flags().BoolP("follow", "f", false, "Follow log output")
}
