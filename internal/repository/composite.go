package repository

import (
	"context"
	"isw_utility/internal/domain"
)

// CompositeRepository combines multiple repositories to provide data from the best available source.
type CompositeRepository struct {
	isw    *ISWRepository
	hwmon  *HwmonRepository
	system *SystemRepository
}

func NewCompositeRepository(isw *ISWRepository, hwmon *HwmonRepository, system *SystemRepository) *CompositeRepository {
	return &CompositeRepository{
		isw:    isw,
		hwmon:  hwmon,
		system: system,
	}
}

func (r *CompositeRepository) GetTelemetry(ctx context.Context) (domain.Telemetry, error) {
	// Try ISW first
	t, _ := r.isw.GetTelemetry(ctx)
	
	// Fallback/Merge for temperatures
	if t.CPUTemp == 0 || t.GPUTemp == 0 {
		h, _ := r.hwmon.GetTelemetry(ctx)
		if t.CPUTemp == 0 { t.CPUTemp = h.CPUTemp }
		if t.GPUTemp == 0 { t.GPUTemp = h.GPUTemp }
	}
	
	// Ensure GPU temp from nvidia-smi if still missing
	if t.GPUTemp == 0 {
		t.GPUTemp = r.system.GetGPUTemp()
	}

	// Always get load from system repo as it's more reliable
	t.CPULoad, _ = r.system.GetCPULoad()
	t.GPULoad = r.system.GetGPULoad()
	
	return t, nil
}

func (r *CompositeRepository) GetFans(ctx context.Context) ([]domain.FanStatus, error) {
	fans, _ := r.isw.GetFans(ctx)
	
	// If ISW returned no fans or all zero RPMs, try Hwmon
	allZero := true
	for _, f := range fans {
		if f.RPM > 0 {
			allZero = false
			break
		}
	}
	
	if len(fans) == 0 || allZero {
		hFans, _ := r.hwmon.GetFans(ctx)
		if len(hFans) > 0 {
			return hFans, nil
		}
	}
	
	return fans, nil
}

func (r *CompositeRepository) GetSystemInfo(ctx context.Context) (domain.SystemInfo, error) {
	return r.system.GetSystemInfo(ctx)
}

func (r *CompositeRepository) SetBoostMode(ctx context.Context, enabled bool) error {
	return r.isw.SetBoostMode(ctx, enabled)
}

func (r *CompositeRepository) GetBoostMode(ctx context.Context) (bool, error) {
	return r.isw.GetBoostMode(ctx)
}
