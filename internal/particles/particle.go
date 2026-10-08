package particles

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Particle represents a single dynamic visual effect particle.
type Particle struct {
	X, Y       float64
	VX, VY     float64
	Life       float64
	MaxLife    float64
	Size       float32
	StartSize  float32
	EndSize    float32
	Color      color.RGBA
	FadeColor  color.RGBA
	IsAdditive bool
	Drag       float64
	IsDead     bool
}

func (p *Particle) Update(dt float64) {
	if p.IsDead {
		return
	}
	p.Life -= dt
	if p.Life <= 0 {
		p.IsDead = true
		return
	}

	// Apply drag and movement
	p.VX *= math.Pow(p.Drag, dt*60.0)
	p.VY *= math.Pow(p.Drag, dt*60.0)
	p.X += p.VX * dt
	p.Y += p.VY * dt

	// Interpolate size
	t := 1.0 - (p.Life / p.MaxLife)
	p.Size = p.StartSize + (p.EndSize-p.StartSize)*float32(t)
}

func (p *Particle) Draw(screen *ebiten.Image) {
	if p.IsDead || p.Size <= 0.1 {
		return
	}
	t := 1.0 - (p.Life / p.MaxLife)
	alpha := uint8(math.Max(0, math.Min(255, float64(p.Color.A)*(1.0-t))))

	c := color.RGBA{
		R: uint8(float64(p.Color.R)*(1-t) + float64(p.FadeColor.R)*t),
		G: uint8(float64(p.Color.G)*(1-t) + float64(p.FadeColor.G)*t),
		B: uint8(float64(p.Color.B)*(1-t) + float64(p.FadeColor.B)*t),
		A: alpha,
	}

	vector.FillCircle(screen, float32(p.X), float32(p.Y), p.Size, c, false)
}
