package particles

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

// ParticleManager coordinates particles, floating numbers, and screen shake.
type ParticleManager struct {
	Particles   []*Particle
	CombatTexts []*CombatText
	Shake       *ScreenShake
}

func NewParticleManager() *ParticleManager {
	return &ParticleManager{
		Particles:   make([]*Particle, 0, 500),
		CombatTexts: make([]*CombatText, 0, 100),
		Shake:       NewScreenShake(),
	}
}

func (pm *ParticleManager) EmitHitSparks(x, y float64, count int, c color.RGBA) {
	for i := 0; i < count; i++ {
		angle := rand.Float64() * math.Pi * 2.0
		speed := 40.0 + rand.Float64()*120.0
		life := 0.2 + rand.Float64()*0.3
		pm.Particles = append(pm.Particles, &Particle{
			X:          x,
			Y:          y,
			VX:         math.Cos(angle) * speed,
			VY:         math.Sin(angle) * speed,
			Life:       life,
			MaxLife:    life,
			StartSize:  3.0 + rand.Float32()*2.0,
			EndSize:    0.5,
			Color:      c,
			FadeColor:  color.RGBA{255, 255, 255, 0},
			IsAdditive: true,
			Drag:       0.90,
		})
	}
}

func (pm *ParticleManager) EmitExplosion(x, y float64, radius float64, c color.RGBA) {
	count := int(radius * 0.8)
	if count < 15 {
		count = 15
	}
	for i := 0; i < count; i++ {
		angle := rand.Float64() * math.Pi * 2.0
		speed := 30.0 + rand.Float64()*radius*2.5
		life := 0.35 + rand.Float64()*0.4
		pm.Particles = append(pm.Particles, &Particle{
			X:          x,
			Y:          y,
			VX:         math.Cos(angle) * speed,
			VY:         math.Sin(angle) * speed,
			Life:       life,
			MaxLife:    life,
			StartSize:  4.0 + rand.Float32()*3.5,
			EndSize:    0.5,
			Color:      c,
			FadeColor:  color.RGBA{R: 255, G: 80, B: 20, A: 0},
			IsAdditive: true,
			Drag:       0.88,
		})
	}
	pm.Shake.AddTrauma(0.25)
}

func (pm *ParticleManager) EmitManaVapor(x, y float64, c color.RGBA) {
	for i := 0; i < 2; i++ {
		angle := -math.Pi/2 + (rand.Float64()-0.5)*0.8
		speed := 20.0 + rand.Float64()*35.0
		life := 0.4 + rand.Float64()*0.3
		pm.Particles = append(pm.Particles, &Particle{
			X:          x + (rand.Float64()-0.5)*10.0,
			Y:          y + (rand.Float64()-0.5)*10.0,
			VX:         math.Cos(angle) * speed,
			VY:         math.Sin(angle) * speed,
			Life:       life,
			MaxLife:    life,
			StartSize:  3.5,
			EndSize:    1.0,
			Color:      c,
			FadeColor:  color.RGBA{220, 160, 255, 0},
			IsAdditive: true,
			Drag:       0.95,
		})
	}
}

func (pm *ParticleManager) AddCombatText(text string, x, y float64, c color.RGBA) {
	pm.CombatTexts = append(pm.CombatTexts, &CombatText{
		Text:    text,
		X:       x + (rand.Float64()-0.5)*8.0,
		Y:       y - 10.0,
		VY:      -35.0 - rand.Float64()*20.0,
		Life:    0.65,
		MaxLife: 0.65,
		Color:   c,
	})
}

func (pm *ParticleManager) AddDamageText(dmg float64, x, y float64, isCrit bool) {
	c := color.RGBA{R: 255, G: 220, B: 100, A: 255}
	text := fmt.Sprintf("%.0f", dmg)
	if isCrit {
		c = color.RGBA{R: 255, G: 50, B: 70, A: 255}
		text = fmt.Sprintf("!%.0f!", dmg)
	}
	pm.AddCombatText(text, x, y, c)
}

func (pm *ParticleManager) Update(dt float64) {
	pm.Shake.Update(dt)

	// Update Particles
	aliveP := pm.Particles[:0]
	for _, p := range pm.Particles {
		p.Update(dt)
		if !p.IsDead {
			aliveP = append(aliveP, p)
		}
	}
	pm.Particles = aliveP

	// Update Combat Texts
	aliveCT := pm.CombatTexts[:0]
	for _, ct := range pm.CombatTexts {
		ct.Update(dt)
		if !ct.IsDead {
			aliveCT = append(aliveCT, ct)
		}
	}
	pm.CombatTexts = aliveCT
}

func (pm *ParticleManager) Draw(screen *ebiten.Image) {
	for _, p := range pm.Particles {
		p.Draw(screen)
	}
	for _, ct := range pm.CombatTexts {
		ct.Draw(screen)
	}
}
