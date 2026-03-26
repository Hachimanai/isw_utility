package domain

// Telemetry represents CPU/GPU temperature and load.
type Telemetry struct {
	CPUTemp float64 `json:"cpu_temp"` // Temperature in °C
	CPULoad float64 `json:"cpu_load"` // Load in %
	GPUTemp float64 `json:"gpu_temp"` // Temperature in °C
	GPULoad float64 `json:"gpu_load"` // Load in %
}
