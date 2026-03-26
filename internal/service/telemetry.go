package service

import (
	"context"
	"fmt"
	"isw_utility/internal/domain"
	"log/slog"
	"time"
)

// DashboardState combines all data for the UI.
type DashboardState struct {
	Telemetry domain.Telemetry   `json:"telemetry"`
	Fans      []domain.FanStatus `json:"fans"`
	System    domain.SystemInfo  `json:"system"`
	Boost     bool               `json:"boost"`
	Timestamp time.Time          `json:"timestamp"`
}

// TelemetryService orchestrates data collection and controls.
type TelemetryService struct {
	sensorRepo  domain.SensorRepository
	controlRepo domain.ControlRepository
	interval    time.Duration
	updates     chan DashboardState
	stop        chan struct{}
}

// NewTelemetryService creates a new TelemetryService.
func NewTelemetryService(sensorRepo domain.SensorRepository, controlRepo domain.ControlRepository, interval time.Duration) *TelemetryService {
	return &TelemetryService{
		sensorRepo:  sensorRepo,
		controlRepo: controlRepo,
		interval:    interval,
		updates:     make(chan DashboardState, 1),
		stop:        make(chan struct{}),
	}
}

// Start runs the polling loop.
// It is intended to be run in its own goroutine.
func (s *TelemetryService) Start(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	slog.Info("Starting telemetry polling loop", "interval", s.interval)

	for {
		select {
		case <-ctx.Done():
			slog.Info("Stopping telemetry polling loop (context cancelled)")
			return
		case <-s.stop:
			slog.Info("Stopping telemetry polling loop (stop signal)")
			return
		case <-ticker.C:
			state, err := s.collectData(ctx)
			if err != nil {
				slog.Error("Failed to collect telemetry data", "error", err)
				continue
			}

			// Send update (non-blocking if UI is slow to consume)
			select {
			case s.updates <- state:
			default:
				// If channel is full, we skip this update to avoid blocking the polling
				slog.Debug("Skipping update, channel full")
			}
		}
	}
}

// Stop sends a signal to stop the polling loop.
func (s *TelemetryService) Stop() {
	close(s.stop)
}

// Updates returns the channel for telemetry updates.
func (s *TelemetryService) Updates() <-chan DashboardState {
	return s.updates
}

// SetBoostMode enables or disables boost mode.
func (s *TelemetryService) SetBoostMode(ctx context.Context, enabled bool) error {
	if err := s.controlRepo.SetBoostMode(ctx, enabled); err != nil {
		return fmt.Errorf("failed to set boost mode: %w", err)
	}
	return nil
}

func (s *TelemetryService) collectData(ctx context.Context) (DashboardState, error) {
	telemetry, err := s.sensorRepo.GetTelemetry(ctx)
	if err != nil {
		return DashboardState{}, fmt.Errorf("failed to get telemetry: %w", err)
	}

	fans, err := s.sensorRepo.GetFans(ctx)
	if err != nil {
		return DashboardState{}, fmt.Errorf("failed to get fans: %w", err)
	}

	system, err := s.sensorRepo.GetSystemInfo(ctx)
	if err != nil {
		return DashboardState{}, fmt.Errorf("failed to get system info: %w", err)
	}

	boost, err := s.controlRepo.GetBoostMode(ctx)
	if err != nil {
		// Log but don't fail, maybe boost info is not critical or not available
		slog.Warn("Failed to get boost mode status", "error", err)
	}

	return DashboardState{
		Telemetry: telemetry,
		Fans:      fans,
		System:    system,
		Boost:     boost,
		Timestamp: time.Now(),
	}, nil
}
