package repository

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"isw_utility/internal/domain"
)

// HwmonRepository provides fallback access to /sys/class/hwmon for temperatures and fans.
type HwmonRepository struct {
	sysPath string
}

func NewHwmonRepository() *HwmonRepository {
	return &HwmonRepository{
		sysPath: "/sys/class/hwmon",
	}
}

func (r *HwmonRepository) GetTelemetry(ctx context.Context) (domain.Telemetry, error) {
	t := domain.Telemetry{}
	hwmons, _ := filepath.Glob(filepath.Join(r.sysPath, "hwmon*"))
	for _, hwmon := range hwmons {
		name, _ := os.ReadFile(filepath.Join(hwmon, "name"))
		nameStr := strings.TrimSpace(string(name))
		
		// CPU temp usually from coretemp
		if nameStr == "coretemp" {
			t.CPUTemp = r.readTemp(hwmon, "temp1_input")
		}
		// GPU temp might be from amdgpu or nouveau/nvidia
		if nameStr == "amdgpu" || nameStr == "nouveau" {
			t.GPUTemp = r.readTemp(hwmon, "temp1_input")
		}
	}
	return t, nil
}

func (r *HwmonRepository) GetFans(ctx context.Context) ([]domain.FanStatus, error) {
	// Fallback implementation if needed, for now empty.
	return nil, nil
}

func (r *HwmonRepository) GetSystemInfo(ctx context.Context) (domain.SystemInfo, error) {
	return domain.SystemInfo{}, nil
}

func (r *HwmonRepository) readTemp(path, file string) float64 {
	data, err := os.ReadFile(filepath.Join(path, file))
	if err != nil {
		return 0
	}
	val, _ := strconv.ParseFloat(strings.TrimSpace(string(data)), 64)
	return val / 1000.0 // Hwmon values are in milli-degrees Celsius
}
