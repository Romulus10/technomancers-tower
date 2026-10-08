package gfx

import (
	"image/color"
	"math"
)

// Cyber-Arcane primary palette colors
var (
	ColorCyanNeon      = color.RGBA{R: 0, G: 220, B: 255, A: 255}
	ColorCyanGlow      = color.RGBA{R: 0, G: 180, B: 240, A: 160}
	ColorPurpleArcane  = color.RGBA{R: 170, G: 70, B: 255, A: 255}
	ColorPurpleGlow    = color.RGBA{R: 140, G: 40, B: 220, A: 160}
	ColorGoldMatrix    = color.RGBA{R: 255, G: 215, B: 30, A: 255}
	ColorEmeraldBio    = color.RGBA{R: 40, G: 240, B: 130, A: 255}
	ColorCrimsonGlitch = color.RGBA{R: 255, G: 50, B: 80, A: 255}
	ColorAmberFire     = color.RGBA{R: 255, G: 130, B: 30, A: 255}
	ColorFrostIce      = color.RGBA{R: 120, G: 220, B: 255, A: 255}

	// Panel backplates
	ColorPanelDark      = color.RGBA{R: 12, G: 16, B: 24, A: 245}
	ColorPanelHeader    = color.RGBA{R: 18, G: 24, B: 36, A: 255}
	ColorPanelBorder    = color.RGBA{R: 0, G: 180, B: 220, A: 220}
	ColorPanelHighlight = color.RGBA{R: 100, G: 240, B: 255, A: 255}
	ColorGridTrace      = color.RGBA{R: 18, G: 45, B: 65, A: 180}
	ColorGridTraceHot   = color.RGBA{R: 0, G: 200, B: 240, A: 230}
)

// LerpColor smoothly interpolates between two RGBA colors.
func LerpColor(c1, c2 color.RGBA, t float64) color.RGBA {
	t = math.Max(0.0, math.Min(1.0, t))
	return color.RGBA{
		R: uint8(float64(c1.R)*(1-t) + float64(c2.R)*t),
		G: uint8(float64(c1.G)*(1-t) + float64(c2.G)*t),
		B: uint8(float64(c1.B)*(1-t) + float64(c2.B)*t),
		A: uint8(float64(c1.A)*(1-t) + float64(c2.A)*t),
	}
}

// FadeAlpha returns a copy of the color with scaled alpha.
func FadeAlpha(c color.RGBA, alpha float64) color.RGBA {
	a := math.Max(0, math.Min(255, float64(c.A)*alpha))
	return color.RGBA{
		R: c.R,
		G: c.G,
		B: c.B,
		A: uint8(a),
	}
}

// Brighten shifts a color towards white by factor (0..1).
func Brighten(c color.RGBA, factor float64) color.RGBA {
	factor = math.Max(0.0, math.Min(1.0, factor))
	return color.RGBA{
		R: uint8(float64(c.R) + (255.0-float64(c.R))*factor),
		G: uint8(float64(c.G) + (255.0-float64(c.G))*factor),
		B: uint8(float64(c.B) + (255.0-float64(c.B))*factor),
		A: c.A,
	}
}

// Darken shifts a color towards black by factor (0..1).
func Darken(c color.RGBA, factor float64) color.RGBA {
	factor = math.Max(0.0, math.Min(1.0, factor))
	return color.RGBA{
		R: uint8(float64(c.R) * (1.0 - factor)),
		G: uint8(float64(c.G) * (1.0 - factor)),
		B: uint8(float64(c.B) * (1.0 - factor)),
		A: c.A,
	}
}
