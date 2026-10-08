package particles

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// CombatText represents a floating message (damage number, crit, level up, status).
type CombatText struct {
	Text    string
	X, Y    float64
	VY      float64
	Life    float64
	MaxLife float64
	Color   color.RGBA
	IsDead  bool
}

func (ct *CombatText) Update(dt float64) {
	if ct.IsDead {
		return
	}
	ct.Life -= dt
	if ct.Life <= 0 {
		ct.IsDead = true
		return
	}
	ct.Y += ct.VY * dt
	ct.VY *= math.Pow(0.92, dt*60.0)
}

func (ct *CombatText) Draw(screen *ebiten.Image) {
	if ct.IsDead {
		return
	}
	ebitenutil.DebugPrintAt(screen, ct.Text, int(ct.X), int(ct.Y))
}
