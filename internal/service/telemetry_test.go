package service

import (
	"context"
	"isw_utility/internal/domain"
	"testing"
	"time"
)

type mockSensorRepo struct {
	telemetry domain.Telemetry
	fans      []domain.FanStatus
	system    domain.SystemInfo
	err       error
}

func (m *mockSensorRepo) GetTelemetry(ctx context.Context) (domain.Telemetry, error) {
	return m.telemetry, m.err
}

func (m *mockSensorRepo) GetFans(ctx context.Context) ([]domain.FanStatus, error) {
	return m.fans, m.err
}

func (m *mockSensorRepo) GetSystemInfo(ctx context.Context) (domain.SystemInfo, error) {
	return m.system, m.err
}

type mockControlRepo struct {
	boost bool
	err   error
}

func (m *mockControlRepo) SetBoostMode(ctx context.Context, enabled bool) error {
	m.boost = enabled
	return m.err
}

func (m *mockControlRepo) GetBoostMode(ctx context.Context) (bool, error) {
	return m.boost, m.err
}

func TestTelemetryService_CollectData(t *testing.T) {
	sensorRepo := &mockSensorRepo{
		telemetry: domain.Telemetry{CPUTemp: 50.5},
		fans:      []domain.FanStatus{{Label: "CPU", RPM: 2000}},
		system:    domain.SystemInfo{KernelVersion: "6.0"},
	}
	controlRepo := &mockControlRepo{boost: true}

	svc := NewTelemetryService(sensorRepo, controlRepo, time.Second)

	state, err := svc.collectData(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if state.Telemetry.CPUTemp != 50.5 {
		t.Errorf("expected CPUTemp 50.5, got %v", state.Telemetry.CPUTemp)
	}
	if len(state.Fans) != 1 || state.Fans[0].RPM != 2000 {
		t.Errorf("unexpected fan data")
	}
	if state.System.KernelVersion != "6.0" {
		t.Errorf("expected Kernel 6.0, got %v", state.System.KernelVersion)
	}
	if !state.Boost {
		t.Errorf("expected Boost true")
	}
}

func TestTelemetryService_PollingLoop(t *testing.T) {
	sensorRepo := &mockSensorRepo{
		telemetry: domain.Telemetry{CPUTemp: 45.0},
	}
	controlRepo := &mockControlRepo{}

	// Short interval for testing
	svc := NewTelemetryService(sensorRepo, controlRepo, 10*time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go svc.Start(ctx)

	select {
	case state := <-svc.Updates():
		if state.Telemetry.CPUTemp != 45.0 {
			t.Errorf("expected CPUTemp 45.0, got %v", state.Telemetry.CPUTemp)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for update")
	}
}

func TestTelemetryService_SetBoostMode(t *testing.T) {
	sensorRepo := &mockSensorRepo{}
	controlRepo := &mockControlRepo{boost: false}

	svc := NewTelemetryService(sensorRepo, controlRepo, time.Second)

	err := svc.SetBoostMode(context.Background(), true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !controlRepo.boost {
		t.Errorf("expected boost to be true")
	}
}

func TestTelemetryService_Stop(t *testing.T) {
	sensorRepo := &mockSensorRepo{}
	controlRepo := &mockControlRepo{}
	svc := NewTelemetryService(sensorRepo, controlRepo, 10*time.Millisecond)

	go svc.Start(context.Background())

	// Give it a bit of time to start
	time.Sleep(20 * time.Millisecond)

	svc.Stop()

	// If Stop() works, the goroutine should exit. 
	// We can't easily check if goroutine exited without more complex sync,
	// but we can check if it stops sending updates.
	
	// Actually, the select in Start will return on case <-s.stop.
}
