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

	// Boost Mode Buttons
	boostOnBtn  *widget.Button
	boostOffBtn *widget.Button

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
	monoBold := fyne.TextStyle{Monospace: true, Bold: true}
	mono := fyne.TextStyle{Monospace: true}

	// Header Widgets
	d.cpuTempLabel = widget.NewLabel("00°C")
	d.cpuLoadLabel = widget.NewLabel("00%")
	d.gpuTempLabel = widget.NewLabel("00°C")
	d.gpuLoadLabel = widget.NewLabel("00%")

	d.cpuTempLabel.TextStyle = mono
	d.cpuLoadLabel.TextStyle = mono
	d.gpuTempLabel.TextStyle = mono
	d.gpuLoadLabel.TextStyle = mono

	// Fan Gauges
	d.cpuFanGauge = NewFanGauge("CPU_FAN", 6000)
	d.gpuFanGauge = NewFanGauge("GPU_FAN", 6000)

	// Analytics
	d.cpuTempHist = NewTemperatureHistogram("CPU_TEMP_HISTORY")
	d.gpuTempHist = NewTemperatureHistogram("GPU_TEMP_HISTORY")

	// Boost Mode Buttons
	d.boostOnBtn = widget.NewButton("BOOST_ON", func() {
		d.setBoost(true)
	})
	d.boostOnBtn.Importance = widget.HighImportance

	d.boostOffBtn = widget.NewButton("BOOST_OFF", func() {
		d.setBoost(false)
	})

	// System Info
	d.kernelLabel = widget.NewLabel("KERNEL: ...")
	d.kernelLabel.TextStyle = mono
	d.uptimeLabel = widget.NewLabel("UPTIME: ...")
	d.uptimeLabel.TextStyle = mono
	d.cpuFreqLabel = widget.NewLabel("FREQ:   ...")
	d.cpuFreqLabel.TextStyle = mono

	// Status
	d.statusLabel = widget.NewLabel("SYSTEM_READY")
	d.statusLabel.TextStyle = monoBold
}

func (d *Dashboard) setBoost(enabled bool) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		stateStr := "ENABLING"
		if !enabled {
			stateStr = "DISABLING"
		}
		
		fyne.Do(func() {
			d.SetStatus(fmt.Sprintf("%s BOOST_MODE...", stateStr))
		})

		if err := d.telemetryService.SetBoostMode(ctx, enabled); err != nil {
			fyne.Do(func() {
				d.SetStatus(fmt.Sprintf("ERR: %v", err))
			})
		} else {
			finalState := "ACTIVE"
			if !enabled {
				finalState = "STANDBY"
			}
			fyne.Do(func() {
				d.SetStatus(fmt.Sprintf("BOOST_MODE: %s", finalState))
			})
		}
	}()
}

func (d *Dashboard) startAnimationLoop() {
	go func() {
		ticker := time.NewTicker(30 * time.Millisecond) // ~33 FPS
		for range ticker.C {
			lerp := 0.1
			d.currentCPURPM += (float64(d.targetCPURPM) - d.currentCPURPM) * lerp
			d.currentGPURPM += (float64(d.targetGPURPM) - d.currentGPURPM) * lerp

			d.cpuFanGauge.SetRPM(int(d.currentCPURPM))
			d.gpuFanGauge.SetRPM(int(d.currentGPURPM))
		}
	}()
}

func createSectionTitle(text string) *widget.Label {
	l := widget.NewLabel(text)
	l.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}
	return l
}

// BuildLayout constructs the Fyne container hierarchy.
func (d *Dashboard) BuildLayout() fyne.CanvasObject {
	headerBg := canvas.NewRectangle(ColorSurfaceLow)
	
	headerContent := container.NewHBox(
		layout.NewSpacer(),
		createSectionTitle("CPU:"), d.cpuTempLabel, d.cpuLoadLabel,
		layout.NewSpacer(),
		createSectionTitle("GPU:"), d.gpuTempLabel, d.gpuLoadLabel,
		layout.NewSpacer(),
	)
	header := container.NewStack(headerBg, container.NewPadded(headerContent))

	fansContent := container.NewGridWithColumns(2, d.cpuFanGauge, d.gpuFanGauge)
	fansCard := container.NewVBox(
		createSectionTitle("FAN_GAUGES"),
		fansContent,
	)
	
	analyticsCard := container.NewVBox(
		createSectionTitle("ANALYTICS"),
		d.cpuTempHist,
		d.gpuTempHist,
	)

	leftSection := container.NewVBox(fansCard, layout.NewSpacer(), analyticsCard)

	sidePanelBg := canvas.NewRectangle(ColorSurfaceLow)
	sidePanelContent := container.NewVBox(
		createSectionTitle("CONTROL"),
		container.NewGridWithColumns(2, d.boostOnBtn, d.boostOffBtn),
		layout.NewSpacer(),
		createSectionTitle("SYSTEM_INFO"),
		d.kernelLabel,
		d.uptimeLabel,
		d.cpuFreqLabel,
	)
	sidePanel := container.NewStack(sidePanelBg, container.NewPadded(sidePanelContent))

	mainArea := container.NewHSplit(container.NewPadded(leftSection), sidePanel)
	mainArea.Offset = 0.7 

	footerBg := canvas.NewRectangle(ColorSurfaceHigh)
	footerContent := container.NewHBox(d.statusLabel, layout.NewSpacer())
	footer := container.NewStack(footerBg, container.NewPadded(footerContent))

	return container.NewBorder(header, footer, nil, nil, mainArea)
}

// SetStatus updates the status label.
func (d *Dashboard) SetStatus(msg string) {
	d.statusLabel.SetText(msg)
}

// Update updates the UI widgets with new data.
func (d *Dashboard) Update(state service.DashboardState) {
	// Use a goroutine to avoid blocking the telemetry loop if UI is busy.
	go fyne.Do(func() {
		d.cpuTempLabel.SetText(fmt.Sprintf("%2.0f°C", state.Telemetry.CPUTemp))
		d.cpuLoadLabel.SetText(fmt.Sprintf("%2.0f%%", state.Telemetry.CPULoad))
		d.gpuTempLabel.SetText(fmt.Sprintf("%2.0f°C", state.Telemetry.GPUTemp))
		d.gpuLoadLabel.SetText(fmt.Sprintf("%2.0f%%", state.Telemetry.GPULoad))

		d.cpuTempHist.AddValue(state.Telemetry.CPUTemp)
		d.gpuTempHist.AddValue(state.Telemetry.GPUTemp)

		for _, fan := range state.Fans {
			if fan.Label == "CPU_FAN" {
				d.targetCPURPM = fan.RPM
			} else if fan.Label == "GPU_FAN" {
				d.targetGPURPM = fan.RPM
			}
		}
		// Fallback if labels are generic
		if len(state.Fans) > 0 && d.targetCPURPM == 0 && d.targetGPURPM == 0 {
			d.targetCPURPM = state.Fans[0].RPM
			if len(state.Fans) > 1 {
				d.targetGPURPM = state.Fans[1].RPM
			}
		}

		d.kernelLabel.SetText(fmt.Sprintf("KERNEL: %s", state.System.KernelVersion))
		d.uptimeLabel.SetText(fmt.Sprintf("UPTIME: %s", formatDuration(state.System.Uptime)))
		d.cpuFreqLabel.SetText(fmt.Sprintf("FREQ:   %.2f GHz", state.System.CPUFreq))
	})
}

func formatDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	return fmt.Sprintf("%02dh %02dm %02ds", h, m, s)
}
