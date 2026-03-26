package repository

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"isw_utility/internal/domain"
)

// ISWRepository implements domain.SensorRepository and domain.ControlRepository using the 'isw' CLI.
type ISWRepository struct{}

func NewISWRepository() *ISWRepository {
	return &ISWRepository{}
}

func (r *ISWRepository) GetTelemetry(ctx context.Context) (domain.Telemetry, error) {
	// isw -r typically provides temperatures.
	output, err := r.runCommand(ctx, "isw", "-r")
	if err != nil {
		return domain.Telemetry{}, err
	}

	return r.parseTelemetry(output), nil
}

func (r *ISWRepository) GetFans(ctx context.Context) ([]domain.FanStatus, error) {
	output, err := r.runCommand(ctx, "isw", "-r")
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
	_, err := r.runCommand(ctx, "isw", "-b", state)
	return err
}

func (r *ISWRepository) GetBoostMode(ctx context.Context) (bool, error) {
	// isw -r often shows boost mode status or we can infer from max RPM.
	// For now, let's assume it's in the status output.
	output, err := r.runCommand(ctx, "isw", "-r")
	if err != nil {
		return false, err
	}
	return strings.Contains(strings.ToLower(output), "cooler boost: on"), nil
}

func (r *ISWRepository) runCommand(ctx context.Context, name string, arg ...string) (string, error) {
	cmd := exec.CommandContext(ctx, name, arg...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to run %s %v: %w (stderr: %s)", name, arg, err, stderr.String())
	}
	return stdout.String(), nil
}

func (r *ISWRepository) parseTelemetry(output string) domain.Telemetry {
	t := domain.Telemetry{}
	// Regex for "CPU Temperature: 45°C"
	reTemp := regexp.MustCompile(`(CPU|GPU) Temperature:\s*(\d+)`)
	matches := reTemp.FindAllStringSubmatch(output, -1)
	for _, m := range matches {
		val, _ := strconv.Atoi(m[2])
		if m[1] == "CPU" {
			t.CPUTemp = float64(val)
		} else {
			t.GPUTemp = float64(val)
		}
	}
	// Note: isw -r might not provide Load. Load is better fetched from /proc or another tool.
	return t
}

func (r *ISWRepository) parseFans(output string) []domain.FanStatus {
	fans := []domain.FanStatus{}
	// Regex for "CPU Fan Speed: 2150 RPM"
	reFan := regexp.MustCompile(`(CPU|GPU) Fan Speed:\s*(\d+)`)
	matches := reFan.FindAllStringSubmatch(output, -1)
	for _, m := range matches {
		val, _ := strconv.Atoi(m[2])
		fans = append(fans, domain.FanStatus{
			Label:   m[1] + "_FAN",
			RPM:     val,
			// IDLE/MAX ranges could be fetched from config or calibrated.
			// Defaulting to reasonable values for a laptop.
			IdleRPM: 0,
			MaxRPM:  6000,
		})
	}
	return fans
}
