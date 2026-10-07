package economy

import (
	"fmt"
	"math"
)

// RunState tracks in-run resources, progression, and combat stats.
type RunState struct {
	// Currencies
	Bytes        float64
	TotalBytes   float64
	Mana         float64
	MaxMana      float64
	BaseManaRate float64 // Mana per second from towers
	BaseByteRate float64 // Bytes per second from passive miners

	// Base/Kernel Health
	KernelHP    float64
	MaxKernelHP float64

	// XP & Roguelite Level
	Level         int
	CurrentXP     float64
	TargetXP      float64
	PendingDrafts int

	// Survival & Combat Stats
	RunTime      float64 // in seconds
	Kills        int
	BossKills    int
	ShardsEarned int

	// Global Modifiers (Applied via Meta Talents & In-run Drafts)
	TowerDamageMult   float64
	TowerRangeMult    float64
	TowerSpeedMult    float64
	ByteBountyMult    float64
	ManaGenMult       float64
	SpellCooldownMult float64
	SpellCostMult     float64
}

func NewRunState(metaShieldLvl, metaBootLvl, metaManaLvl, metaSpeedLvl, metaSpellLvl, metaScrapLvl, metaDamageLvl, metaRangeLvl, metaNaniteLvl int) *RunState {
	maxHP := 100.0 + float64(metaShieldLvl)*25.0
	startingBytes := 120.0 + float64(metaBootLvl)*40.0
	manaGenMult := 1.0 + float64(metaManaLvl)*0.20
	speedMult := 1.0 + float64(metaSpeedLvl)*0.05
	dmgMult := 1.0 + float64(metaDamageLvl)*0.06
	rangeMult := 1.0 + float64(metaRangeLvl)*0.06
	spellCdMult := math.Max(0.5, 1.0-float64(metaSpellLvl)*0.08)
	spellCostMult := math.Max(0.5, 1.0-float64(metaSpellLvl)*0.10)
	scrapMult := 1.0 + float64(metaScrapLvl)*0.10

	return &RunState{
		Bytes:             startingBytes,
		TotalBytes:        startingBytes,
		Mana:              20.0,
		MaxMana:           100.0,
		BaseManaRate:      0.0,
		BaseByteRate:      0.0,
		KernelHP:          maxHP,
		MaxKernelHP:       maxHP,
		Level:             1,
		CurrentXP:         0,
		TargetXP:          40.0,
		PendingDrafts:     0,
		TowerDamageMult:   dmgMult,
		TowerRangeMult:    rangeMult,
		TowerSpeedMult:    speedMult,
		ByteBountyMult:    scrapMult,
		ManaGenMult:       manaGenMult,
		SpellCooldownMult: spellCdMult,
		SpellCostMult:     spellCostMult,
	}
}

func (r *RunState) Update(dt float64) {
	r.RunTime += dt

	// Passive mana generation from towers
	totalManaRate := r.BaseManaRate * r.ManaGenMult
	if totalManaRate > 0 {
		r.Mana += totalManaRate * dt
		if r.Mana > r.MaxMana {
			r.Mana = r.MaxMana
		}
	}

	// Passive byte generation from crypto miners
	totalByteRate := r.BaseByteRate * r.ByteBountyMult
	if totalByteRate > 0 {
		gained := totalByteRate * dt
		r.Bytes += gained
		r.TotalBytes += gained
	}
}

// AddXP awards runtime XP and queues draft choices upon leveling up.
func (r *RunState) AddXP(amount float64) bool {
	if amount <= 0 {
		return false
	}
	r.CurrentXP += amount
	leveledUp := false
	for r.CurrentXP >= r.TargetXP && r.TargetXP > 0 {
		r.CurrentXP -= r.TargetXP
		r.Level++
		r.TargetXP = math.Round(r.TargetXP * 1.35)
		r.PendingDrafts++
		leveledUp = true
	}
	return leveledUp
}

// AddKill awards bytes and registers a kill.
func (r *RunState) AddKill(bounty float64, xp float64, isBoss bool) bool {
	actualBounty := bounty * r.ByteBountyMult
	r.Bytes += actualBounty
	r.TotalBytes += actualBounty
	r.Kills++
	if isBoss {
		r.BossKills++
	}
	return r.AddXP(xp)
}

// DamageKernel inflicts damage on the CPU core. Returns true if core is destroyed.
func (r *RunState) DamageKernel(amount float64) bool {
	r.KernelHP -= amount
	if r.KernelHP <= 0 {
		r.KernelHP = 0
		r.CalculateFinalShards()
		return true
	}
	return false
}

func (r *RunState) CalculateFinalShards() int {
	// Base shards from survival time (1 shard per 10s survived)
	timeShards := int(r.RunTime / 10.0)
	// Shards from kills (1 shard per 20 kills)
	killShards := r.Kills / 20
	// Shards from bosses (15 shards per boss kill)
	bossShards := r.BossKills * 15

	total := timeShards + killShards + bossShards
	if total < 1 {
		total = 1
	}
	r.ShardsEarned = total
	return total
}

func (r *RunState) CanAffordBytes(amount float64) bool {
	return r.Bytes >= amount
}

func (r *RunState) SpendBytes(amount float64) bool {
	if r.Bytes >= amount {
		r.Bytes -= amount
		return true
	}
	return false
}

func (r *RunState) CanAffordMana(amount float64) bool {
	return r.Mana >= (amount * r.SpellCostMult)
}

func (r *RunState) SpendMana(amount float64) bool {
	cost := amount * r.SpellCostMult
	if r.Mana >= cost {
		r.Mana -= cost
		return true
	}
	return false
}

func FormatNumber(n float64) string {
	if n < 1000 {
		if n == math.Trunc(n) {
			return fmt.Sprintf("%.0f", n)
		}
		return fmt.Sprintf("%.1f", n)
	}
	units := []string{"k", "M", "B", "T"}
	unitIndex := -1
	val := n
	for val >= 1000 && unitIndex < len(units)-1 {
		val /= 1000
		unitIndex++
	}
	return fmt.Sprintf("%.1f%s", val, units[unitIndex])
}
