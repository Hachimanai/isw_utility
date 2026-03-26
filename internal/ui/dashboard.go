package ui

import (
	"context"
	"fmt"
	"isw_utility/internal/service"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// Dashboard represents the main UI component.
type Dashboard struct {
	window           fyne.Window
	telemetryService *service.TelemetryService

	// Header Widgets
	cpuTempLabel *widget.Label
	cpuLoadLabel *widget.Label
	gpuTempLabel *widget.Label
	gpuLoadLabel *widget.Label

	// Fan Gauges
	cpuFanGauge *FanGauge
	gpuFanGauge *FanGauge

	// Analytics
	cpuTempHist *TemperatureHistogram
	gpuTempHist *TemperatureHistogram

	// Boost Mode
	boostToggle *widget.Check

	// System Info
	kernelLabel  *widget.Label
	uptimeLabel  *widget.Label
	cpuFreqLabel *widget.Label
}

// NewDashboard creates a new Dashboard instance.
func NewDashboard(window fyne.Window, telemetryService *service.TelemetryService) *Dashboard {
	d := &Dashboard{
		window:           window,
		telemetryService: telemetryService,
	}
	d.initWidgets()
	return d
}

func (d *Dashboard) initWidgets() {
	// Header Widgets
	d.cpuTempLabel = widget.NewLabel("00°C")
	d.cpuLoadLabel = widget.NewLabel("00%")
	d.gpuTempLabel = widget.NewLabel("00°C")
	d.gpuLoadLabel = widget.NewLabel("00%")

	// Monospace for technical data
	d.cpuTempLabel.TextStyle = fyne.TextStyle{Monospace: true}
	d.cpuLoadLabel.TextStyle = fyne.TextStyle{Monospace: true}
	d.gpuTempLabel.TextStyle = fyne.TextStyle{Monospace: true}
	d.gpuLoadLabel.TextStyle = fyne.TextStyle{Monospace: true}

	// Fan Gauges
	d.cpuFanGauge = NewFanGauge("CPU FAN", 6000)
	d.gpuFanGauge = NewFanGauge("GPU FAN", 6000)

	// Analytics
	d.cpuTempHist = NewTemperatureHistogram("CPU TEMP HISTORY")
	d.gpuTempHist = NewTemperatureHistogram("GPU TEMP HISTORY")

	// Boost Toggle
	d.boostToggle = widget.NewCheck("BOOST MODE", func(enabled bool) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := d.telemetryService.SetBoostMode(ctx, enabled); err != nil {
			fmt.Printf("Error setting boost mode: %v\n", err)
		}
	})

	// System Info
	d.kernelLabel = widget.NewLabel("Kernel: ...")
	d.kernelLabel.TextStyle = fyne.TextStyle{Monospace: true}
	d.uptimeLabel = widget.NewLabel("Uptime: ...")
	d.uptimeLabel.TextStyle = fyne.TextStyle{Monospace: true}
	d.cpuFreqLabel = widget.NewLabel("Freq: ...")
	d.cpuFreqLabel.TextStyle = fyne.TextStyle{Monospace: true}
}

// BuildLayout constructs the Fyne container hierarchy.
func (d *Dashboard) BuildLayout() fyne.CanvasObject {
	// 1. Header (Real-time Telemetry)
	header := container.NewHBox(
		layout.NewSpacer(),
		widget.NewLabel("CPU:"), d.cpuTempLabel, d.cpuLoadLabel,
		layout.NewSpacer(),
		widget.NewLabel("GPU:"), d.gpuTempLabel, d.gpuLoadLabel,
		layout.NewSpacer(),
	)

	// 2. Main Area (Asymmetric 2/3 - 1/3)
	// Left Area: Fans and Analytics
	fansSection := container.NewVBox(
		widget.NewLabel("FANS"),
		d.cpuFanGauge,
		d.gpuFanGauge,
		layout.NewSpacer(),
		d.cpuTempHist,
		d.gpuTempHist,
	)

	// Right Area: Boost Mode and System Info
	sidePanel := container.NewVBox(
		widget.NewLabel("CONTROL"),
		d.boostToggle,
		widget.NewLabel("SYSTEM"),
		d.kernelLabel,
		d.uptimeLabel,
		d.cpuFreqLabel,
	)

	mainArea := container.New(layout.NewGridLayoutWithColumns(2), fansSection, sidePanel)

	// Final assembly
	return container.NewBorder(header, nil, nil, nil, mainArea)
}

// Update updates the UI widgets with new data.
func (d *Dashboard) Update(state service.DashboardState) {
	d.cpuTempLabel.SetText(fmt.Sprintf("%2.0f°C", state.Telemetry.CPUTemp))
	d.cpuLoadLabel.SetText(fmt.Sprintf("%2.0f%%", state.Telemetry.CPULoad))
	d.gpuTempLabel.SetText(fmt.Sprintf("%2.0f°C", state.Telemetry.GPUTemp))
	d.gpuLoadLabel.SetText(fmt.Sprintf("%2.0f%%", state.Telemetry.GPULoad))

	d.cpuTempHist.AddValue(state.Telemetry.CPUTemp)
	d.gpuTempHist.AddValue(state.Telemetry.GPUTemp)

	if len(state.Fans) >= 2 {
		d.cpuFanGauge.SetRPM(state.Fans[0].RPM)
		d.gpuFanGauge.SetRPM(state.Fans[1].RPM)
	} else if len(state.Fans) == 1 {
		d.cpuFanGauge.SetRPM(state.Fans[0].RPM)
	}

	d.boostToggle.SetChecked(state.Boost)

	d.kernelLabel.SetText(fmt.Sprintf("Kernel: %s", state.System.KernelVersion))
	d.uptimeLabel.SetText(fmt.Sprintf("Uptime: %s", formatDuration(state.System.Uptime)))
	d.cpuFreqLabel.SetText(fmt.Sprintf("Freq:   %.2f GHz", state.System.CPUFreq))
}

func formatDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	return fmt.Sprintf("%02dh %02dm %02ds", h, m, s)
}
