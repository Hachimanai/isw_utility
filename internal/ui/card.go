package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

// TonalCard is a container with a themed background and padding.
func TonalCard(objects ...fyne.CanvasObject) fyne.CanvasObject {
	bg := canvas.NewRectangle(ColorSurface)
	content := container.NewPadded(objects...)
	return container.NewStack(bg, content)
}

func TonalCardLow(objects ...fyne.CanvasObject) fyne.CanvasObject {
	bg := canvas.NewRectangle(ColorSurfaceLow)
	content := container.NewPadded(objects...)
	return container.NewStack(bg, content)
}

func TonalCardHigh(objects ...fyne.CanvasObject) fyne.CanvasObject {
	bg := canvas.NewRectangle(ColorSurfaceHigh)
	content := container.NewPadded(objects...)
	return container.NewStack(bg, content)
}
