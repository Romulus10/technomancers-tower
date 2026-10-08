package gfx

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// GenerateTowerSprite generates a procedural 32x32 sprite for a specific tower ID and tier.
func GenerateTowerSprite(id string, tier int, baseColor color.RGBA) *ebiten.Image {
	const size = 32
	img := ebiten.NewImage(size, size)

	cx := float32(size / 2)
	cy := float32(size / 2)

	// 1. Silicon Substrate Chip Base (Rotated 45 deg or beveled square)
	chipColor := color.RGBA{R: 20, G: 26, B: 38, A: 255}
	accentColor := baseColor
	glowColor := FadeAlpha(baseColor, 0.4)

	// Base plate
	vector.FillRect(img, 4, 4, 24, 24, chipColor, false)

	// Corner gold solder pads / pins
	padColor := ColorGoldMatrix
	vector.FillRect(img, 2, 6, 2, 3, padColor, false)
	vector.FillRect(img, 2, 14, 2, 4, padColor, false)
	vector.FillRect(img, 2, 23, 2, 3, padColor, false)

	vector.FillRect(img, 28, 6, 2, 3, padColor, false)
	vector.FillRect(img, 28, 14, 2, 4, padColor, false)
	vector.FillRect(img, 28, 23, 2, 3, padColor, false)

	vector.FillRect(img, 6, 2, 3, 2, padColor, false)
	vector.FillRect(img, 14, 2, 4, 2, padColor, false)
	vector.FillRect(img, 23, 2, 3, 2, padColor, false)

	vector.FillRect(img, 6, 28, 3, 2, padColor, false)
	vector.FillRect(img, 14, 28, 4, 2, padColor, false)
	vector.FillRect(img, 23, 28, 3, 2, padColor, false)

	// Chip Border & Inner Chamfer
	vector.StrokeRect(img, 4, 4, 24, 24, 1.2, Darken(accentColor, 0.2), false)
	vector.StrokeRect(img, 6, 6, 20, 20, 1.0, color.RGBA{R: 35, G: 45, B: 65, A: 255}, false)

	// 2. Specific Archetype Circuit Geometry
	switch id {
	case "bit_driver":
		// High-speed dual laser barrel emitter
		vector.FillRect(img, 8, 8, 16, 16, color.RGBA{R: 15, G: 20, B: 30, A: 255}, false)
		vector.FillRect(img, 10, 12, 12, 8, Darken(accentColor, 0.3), false)
		// Dual focus rails
		vector.FillRect(img, 11, 6, 3, 10, accentColor, false)
		vector.FillRect(img, 18, 6, 3, 10, accentColor, false)
		// Glowing emitter core
		vector.FillCircle(img, cx, cy+2, 4, Brighten(accentColor, 0.3), false)
		vector.FillCircle(img, cx, cy+2, 2, color.RGBA{255, 255, 255, 255}, false)

	case "mana_siphon":
		// Arcane siphon crystal & containment rings
		vector.FillCircle(img, cx, cy, 8, color.RGBA{R: 30, G: 15, B: 45, A: 255}, false)
		vector.StrokeCircle(img, cx, cy, 9, 1.5, accentColor, false)
		// Crystal diamond in center
		vector.FillRect(img, 13, 9, 6, 14, accentColor, false)
		vector.FillRect(img, 9, 13, 14, 6, accentColor, false)
		vector.FillCircle(img, cx, cy, 3.5, color.RGBA{255, 220, 255, 255}, false)
		// Energy conduit lines to 4 corners
		vector.StrokeLine(img, 6, 6, 11, 11, 1.2, glowColor, false)
		vector.StrokeLine(img, 26, 6, 21, 11, 1.2, glowColor, false)
		vector.StrokeLine(img, 6, 26, 11, 21, 1.2, glowColor, false)
		vector.StrokeLine(img, 26, 26, 21, 21, 1.2, glowColor, false)

	case "crypto_miner", "quantum_miner":
		// Matrix computation array / heatsink fins
		vector.FillRect(img, 7, 7, 18, 18, color.RGBA{R: 25, G: 35, B: 30, A: 255}, false)
		for i := 0; i < 4; i++ {
			y := float32(9 + i*4)
			vector.FillRect(img, 9, y, 14, 2, accentColor, false)
		}
		// Central glowing microprocessor die
		vector.FillRect(img, 13, 13, 6, 6, Brighten(accentColor, 0.5), false)
		vector.FillRect(img, 14, 14, 4, 4, color.RGBA{255, 255, 255, 255}, false)

	case "tesla_bus":
		// Tesla logic coil & capacitor banks
		vector.FillCircle(img, cx, cy, 9, color.RGBA{R: 35, G: 30, B: 15, A: 255}, false)
		vector.StrokeCircle(img, cx, cy, 9, 1.2, Darken(accentColor, 0.3), false)
		// Central Tesla sphere
		vector.FillCircle(img, cx, cy, 5.5, accentColor, false)
		vector.FillCircle(img, cx, cy, 3.0, color.RGBA{255, 255, 200, 255}, false)
		// 4 peripheral mini electrode nodes
		vector.FillCircle(img, 8, 8, 2, padColor, false)
		vector.FillCircle(img, 24, 8, 2, padColor, false)
		vector.FillCircle(img, 8, 24, 2, padColor, false)
		vector.FillCircle(img, 24, 24, 2, padColor, false)
		// Lightning cross
		vector.StrokeLine(img, cx, 6, cx, 26, 1.2, glowColor, false)
		vector.StrokeLine(img, 6, cy, 26, cy, 1.2, glowColor, false)

	case "cryo_cache":
		// Cryogenic chiller & snowflake radiator
		vector.FillRect(img, 8, 8, 16, 16, color.RGBA{R: 15, G: 30, B: 45, A: 255}, false)
		vector.StrokeRect(img, 8, 8, 16, 16, 1.2, accentColor, false)
		// Snowflake diamond core
		vector.StrokeLine(img, cx, 8, cx, 24, 2.0, Brighten(accentColor, 0.4), false)
		vector.StrokeLine(img, 8, cy, 24, cy, 2.0, Brighten(accentColor, 0.4), false)
		vector.StrokeLine(img, 11, 11, 21, 21, 1.5, Brighten(accentColor, 0.2), false)
		vector.StrokeLine(img, 21, 11, 11, 21, 1.5, Brighten(accentColor, 0.2), false)
		vector.FillCircle(img, cx, cy, 3.0, color.RGBA{230, 250, 255, 255}, false)

	case "logic_mortar":
		// Heavy ballistic turret housing
		vector.FillCircle(img, cx, cy, 10, color.RGBA{R: 35, G: 20, B: 15, A: 255}, false)
		vector.StrokeCircle(img, cx, cy, 10, 1.8, Darken(accentColor, 0.3), false)
		// Heavy mortar barrel muzzle
		vector.FillCircle(img, cx, cy, 6.0, accentColor, false)
		vector.FillCircle(img, cx, cy, 3.5, color.RGBA{R: 10, G: 10, B: 15, A: 255}, false)
		// Armored shroud bevels
		vector.FillRect(img, 5, 13, 4, 6, Darken(accentColor, 0.4), false)
		vector.FillRect(img, 23, 13, 4, 6, Darken(accentColor, 0.4), false)

	case "emp_railgun":
		// High-voltage magnetic rail accelerators
		vector.FillRect(img, 7, 7, 18, 18, color.RGBA{R: 10, G: 30, B: 35, A: 255}, false)
		// Parallel magnetic coils
		vector.FillRect(img, 9, 8, 4, 16, accentColor, false)
		vector.FillRect(img, 19, 8, 4, 16, accentColor, false)
		// Central plasma acceleration channel
		vector.FillRect(img, 14, 6, 4, 20, Brighten(accentColor, 0.5), false)
		vector.FillRect(img, 15, 6, 2, 20, color.RGBA{255, 255, 255, 255}, false)

	case "plasma_flak":
		// Quad-barrel flak matrix
		vector.FillRect(img, 8, 8, 16, 16, color.RGBA{R: 35, G: 25, B: 15, A: 255}, false)
		vector.FillCircle(img, 11, 11, 3.0, accentColor, false)
		vector.FillCircle(img, 21, 11, 3.0, accentColor, false)
		vector.FillCircle(img, 11, 21, 3.0, accentColor, false)
		vector.FillCircle(img, 21, 21, 3.0, accentColor, false)
		vector.FillCircle(img, cx, cy, 3.0, Brighten(accentColor, 0.4), false)

	case "overclock_beacon":
		// Omni-directional pulse transmitter
		vector.FillCircle(img, cx, cy, 10, color.RGBA{R: 35, G: 15, B: 40, A: 255}, false)
		vector.StrokeCircle(img, cx, cy, 10, 1.5, accentColor, false)
		vector.StrokeCircle(img, cx, cy, 6, 1.2, Brighten(accentColor, 0.3), false)
		vector.FillCircle(img, cx, cy, 3.5, color.RGBA{255, 230, 255, 255}, false)

	case "singularity_vortex":
		// Dark matter gravity containment node
		vector.FillCircle(img, cx, cy, 10, color.RGBA{R: 15, G: 10, B: 30, A: 255}, false)
		vector.StrokeCircle(img, cx, cy, 10, 2.0, accentColor, false)
		// Gravitational vortex spiral blades
		vector.FillRect(img, 14, 6, 4, 20, Darken(accentColor, 0.2), false)
		vector.FillRect(img, 6, 14, 20, 4, Darken(accentColor, 0.2), false)
		// Event horizon core
		vector.FillCircle(img, cx, cy, 4.5, color.RGBA{5, 5, 10, 255}, false)
		vector.StrokeCircle(img, cx, cy, 4.5, 1.5, Brighten(accentColor, 0.6), false)

	case "bit_leech":
		// Vampiric hex-mesh siphon
		vector.FillRect(img, 8, 8, 16, 16, color.RGBA{R: 15, G: 35, B: 20, A: 255}, false)
		vector.StrokeRect(img, 8, 8, 16, 16, 1.5, accentColor, false)
		// Leech fangs / energy drainers
		vector.FillRect(img, 10, 9, 4, 14, Brighten(accentColor, 0.3), false)
		vector.FillRect(img, 18, 9, 4, 14, Brighten(accentColor, 0.3), false)
		vector.FillCircle(img, cx, cy, 3.0, color.RGBA{200, 255, 220, 255}, false)

	default:
		// Generic cyber turret
		vector.FillCircle(img, cx, cy, 8, accentColor, false)
		vector.FillCircle(img, cx, cy, 4, color.RGBA{255, 255, 255, 255}, false)
	}

	// 3. Procedural Tier Evolution Adornments (Tier 2..5)
	if tier >= 2 {
		// Tier 2: Reinforced gold corner braces
		vector.StrokeLine(img, 4, 8, 8, 4, 1.5, ColorGoldMatrix, false)
		vector.StrokeLine(img, 28, 8, 24, 4, 1.5, ColorGoldMatrix, false)
		vector.StrokeLine(img, 4, 24, 8, 28, 1.5, ColorGoldMatrix, false)
		vector.StrokeLine(img, 28, 24, 24, 28, 1.5, ColorGoldMatrix, false)
	}
	if tier >= 3 {
		// Tier 3: Radiant energized sub-chassis glow
		vector.StrokeRect(img, 5, 5, 22, 22, 1.0, Brighten(accentColor, 0.4), false)
	}
	if tier >= 4 {
		// Tier 4: Arcane runic orbital nodes
		vector.FillCircle(img, 6, 6, 1.8, ColorCyanNeon, false)
		vector.FillCircle(img, 26, 6, 1.8, ColorCyanNeon, false)
		vector.FillCircle(img, 6, 26, 1.8, ColorCyanNeon, false)
		vector.FillCircle(img, 26, 26, 1.8, ColorCyanNeon, false)
	}
	if tier >= 5 {
		// Tier 5: Overcharged Master Crown Sigil
		vector.StrokeCircle(img, cx, cy, 13, 1.5, ColorGoldMatrix, false)
		vector.StrokeRect(img, 3, 3, 26, 26, 1.5, ColorGoldMatrix, false)
	}

	return img
}

// GenerateEnemySprite generates a procedural sprite for malware packets.
func GenerateEnemySprite(id string, isBoss bool, baseColor color.RGBA, radius float32) *ebiten.Image {
	size := int(math.Ceil(float64(radius*2.0 + 8.0)))
	if size%2 != 0 {
		size++
	}
	if size < 20 {
		size = 20
	}
	if isBoss && size < 44 {
		size = 44
	}

	img := ebiten.NewImage(size, size)
	cx := float32(size / 2)
	cy := float32(size / 2)

	r := radius
	if r < 4 {
		r = 4
	}

	switch id {
	case "swarmer":
		// Fast diamond data packet with trailing thrusters
		vector.FillCircle(img, cx, cy, r, baseColor, false)
		// Inner sharp core
		vector.FillRect(img, cx-r*0.5, cy-r*0.5, r, r, Brighten(baseColor, 0.4), false)
		vector.StrokeCircle(img, cx, cy, r, 1.2, color.RGBA{255, 255, 255, 200}, false)

	case "sprite":
		// Shielded crystalline hexagon
		vector.FillCircle(img, cx, cy, r, baseColor, false)
		vector.StrokeCircle(img, cx, cy, r+1.5, 1.5, ColorCyanNeon, false)
		// Internal crystal cross
		vector.StrokeLine(img, cx-r*0.6, cy, cx+r*0.6, cy, 1.5, color.RGBA{255, 255, 255, 255}, false)
		vector.StrokeLine(img, cx, cy-r*0.6, cx, cy+r*0.6, 1.5, color.RGBA{255, 255, 255, 255}, false)

	case "golem":
		// Heavy armored trojan tank - faceted hexagonal armor plates
		vector.FillCircle(img, cx, cy, r, color.RGBA{R: 45, G: 30, B: 15, A: 255}, false)
		vector.StrokeCircle(img, cx, cy, r, 2.5, baseColor, false)
		// Thick armor chevron plates
		vector.FillRect(img, cx-r*0.7, cy-r*0.7, r*1.4, r*1.4, Darken(baseColor, 0.2), false)
		vector.FillCircle(img, cx, cy, r*0.4, ColorGoldMatrix, false)
		vector.FillCircle(img, cx, cy, r*0.2, color.RGBA{255, 255, 255, 255}, false)

	case "trojan_shard":
		// Jagged shard fragment
		vector.FillCircle(img, cx, cy, r, baseColor, false)
		vector.FillRect(img, cx-2, cy-2, 4, 4, color.RGBA{255, 255, 255, 255}, false)
		vector.StrokeCircle(img, cx, cy, r, 1.0, ColorAmberFire, false)

	case "healer_worm":
		// Bio-digital parasite segment with emerald medical cross
		vector.FillCircle(img, cx, cy, r, color.RGBA{R: 15, G: 45, B: 25, A: 255}, false)
		vector.StrokeCircle(img, cx, cy, r, 1.8, baseColor, false)
		// Emerald Medical Cross
		crossW := r * 0.4
		crossL := r * 1.2
		vector.FillRect(img, cx-crossW/2, cy-crossL/2, crossW, crossL, Brighten(baseColor, 0.3), false)
		vector.FillRect(img, cx-crossL/2, cy-crossW/2, crossL, crossW, Brighten(baseColor, 0.3), false)
		vector.FillCircle(img, cx, cy, crossW*0.6, color.RGBA{255, 255, 255, 255}, false)

	case "phantom_glitch":
		// Shifting ethereal glitch packet
		vector.FillCircle(img, cx, cy, r, color.RGBA{R: 40, G: 15, B: 55, A: 200}, false)
		vector.StrokeCircle(img, cx, cy, r, 1.5, baseColor, false)
		// Glitch scanlines
		vector.FillRect(img, cx-r*0.8, cy-2, r*1.6, 1.5, ColorPurpleArcane, false)
		vector.FillRect(img, cx-r*0.5, cy+3, r*1.0, 1.5, ColorCyanNeon, false)
		vector.FillCircle(img, cx, cy, r*0.35, color.RGBA{255, 220, 255, 255}, false)

	case "buffer_overflower":
		// Unstable volatility bomb
		vector.FillCircle(img, cx, cy, r, color.RGBA{R: 50, G: 10, B: 10, A: 255}, false)
		vector.StrokeCircle(img, cx, cy, r, 2.0, baseColor, false)
		// Radiation / Hazard symbol
		vector.FillCircle(img, cx, cy, r*0.45, ColorAmberFire, false)
		vector.FillCircle(img, cx, cy, r*0.25, ColorGoldMatrix, false)
		vector.StrokeLine(img, cx-r*0.7, cy-r*0.7, cx+r*0.7, cy+r*0.7, 1.5, ColorCrimsonGlitch, false)
		vector.StrokeLine(img, cx+r*0.7, cy-r*0.7, cx-r*0.7, cy+r*0.7, 1.5, ColorCrimsonGlitch, false)

	case "boss_daemon":
		// Tier 1 Boss: Colossal Rootkit Mainframe
		vector.FillCircle(img, cx, cy, r, color.RGBA{R: 35, G: 8, B: 18, A: 255}, false)
		vector.StrokeCircle(img, cx, cy, r, 3.0, baseColor, false)
		vector.StrokeCircle(img, cx, cy, r-4, 1.5, ColorGoldMatrix, false)
		// Skull / Demon Horn Matrix geometry
		vector.FillRect(img, cx-r*0.6, cy-r*0.6, r*1.2, r*1.2, Darken(baseColor, 0.3), false)
		// Glowing red demonic ocular core
		vector.FillCircle(img, cx-5, cy-2, 3.0, ColorCrimsonGlitch, false)
		vector.FillCircle(img, cx+5, cy-2, 3.0, ColorCrimsonGlitch, false)
		vector.FillCircle(img, cx, cy+4, 4.0, Brighten(baseColor, 0.4), false)
		vector.FillCircle(img, cx, cy+4, 2.0, color.RGBA{255, 255, 255, 255}, false)

	case "boss_zeroday":
		// Tier 2 Boss: Zero-Day Harbinger Void Core
		vector.FillCircle(img, cx, cy, r, color.RGBA{R: 25, G: 5, B: 30, A: 255}, false)
		vector.StrokeCircle(img, cx, cy, r, 3.0, baseColor, false)
		// Triple orbiting arcane containment blades
		vector.StrokeCircle(img, cx, cy, r-3, 2.0, ColorCyanNeon, false)
		vector.StrokeCircle(img, cx, cy, r-7, 1.5, ColorPurpleArcane, false)
		// Pulsing zero-point singularity core
		vector.FillCircle(img, cx, cy, r*0.4, color.RGBA{5, 5, 10, 255}, false)
		vector.FillCircle(img, cx, cy, r*0.2, color.RGBA{255, 255, 255, 255}, false)

	default:
		vector.FillCircle(img, cx, cy, r, baseColor, false)
		vector.StrokeCircle(img, cx, cy, r, 1.5, color.RGBA{255, 255, 255, 255}, false)
	}

	return img
}

// GenerateCPUCoreSprite generates the CPU Kernel Core sprite (48x48).
func GenerateCPUCoreSprite() *ebiten.Image {
	const size = 48
	img := ebiten.NewImage(size, size)
	cx := float32(size / 2)
	cy := float32(size / 2)

	// Outer motherboard socket base
	vector.FillRect(img, 4, 4, 40, 40, color.RGBA{R: 16, G: 24, B: 36, A: 255}, false)
	vector.StrokeRect(img, 4, 4, 40, 40, 2.0, ColorCyanNeon, false)

	// Gold pin array along socket perimeter
	for i := 6; i <= 42; i += 4 {
		vector.FillRect(img, float32(i), 2, 2, 2, ColorGoldMatrix, false)
		vector.FillRect(img, float32(i), 44, 2, 2, ColorGoldMatrix, false)
		vector.FillRect(img, 2, float32(i), 2, 2, ColorGoldMatrix, false)
		vector.FillRect(img, 44, float32(i), 2, 2, ColorGoldMatrix, false)
	}

	// Heat spreader lid
	vector.FillRect(img, 8, 8, 32, 32, color.RGBA{R: 24, G: 34, B: 50, A: 255}, false)
	vector.StrokeRect(img, 8, 8, 32, 32, 1.5, ColorFrostIce, false)

	// Central Silicon Die
	vector.FillRect(img, 14, 14, 20, 20, color.RGBA{R: 12, G: 18, B: 28, A: 255}, false)
	vector.StrokeRect(img, 14, 14, 20, 20, 1.2, ColorCyanNeon, false)

	// Glowing Cyber Arcane Heart
	vector.FillCircle(img, cx, cy, 7.0, ColorCyanNeon, false)
	vector.FillCircle(img, cx, cy, 4.0, ColorFrostIce, false)
	vector.FillCircle(img, cx, cy, 2.0, color.RGBA{255, 255, 255, 255}, false)

	return img
}

// GenerateSpawnerPortSprite generates the Perimeter Spawner socket (32x32).
func GenerateSpawnerPortSprite() *ebiten.Image {
	const size = 32
	img := ebiten.NewImage(size, size)
	cx := float32(size / 2)
	cy := float32(size / 2)

	// Hazard warning chassis
	vector.FillRect(img, 3, 3, 26, 26, color.RGBA{R: 35, G: 15, B: 20, A: 255}, false)
	vector.StrokeRect(img, 3, 3, 26, 26, 1.5, ColorCrimsonGlitch, false)

	// Danger diagonal stripes on border
	vector.StrokeLine(img, 4, 8, 8, 4, 1.2, ColorAmberFire, false)
	vector.StrokeLine(img, 24, 28, 28, 24, 1.2, ColorAmberFire, false)
	vector.StrokeLine(img, 4, 24, 8, 28, 1.2, ColorAmberFire, false)
	vector.StrokeLine(img, 24, 4, 28, 8, 1.2, ColorAmberFire, false)

	// Intake breach aperture
	vector.FillCircle(img, cx, cy, 8.0, color.RGBA{R: 15, G: 5, B: 10, A: 255}, false)
	vector.StrokeCircle(img, cx, cy, 8.0, 1.8, ColorCrimsonGlitch, false)
	vector.FillCircle(img, cx, cy, 3.5, ColorAmberFire, false)

	return img
}

// DrawHolographicPanel draws a procedural holographic beveled cyber-panel directly to the target.
func DrawHolographicPanel(screen *ebiten.Image, x, y, w, h float32, borderColor, bgColor color.RGBA) {
	// Semi-transparent dark backplate
	vector.FillRect(screen, x, y, w, h, bgColor, false)

	// Outer border frame
	vector.StrokeRect(screen, x, y, w, h, 1.5, borderColor, false)

	// Glowing sci-fi corner bracket cuts
	cLen := float32(10.0)
	brightBorder := Brighten(borderColor, 0.4)

	// Top-left
	vector.StrokeLine(screen, x, y, x+cLen, y, 2.5, brightBorder, false)
	vector.StrokeLine(screen, x, y, x, y+cLen, 2.5, brightBorder, false)

	// Top-right
	vector.StrokeLine(screen, x+w, y, x+w-cLen, y, 2.5, brightBorder, false)
	vector.StrokeLine(screen, x+w, y, x+w, y+cLen, 2.5, brightBorder, false)

	// Bottom-left
	vector.StrokeLine(screen, x, y+h, x+cLen, y+h, 2.5, brightBorder, false)
	vector.StrokeLine(screen, x, y+h, x, y+h-cLen, 2.5, brightBorder, false)

	// Bottom-right
	vector.StrokeLine(screen, x+w, y+h, x+w-cLen, y+h, 2.5, brightBorder, false)
	vector.StrokeLine(screen, x+w, y+h, x+w, y+h-cLen, 2.5, brightBorder, false)
}

// DrawGradientBar renders an enhanced multi-layer HUD status bar with glowing border and segments.
func DrawGradientBar(screen *ebiten.Image, x, y, w, h float32, ratio float32, primaryColor, secondaryColor color.RGBA) {
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}

	// Bar background track
	vector.FillRect(screen, x, y, w, h, color.RGBA{R: 16, G: 20, B: 28, A: 255}, false)

	// Bar fill with gradient look
	fillW := w * ratio
	if fillW > 0 {
		vector.FillRect(screen, x, y, fillW, h, primaryColor, false)
		// Top highlight sliver
		vector.FillRect(screen, x, y, fillW, h*0.35, Brighten(primaryColor, 0.35), false)
	}

	// Track border
	vector.StrokeRect(screen, x, y, w, h, 1.2, secondaryColor, false)

	// Segment dividers (every 25%)
	for i := 1; i <= 3; i++ {
		divX := x + w*(float32(i)*0.25)
		vector.StrokeLine(screen, divX, y, divX, y+h, 1.0, color.RGBA{R: 20, G: 30, B: 45, A: 200}, false)
	}
}

// GenerateIconBadge creates a stylized 28x28 icon for a tower or spell.
func GenerateIconBadge(id string, isTower bool, baseColor color.RGBA) *ebiten.Image {
	const size = 28
	img := ebiten.NewImage(size, size)
	cx := float32(size / 2)
	cy := float32(size / 2)

	// Background plate
	vector.FillRect(img, 2, 2, 24, 24, color.RGBA{R: 18, G: 24, B: 36, A: 255}, false)
	vector.StrokeRect(img, 2, 2, 24, 24, 1.2, baseColor, false)

	if isTower {
		// Render mini tower symbol
		vector.FillRect(img, 7, 7, 14, 14, Darken(baseColor, 0.2), false)
		vector.FillCircle(img, cx, cy, 3.5, Brighten(baseColor, 0.4), false)
		vector.FillCircle(img, cx, cy, 1.5, color.RGBA{255, 255, 255, 255}, false)
	} else {
		// Render mini spell arcane sigil
		vector.StrokeCircle(img, cx, cy, 7.0, 1.5, baseColor, false)
		vector.FillCircle(img, cx, cy, 3.0, Brighten(baseColor, 0.5), false)
	}

	return img
}
