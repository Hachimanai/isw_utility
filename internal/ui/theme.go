package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

var (
	// Terminal Architect Palette
	ColorBackground       = color.RGBA{R: 0x0c, G: 0x0d, B: 0x18, A: 0xff} // Deep-space background
	ColorPrimary          = color.RGBA{R: 0x8b, G: 0xce, B: 0xff, A: 0xff} // Primary blue
	ColorSurface          = color.RGBA{R: 0x16, G: 0x18, B: 0x2c, A: 0xff} // Surface container
	ColorSurfaceLow       = color.RGBA{R: 0x11, G: 0x12, B: 0x21, A: 0xff} // Lower surface
	ColorSurfaceHigh      = color.RGBA{R: 0x20, G: 0x23, B: 0x41, A: 0xff} // Higher surface
	ColorText             = color.RGBA{R: 0xe3, G: 0xe3, B: 0xff, A: 0xff} // Main text
	ColorError            = color.RGBA{R: 0xee, G: 0x7d, B: 0x77, A: 0xff} // Critical error/temp
	ColorSecondaryGPU     = color.RGBA{R: 0x17, G: 0x93, B: 0xd1, A: 0xff} // GPU secondary color
)

// ArchitectTheme implements a custom fyne.Theme following the "Terminal Architect" design system.
type ArchitectTheme struct{}

var _ fyne.Theme = (*ArchitectTheme)(nil)

// Color returns the color for a specific theme name and variant.
func (m *ArchitectTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	fallbackVariant := theme.VariantDark

	switch name {
	case theme.ColorNameBackground:
		return ColorBackground
	case theme.ColorNamePrimary:
		return ColorPrimary
	case theme.ColorNameInputBackground:
		return ColorSurface
	case theme.ColorNameForeground:
		return ColorText
	case theme.ColorNameError:
		return ColorError
	case theme.ColorNameButton:
		return ColorSurface
	case theme.ColorNameSeparator:
		return ColorBackground // "No-Line" Rule
	case theme.ColorNameShadow:
		return color.Transparent
	}

	return theme.DefaultTheme().Color(name, fallbackVariant)
}

// Font returns the resource for the requested font style.
func (m *ArchitectTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

// Icon returns the resource for the requested icon name.
func (m *ArchitectTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

// Size returns the standard size for the requested size name.
func (m *ArchitectTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNamePadding:
		return 8
	case theme.SizeNameInnerPadding:
		return 4
	case theme.SizeNameScrollBarRadius:
		return 5
	case theme.SizeNameSelectionRadius:
		return 2
	}
	return theme.DefaultTheme().Size(name)
}

// NewArchitectTheme returns a new instance of ArchitectTheme.
func NewArchitectTheme() fyne.Theme {
	return &ArchitectTheme{}
}
