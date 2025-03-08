package monitoring

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type SystemStats struct {
	CPU      float64
	Memory   MemoryStats
	Disk     DiskStats
	Uptime   time.Duration
	Services map[string]string // service name -> status
}

type MemoryStats struct {
	Total     uint64
	Used      uint64
	Free      uint64
	UsagePerc float64
}

type DiskStats struct {
	Total     uint64
	Used      uint64
	Free      uint64
	UsagePerc float64
}

// GetSystemStats returns current system statistics
func GetSystemStats() (*SystemStats, error) {
	stats := &SystemStats{
		Services: make(map[string]string),
	}

	var err error

	// Get CPU usage
	if stats.CPU, err = getCPUUsage(); err != nil {
		return nil, fmt.Errorf("failed to get CPU usage: %v", err)
	}

	// Get memory stats
	if stats.Memory, err = getMemoryStats(); err != nil {
		return nil, fmt.Errorf("failed to get memory stats: %v", err)
	}

	// Get disk stats
	if stats.Disk, err = getDiskStats(); err != nil {
		return nil, fmt.Errorf("failed to get disk stats: %v", err)
	}

	// Get uptime
	if stats.Uptime, err = getUptime(); err != nil {
		return nil, fmt.Errorf("failed to get uptime: %v", err)
	}

	// Get service status
	services := []string{"caddy", "mariadb"}
	for _, service := range services {
		status, _ := getServiceStatus(service)
		stats.Services[service] = status
	}

	return stats, nil
}

// getCPUUsage returns current CPU usage percentage
func getCPUUsage() (float64, error) {
	cmd := exec.Command("top", "-bn1")
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Cpu(s)") {
			fields := strings.Fields(line)
			for _, field := range fields {
				if strings.HasSuffix(field, "id,") {
					idle, err := strconv.ParseFloat(strings.TrimSuffix(field, "id,"), 64)
					if err != nil {
						return 0, err
					}
					return 100.0 - idle, nil
				}
			}
		}
	}

	return 0, fmt.Errorf("CPU usage not found in top output")
}

// getMemoryStats returns current memory statistics
func getMemoryStats() (MemoryStats, error) {
	stats := MemoryStats{}
	cmd := exec.Command("free", "-b")
	output, err := cmd.Output()
	if err != nil {
		return stats, err
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Mem:") {
			fields := strings.Fields(line)
			if len(fields) >= 4 {
				stats.Total, _ = strconv.ParseUint(fields[1], 10, 64)
				stats.Used, _ = strconv.ParseUint(fields[2], 10, 64)
				stats.Free, _ = strconv.ParseUint(fields[3], 10, 64)
				stats.UsagePerc = float64(stats.Used) / float64(stats.Total) * 100
				return stats, nil
			}
		}
	}

	return stats, fmt.Errorf("memory stats not found in free output")
}

// getDiskStats returns current disk usage statistics
func getDiskStats() (DiskStats, error) {
	stats := DiskStats{}
	cmd := exec.Command("df", "-B1", "/")
	output, err := cmd.Output()
	if err != nil {
		return stats, err
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) >= 2 {
		fields := strings.Fields(lines[1])
		if len(fields) >= 4 {
			stats.Total, _ = strconv.ParseUint(fields[1], 10, 64)
			stats.Used, _ = strconv.ParseUint(fields[2], 10, 64)
			stats.Free, _ = strconv.ParseUint(fields[3], 10, 64)
			stats.UsagePerc = float64(stats.Used) / float64(stats.Total) * 100
			return stats, nil
		}
	}

	return stats, fmt.Errorf("disk stats not found in df output")
}

// getUptime returns system uptime
func getUptime() (time.Duration, error) {
	cmd := exec.Command("cat", "/proc/uptime")
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	fields := strings.Fields(string(output))
	if len(fields) > 0 {
		uptime, err := strconv.ParseFloat(fields[0], 64)
		if err != nil {
			return 0, err
		}
		return time.Duration(uptime) * time.Second, nil
	}

	return 0, fmt.Errorf("uptime not found")
}

// getServiceStatus returns the status of a system service
func getServiceStatus(service string) (string, error) {
	cmd := exec.Command("systemctl", "is-active", service)
	output, err := cmd.Output()
	if err != nil {
		return "inactive", nil
	}
	return strings.TrimSpace(string(output)), nil
}

// FormatBytes formats bytes into human readable format
func FormatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := unit, 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// FormatUptime formats uptime duration into human readable format
func FormatUptime(d time.Duration) string {
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}
