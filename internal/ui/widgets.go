package ui

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
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

// FanGauge is a custom widget for displaying fan RPM.
type FanGauge struct {
	widget.BaseWidget
	Label string
	RPM   int
	Max   int

	isHovered bool
	rpmLabel  *widget.Label
	bar       *canvas.Rectangle
	bgBar     *canvas.Rectangle
	glow      *canvas.Rectangle
}

func NewFanGauge(label string, maxRPM int) *FanGauge {
	g := &FanGauge{
		Label: label,
		Max:   maxRPM,
	}
	g.ExtendBaseWidget(g)

	g.rpmLabel = widget.NewLabel("0000 RPM")
	g.rpmLabel.TextStyle = fyne.TextStyle{Monospace: true}
	g.rpmLabel.Alignment = fyne.TextAlignTrailing

	g.bgBar = canvas.NewRectangle(ColorSurfaceLow)
	g.bgBar.SetMinSize(fyne.NewSize(200, 4))

	g.bar = canvas.NewRectangle(ColorPrimary)
	g.bar.SetMinSize(fyne.NewSize(0, 4))

	g.glow = canvas.NewRectangle(color.Transparent)
	g.glow.Hide()

	return g
}

func (g *FanGauge) SetRPM(rpm int) {
	g.RPM = rpm
	g.rpmLabel.SetText(fmt.Sprintf("%d RPM", rpm))
	g.Refresh()
}

func (g *FanGauge) MouseIn(*desktop.MouseEvent) {
	g.isHovered = true
	g.glow.FillColor = color.RGBA{R: ColorPrimary.R, G: ColorPrimary.G, B: ColorPrimary.B, A: 40}
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
	
	header := container.NewHBox(title, layout.NewSpacer(), g.rpmLabel)
	stack := container.NewStack(g.bgBar, g.glow, g.bar)
	
	content := container.NewVBox(header, stack)
	return &fanGaugeRenderer{
		gauge:   g,
		content: content,
	}
}

type fanGaugeRenderer struct {
	gauge   *FanGauge
	content fyne.CanvasObject
}

func (r *fanGaugeRenderer) Layout(size fyne.Size) {
	r.content.Resize(size)
	r.gauge.bgBar.SetMinSize(fyne.NewSize(size.Width, 4))
	r.gauge.glow.Resize(fyne.NewSize(size.Width, 12))
	r.gauge.glow.Move(fyne.NewPos(0, -4))
	
	percent := float32(r.gauge.RPM) / float32(r.gauge.Max)
	if percent > 1.0 {
		percent = 1.0
	}
	r.gauge.bar.Resize(fyne.NewSize(size.Width*percent, 4))
}

func (r *fanGaugeRenderer) MinSize() fyne.Size {
	return r.content.MinSize()
}

func (r *fanGaugeRenderer) Refresh() {
	if r.gauge.isHovered {
		r.gauge.bar.FillColor = brighterColor(ColorPrimary)
	} else {
		r.gauge.bar.FillColor = ColorPrimary
	}
	
	r.gauge.bgBar.FillColor = ColorSurfaceLow
	
	percent := float32(r.gauge.RPM) / float32(r.gauge.Max)
	if percent > 1.0 {
		percent = 1.0
	}
	r.gauge.bar.Resize(fyne.NewSize(r.gauge.bgBar.Size().Width*percent, 4))
	
	canvas.Refresh(r.content)
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
}

func (r *tempHistogramRenderer) MinSize() fyne.Size {
	return fyne.NewSize(200, 100)
}

func (r *tempHistogramRenderer) Refresh() {
	r.Layout(r.content.Size())
	canvas.Refresh(r.content)
}

func (r *tempHistogramRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.content}
}

func (r *tempHistogramRenderer) Destroy() {}
