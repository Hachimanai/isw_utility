package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

var (
	// Terminal Architect Palette
	colorBackground    = color.RGBA{R: 0x0c, G: 0x0d, B: 0x18, A: 0xff} // Deep-space background
	colorPrimary      = color.RGBA{R: 0x8b, G: 0xce, B: 0xff, A: 0xff} // Primary blue
	colorSurface       = color.RGBA{R: 0x16, G: 0x18, B: 0x2c, A: 0xff} // Surface container
	colorText          = color.RGBA{R: 0xe3, G: 0xe3, B: 0xff, A: 0xff} // Main text
	colorError         = color.RGBA{R: 0xee, G: 0x7d, B: 0x77, A: 0xff} // Critical error/temp
	colorSecondaryGPU = color.RGBA{R: 0x17, G: 0x93, B: 0xd1, A: 0xff} // GPU secondary color
)

// ArchitectTheme implements a custom fyne.Theme following the "Terminal Architect" design system.
type ArchitectTheme struct{}

var _ fyne.Theme = (*ArchitectTheme)(nil)

// Color returns the color for a specific theme name and variant.
// This theme is optimized for a dark appearance regardless of the system variant.
func (m *ArchitectTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	// Default to dark variant for fallback colors as this is a dark-only theme design
	fallbackVariant := theme.VariantDark

	switch name {
	case theme.ColorNameBackground:
		return colorBackground
	case theme.ColorNamePrimary:
		return colorPrimary
	case theme.ColorNameInputBackground:
		return colorSurface
	case theme.ColorNameForeground:
		return colorText
	case theme.ColorNameError:
		return colorError
	case theme.ColorNameButton:
		return colorSurface
	// "No-Line" Rule: Use background color for separators to avoid visible 1px borders
	case theme.ColorNameSeparator:
		return colorBackground
	case theme.ColorNameShadow:
		return color.Transparent
	}

	return theme.DefaultTheme().Color(name, fallbackVariant)
}

// Font returns the resource for the requested font style.
// Currently returns Fyne's default font as a placeholder for "Space Grotesk".
func (m *ArchitectTheme) Font(style fyne.TextStyle) fyne.Resource {
	// TODO: Replace with "Space Grotesk" / "Monospace" as per design spec
	return theme.DefaultTheme().Font(style)
}

// Icon returns the resource for the requested icon name.
func (m *ArchitectTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

// Size returns the standard size for the requested size name.
func (m *ArchitectTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

// NewArchitectTheme returns a new instance of ArchitectTheme.
func NewArchitectTheme() fyne.Theme {
	return &ArchitectTheme{}
}
