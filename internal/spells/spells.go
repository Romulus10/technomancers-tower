package spells

import (
	"image/color"
	"math"

	"technomancers-tower/internal/data"
	"technomancers-tower/internal/economy"
	"technomancers-tower/internal/gfx"
	"technomancers-tower/internal/malware"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type SpellEffectVFX struct {
	X, Y      float64
	Radius    float64
	Life      float64
	MaxLife   float64
	IsThunder bool
	TargetX   float64
	TargetY   float64
	Color     color.RGBA
}

type SpellManager struct {
	Registry       *data.Registry
	Cooldowns      map[string]float64
	Flashes        map[string]float64
	SelectedSpell  string
	HasSelection   bool
	OverclockTimer float64
	ActiveEffects  []*SpellEffectVFX
}

func NewSpellManager(reg *data.Registry) *SpellManager {
	return &SpellManager{
		Registry:      reg,
		Cooldowns:     make(map[string]float64),
		Flashes:       make(map[string]float64),
		ActiveEffects: make([]*SpellEffectVFX, 0),
	}
}

func (sm *SpellManager) TriggerFlash(id string) {
	sm.Flashes[id] = 0.35
}

func (sm *SpellManager) SelectSpell(id string) {
	def := sm.Registry.GetSpell(id)
	if def == nil || !def.Unlocked {
		sm.TriggerFlash(id)
		return
	}
	if sm.Cooldowns[id] > 0 {
		sm.TriggerFlash(id)
		return
	}

	if sm.HasSelection && sm.SelectedSpell == id {
		sm.HasSelection = false
	} else {
		sm.SelectedSpell = id
		sm.HasSelection = true
	}
}

func (sm *SpellManager) Deselect() {
	sm.HasSelection = false
}

func (sm *SpellManager) Update(dt float64, run *economy.RunState) {
	// Update Cooldowns
	for id, cd := range sm.Cooldowns {
		if cd > 0 {
			sm.Cooldowns[id] = cd - dt
		}
	}

	// Update Flash Timers
	for id, fl := range sm.Flashes {
		if fl > 0 {
			sm.Flashes[id] = fl - dt
		}
	}

	// Overclock effect
	if sm.OverclockTimer > 0 {
		sm.OverclockTimer -= dt
		run.TowerSpeedMult = 1.70
	} else {
		run.TowerSpeedMult = 1.0
	}

	// Update VFX
	alive := sm.ActiveEffects[:0]
	for _, ef := range sm.ActiveEffects {
		ef.Life -= dt
		if ef.Life > 0 {
			alive = append(alive, ef)
		}
	}
	sm.ActiveEffects = alive
}

func (sm *SpellManager) CastSpell(id string, targetX, targetY float64, enemies []*malware.Enemy, run *economy.RunState, spawner *malware.Spawner) bool {
	def := sm.Registry.GetSpell(id)
	if def == nil || !def.Unlocked {
		sm.TriggerFlash(id)
		return false
	}
	if sm.Cooldowns[id] > 0 {
		sm.TriggerFlash(id)
		return false
	}
	if !run.CanAffordMana(def.ManaCost) {
		sm.TriggerFlash(id)
		return false
	}

	run.SpendMana(def.ManaCost)
	sm.Cooldowns[id] = def.Cooldown * run.SpellCooldownMult

	// Data-driven execution of all modular payload effects
	for _, eff := range def.Effects {
		switch eff.Type {
		case data.SpellEffectChainLightning:
			var firstHit *malware.Enemy
			minDist := eff.Radius
			for _, e := range enemies {
				if e.IsDead || e.IsPhased {
					continue
				}
				dist := math.Hypot(e.X-targetX, e.Y-targetY)
				if dist <= minDist {
					minDist = dist
					firstHit = e
				}
			}

			if firstHit != nil {
				if firstHit.TakeDamage(eff.Damage, spawner) {
					run.AddKill(firstHit.Bounty, firstHit.XP, firstHit.Type == malware.TypeBoss)
				}
				sm.ActiveEffects = append(sm.ActiveEffects, &SpellEffectVFX{
					X: targetX, Y: targetY, TargetX: firstHit.X, TargetY: firstHit.Y,
					Life: 0.25, MaxLife: 0.25, IsThunder: true, Color: def.Color,
				})

				curr := firstHit
				chained := 0
				for _, other := range enemies {
					if other == curr || other.IsDead || other.IsPhased {
						continue
					}
					dist := math.Hypot(other.X-curr.X, other.Y-curr.Y)
					if dist <= eff.ChainRange {
						if other.TakeDamage(eff.ChainDamage, spawner) {
							run.AddKill(other.Bounty, other.XP, other.Type == malware.TypeBoss)
						}
						sm.ActiveEffects = append(sm.ActiveEffects, &SpellEffectVFX{
							X: curr.X, Y: curr.Y, TargetX: other.X, TargetY: other.Y,
							Life: 0.2, MaxLife: 0.2, IsThunder: true, Color: def.Color,
						})
						curr = other
						chained++
						if chained >= eff.ChainCount {
							break
						}
					}
				}
			}

		case data.SpellEffectAreaFreeze:
			for _, e := range enemies {
				if e.IsDead {
					continue
				}
				dist := math.Hypot(e.X-targetX, e.Y-targetY)
				if dist <= eff.Radius {
					e.ApplyFreeze(eff.FreezeSeconds)
					if eff.Damage > 0 {
						if e.TakeDamage(eff.Damage, spawner) {
							run.AddKill(e.Bounty, e.XP, e.Type == malware.TypeBoss)
						}
					}
				}
			}
			sm.ActiveEffects = append(sm.ActiveEffects, &SpellEffectVFX{
				X: targetX, Y: targetY, Radius: eff.Radius,
				Life: 0.6, MaxLife: 0.6, Color: def.Color,
			})

		case data.SpellEffectAreaDamage:
			for _, e := range enemies {
				if e.IsDead {
					continue
				}
				dist := math.Hypot(e.X-targetX, e.Y-targetY)
				if dist <= eff.Radius {
					falloff := 1.0 - (dist / (eff.Radius * 1.25))
					if falloff < 0.4 {
						falloff = 0.4
					}
					if e.TakeDamage(eff.Damage*falloff, spawner) {
						run.AddKill(e.Bounty, e.XP, e.Type == malware.TypeBoss)
					}
				}
			}
			sm.ActiveEffects = append(sm.ActiveEffects, &SpellEffectVFX{
				X: targetX, Y: targetY, Radius: eff.Radius,
				Life: 0.5, MaxLife: 0.5, Color: def.Color,
			})

		case data.SpellEffectOverclock:
			sm.OverclockTimer = eff.BuffDuration
			sm.ActiveEffects = append(sm.ActiveEffects, &SpellEffectVFX{
				X: 400, Y: 300, Radius: 350,
				Life: 0.8, MaxLife: 0.8, Color: def.Color,
			})
		}
	}

	sm.HasSelection = false
	return true
}

func (sm *SpellManager) Draw(screen *ebiten.Image) {
	for _, ef := range sm.ActiveEffects {
		ratio := float32(ef.Life / ef.MaxLife)
		if ef.IsThunder {
			// Outer electric aura
			vector.StrokeLine(screen, float32(ef.X), float32(ef.Y), float32(ef.TargetX), float32(ef.TargetY), 4.5*ratio, ef.Color, false)
			// Core intense white plasma bolt
			vector.StrokeLine(screen, float32(ef.X), float32(ef.Y), float32(ef.TargetX), float32(ef.TargetY), 1.8*ratio, color.RGBA{255, 255, 255, 255}, false)
		} else {
			// Expanding shockwave ring
			r := float32(ef.Radius * (1.15 - float64(ratio)*0.15))
			vector.StrokeCircle(screen, float32(ef.X), float32(ef.Y), r, 3.5*ratio, ef.Color, false)
			vector.StrokeCircle(screen, float32(ef.X), float32(ef.Y), r*0.9, 1.5*ratio, gfx.Brighten(ef.Color, 0.4), false)
			vector.FillCircle(screen, float32(ef.X), float32(ef.Y), r*0.85*ratio, color.RGBA{
				R: ef.Color.R, G: ef.Color.G, B: ef.Color.B, A: uint8(120 * ratio),
			}, false)
		}
	}
}
