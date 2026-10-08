package economy

import (
	"fmt"
	"math"
)

type TalentGetter interface {
	GetTalentLevel(id string) int
}

type MapTalents map[string]int

func (m MapTalents) GetTalentLevel(id string) int {
	return m[id]
}

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
	TowerDamageMult    float64
	TowerRangeMult     float64
	TowerSpeedMult     float64
	BaseTowerSpeedMult float64
	ByteBountyMult     float64
	ManaGenMult        float64
	SpellCooldownMult  float64
	SpellCostMult      float64

	// Expanded Meta-Progression Modifiers
	TowerCritChance    float64
	TowerCritMult      float64
	StatusDurationMult float64
	NaturalManaRegen   float64
	SpellCritChance    float64
	SpellCritMult      float64
	SpellRadiusMult    float64
	ByteInterestRate   float64
	ByteInterestTimer  float64
	SalvageRefundRate  float64
	MinerBoostMult     float64
	ShardsMult         float64
	DraftRerolls       int
	DraftCardChoices   int
	XPMult             float64
	NaniteRegenRate    float64
	NaniteRegenTimer   float64
	NaniteBossHeal     float64
	HardenedArmor      float64
	KernelEMPDamage    float64
}

func NewRunState(meta TalentGetter) *RunState {
	get := func(id string) int {
		if meta == nil {
			return 0
		}
		return meta.GetTalentLevel(id)
	}

	maxHP := 100.0 + float64(get("kernel_shield"))*25.0
	startingBytes := 120.0 + float64(get("boot_bytes"))*40.0
	maxMana := 100.0 + float64(get("expanded_mana_pool"))*25.0
	manaGenMult := 1.0 + float64(get("mana_conductor"))*0.20
	naturalManaRegen := float64(get("natural_mana_regen")) * 0.50
	speedMult := 1.0 + float64(get("overclock_nodes"))*0.05
	dmgMult := 1.0 + float64(get("tower_potency"))*0.06
	rangeMult := 1.0 + float64(get("sensor_array"))*0.06
	towerCritChance := float64(get("critical_subroutines")) * 0.04
	statusDurationMult := 1.0 + float64(get("status_amplification"))*0.12
	spellCdMult := math.Max(0.5, 1.0-float64(get("spell_efficiency"))*0.08)
	spellCostMult := math.Max(0.5, 1.0-float64(get("spell_efficiency"))*0.10)
	spellCritChance := float64(get("spell_critical")) * 0.05
	spellRadiusMult := 1.0 + float64(get("spell_critical"))*0.08
	scrapMult := 1.0 + float64(get("scrap_leech"))*0.10
	byteInterestRate := float64(get("byte_interest")) * 0.04
	salvageRefundRate := 0.50 + float64(get("salvage_efficiency"))*0.10
	minerBoostMult := 1.0 + float64(get("miner_boost"))*0.15
	shardsMult := 1.0 + float64(get("shards_harvest"))*0.15
	draftRerolls := get("draft_rerolls")
	draftCardChoices := 3 + get("extra_draft_slot")
	startLevelTalent := get("starting_level")
	xpMult := 1.0 + float64(get("xp_multiplier"))*0.08
	naniteRegenRate := float64(get("nanite_regen")) * 1.0
	naniteBossHeal := float64(get("nanite_resilience")) * 15.0
	hardenedArmor := float64(get("hardened_firewall")) * 1.0
	kernelEMPDamage := float64(get("kernel_emp")) * 30.0

	initLevel := 1 + startLevelTalent
	initTargetXP := 40.0
	for i := 1; i < initLevel; i++ {
		initTargetXP = math.Round(initTargetXP * 1.35)
	}

	return &RunState{
		Bytes:              startingBytes,
		TotalBytes:         startingBytes,
		Mana:               20.0,
		MaxMana:            maxMana,
		BaseManaRate:       0.0,
		BaseByteRate:       0.0,
		KernelHP:           maxHP,
		MaxKernelHP:        maxHP,
		Level:              initLevel,
		CurrentXP:          0,
		TargetXP:           initTargetXP,
		PendingDrafts:      startLevelTalent,
		TowerDamageMult:    dmgMult,
		TowerRangeMult:     rangeMult,
		TowerSpeedMult:     speedMult,
		BaseTowerSpeedMult: speedMult,
		ByteBountyMult:     scrapMult,
		ManaGenMult:        manaGenMult,
		SpellCooldownMult:  spellCdMult,
		SpellCostMult:      spellCostMult,
		TowerCritChance:    towerCritChance,
		TowerCritMult:      2.0,
		StatusDurationMult: statusDurationMult,
		NaturalManaRegen:   naturalManaRegen,
		SpellCritChance:    spellCritChance,
		SpellCritMult:      2.0,
		SpellRadiusMult:    spellRadiusMult,
		ByteInterestRate:   byteInterestRate,
		ByteInterestTimer:  30.0,
		SalvageRefundRate:  salvageRefundRate,
		MinerBoostMult:     minerBoostMult,
		ShardsMult:         shardsMult,
		DraftRerolls:       draftRerolls,
		DraftCardChoices:   draftCardChoices,
		XPMult:             xpMult,
		NaniteRegenRate:    naniteRegenRate,
		NaniteRegenTimer:   15.0,
		NaniteBossHeal:     naniteBossHeal,
		HardenedArmor:      hardenedArmor,
		KernelEMPDamage:    kernelEMPDamage,
	}
}

func (r *RunState) Update(dt float64) {
	r.RunTime += dt

	// Passive mana generation from towers + natural ambient ether siphon
	totalManaRate := (r.BaseManaRate * r.ManaGenMult) + r.NaturalManaRegen
	if totalManaRate > 0 {
		r.Mana += totalManaRate * dt
		if r.Mana > r.MaxMana {
			r.Mana = r.MaxMana
		}
	}

	// Passive byte generation from crypto miners with miner boost
	totalByteRate := r.BaseByteRate * r.ByteBountyMult * r.MinerBoostMult
	if totalByteRate > 0 {
		gained := totalByteRate * dt
		r.Bytes += gained
		r.TotalBytes += gained
	}

	// Algorithmic Interest tick
	if r.ByteInterestRate > 0 {
		r.ByteInterestTimer -= dt
		if r.ByteInterestTimer <= 0 {
			r.ByteInterestTimer = 30.0
			if r.Bytes > 0 {
				interest := math.Round(r.Bytes * r.ByteInterestRate)
				r.Bytes += interest
				r.TotalBytes += interest
			}
		}
	}

	// Nanite Regeneration tick
	if r.NaniteRegenRate > 0 && r.KernelHP < r.MaxKernelHP {
		r.NaniteRegenTimer -= dt
		if r.NaniteRegenTimer <= 0 {
			r.NaniteRegenTimer = 15.0
			r.KernelHP = math.Min(r.MaxKernelHP, r.KernelHP+r.NaniteRegenRate)
		}
	}
}

// AddXP awards runtime XP and queues draft choices upon leveling up.
func (r *RunState) AddXP(amount float64) bool {
	if amount <= 0 {
		return false
	}
	actualXP := amount * r.XPMult
	r.CurrentXP += actualXP
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
		if r.NaniteBossHeal > 0 {
			r.KernelHP = math.Min(r.MaxKernelHP, r.KernelHP+r.NaniteBossHeal)
		}
	}
	return r.AddXP(xp)
}

// DamageKernel inflicts damage on the CPU core with Hardened Firewall reduction. Returns true if core is destroyed.
func (r *RunState) DamageKernel(amount float64) bool {
	effectiveDmg := amount - r.HardenedArmor
	if effectiveDmg < 1.0 {
		effectiveDmg = 1.0
	}
	r.KernelHP -= effectiveDmg
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

	baseTotal := float64(timeShards + killShards + bossShards)
	total := int(math.Round(baseTotal * r.ShardsMult))
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
