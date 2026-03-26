package repository

import (
	"reflect"
	"testing"

	"isw_utility/internal/domain"
)

func TestISWRepository_parseTelemetry(t *testing.T) {
	repo := &ISWRepository{}
	output := `
CPU Temperature: 45
GPU Temperature: 52
Cooler Boost: off
`
	expected := domain.Telemetry{
		CPUTemp: 45,
		GPUTemp: 52,
	}
	got := repo.parseTelemetry(output)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("parseTelemetry() = %v, want %v", got, expected)
	}
}

func TestISWRepository_parseFans(t *testing.T) {
	repo := &ISWRepository{}
	output := `
CPU Fan Speed: 2150 RPM
GPU Fan Speed: 2300 RPM
`
	expected := []domain.FanStatus{
		{
			Label:   "CPU_FAN",
			RPM:     2150,
			IdleRPM: 0,
			MaxRPM:  6000,
		},
		{
			Label:   "GPU_FAN",
			RPM:     2300,
			IdleRPM: 0,
			MaxRPM:  6000,
		},
	}
	got := repo.parseFans(output)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("parseFans() = %v, want %v", got, expected)
	}
}
