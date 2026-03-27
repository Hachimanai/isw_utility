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

// FanGauge is a custom widget for displaying fan RPM as a circular gauge.
type FanGauge struct {
	widget.BaseWidget
	Label    string
	SubLabel string
	Icon     fyne.Resource
	RPM      int
	Max      int

	isHovered bool
	rpmLabel  *widget.Label
	subLabel  *widget.Label
	icon      *canvas.Image
	track     *canvas.Arc
	progress  *canvas.Arc
	glow      *canvas.Circle

	// Bottom bar
	bottomTrack    *canvas.Rectangle
	bottomProgress *canvas.Rectangle
}

func NewFanGauge(label, subLabel string, icon fyne.Resource, maxRPM int) *FanGauge {
	g := &FanGauge{
		Label:    label,
		SubLabel: subLabel,
		Icon:     icon,
		Max:      maxRPM,
	}
	g.ExtendBaseWidget(g)

	g.rpmLabel = widget.NewLabel("0000")
	g.rpmLabel.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}
	g.rpmLabel.Alignment = fyne.TextAlignCenter

	g.subLabel = widget.NewLabel(subLabel)
	g.subLabel.TextStyle = fyne.TextStyle{Monospace: true}
	g.subLabel.Alignment = fyne.TextAlignCenter

	g.icon = canvas.NewImageFromResource(icon)
	g.icon.FillMode = canvas.ImageFillContain

	// Background track
	g.track = &canvas.Arc{
		StrokeColor: ColorSurfaceLow,
		StartAngle:  0,
		EndAngle:    360,
		StrokeWidth: 10,
		CutoutRatio: 1.0,
	}

	// Progress arc
	g.progress = &canvas.Arc{
		StrokeColor: ColorPrimary,
		StartAngle:  0, // Top
		EndAngle:    0,
		StrokeWidth: 12,
		CutoutRatio: 1.0,
	}

	g.glow = canvas.NewCircle(color.Transparent)
	g.glow.StrokeWidth = 0
	g.glow.Hide()

	g.bottomTrack = canvas.NewRectangle(ColorSurfaceLow)
	g.bottomProgress = canvas.NewRectangle(ColorPrimary)

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
	title.Hide() // We hide the top title as it's not in the card in screen.png, but we keep the field

	centerContent := container.NewVBox(
		container.NewCenter(container.NewGridWithRows(1, container.NewPadded(g.icon))),
		g.rpmLabel,
		g.subLabel,
	)

	gaugeStack := container.NewStack(g.glow, g.track, g.progress, container.NewCenter(centerContent))

	idleLabel := widget.NewLabel("IDLE")
	idleLabel.TextStyle = fyne.TextStyle{Monospace: true}
	maxLabel := widget.NewLabel("MAX")
	maxLabel.TextStyle = fyne.TextStyle{Monospace: true}

	bottomBar := container.NewStack(g.bottomTrack, g.bottomProgress)

	bottomSection := container.NewBorder(nil, nil, idleLabel, maxLabel, container.NewPadded(bottomBar))

	content := container.NewVBox(
		layout.NewSpacer(),
		gaugeStack,
		layout.NewSpacer(),
		bottomSection,
	)

	return &fanGaugeRenderer{
		gauge:   g,
		content: content,
		center:  gaugeStack,
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

		// Circular gauge area
		gaugeSize := size.Width * 0.7
		if size.Height*0.6 < gaugeSize {
			gaugeSize = size.Height * 0.6
		}

		r.center.Resize(fyne.NewSize(gaugeSize, gaugeSize))
		r.center.Move(fyne.NewPos((size.Width-gaugeSize)/2, (size.Height*0.7-gaugeSize)/2))

		r.gauge.track.Resize(fyne.NewSize(gaugeSize, gaugeSize))
		r.gauge.progress.Resize(fyne.NewSize(gaugeSize, gaugeSize))
		r.gauge.glow.Resize(fyne.NewSize(gaugeSize+20, gaugeSize+20))
		r.gauge.glow.Move(fyne.NewPos(-10, -10))

		r.gauge.icon.SetMinSize(fyne.NewSize(gaugeSize*0.2, gaugeSize*0.2))

		percent := float32(r.gauge.RPM) / float32(r.gauge.Max)
		if percent > 1.0 {
			percent = 1.0
		}
		r.gauge.progress.EndAngle = 360 * percent

		// Bottom bar
		barWidth := size.Width - 40
		r.gauge.bottomTrack.Resize(fyne.NewSize(barWidth, 4))
		r.gauge.bottomProgress.Resize(fyne.NewSize(barWidth*percent, 4))
	})
}

func (r *fanGaugeRenderer) MinSize() fyne.Size {
	return fyne.NewSize(240, 320)
}

func (r *fanGaugeRenderer) Refresh() {
	fyne.Do(func() {
		r.gauge.rpmLabel.SetText(fmt.Sprintf("%d", r.gauge.RPM))

		if r.gauge.isHovered {
			r.gauge.progress.StrokeColor = brighterColor(ColorPrimary)
			r.gauge.bottomProgress.FillColor = brighterColor(ColorPrimary)
		} else {
			r.gauge.progress.StrokeColor = ColorPrimary
			r.gauge.bottomProgress.FillColor = ColorPrimary
		}

		r.gauge.track.StrokeColor = ColorSurfaceLow
		r.gauge.bottomTrack.FillColor = ColorSurfaceLow

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
	Icon    fyne.Resource
	History []float64
	MaxVal  float64

	isHovered bool
}

func NewTemperatureHistogram(label string, icon fyne.Resource) *TemperatureHistogram {
	h := &TemperatureHistogram{
		Label:   label,
		Icon:    icon,
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
	title.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}

	icon := canvas.NewImageFromResource(h.Icon)
	icon.FillMode = canvas.ImageFillContain
	icon.SetMinSize(fyne.NewSize(16, 16))

	rangeLabel := widget.NewLabel("RANGE: 0-100°C")
	rangeLabel.TextStyle = fyne.TextStyle{Monospace: true}
	rangeLabel.Alignment = fyne.TextAlignTrailing

	header := container.NewHBox(icon, title, layout.NewSpacer(), rangeLabel)

	bars := container.NewMax() // Container for bars

	t1 := widget.NewLabel("-15m")
	t2 := widget.NewLabel("-10m")
	t3 := widget.NewLabel("-5m")
	t4 := widget.NewLabel("NOW")
	for _, l := range []*widget.Label{t1, t2, t3, t4} {
		l.TextStyle = fyne.TextStyle{Monospace: true}
		l.Alignment = fyne.TextAlignCenter
	}

	footer := container.NewHBox(t1, layout.NewSpacer(), t2, layout.NewSpacer(), t3, layout.NewSpacer(), t4)

	content := container.NewVBox(header, container.NewPadded(bars), footer)

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
		maxH := size.Height - 60 // Subtract header and footer height
		
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

			rect.Resize(fyne.NewSize(barWidth-2, h))
			rect.Move(fyne.NewPos(float32(i)*barWidth, maxH-h))
			r.bars.Add(rect)
		}
	})
}

func (r *tempHistogramRenderer) MinSize() fyne.Size {
	return fyne.NewSize(400, 200)
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
