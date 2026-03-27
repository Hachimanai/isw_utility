package repository

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"isw_utility/internal/domain"
)

// SystemRepository provides generic system info from Linux /proc and other tools.
type SystemRepository struct {
	prevIdle  uint64
	prevTotal uint64
}

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
	data, err := os.ReadFile("/proc/stat")
	if err != nil {
		return 0, err
	}

	lines := strings.Split(string(data), "\n")
	if len(lines) == 0 {
		return 0, fmt.Errorf("empty /proc/stat")
	}

	// First line is "cpu  ..."
	fields := strings.Fields(lines[0])
	if len(fields) < 5 {
		return 0, fmt.Errorf("invalid /proc/stat format")
	}

	var total uint64
	for i := 1; i < len(fields); i++ {
		val, _ := strconv.ParseUint(fields[i], 10, 64)
		total += val
	}
	idle, _ := strconv.ParseUint(fields[4], 10, 64)

	diffIdle := idle - r.prevIdle
	diffTotal := total - r.prevTotal

	r.prevIdle = idle
	r.prevTotal = total

	if diffTotal == 0 {
		return 0, nil
	}

	return 100 * (1 - float64(diffIdle)/float64(diffTotal)), nil
}

func (r *SystemRepository) GetGPULoad() float64 {
	// 1. Try AMD GPU path
	data, err := os.ReadFile("/sys/class/drm/card0/device/gpu_busy_percent")
	if err == nil {
		val, _ := strconv.ParseFloat(strings.TrimSpace(string(data)), 64)
		return val
	}
	
	// 2. Try secondary AMD card (common in laptops with iGPU)
	data, err = os.ReadFile("/sys/class/drm/card1/device/gpu_busy_percent")
	if err == nil {
		val, _ := strconv.ParseFloat(strings.TrimSpace(string(data)), 64)
		return val
	}

	// 3. Try NVIDIA via nvidia-smi (if installed)
	out, err := exec.Command("nvidia-smi", "--query-gpu=utilization.gpu", "--format=csv,noheader,nounits").Output()
	if err == nil {
		val, _ := strconv.ParseFloat(strings.TrimSpace(string(out)), 64)
		return val
	}

	return 0
}

func (r *SystemRepository) GetGPUTemp() float64 {
	// 1. Try NVIDIA via nvidia-smi (reliable for NVIDIA)
	out, err := exec.Command("nvidia-smi", "--query-gpu=temperature.gpu", "--format=csv,noheader,nounits").Output()
	if err == nil {
		val, _ := strconv.ParseFloat(strings.TrimSpace(string(out)), 64)
		return val
	}

	// 2. Try common AMD temperature path
	hwmons, _ := filepath.Glob("/sys/class/hwmon/hwmon*/name")
	for _, hwmonNameFile := range hwmons {
		name, _ := os.ReadFile(hwmonNameFile)
		nameStr := strings.TrimSpace(string(name))
		if nameStr == "amdgpu" {
			hwmonDir := filepath.Dir(hwmonNameFile)
			tempData, err := os.ReadFile(filepath.Join(hwmonDir, "temp1_input"))
			if err == nil {
				val, _ := strconv.ParseFloat(strings.TrimSpace(string(tempData)), 64)
				return val / 1000.0
			}
		}
	}

	return 0
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
