package repository

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"isw_utility/internal/domain"
)

// ISWRepository implements domain.SensorRepository and domain.ControlRepository using the 'isw' CLI.
type ISWRepository struct{}

func NewISWRepository() *ISWRepository {
	return &ISWRepository{}
}

func (r *ISWRepository) GetTelemetry(ctx context.Context) (domain.Telemetry, error) {
	// Use a short timeout for ISW calls to avoid blocking the main loop
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	
	output, err := r.runCommand(ctx, "isw", "-r", "MSI_ADDRESS_DEFAULT")
	if err != nil {
		return domain.Telemetry{}, err
	}

	return r.parseTelemetry(output), nil
}

func (r *ISWRepository) GetFans(ctx context.Context) ([]domain.FanStatus, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	
	output, err := r.runCommand(ctx, "isw", "-r", "MSI_ADDRESS_DEFAULT")
	if err != nil {
		return nil, err
	}

	return r.parseFans(output), nil
}

func (r *ISWRepository) GetSystemInfo(ctx context.Context) (domain.SystemInfo, error) {
	// isw doesn't provide kernel info, but it might have CPU frequency in some versions.
	// We'll leave this to the SystemRepository for most parts.
	return domain.SystemInfo{}, nil
}

func (r *ISWRepository) SetBoostMode(ctx context.Context, enabled bool) error {
	state := "off"
	if enabled {
		state = "on"
	}
	// Use pkexec with absolute path for better reliability
	_, err := r.runCommand(ctx, "pkexec", "/usr/bin/isw", "-b", state)
	return err
}

func (r *ISWRepository) GetBoostMode(ctx context.Context) (bool, error) {
	// Use a short timeout to prevent hanging the whole app
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	output, err := r.runCommand(ctx, "isw", "-r", "MSI_ADDRESS_DEFAULT")
	if err != nil {
		return false, err
	}
	// Fallback check: if fans are at very high RPM, boost is likely ON
	fans := r.parseFans(output)
	for _, f := range fans {
		if f.RPM > 5500 {
			return true, nil
		}
	}
	return strings.Contains(strings.ToLower(output), "cooler boost: on") || 
	       strings.Contains(strings.ToLower(output), "cooler boost: 1"), nil
}

func (r *ISWRepository) runCommand(ctx context.Context, name string, arg ...string) (string, error) {
	cmd := exec.CommandContext(ctx, name, arg...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		stderrStr := strings.TrimSpace(stderr.String())
		if stderrStr == "" {
			return "", fmt.Errorf("failed to run %s %v: %w", name, arg, err)
		}
		return "", fmt.Errorf("%s", stderrStr)
	}
	return stdout.String(), nil
}

func (r *ISWRepository) parseTelemetry(output string) domain.Telemetry {
	t := domain.Telemetry{}
	// Look for any number followed by °C
	reTemp := regexp.MustCompile(`(\d+)°C`)
	matches := reTemp.FindAllStringSubmatch(output, -1)
	if len(matches) >= 1 {
		val, _ := strconv.Atoi(matches[0][1])
		t.CPUTemp = float64(val)
	}
	if len(matches) >= 2 {
		val, _ := strconv.Atoi(matches[1][1])
		t.GPUTemp = float64(val)
	}
	
	// Also support the original label format as fallback
	reLabel := regexp.MustCompile(`(?i)(CPU|GPU)\s*(?:Temperature|Temp):\s*(\d+)`)
	matchesLabel := reLabel.FindAllStringSubmatch(output, -1)
	for _, m := range matchesLabel {
		val, _ := strconv.Atoi(m[2])
		label := strings.ToUpper(m[1])
		if label == "CPU" && t.CPUTemp == 0 {
			t.CPUTemp = float64(val)
		} else if label == "GPU" && t.GPUTemp == 0 {
			t.GPUTemp = float64(val)
		}
	}
	return t
}

func (r *ISWRepository) parseFans(output string) []domain.FanStatus {
	fans := []domain.FanStatus{}
	// Look for any number followed by RPM
	reRPM := regexp.MustCompile(`(\d+)\s*RPM`)
	matches := reRPM.FindAllStringSubmatch(output, -1)
	
	if len(matches) >= 1 {
		val, _ := strconv.Atoi(matches[0][1])
		fans = append(fans, domain.FanStatus{Label: "CPU_FAN", RPM: val, MaxRPM: 6000})
	}
	if len(matches) >= 2 {
		val, _ := strconv.Atoi(matches[1][1])
		fans = append(fans, domain.FanStatus{Label: "GPU_FAN", RPM: val, MaxRPM: 6000})
	}

	// Also support the original label format as fallback
	reLabel := regexp.MustCompile(`(?i)(CPU|GPU)\s*(?:Fan Speed|Fan|FAN):\s*(\d+)`)
	matchesLabel := reLabel.FindAllStringSubmatch(output, -1)
	for _, m := range matchesLabel {
		val, _ := strconv.Atoi(m[2])
		label := strings.ToUpper(m[1]) + "_FAN"
		// Check if we already have it
		found := false
		for _, f := range fans {
			if f.Label == label {
				found = true
				break
			}
		}
		if !found {
			fans = append(fans, domain.FanStatus{Label: label, RPM: val, MaxRPM: 6000})
		}
	}
	return fans
}
