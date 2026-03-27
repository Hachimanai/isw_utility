package repository

import (
	"context"
	"fmt"
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
	
	var maxTemp float64
	for _, hwmon := range hwmons {
		name, _ := os.ReadFile(filepath.Join(hwmon, "name"))
		nameStr := strings.TrimSpace(string(name))
		
		// 1. Explicit CPU temps
		if nameStr == "coretemp" || nameStr == "k10temp" {
			temp := r.readTemp(hwmon, "temp1_input")
			if temp > 0 { t.CPUTemp = temp }
		}
		
		// 2. Explicit GPU temps
		isGPU := nameStr == "amdgpu" || nameStr == "nouveau" || nameStr == "nvidia" || 
		         nameStr == "radeon" || strings.Contains(nameStr, "gpu")
		if isGPU {
			temp := r.readTemp(hwmon, "temp1_input")
			if temp == 0 { temp = r.readTemp(hwmon, "temp2_input") }
			if temp > 0 { t.GPUTemp = temp }
		}
		
		// 3. Keep track of highest temperature found anywhere as a fallback for CPU
		for i := 1; i < 5; i++ {
			temp := r.readTemp(hwmon, fmt.Sprintf("temp%d_input", i))
			if temp > maxTemp {
				maxTemp = temp
			}
		}
	}
	
	if t.CPUTemp == 0 {
		t.CPUTemp = maxTemp
	}
	
	return t, nil
}

func (r *HwmonRepository) GetFans(ctx context.Context) ([]domain.FanStatus, error) {
	fans := []domain.FanStatus{}
	hwmons, _ := filepath.Glob(filepath.Join(r.sysPath, "hwmon*"))
	
	for _, hwmon := range hwmons {
		name, _ := os.ReadFile(filepath.Join(hwmon, "name"))
		nameStr := strings.TrimSpace(string(name))
		
		// Common fan controllers: thinkpad, asus-nb, dell-smm, etc.
		// We'll search for fan*_input files
		fanFiles, _ := filepath.Glob(filepath.Join(hwmon, "fan*_input"))
		for i, fanFile := range fanFiles {
			data, err := os.ReadFile(fanFile)
			if err != nil {
				continue
			}
			rpm, _ := strconv.Atoi(strings.TrimSpace(string(data)))
			
			label := "FAN_" + strconv.Itoa(i+1)
			if nameStr == "thinkpad" || nameStr == "asus" || nameStr == "msi" || nameStr == "msi_wmi_platform" {
				if i == 0 {
					label = "CPU_FAN"
				} else if i == 1 {
					label = "GPU_FAN"
				}
			}

			fans = append(fans, domain.FanStatus{
				Label:   label,
				RPM:     rpm,
				IdleRPM: 0,
				MaxRPM:  6000,
			})
		}
	}
	return fans, nil
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
