package repository

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"isw_utility/internal/domain"
)

// SystemRepository provides generic system info from Linux /proc and other tools.
type SystemRepository struct{}

func NewSystemRepository() *SystemRepository {
	return &SystemRepository{}
}

func (r *SystemRepository) GetSystemInfo(ctx context.Context) (domain.SystemInfo, error) {
	info := domain.SystemInfo{}

	// Kernel version from uname
	out, _ := exec.Command("uname", "-r").Output()
	info.KernelVersion = strings.TrimSpace(string(out))

	// Uptime from /proc/uptime
	uptimeSec, _ := r.getUptime()
	info.Uptime = time.Duration(uptimeSec) * time.Second

	// CPU Frequency from /proc/cpuinfo
	info.CPUFreq = r.getCPUFreq()

	return info, nil
}

func (r *SystemRepository) GetCPULoad() (float64, error) {
	// Simple CPU load calculation from /proc/stat
	// This is a placeholder for a more robust calculation
	return 0.0, nil
}

func (r *SystemRepository) getUptime() (float64, error) {
	data, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return 0, err
	}
	fields := strings.Fields(string(data))
	if len(fields) > 0 {
		return strconv.ParseFloat(fields[0], 64)
	}
	return 0, fmt.Errorf("invalid /proc/uptime format")
}

func (r *SystemRepository) getCPUFreq() float64 {
	data, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return 0
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.Contains(line, "cpu MHz") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				mhz, _ := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
				return mhz / 1000.0 // Returning GHz
			}
		}
	}
	return 0
}
