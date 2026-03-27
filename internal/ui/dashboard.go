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
	"fyne.io/fyne/v2/theme"
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

	// Header Widgets (Telemetry Card)
	d.cpuTempLabel = widget.NewLabel("00°C")
	d.cpuLoadLabel = widget.NewLabel("00%")
	d.gpuTempLabel = widget.NewLabel("00°C")
	d.gpuLoadLabel = widget.NewLabel("00%")

	d.cpuTempLabel.TextStyle = monoBold
	d.cpuLoadLabel.TextStyle = monoBold
	d.gpuTempLabel.TextStyle = monoBold
	d.gpuLoadLabel.TextStyle = monoBold

	// Fan Gauges (using HistoryIcon as a circular placeholder)
	d.cpuFanGauge = NewFanGauge("CPU_FAN", "CPU_RPM", theme.HistoryIcon(), 6000, ColorPrimary)
	d.gpuFanGauge = NewFanGauge("GPU_FAN", "GPU_RPM", theme.HistoryIcon(), 6000, ColorSecondaryGPU)

	// Analytics
	d.cpuTempHist = NewTemperatureHistogram("CPU_TEMP_HISTORY", theme.HistoryIcon())
	d.gpuTempHist = NewTemperatureHistogram("GPU_TEMP_HISTORY", theme.HistoryIcon())

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

func createTelemetryItem(label string, value *widget.Label) fyne.CanvasObject {
	l := widget.NewLabel(label)
	l.TextStyle = fyne.TextStyle{Monospace: true}
	return container.NewVBox(l, value)
}

// BuildLayout constructs the Fyne container hierarchy.
func (d *Dashboard) BuildLayout() fyne.CanvasObject {
	// 1. Header
	title := widget.NewLabel("THERMAL_ARCHITECT")
	title.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}
	title.Alignment = fyne.TextAlignLeading

	telemetryCard := TonalCard(container.NewHBox(
		createTelemetryItem("CPU_TEMP", d.cpuTempLabel),
		layout.NewSpacer(),
		createTelemetryItem("CPU_LOAD", d.cpuLoadLabel),
		layout.NewSpacer(),
		createTelemetryItem("GPU_TEMP", d.gpuTempLabel),
		layout.NewSpacer(),
		createTelemetryItem("GPU_LOAD", d.gpuLoadLabel),
	))

	header := container.NewHBox(title, layout.NewSpacer(), telemetryCard)

	// 2. Main Section (3 columns)
	boostCard := TonalCard(container.NewVBox(
		container.NewHBox(widget.NewLabel("⚡"), createSectionTitle("BOOST_MODE")),
		widget.NewLabel("Override EC fan curves for maximum\nthermal dissipation (isw -b on)."),
		container.NewGridWithColumns(2, d.boostOnBtn, d.boostOffBtn),
	))

	sysInfoCard := TonalCard(container.NewVBox(
		createSectionTitle("SYSTEM_INFO"),
		container.NewHBox(widget.NewLabel("KERNEL"), layout.NewSpacer(), d.kernelLabel),
		container.NewHBox(widget.NewLabel("UPTIME"), layout.NewSpacer(), d.uptimeLabel),
		container.NewHBox(widget.NewLabel("CPU_FREQ"), layout.NewSpacer(), d.cpuFreqLabel),
	))

	mainGrid := container.NewGridWithColumns(3,
		TonalCard(d.cpuFanGauge),
		TonalCard(d.gpuFanGauge),
		container.NewVBox(boostCard, sysInfoCard),
	)

	// 3. Analytics Section (2 columns)
	analyticsGrid := container.NewGridWithColumns(2,
		TonalCard(d.cpuTempHist),
		TonalCard(d.gpuTempHist),
	)

	// Footer
	footer := TonalCardHigh(container.NewHBox(d.statusLabel, layout.NewSpacer()))

	content := container.NewVBox(
		container.NewPadded(header),
		container.NewPadded(mainGrid),
		container.NewPadded(analyticsGrid),
		layout.NewSpacer(),
		footer,
	)

	return container.NewStack(canvas.NewRectangle(ColorBackground), content)
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
