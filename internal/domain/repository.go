package domain

import "context"

// SensorRepository defines the port to read telemetry and fan data.
type SensorRepository interface {
	GetTelemetry(ctx context.Context) (Telemetry, error)
	GetFans(ctx context.Context) ([]FanStatus, error)
	GetSystemInfo(ctx context.Context) (SystemInfo, error)
}

// ControlRepository defines the port to execute controls like "Boost Mode".
type ControlRepository interface {
	SetBoostMode(ctx context.Context, enabled bool) error
	GetBoostMode(ctx context.Context) (bool, error)
}
