package data

import (
	"image/color"
)

// --- TOWERS & EFFECTS ---

type DeliveryType int

const (
	DeliveryDirectProjectile DeliveryType = iota
	DeliveryBallisticMortar
	DeliveryInstantBeam
	DeliveryRadialPulse
	DeliveryPassiveGenerator
)

type DeliveryConfig struct {
	Type         DeliveryType
	Speed        float64
	SplashRadius float64
	BeamDuration float64
}

type OnHitType int

const (
	OnHitDamage OnHitType = iota
	OnHitSplashDamage
	OnHitSlow
	OnHitFreeze
	OnHitChain
)

type OnHitEffect struct {
	Type          OnHitType
	Damage        float64
	SplashRadius  float64
	SlowFactor    float64
	SlowDuration  float64
	FreezeSeconds float64
	ChainCount    int
	ChainRange    float64
	ChainFalloff  float64
}

type TowerDef struct {
	ID          string
	Name        string
	Description string
	BaseCost    float64
	Range       float64
	Cooldown    float64
	ManaRate    float64 // For economy / Mana Siphons
	ByteRate    float64 // For economy / Passive Crypto Miners
	Color       color.RGBA
	Unlocked    bool
	Delivery    DeliveryConfig
	OnHit       []OnHitEffect
}

// --- SPELLS & EFFECTS ---

type SpellCastType int

const (
	CastPointArea SpellCastType = iota
	CastSingleTarget
	CastGlobalBuff
)

type SpellVisualType int

const (
	VisualLightning SpellVisualType = iota
	VisualFreezeRing
	VisualExplosion
	VisualGlobalPulse
)

type SpellEffectType int

const (
	SpellEffectAreaDamage SpellEffectType = iota
	SpellEffectAreaFreeze
	SpellEffectChainLightning
	SpellEffectOverclock
)

type SpellEffectDef struct {
	Type          SpellEffectType
	Damage        float64
	Radius        float64
	FreezeSeconds float64
	ChainCount    int
	ChainRange    float64
	ChainDamage   float64
	BuffDuration  float64
	BuffSpeedMult float64
}

type SpellDef struct {
	ID           string
	Name         string
	Description  string
	ManaCost     float64
	Cooldown     float64
	TargetRadius float64
	HotKey       string
	Color        color.RGBA
	Unlocked     bool
	CastType     SpellCastType
	Visual       SpellVisualType
	Effects      []SpellEffectDef
}

// --- ENEMIES ---

type EnemyArchetype int

const (
	ArchetypeNormal EnemyArchetype = iota
	ArchetypeShielded
	ArchetypeTankSplitter
	ArchetypeSplitterSub
	ArchetypeHealer
	ArchetypePhaseCloaker
	ArchetypeExploder
	ArchetypeBoss
)

type EnemyDef struct {
	ID           string
	Name         string
	Description  string
	Archetype    EnemyArchetype
	BaseHP       float64
	BaseShield   float64
	BaseSpeed    float64
	BaseBounty   float64
	BaseXP       float64
	CoreDamage   float64
	Radius       float32
	Color        color.RGBA
	IsBoss       bool
	MinSpawnTime float64
	SpawnWeight  float64
	// Archetype specific parameters
	HealAuraRadius float64
	HealAuraAmount float64
	HealAuraPeriod float64
	PhasePeriod    float64
	PhaseDuration  float64
	ExplodeRadius  float64
	ExplodeDamage  float64
	SplitSpawnID   string
	SplitCount     int
}

// --- DRAFT CARDS & UPGRADES ---

type ModifierType int

const (
	ModUnlockTower ModifierType = iota
	ModUnlockSpell
	ModTowerDamageMult
	ModTowerRangeMult
	ModTowerSpeedMult
	ModManaGenMult
	ModByteBountyMult
	ModKernelHeal
)

type CardEffect struct {
	Type     ModifierType
	TargetID string  // Tower ID or Spell ID to unlock
	Value    float64 // Multiplier or flat amount
}

type CardDef struct {
	ID          string
	Title       string
	Subtitle    string
	Description string
	Color       color.RGBA
	IsUnlock    bool
	Effects     []CardEffect
}
