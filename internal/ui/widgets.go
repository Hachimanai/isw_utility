package ui

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

var _ desktop.Hoverable = (*FanGauge)(nil)
var _ desktop.Hoverable = (*TemperatureHistogram)(nil)

// brighterColor returns a color that is slightly more luminous.
func brighterColor(c color.Color) color.Color {
	r, g, b, a := c.RGBA()
	// Increase brightness by 20%
	return color.RGBA{
		R: uint8(min(int(r>>8)+40, 255)),
		G: uint8(min(int(g>>8)+40, 255)),
		B: uint8(min(int(b>>8)+40, 255)),
		A: uint8(a >> 8),
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// FanGauge is a custom widget for displaying fan RPM as a circular gauge.
type FanGauge struct {
	widget.BaseWidget
	Label string
	RPM   int
	Max   int

	isHovered bool
	rpmLabel  *widget.Label
	track     *canvas.Arc
	progress  *canvas.Arc
	glow      *canvas.Circle
}

func NewFanGauge(label string, maxRPM int) *FanGauge {
	g := &FanGauge{
		Label: label,
		Max:   maxRPM,
	}
	g.ExtendBaseWidget(g)

	g.rpmLabel = widget.NewLabel("0000 RPM")
	g.rpmLabel.TextStyle = fyne.TextStyle{Monospace: true}
	g.rpmLabel.Alignment = fyne.TextAlignCenter

	// Background track
	g.track = &canvas.Arc{
		StrokeColor: ColorSurfaceLow,
		StartAngle:  0,
		EndAngle:    360,
		StrokeWidth: 6,
		CutoutRatio: 1.0,
	}

	// Progress arc
	g.progress = &canvas.Arc{
		StrokeColor: ColorPrimary,
		StartAngle:  0, // Top
		EndAngle:    0,
		StrokeWidth: 8,
		CutoutRatio: 1.0,
	}

	g.glow = canvas.NewCircle(color.Transparent)
	g.glow.StrokeWidth = 0
	g.glow.Hide()

	return g
}

func (g *FanGauge) SetRPM(rpm int) {
	g.RPM = rpm
	g.Refresh()
}

func (g *FanGauge) MouseIn(*desktop.MouseEvent) {
	g.isHovered = true
	g.glow.FillColor = color.RGBA{R: ColorPrimary.R, G: ColorPrimary.G, B: ColorPrimary.B, A: 20}
	g.glow.Show()
	g.Refresh()
}

func (g *FanGauge) MouseOut() {
	g.isHovered = false
	g.glow.Hide()
	g.Refresh()
}

func (g *FanGauge) MouseMoved(*desktop.MouseEvent) {}

func (g *FanGauge) CreateRenderer() fyne.WidgetRenderer {
	title := widget.NewLabel(g.Label)
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter
	
	center := container.NewStack(g.glow, g.track, g.progress, g.rpmLabel)
	content := container.NewVBox(title, center)
	
	return &fanGaugeRenderer{
		gauge:   g,
		content: content,
		center:  center,
	}
}

type fanGaugeRenderer struct {
	gauge   *FanGauge
	content fyne.CanvasObject
	center  *fyne.Container
}

func (r *fanGaugeRenderer) Layout(size fyne.Size) {
	fyne.Do(func() {
		r.content.Resize(size)
		
		// Square gauge in the center
		titleHeight := r.content.(*fyne.Container).Objects[0].MinSize().Height
		gaugeSize := size.Width
		if size.Height-titleHeight < gaugeSize {
			gaugeSize = size.Height - titleHeight
		}
		gaugeSize -= 10 // Padding
		
		r.center.Resize(fyne.NewSize(gaugeSize, gaugeSize))
		r.center.Move(fyne.NewPos((size.Width-gaugeSize)/2, titleHeight))
		
		r.gauge.track.Resize(fyne.NewSize(gaugeSize, gaugeSize))
		r.gauge.progress.Resize(fyne.NewSize(gaugeSize, gaugeSize))
		r.gauge.glow.Resize(fyne.NewSize(gaugeSize+10, gaugeSize+10))
		r.gauge.glow.Move(fyne.NewPos(-5, -5))

		percent := float32(r.gauge.RPM) / float32(r.gauge.Max)
		if percent > 1.0 {
			percent = 1.0
		}
		r.gauge.progress.EndAngle = 360 * percent
	})
}

func (r *fanGaugeRenderer) MinSize() fyne.Size {
	return fyne.NewSize(120, 140)
}

func (r *fanGaugeRenderer) Refresh() {
	fyne.Do(func() {
		r.gauge.rpmLabel.SetText(fmt.Sprintf("%d\nRPM", r.gauge.RPM))
		
		if r.gauge.isHovered {
			r.gauge.progress.StrokeColor = brighterColor(ColorPrimary)
		} else {
			r.gauge.progress.StrokeColor = ColorPrimary
		}
		
		r.gauge.track.StrokeColor = ColorSurfaceLow
		
		percent := float32(r.gauge.RPM) / float32(r.gauge.Max)
		if percent > 1.0 {
			percent = 1.0
		}
		r.gauge.progress.EndAngle = 360 * percent
		
		canvas.Refresh(r.content)
	})
}

func (r *fanGaugeRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.content}
}

func (r *fanGaugeRenderer) Destroy() {}

// TemperatureHistogram displays a bar chart of recent temperatures.
type TemperatureHistogram struct {
	widget.BaseWidget
	Label   string
	History []float64
	MaxVal  float64

	isHovered bool
}

func NewTemperatureHistogram(label string) *TemperatureHistogram {
	h := &TemperatureHistogram{
		Label:   label,
		History: make([]float64, 40),
		MaxVal:  100.0,
	}
	h.ExtendBaseWidget(h)
	return h
}

func (h *TemperatureHistogram) AddValue(val float64) {
	h.History = append(h.History[1:], val)
	h.Refresh()
}

func (h *TemperatureHistogram) MouseIn(*desktop.MouseEvent) {
	h.isHovered = true
	h.Refresh()
}

func (h *TemperatureHistogram) MouseOut() {
	h.isHovered = false
	h.Refresh()
}

func (h *TemperatureHistogram) MouseMoved(*desktop.MouseEvent) {}

func (h *TemperatureHistogram) CreateRenderer() fyne.WidgetRenderer {
	title := widget.NewLabel(h.Label)
	title.TextStyle = fyne.TextStyle{Bold: true}
	
	bars := container.NewMax() // Container for bars
	
	content := container.NewVBox(title, bars)
	
	return &tempHistogramRenderer{
		hist:    h,
		content: content,
		bars:    bars,
	}
}

type tempHistogramRenderer struct {
	hist    *TemperatureHistogram
	content fyne.CanvasObject
	bars    *fyne.Container
}

func (r *tempHistogramRenderer) Layout(size fyne.Size) {
	fyne.Do(func() {
		r.content.Resize(size)
		
		// Layout bars within the bars container
		barWidth := size.Width / float32(len(r.hist.History))
		maxH := size.Height - 30 // Subtract title height
		
		r.bars.Objects = nil
		for i, val := range r.hist.History {
			if val <= 0 {
				continue
			}
			h := float32(val/r.hist.MaxVal) * maxH
			rect := canvas.NewRectangle(ColorPrimary)
			if val > 80 {
				rect.FillColor = ColorError
			}
			
			if r.hist.isHovered {
				rect.FillColor = brighterColor(rect.FillColor)
			}

			rect.Resize(fyne.NewSize(barWidth-1, h))
			rect.Move(fyne.NewPos(float32(i)*barWidth, maxH-h))
			r.bars.Add(rect)
		}
	})
}

func (r *tempHistogramRenderer) MinSize() fyne.Size {
	return fyne.NewSize(200, 100)
}

func (r *tempHistogramRenderer) Refresh() {
	fyne.Do(func() {
		r.Layout(r.content.Size())
		canvas.Refresh(r.content)
	})
}

func (r *tempHistogramRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.content}
}

func (r *tempHistogramRenderer) Destroy() {}
