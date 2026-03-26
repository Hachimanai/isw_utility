package ui

import (
	"context"
	"fmt"
	"isw_utility/internal/service"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
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

	// Status/Error Bar
	statusLabel *widget.Label

	// Interpolation targets
	targetCPURPM int
	targetGPURPM int
	currentCPURPM float64
	currentGPURPM float64
}

// NewDashboard creates a new Dashboard instance.
func NewDashboard(window fyne.Window, telemetryService *service.TelemetryService) *Dashboard {
	d := &Dashboard{
		window:           window,
		telemetryService: telemetryService,
	}
	d.initWidgets()
	d.startAnimationLoop()
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
			d.SetStatus(fmt.Sprintf("ERR: %v", err))
		} else {
			d.SetStatus("BOOST MODE UPDATED")
			go func() {
				time.Sleep(3 * time.Second)
				d.SetStatus("SYSTEM READY")
			}()
		}
	})

	// System Info
	d.kernelLabel = widget.NewLabel("Kernel: ...")
	d.kernelLabel.TextStyle = fyne.TextStyle{Monospace: true}
	d.uptimeLabel = widget.NewLabel("Uptime: ...")
	d.uptimeLabel.TextStyle = fyne.TextStyle{Monospace: true}
	d.cpuFreqLabel = widget.NewLabel("Freq: ...")
	d.cpuFreqLabel.TextStyle = fyne.TextStyle{Monospace: true}

	// Status
	d.statusLabel = widget.NewLabel("SYSTEM READY")
	d.statusLabel.TextStyle = fyne.TextStyle{Monospace: true}
}

func (d *Dashboard) startAnimationLoop() {
	go func() {
		ticker := time.NewTicker(30 * time.Millisecond) // ~33 FPS
		for range ticker.C {
			// Linear interpolation for smooth RPM bars
			lerp := 0.1
			d.currentCPURPM += (float64(d.targetCPURPM) - d.currentCPURPM) * lerp
			d.currentGPURPM += (float64(d.targetGPURPM) - d.currentGPURPM) * lerp

			d.cpuFanGauge.SetRPM(int(d.currentCPURPM))
			d.gpuFanGauge.SetRPM(int(d.currentGPURPM))
		}
	}()
}

// BuildLayout constructs the Fyne container hierarchy.
func (d *Dashboard) BuildLayout() fyne.CanvasObject {
	headerBg := canvas.NewRectangle(ColorSurfaceLow)
	headerContent := container.NewHBox(
		layout.NewSpacer(),
		widget.NewLabel("CPU:"), d.cpuTempLabel, d.cpuLoadLabel,
		layout.NewSpacer(),
		widget.NewLabel("GPU:"), d.gpuTempLabel, d.gpuLoadLabel,
		layout.NewSpacer(),
	)
	header := container.NewStack(headerBg, headerContent)

	fansCard := container.NewVBox(
		widget.NewLabel("FANS"),
		d.cpuFanGauge,
		d.gpuFanGauge,
	)
	
	analyticsCard := container.NewVBox(
		widget.NewLabel("ANALYTICS"),
		d.cpuTempHist,
		d.gpuTempHist,
	)

	leftSection := container.NewVBox(fansCard, layout.NewSpacer(), analyticsCard)

	sidePanelBg := canvas.NewRectangle(ColorSurfaceLow)
	sidePanelContent := container.NewVBox(
		widget.NewLabel("CONTROL"),
		d.boostToggle,
		layout.NewSpacer(),
		widget.NewLabel("SYSTEM"),
		d.kernelLabel,
		d.uptimeLabel,
		d.cpuFreqLabel,
	)
	sidePanel := container.NewStack(sidePanelBg, sidePanelContent)

	mainArea := container.NewHSplit(leftSection, sidePanel)
	mainArea.Offset = 0.7 

	footerBg := canvas.NewRectangle(ColorSurfaceHigh)
	footerContent := container.NewHBox(d.statusLabel, layout.NewSpacer())
	footer := container.NewStack(footerBg, footerContent)

	return container.NewBorder(header, footer, nil, nil, mainArea)
}

// SetStatus updates the status label.
func (d *Dashboard) SetStatus(msg string) {
	d.statusLabel.SetText(msg)
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
		d.targetCPURPM = state.Fans[0].RPM
		d.targetGPURPM = state.Fans[1].RPM
	} else if len(state.Fans) == 1 {
		d.targetCPURPM = state.Fans[0].RPM
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
