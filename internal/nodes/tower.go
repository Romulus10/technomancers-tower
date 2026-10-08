package nodes

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"technomancers-tower/internal/data"
	"technomancers-tower/internal/economy"
	"technomancers-tower/internal/gfx"
	"technomancers-tower/internal/malware"
	"technomancers-tower/internal/motherboard"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const MaxTowerLevel = 5

// NextLevelXP returns the required XP to advance from the given level to the next.
func NextLevelXP(level int) float64 {
	switch level {
	case 1:
		return 100.0
	case 2:
		return 250.0
	case 3:
		return 550.0
	case 4:
		return 1100.0
	default:
		return 0.0 // Max Level
	}
}

// Tower represents an active defense structure built on the motherboard.
type Tower struct {
	ID       int
	Def      *data.TowerDef
	GridX    int
	GridY    int
	WorldX   float64
	WorldY   float64
	Timer    float64
	Tier     int // Overclock Tier (starts at 1, promoted at Level 5)
	Level    int // Level within current Tier (1..5)
	XP       float64
	TargetXP float64

	// Visual effect states
	BeamTargetX float64
	BeamTargetY float64
	BeamLife    float64
	BeamMaxLife float64

	// Level-up & Promotion visual effects
	LevelUpAnim     float64
	LevelUpTextLife float64
}

// GetDamageMult returns the cumulative damage multiplier from Tier and Level.
func (t *Tower) GetDamageMult() float64 {
	return (1.0 + float64(t.Tier-1)*0.35) * (1.0 + float64(t.Level-1)*0.10)
}

// GetSpeedMult returns the cumulative fire rate multiplier from Tier and Level.
func (t *Tower) GetSpeedMult() float64 {
	return (1.0 + float64(t.Tier-1)*0.20) * (1.0 + float64(t.Level-1)*0.05)
}

// GetRangeMult returns the cumulative range multiplier from Tier and Level.
func (t *Tower) GetRangeMult() float64 {
	return 1.0 + float64(t.Tier-1)*0.05 + float64(t.Level-1)*0.03
}

// GetManaRate returns the scaled mana harvest rate from Tier and Level.
func (t *Tower) GetManaRate() float64 {
	return t.Def.ManaRate * (1.0 + float64(t.Tier-1)*0.35) * (1.0 + float64(t.Level-1)*0.10)
}

// GetByteRate returns the scaled byte generation rate from Tier and Level.
func (t *Tower) GetByteRate() float64 {
	return t.Def.ByteRate * (1.0 + float64(t.Tier-1)*0.35) * (1.0 + float64(t.Level-1)*0.10)
}

// GetUpgradeCost calculates the Byte cost to promote the tower to the next Overclock Tier.
func (t *Tower) GetUpgradeCost() float64 {
	return math.Round(t.Def.BaseCost * 0.50 * float64(t.Tier))
}

// CanPromote returns true if the tower is at Max Level and the player can afford the upgrade cost.
func (t *Tower) CanPromote(run *economy.RunState) bool {
	return t.Level >= MaxTowerLevel && run.CanAffordBytes(t.GetUpgradeCost())
}

// Promote consumes the Byte cost, advances to the next Tier, and resets the Level to 1.
func (t *Tower) Promote(run *economy.RunState) bool {
	if t.Level < MaxTowerLevel {
		return false
	}
	cost := t.GetUpgradeCost()
	if !run.CanAffordBytes(cost) {
		return false
	}

	run.SpendBytes(cost)
	t.Tier++
	t.Level = 1
	t.XP = 0
	t.TargetXP = NextLevelXP(1)
	t.LevelUpAnim = 0.8
	t.LevelUpTextLife = 1.5
	return true
}

// AddXP adds combat or harvest experience to the tower, triggering level-ups up to MaxTowerLevel.
func (t *Tower) AddXP(amount float64) bool {
	if t.Level >= MaxTowerLevel || amount <= 0 {
		return false
	}

	t.XP += amount
	leveledUp := false

	for t.Level < MaxTowerLevel && t.XP >= t.TargetXP && t.TargetXP > 0 {
		t.XP -= t.TargetXP
		t.Level++
		t.TargetXP = NextLevelXP(t.Level)
		t.LevelUpAnim = 0.6
		t.LevelUpTextLife = 1.2
		leveledUp = true
	}

	if t.Level >= MaxTowerLevel {
		t.XP = 0
		t.TargetXP = 0
	}

	return leveledUp
}

// Projectile handles both direct homing bolts and ballistic artillery mortar shells.
type Projectile struct {
	SourceTowerID  int
	StartX, StartY float64
	X, Y           float64
	TargetX        float64
	TargetY        float64
	TargetEnemy    *malware.Enemy
	Speed          float64
	IsMortar       bool
	SplashRadius   float64
	OnHits         []data.OnHitEffect
	DamageMult     float64
	IsDead         bool
}

// TowerManager handles tower logic, attacks, upgrades, and projectile lifecycles data-driven.
type TowerManager struct {
	nextID      int
	Towers      []*Tower
	Projectiles []*Projectile
	Registry    *data.Registry
}

func NewTowerManager(reg *data.Registry) *TowerManager {
	return &TowerManager{
		Towers:      make([]*Tower, 0),
		Projectiles: make([]*Projectile, 0),
		Registry:    reg,
	}
}

func (tm *TowerManager) BuildTower(towerID string, gx, gy int, grid *motherboard.Grid, run *economy.RunState) bool {
	def := tm.Registry.GetTower(towerID)
	if def == nil || !def.Unlocked {
		return false
	}
	if !run.CanAffordBytes(def.BaseCost) {
		return false
	}
	if !grid.SetTower(gx, gy) {
		return false
	}

	run.SpendBytes(def.BaseCost)
	wx, wy := grid.GridToScreenCenter(gx, gy)
	tm.nextID++

	tower := &Tower{
		ID:       tm.nextID,
		Def:      def,
		GridX:    gx,
		GridY:    gy,
		WorldX:   wx,
		WorldY:   wy,
		Tier:     1,
		Level:    1,
		XP:       0,
		TargetXP: NextLevelXP(1),
	}

	tm.Towers = append(tm.Towers, tower)
	tm.RecalculateBaseManaRate(run)
	return true
}

func (tm *TowerManager) PromoteTowerAt(gx, gy int, run *economy.RunState) bool {
	t := tm.GetTowerAt(gx, gy)
	if t == nil {
		return false
	}
	promoted := t.Promote(run)
	if promoted {
		tm.RecalculateBaseManaRate(run)
	}
	return promoted
}

// GetPromotableInfo returns the number of max-level towers awaiting promotion and their total Byte cost.
func (tm *TowerManager) GetPromotableInfo() (count int, totalCost float64) {
	for _, t := range tm.Towers {
		if t.Level >= MaxTowerLevel {
			count++
			totalCost += t.GetUpgradeCost()
		}
	}
	return count, totalCost
}

// PromoteAll upgrades all eligible max-level towers that the player can afford in one action.
func (tm *TowerManager) PromoteAll(run *economy.RunState) int {
	promotedCount := 0
	for _, t := range tm.Towers {
		if t.Level >= MaxTowerLevel && run.CanAffordBytes(t.GetUpgradeCost()) {
			if t.Promote(run) {
				promotedCount++
			}
		}
	}
	if promotedCount > 0 {
		tm.RecalculateBaseManaRate(run)
	}
	return promotedCount
}

// GetTotalOverclockTiers calculates the cumulative overclock tiers across all placed towers (sum of Tier - 1).
func (tm *TowerManager) GetTotalOverclockTiers() int {
	total := 0
	for _, t := range tm.Towers {
		if t.Tier > 1 {
			total += (t.Tier - 1)
		}
	}
	return total
}

func (tm *TowerManager) AddTowerXP(sourceTowerID int, amount float64, run *economy.RunState) {
	if sourceTowerID <= 0 || amount <= 0 {
		return
	}
	for _, t := range tm.Towers {
		if t.ID == sourceTowerID {
			if t.AddXP(amount) {
				tm.RecalculateBaseManaRate(run)
			}
			break
		}
	}
}

func (tm *TowerManager) GetTowerAt(gx, gy int) *Tower {
	for _, t := range tm.Towers {
		if t.GridX == gx && t.GridY == gy {
			return t
		}
	}
	return nil
}

func (tm *TowerManager) RecalculateBaseManaRate(run *economy.RunState) {
	totalManaRate := 0.0
	totalByteRate := 0.0
	for _, t := range tm.Towers {
		totalManaRate += t.GetManaRate()
		totalByteRate += t.GetByteRate()
	}
	run.BaseManaRate = totalManaRate
	run.BaseByteRate = totalByteRate
}

func (tm *TowerManager) Update(dt float64, enemies []*malware.Enemy, run *economy.RunState, spawner *malware.Spawner) {
	// 1. Update Towers
	for _, t := range tm.Towers {
		if t.BeamLife > 0 {
			t.BeamLife -= dt
		}
		if t.LevelUpAnim > 0 {
			t.LevelUpAnim -= dt
		}
		if t.LevelUpTextLife > 0 {
			t.LevelUpTextLife -= dt
		}

		if t.Def.Delivery.Type == data.DeliveryPassiveGenerator {
			// Passive generators gain 15 XP per second actively operating
			if t.AddXP(15.0 * dt) {
				tm.RecalculateBaseManaRate(run)
			}
			continue
		}

		t.Timer -= dt
		if t.Timer <= 0 {
			effectiveRange := t.Def.Range * run.TowerRangeMult * t.GetRangeMult()
			effectiveCooldown := (t.Def.Cooldown / run.TowerSpeedMult) / t.GetSpeedMult()

			target := tm.findBestTarget(t.WorldX, t.WorldY, effectiveRange, enemies)
			if target != nil {
				t.Timer = effectiveCooldown
				totalDamageMult := run.TowerDamageMult * t.GetDamageMult()

				// Data-driven attack delivery execution
				switch t.Def.Delivery.Type {
				case data.DeliveryDirectProjectile:
					tm.Projectiles = append(tm.Projectiles, &Projectile{
						SourceTowerID: t.ID,
						StartX:        t.WorldX,
						StartY:        t.WorldY,
						X:             t.WorldX,
						Y:             t.WorldY,
						TargetEnemy:   target,
						Speed:         t.Def.Delivery.Speed,
						OnHits:        t.Def.OnHit,
						DamageMult:    totalDamageMult,
					})

				case data.DeliveryBallisticMortar:
					tm.Projectiles = append(tm.Projectiles, &Projectile{
						SourceTowerID: t.ID,
						StartX:        t.WorldX,
						StartY:        t.WorldY,
						X:             t.WorldX,
						Y:             t.WorldY,
						TargetX:       target.X,
						TargetY:       target.Y,
						Speed:         t.Def.Delivery.Speed,
						SplashRadius:  t.Def.Delivery.SplashRadius,
						IsMortar:      true,
						OnHits:        t.Def.OnHit,
						DamageMult:    totalDamageMult,
					})

				case data.DeliveryInstantBeam:
					t.BeamTargetX = target.X
					t.BeamTargetY = target.Y
					t.BeamLife = t.Def.Delivery.BeamDuration
					t.BeamMaxLife = t.Def.Delivery.BeamDuration
					tm.ApplyOnHits(t.ID, t.Def.OnHit, target, enemies, run, totalDamageMult, target.X, target.Y, spawner)

				case data.DeliveryRadialPulse:
					t.BeamLife = t.Def.Delivery.BeamDuration
					t.BeamMaxLife = t.Def.Delivery.BeamDuration
					// Pulse all enemies within effective range
					for _, e := range enemies {
						if e.IsDead {
							continue
						}
						dist := math.Hypot(e.X-t.WorldX, e.Y-t.WorldY)
						if dist <= effectiveRange {
							tm.ApplyOnHits(t.ID, t.Def.OnHit, e, enemies, run, totalDamageMult, t.WorldX, t.WorldY, spawner)
						}
					}
				}
			}
		}
	}

	// 2. Update Projectiles
	aliveProj := tm.Projectiles[:0]
	for _, p := range tm.Projectiles {
		p.Update(dt, enemies, run, tm, spawner)
		if !p.IsDead {
			aliveProj = append(aliveProj, p)
		}
	}
	tm.Projectiles = aliveProj
}

// ApplyOnHits executes all configured OnHitEffects (Damage, Splash, Slow, Freeze, Chain) and attributes XP.
func (tm *TowerManager) ApplyOnHits(sourceTowerID int, effects []data.OnHitEffect, target *malware.Enemy, enemies []*malware.Enemy, run *economy.RunState, dmgMult float64, originX, originY float64, spawner *malware.Spawner) {
	effectiveDmgMult := dmgMult
	if run != nil && run.TowerCritChance > 0 && rand.Float64() < run.TowerCritChance {
		effectiveDmgMult *= run.TowerCritMult
	}

	for _, eff := range effects {
		switch eff.Type {
		case data.OnHitDamage:
			if target != nil && !target.IsDead {
				dmg := eff.Damage * effectiveDmgMult
				if target.TakeDamage(dmg, spawner) {
					run.AddKill(target.Bounty, target.XP, target.Type == malware.TypeBoss)
					tm.AddTowerXP(sourceTowerID, dmg+float64(target.XP)*5.0, run)
				} else {
					tm.AddTowerXP(sourceTowerID, dmg, run)
				}
			}

		case data.OnHitSplashDamage:
			for _, e := range enemies {
				if e.IsDead {
					continue
				}
				dist := math.Hypot(e.X-originX, e.Y-originY)
				if dist <= eff.SplashRadius {
					falloff := 1.0 - (dist / (eff.SplashRadius * 1.25))
					if falloff < 0.35 {
						falloff = 0.35
					}
					dmg := eff.Damage * effectiveDmgMult * falloff
					if e.TakeDamage(dmg, spawner) {
						run.AddKill(e.Bounty, e.XP, e.Type == malware.TypeBoss)
						tm.AddTowerXP(sourceTowerID, dmg+float64(e.XP)*5.0, run)
					} else {
						tm.AddTowerXP(sourceTowerID, dmg, run)
					}
				}
			}

		case data.OnHitSlow:
			if target != nil && !target.IsDead {
				slowDur := eff.SlowDuration
				if run != nil && run.StatusDurationMult > 0 {
					slowDur *= run.StatusDurationMult
				}
				target.ApplySlow(eff.SlowFactor, slowDur)
			}

		case data.OnHitFreeze:
			if target != nil && !target.IsDead {
				freezeSec := eff.FreezeSeconds
				if run != nil && run.StatusDurationMult > 0 {
					freezeSec *= run.StatusDurationMult
				}
				target.ApplyFreeze(freezeSec)
			}

		case data.OnHitChain:
			if target == nil {
				continue
			}
			curr := target
			chained := 0
			for _, other := range enemies {
				if other == curr || other.IsDead || other.IsPhased {
					continue
				}
				dist := math.Hypot(other.X-curr.X, other.Y-curr.Y)
				if dist <= eff.ChainRange {
					falloff := eff.ChainFalloff
					if falloff <= 0 {
						falloff = 0.75
					}
					dmg := eff.Damage * effectiveDmgMult * falloff
					if other.TakeDamage(dmg, spawner) {
						run.AddKill(other.Bounty, other.XP, other.Type == malware.TypeBoss)
						tm.AddTowerXP(sourceTowerID, dmg+float64(other.XP)*5.0, run)
					} else {
						tm.AddTowerXP(sourceTowerID, dmg, run)
					}
					curr = other
					chained++
					if chained >= eff.ChainCount {
						break
					}
				}
			}
		}
	}
}

func (tm *TowerManager) findBestTarget(tx, ty, rng float64, enemies []*malware.Enemy) *malware.Enemy {
	var best *malware.Enemy
	minDist := rng

	for _, e := range enemies {
		if e.IsDead || e.IsPhased {
			continue
		}
		dist := math.Hypot(e.X-tx, e.Y-ty)
		if dist <= minDist {
			minDist = dist
			best = e
		}
	}
	return best
}

func (p *Projectile) Update(dt float64, enemies []*malware.Enemy, run *economy.RunState, tm *TowerManager, spawner *malware.Spawner) {
	if p.IsMortar {
		dx := p.TargetX - p.X
		dy := p.TargetY - p.Y
		dist := math.Hypot(dx, dy)
		step := p.Speed * dt

		if dist <= step || dist < 4 {
			p.IsDead = true
			tm.ApplyOnHits(p.SourceTowerID, p.OnHits, nil, enemies, run, p.DamageMult, p.TargetX, p.TargetY, spawner)
			return
		}

		p.X += (dx / dist) * step
		p.Y += (dy / dist) * step
		return
	}

	// Direct homing projectile
	if p.TargetEnemy == nil || p.TargetEnemy.IsDead {
		p.IsDead = true
		return
	}

	dx := p.TargetEnemy.X - p.X
	dy := p.TargetEnemy.Y - p.Y
	dist := math.Hypot(dx, dy)
	step := p.Speed * dt

	if dist <= step || dist < 6 {
		p.IsDead = true
		tm.ApplyOnHits(p.SourceTowerID, p.OnHits, p.TargetEnemy, enemies, run, p.DamageMult, p.TargetEnemy.X, p.TargetEnemy.Y, spawner)
		return
	}

	p.X += (dx / dist) * step
	p.Y += (dy / dist) * step
}

func (tm *TowerManager) Draw(screen *ebiten.Image) {
	cache := gfx.GetCache()

	for _, t := range tm.Towers {
		tx := float32(t.WorldX)
		ty := float32(t.WorldY)

		// 1. Draw Cached Procedural Sprite
		sprite := cache.GetTowerSprite(t.Def.ID, t.Tier, t.Def.Color)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(tx-16), float64(ty-16))
		screen.DrawImage(sprite, op)

		// 2. Overclock Ready / Max Rank Pulsing Aura
		if t.Level >= MaxTowerLevel {
			vector.StrokeRect(screen, tx-16, ty-16, 32, 32, 2.0, gfx.ColorGoldMatrix, false)
		}

		// 3. Rank Pips (Level 1..5)
		pipW := float32(3.5)
		pipGap := float32(1.5)
		totalPipsW := float32(MaxTowerLevel)*pipW + float32(MaxTowerLevel-1)*pipGap
		pipStartX := tx - totalPipsW/2
		pipY := ty - 12
		for lvl := 1; lvl <= MaxTowerLevel; lvl++ {
			px := pipStartX + float32(lvl-1)*(pipW+pipGap)
			if lvl <= t.Level {
				vector.FillRect(screen, px, pipY, pipW, pipW, gfx.ColorGoldMatrix, false)
			} else {
				vector.FillRect(screen, px, pipY, pipW, pipW, color.RGBA{R: 30, G: 40, B: 55, A: 160}, false)
			}
		}

		// 4. Floating Upgrade badge if tower reached Max Level
		if t.Level >= MaxTowerLevel {
			upgCost := t.GetUpgradeCost()
			badgeText := fmt.Sprintf("^%sB", economy.FormatNumber(upgCost))
			vector.FillRect(screen, tx-18, ty-24, 36, 11, color.RGBA{R: 15, G: 30, B: 20, A: 230}, false)
			vector.StrokeRect(screen, tx-18, ty-24, 36, 11, 1.0, gfx.ColorEmeraldBio, false)
			ebitenutil.DebugPrintAt(screen, badgeText, int(tx)-16, int(ty)-25)
		}

		// 5. Radial Pulse visual effect
		if t.Def.Delivery.Type == data.DeliveryRadialPulse && t.BeamLife > 0 {
			ratio := float32(t.BeamLife / t.BeamMaxLife)
			alpha := uint8(200 * ratio)
			vector.StrokeCircle(screen, tx, ty, float32(t.Def.Range*t.GetRangeMult()), 2.5*ratio, color.RGBA{R: 120, G: 220, B: 255, A: alpha}, false)
			vector.FillCircle(screen, tx, ty, float32(t.Def.Range*t.GetRangeMult()*0.9), color.RGBA{R: 80, G: 200, B: 255, A: uint8(40 * ratio)}, false)
		}

		// 6. Instant Beam visual effect
		if t.Def.Delivery.Type == data.DeliveryInstantBeam && t.BeamLife > 0 {
			ratio := float32(t.BeamLife / t.BeamMaxLife)
			vector.StrokeLine(screen, tx, ty, float32(t.BeamTargetX), float32(t.BeamTargetY), 3.0*ratio, gfx.Brighten(t.Def.Color, 0.4), false)
			vector.StrokeLine(screen, tx, ty, float32(t.BeamTargetX), float32(t.BeamTargetY), 1.2*ratio, color.RGBA{255, 255, 255, 255}, false)
		}

		// 7. Level up expanding ring animation
		if t.LevelUpAnim > 0 {
			progress := float32(1.0 - (t.LevelUpAnim / 0.6))
			ringRadius := 12.0 + progress*30.0
			alpha := uint8(255 * (t.LevelUpAnim / 0.6))
			vector.StrokeCircle(screen, tx, ty, ringRadius, 2.5, color.RGBA{R: 255, G: 230, B: 80, A: alpha}, false)
		}

		// 8. Floating "+LV.X!" or "+TIER X!" text
		if t.LevelUpTextLife > 0 {
			offsetY := int(ty - 18 - float32((1.2-t.LevelUpTextLife)*20.0))
			if t.Tier > 1 && t.Level == 1 {
				ebitenutil.DebugPrintAt(screen, fmt.Sprintf("+TIER %d!", t.Tier), int(tx)-22, offsetY)
			} else {
				ebitenutil.DebugPrintAt(screen, fmt.Sprintf("+LV.%d!", t.Level), int(tx)-14, offsetY)
			}
		}
	}

	// Draw Projectiles with glowing layered cores
	for _, p := range tm.Projectiles {
		px := float32(p.X)
		py := float32(p.Y)
		if p.IsMortar {
			// Ballistic mortar shell with fiery tail
			vector.FillCircle(screen, px, py, 5.5, gfx.ColorAmberFire, false)
			vector.FillCircle(screen, px, py, 3.0, gfx.ColorGoldMatrix, false)
			vector.FillCircle(screen, px, py, 1.5, color.RGBA{255, 255, 255, 255}, false)
		} else {
			// High-velocity energy bolt
			vector.FillCircle(screen, px, py, 4.0, gfx.ColorCyanGlow, false)
			vector.FillCircle(screen, px, py, 2.5, gfx.ColorCyanNeon, false)
			vector.FillCircle(screen, px, py, 1.2, color.RGBA{255, 255, 255, 255}, false)
		}
	}
}
