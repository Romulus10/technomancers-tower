package data

import (
	"image/color"
)

type Registry struct {
	Towers  []*TowerDef
	Spells  []*SpellDef
	Enemies []*EnemyDef
	Cards   []*CardDef

	towerMap map[string]*TowerDef
	spellMap map[string]*SpellDef
	enemyMap map[string]*EnemyDef
	cardMap  map[string]*CardDef
}

func NewRegistry() *Registry {
	r := &Registry{
		towerMap: make(map[string]*TowerDef),
		spellMap: make(map[string]*SpellDef),
		enemyMap: make(map[string]*EnemyDef),
		cardMap:  make(map[string]*CardDef),
	}

	// 1. EXPANDED TOWERS REGISTRY (10 Unique Defense Nodes)
	r.Towers = []*TowerDef{
		{
			ID:          "bit_driver",
			Name:        "Bit-Driver",
			Description: "Rapid single-target packet cannon with high firing rate.",
			BaseCost:    30,
			Range:       100,
			Cooldown:    0.35,
			Color:       color.RGBA{R: 0, G: 220, B: 255, A: 255},
			Unlocked:    true,
			Delivery:    DeliveryConfig{Type: DeliveryDirectProjectile, Speed: 360.0},
			OnHit:       []OnHitEffect{{Type: OnHitDamage, Damage: 10.0}},
		},
		{
			ID:          "mana_siphon",
			Name:        "Mana Siphon",
			Description: "Generates +4.0 Cyber-Mana/sec to power active spells.",
			BaseCost:    45,
			Range:       0,
			Cooldown:    0,
			ManaRate:    4.0,
			Color:       color.RGBA{R: 160, G: 80, B: 255, A: 255},
			Unlocked:    true,
			Delivery:    DeliveryConfig{Type: DeliveryPassiveGenerator},
		},
		{
			ID:          "crypto_miner",
			Name:        "Crypto Miner",
			Description: "Passive economy rig generating continuous +2.0 Bytes/sec.",
			BaseCost:    35,
			Range:       0,
			Cooldown:    0,
			ByteRate:    2.0,
			Color:       color.RGBA{R: 40, G: 230, B: 120, A: 255},
			Unlocked:    true,
			Delivery:    DeliveryConfig{Type: DeliveryPassiveGenerator},
		},
		{
			ID:          "quantum_miner",
			Name:        "Quantum Miner",
			Description: "High-throughput quantum economy node generating +6.0 Bytes/sec.",
			BaseCost:    85,
			Range:       0,
			Cooldown:    0,
			ByteRate:    6.0,
			Color:       color.RGBA{R: 255, G: 215, B: 30, A: 255},
			Unlocked:    false,
			Delivery:    DeliveryConfig{Type: DeliveryPassiveGenerator},
		},
		{
			ID:          "tesla_bus",
			Name:        "Tesla Bus",
			Description: "Arcs high-voltage logic current across up to 4 clustered enemies.",
			BaseCost:    75,
			Range:       95,
			Cooldown:    1.1,
			Color:       color.RGBA{R: 255, G: 220, B: 50, A: 255},
			Unlocked:    false,
			Delivery:    DeliveryConfig{Type: DeliveryInstantBeam, BeamDuration: 0.16},
			OnHit: []OnHitEffect{
				{Type: OnHitDamage, Damage: 22.0},
				{Type: OnHitChain, ChainCount: 3, ChainRange: 75.0, ChainFalloff: 0.70, Damage: 22.0},
			},
		},
		{
			ID:          "cryo_cache",
			Name:        "Cryo-Cache",
			Description: "Emits cryogenic pulses that slow passing malware by 50%.",
			BaseCost:    60,
			Range:       80,
			Cooldown:    1.4,
			Color:       color.RGBA{R: 60, G: 190, B: 255, A: 255},
			Unlocked:    false,
			Delivery:    DeliveryConfig{Type: DeliveryRadialPulse, BeamDuration: 0.25},
			OnHit: []OnHitEffect{
				{Type: OnHitDamage, Damage: 5.0},
				{Type: OnHitSlow, SlowFactor: 0.50, SlowDuration: 2.2},
			},
		},
		{
			ID:          "logic_mortar",
			Name:        "Logic Mortar",
			Description: "Long-range ballistic artillery dealing heavy splash AoE.",
			BaseCost:    110,
			Range:       150,
			Cooldown:    2.4,
			Color:       color.RGBA{R: 255, G: 120, B: 40, A: 255},
			Unlocked:    false,
			Delivery:    DeliveryConfig{Type: DeliveryBallisticMortar, Speed: 180.0, SplashRadius: 55.0},
			OnHit:       []OnHitEffect{{Type: OnHitSplashDamage, Damage: 55.0, SplashRadius: 55.0}},
		},
		{
			ID:          "emp_railgun",
			Name:        "EMP Railgun",
			Description: "High-voltage magnetic beam piercing targets and stunning them for 2.0s.",
			BaseCost:    130,
			Range:       170,
			Cooldown:    2.6,
			Color:       color.RGBA{R: 0, G: 255, B: 210, A: 255},
			Unlocked:    false,
			Delivery:    DeliveryConfig{Type: DeliveryInstantBeam, BeamDuration: 0.22},
			OnHit: []OnHitEffect{
				{Type: OnHitDamage, Damage: 65.0},
				{Type: OnHitFreeze, FreezeSeconds: 2.0},
			},
		},
		{
			ID:          "plasma_flak",
			Name:        "Plasma Flak",
			Description: "Rapid anti-swarm flak battery firing explosive shrapnel rounds.",
			BaseCost:    95,
			Range:       115,
			Cooldown:    0.85,
			Color:       color.RGBA{R: 255, G: 150, B: 50, A: 255},
			Unlocked:    false,
			Delivery:    DeliveryConfig{Type: DeliveryBallisticMortar, Speed: 230.0, SplashRadius: 42.0},
			OnHit: []OnHitEffect{
				{Type: OnHitSplashDamage, Damage: 28.0, SplashRadius: 42.0},
				{Type: OnHitSlow, SlowFactor: 0.35, SlowDuration: 1.5},
			},
		},
		{
			ID:          "overclock_beacon",
			Name:        "Overclock Beacon",
			Description: "Continuous pulse transmitter dealing periodic damage to all nearby threats.",
			BaseCost:    100,
			Range:       85,
			Cooldown:    0.60,
			Color:       color.RGBA{R: 220, G: 70, B: 255, A: 255},
			Unlocked:    false,
			Delivery:    DeliveryConfig{Type: DeliveryRadialPulse, BeamDuration: 0.20},
			OnHit:       []OnHitEffect{{Type: OnHitDamage, Damage: 14.0}},
		},
		{
			ID:          "singularity_vortex",
			Name:        "Singularity Node",
			Description: "Gravitational node that severely crushes and slows heavy malware.",
			BaseCost:    160,
			Range:       135,
			Cooldown:    2.8,
			Color:       color.RGBA{R: 170, G: 40, B: 240, A: 255},
			Unlocked:    false,
			Delivery:    DeliveryConfig{Type: DeliveryDirectProjectile, Speed: 300.0},
			OnHit: []OnHitEffect{
				{Type: OnHitDamage, Damage: 95.0},
				{Type: OnHitSlow, SlowFactor: 0.65, SlowDuration: 3.0},
			},
		},
		{
			ID:          "bit_leech",
			Name:        "Bit-Leech Firewall",
			Description: "Vampiric node harvesting extra bytes and ambient mana on each hit.",
			BaseCost:    85,
			Range:       90,
			Cooldown:    0.75,
			ManaRate:    1.5,
			Color:       color.RGBA{R: 80, G: 255, B: 140, A: 255},
			Unlocked:    false,
			Delivery:    DeliveryConfig{Type: DeliveryDirectProjectile, Speed: 340.0},
			OnHit:       []OnHitEffect{{Type: OnHitDamage, Damage: 20.0}},
		},
	}

	// 2. EXPANDED ACTIVE SPELLS (9 Unique Cyber-Spells)
	r.Spells = []*SpellDef{
		{
			ID:           "chain_lightning",
			Name:         "Chain Lightning",
			Description:  "High-voltage burst arcing across up to 5 adjacent malware.",
			ManaCost:     25,
			Cooldown:     4.0,
			TargetRadius: 120,
			HotKey:       "Q",
			Color:        color.RGBA{R: 255, G: 230, B: 60, A: 255},
			Unlocked:     true,
			CastType:     CastPointArea,
			Visual:       VisualLightning,
			Effects: []SpellEffectDef{
				{
					Type:        SpellEffectChainLightning,
					Damage:      60.0,
					Radius:      120.0,
					ChainCount:  5,
					ChainRange:  110.0,
					ChainDamage: 48.0,
				},
			},
		},
		{
			ID:           "chrono_freeze",
			Name:         "Chrono Freeze",
			Description:  "Cryogenic glitch wave freezing all malware in target zone for 4.5s.",
			ManaCost:     40,
			Cooldown:     12.0,
			TargetRadius: 110,
			HotKey:       "W",
			Color:        color.RGBA{R: 80, G: 220, B: 255, A: 255},
			Unlocked:     false,
			CastType:     CastPointArea,
			Visual:       VisualFreezeRing,
			Effects: []SpellEffectDef{
				{Type: SpellEffectAreaFreeze, Damage: 10.0, Radius: 110.0, FreezeSeconds: 4.5},
			},
		},
		{
			ID:           "logic_bomb",
			Name:         "Logic Bomb",
			Description:  "High-yield thermal explosion wiping out heavy swarms.",
			ManaCost:     60,
			Cooldown:     15.0,
			TargetRadius: 130,
			HotKey:       "E",
			Color:        color.RGBA{R: 255, G: 90, B: 40, A: 255},
			Unlocked:     false,
			CastType:     CastPointArea,
			Visual:       VisualExplosion,
			Effects: []SpellEffectDef{
				{Type: SpellEffectAreaDamage, Damage: 150.0, Radius: 130.0},
			},
		},
		{
			ID:           "kernel_overclock",
			Name:         "Overclock",
			Description:  "Surges motherboard: +70% firing rate & double Mana for 6s.",
			ManaCost:     50,
			Cooldown:     18.0,
			TargetRadius: 0,
			HotKey:       "R",
			Color:        color.RGBA{R: 200, G: 80, B: 255, A: 255},
			Unlocked:     false,
			CastType:     CastGlobalBuff,
			Visual:       VisualGlobalPulse,
			Effects: []SpellEffectDef{
				{Type: SpellEffectOverclock, BuffDuration: 6.0, BuffSpeedMult: 1.70},
			},
		},
		{
			ID:           "emp_disruption",
			Name:         "EMP Disruption",
			Description:  "Strips all enemy shields and stuns an entire sector for 3.5s.",
			ManaCost:     55,
			Cooldown:     14.0,
			TargetRadius: 140,
			HotKey:       "F",
			Color:        color.RGBA{R: 0, G: 255, B: 220, A: 255},
			Unlocked:     false,
			CastType:     CastPointArea,
			Visual:       VisualLightning,
			Effects: []SpellEffectDef{
				{Type: SpellEffectAreaFreeze, Damage: 40.0, Radius: 140.0, FreezeSeconds: 3.5},
			},
		},
		{
			ID:           "thermal_supernova",
			Name:         "Supernova",
			Description:  "Apocalyptic thermal core overload dealing 300 damage in a massive blast.",
			ManaCost:     85,
			Cooldown:     24.0,
			TargetRadius: 180,
			HotKey:       "G",
			Color:        color.RGBA{R: 255, G: 40, B: 40, A: 255},
			Unlocked:     false,
			CastType:     CastPointArea,
			Visual:       VisualExplosion,
			Effects: []SpellEffectDef{
				{Type: SpellEffectAreaDamage, Damage: 300.0, Radius: 180.0},
			},
		},
		{
			ID:           "time_warp",
			Name:         "Time Warp",
			Description:  "Chronometric distortion slowing all enemies by 85% for 6.0 seconds.",
			ManaCost:     65,
			Cooldown:     20.0,
			TargetRadius: 0,
			HotKey:       "T",
			Color:        color.RGBA{R: 120, G: 160, B: 255, A: 255},
			Unlocked:     false,
			CastType:     CastGlobalBuff,
			Visual:       VisualGlobalPulse,
			Effects: []SpellEffectDef{
				{Type: SpellEffectAreaFreeze, FreezeSeconds: 3.0},
			},
		},
		{
			ID:           "nanite_repair",
			Name:         "Nanite Pulse",
			Description:  "Emergency subroutine instantly restoring +50 Kernel Integrity.",
			ManaCost:     50,
			Cooldown:     22.0,
			TargetRadius: 0,
			HotKey:       "Y",
			Color:        color.RGBA{R: 80, G: 255, B: 150, A: 255},
			Unlocked:     false,
			CastType:     CastGlobalBuff,
			Visual:       VisualGlobalPulse,
			Effects: []SpellEffectDef{
				{Type: SpellEffectOverclock, BuffDuration: 2.0, BuffSpeedMult: 1.0},
			},
		},
	}

	// 3. ENEMIES REGISTRY
	r.Enemies = []*EnemyDef{
		{
			ID:           "swarmer",
			Name:         "Packet Swarmer",
			Description:  "Fast, low-integrity basic packets. Swarms through open maze paths rapidly.",
			Archetype:    ArchetypeNormal,
			BaseHP:       22.0,
			BaseSpeed:    82.0,
			BaseBounty:   5.0,
			BaseXP:       6.0,
			CoreDamage:   4.0,
			Radius:       6.0,
			Color:        color.RGBA{R: 255, G: 70, B: 90, A: 255},
			MinSpawnTime: 0.0,
			SpawnWeight:  100.0,
		},
		{
			ID:           "sprite",
			Name:         "Spyware Sprite",
			Description:  "Shielded malware unit protected by high-frequency energy barriers.",
			Archetype:    ArchetypeShielded,
			BaseHP:       45.0,
			BaseShield:   40.0,
			BaseSpeed:    66.0,
			BaseBounty:   12.0,
			BaseXP:       14.0,
			CoreDamage:   8.0,
			Radius:       8.0,
			Color:        color.RGBA{R: 80, G: 180, B: 255, A: 255},
			MinSpawnTime: 18.0,
			SpawnWeight:  55.0,
		},
		{
			ID:           "golem",
			Name:         "Trojan Golem",
			Description:  "Heavily armored container. Splits into two Trojan Shards upon destruction.",
			Archetype:    ArchetypeTankSplitter,
			BaseHP:       140.0,
			BaseSpeed:    38.0,
			BaseBounty:   20.0,
			BaseXP:       22.0,
			CoreDamage:   16.0,
			Radius:       12.0,
			Color:        color.RGBA{R: 220, G: 140, B: 40, A: 255},
			MinSpawnTime: 35.0,
			SpawnWeight:  40.0,
			SplitSpawnID: "trojan_shard",
			SplitCount:   2,
		},
		{
			ID:           "trojan_shard",
			Name:         "Trojan Shard",
			Description:  "Sub-packet spawned from cracked Trojan Golems. Fast and erratic.",
			Archetype:    ArchetypeSplitterSub,
			BaseHP:       35.0,
			BaseSpeed:    95.0,
			BaseBounty:   6.0,
			BaseXP:       8.0,
			CoreDamage:   5.0,
			Radius:       6.0,
			Color:        color.RGBA{R: 240, G: 180, B: 70, A: 255},
			MinSpawnTime: 99999.0, // Only spawned via splits
			SpawnWeight:  0.0,
		},
		{
			ID:             "healer_worm",
			Name:           "Polymorphic Worm",
			Description:    "Bio-digital parasite that broadcasts self-repair code, healing nearby malware.",
			Archetype:      ArchetypeHealer,
			BaseHP:         85.0,
			BaseSpeed:      52.0,
			BaseBounty:     18.0,
			BaseXP:         20.0,
			CoreDamage:     10.0,
			Radius:         9.0,
			Color:          color.RGBA{R: 60, G: 240, B: 120, A: 255},
			MinSpawnTime:   45.0,
			SpawnWeight:    35.0,
			HealAuraRadius: 75.0,
			HealAuraAmount: 12.0,
			HealAuraPeriod: 1.8,
		},
		{
			ID:            "phantom_glitch",
			Name:          "Phantom Glitch",
			Description:   "Encrypted phantom packet that shifts out of phase, becoming untargetable.",
			Archetype:     ArchetypePhaseCloaker,
			BaseHP:        55.0,
			BaseSpeed:     70.0,
			BaseBounty:    16.0,
			BaseXP:        18.0,
			CoreDamage:    12.0,
			Radius:        7.5,
			Color:         color.RGBA{R: 200, G: 80, B: 255, A: 255},
			MinSpawnTime:  60.0,
			SpawnWeight:   30.0,
			PhasePeriod:   4.5,
			PhaseDuration: 1.8,
		},
		{
			ID:            "buffer_overflower",
			Name:          "Buffer Overflower",
			Description:   "Unstable volatility packet that accelerates and explodes violently upon death.",
			Archetype:     ArchetypeExploder,
			BaseHP:        60.0,
			BaseSpeed:     88.0,
			BaseBounty:    15.0,
			BaseXP:        16.0,
			CoreDamage:    15.0,
			Radius:        8.0,
			Color:         color.RGBA{R: 255, G: 30, B: 30, A: 255},
			MinSpawnTime:  50.0,
			SpawnWeight:   32.0,
			ExplodeRadius: 65.0,
			ExplodeDamage: 35.0,
		},
		{
			ID:           "boss_daemon",
			Name:         "Daemon Overlord",
			Description:  "Tier 1 Boss: Colossal rootkit with thick shields and high-impact kernel breach power.",
			Archetype:    ArchetypeBoss,
			BaseHP:       550.0,
			BaseShield:   250.0,
			BaseSpeed:    30.0,
			BaseBounty:   100.0,
			BaseXP:       120.0,
			CoreDamage:   45.0,
			Radius:       16.0,
			Color:        color.RGBA{R: 255, G: 0, B: 80, A: 255},
			IsBoss:       true,
			MinSpawnTime: 60.0,
			SpawnWeight:  0.0,
		},
		{
			ID:             "boss_zeroday",
			Name:           "Zero-Day Harbinger",
			Description:    "Tier 2+ Boss: Apocalyptic stealth threat that heals allies and phases through firewalls.",
			Archetype:      ArchetypeBoss,
			BaseHP:         950.0,
			BaseShield:     450.0,
			BaseSpeed:      32.0,
			BaseBounty:     200.0,
			BaseXP:         250.0,
			CoreDamage:     60.0,
			Radius:         18.0,
			Color:          color.RGBA{R: 255, G: 40, B: 200, A: 255},
			IsBoss:         true,
			MinSpawnTime:   180.0,
			SpawnWeight:    0.0,
			HealAuraRadius: 90.0,
			HealAuraAmount: 20.0,
			HealAuraPeriod: 2.0,
			PhasePeriod:    6.0,
			PhaseDuration:  1.5,
		},
	}

	// 4. DRAFT CARDS REGISTRY
	r.Cards = []*CardDef{
		// Tower Blueprint Unlocks
		{
			ID: "unlock_tesla", Title: "NEW BLUEPRINT", Subtitle: "Arcane Tesla Bus",
			Description: "Unlocks the Tesla Bus node. Arcs high-voltage chain lightning.",
			Color:       color.RGBA{R: 255, G: 220, B: 50, A: 255}, IsUnlock: true,
			Effects: []CardEffect{{Type: ModUnlockTower, TargetID: "tesla_bus"}},
		},
		{
			ID: "unlock_cryo", Title: "NEW BLUEPRINT", Subtitle: "Cryo-Stasis Cache",
			Description: "Unlocks Cryo-Cache node. Slows passing malware by 50%.",
			Color:       color.RGBA{R: 80, G: 200, B: 255, A: 255}, IsUnlock: true,
			Effects: []CardEffect{{Type: ModUnlockTower, TargetID: "cryo_cache"}},
		},
		{
			ID: "unlock_mortar", Title: "NEW BLUEPRINT", Subtitle: "Logic Mortar",
			Description: "Unlocks long-range artillery node dealing high explosive splash AoE.",
			Color:       color.RGBA{R: 255, G: 120, B: 40, A: 255}, IsUnlock: true,
			Effects: []CardEffect{{Type: ModUnlockTower, TargetID: "logic_mortar"}},
		},
		{
			ID: "unlock_emp_rail", Title: "NEW BLUEPRINT", Subtitle: "EMP Railgun",
			Description: "Unlocks piercing EMP beam node dealing heavy damage and 2s stuns.",
			Color:       color.RGBA{R: 0, G: 255, B: 210, A: 255}, IsUnlock: true,
			Effects: []CardEffect{{Type: ModUnlockTower, TargetID: "emp_railgun"}},
		},
		{
			ID: "unlock_flak", Title: "NEW BLUEPRINT", Subtitle: "Plasma Flak Battery",
			Description: "Unlocks rapid shrapnel cannon that shreds packed swarms.",
			Color:       color.RGBA{R: 255, G: 150, B: 50, A: 255}, IsUnlock: true,
			Effects: []CardEffect{{Type: ModUnlockTower, TargetID: "plasma_flak"}},
		},
		{
			ID: "unlock_beacon", Title: "NEW BLUEPRINT", Subtitle: "Overclock Beacon",
			Description: "Unlocks continuous pulse transmitter damaging all surrounding malware.",
			Color:       color.RGBA{R: 220, G: 70, B: 255, A: 255}, IsUnlock: true,
			Effects: []CardEffect{{Type: ModUnlockTower, TargetID: "overclock_beacon"}},
		},
		{
			ID: "unlock_singularity", Title: "NEW BLUEPRINT", Subtitle: "Singularity Node",
			Description: "Unlocks heavy gravity well node with massive damage and 65% slow.",
			Color:       color.RGBA{R: 170, G: 40, B: 240, A: 255}, IsUnlock: true,
			Effects: []CardEffect{{Type: ModUnlockTower, TargetID: "singularity_vortex"}},
		},
		{
			ID: "unlock_bit_leech", Title: "NEW BLUEPRINT", Subtitle: "Bit-Leech Firewall",
			Description: "Unlocks vampiric node generating bonus mana and byte bounties.",
			Color:       color.RGBA{R: 80, G: 255, B: 140, A: 255}, IsUnlock: true,
			Effects: []CardEffect{{Type: ModUnlockTower, TargetID: "bit_leech"}},
		},
		{
			ID: "unlock_quantum_miner", Title: "NEW BLUEPRINT", Subtitle: "Quantum Miner",
			Description: "Unlocks high-throughput quantum economy node generating +6.0 Bytes/sec.",
			Color:       color.RGBA{R: 255, G: 215, B: 30, A: 255}, IsUnlock: true,
			Effects: []CardEffect{{Type: ModUnlockTower, TargetID: "quantum_miner"}},
		},
		// Spell Unlocks
		{
			ID: "unlock_freeze", Title: "NEW SPELL", Subtitle: "Chrono Freeze [W]",
			Description: "Freezes all malware in target zone for 4.5s. Costs Cyber-Mana.",
			Color:       color.RGBA{R: 100, G: 220, B: 255, A: 255}, IsUnlock: true,
			Effects: []CardEffect{{Type: ModUnlockSpell, TargetID: "chrono_freeze"}},
		},
		{
			ID: "unlock_bomb", Title: "NEW SPELL", Subtitle: "Logic Bomb [E]",
			Description: "Catastrophic thermal payload wiping out swarms in a wide radius.",
			Color:       color.RGBA{R: 255, G: 80, B: 60, A: 255}, IsUnlock: true,
			Effects: []CardEffect{{Type: ModUnlockSpell, TargetID: "logic_bomb"}},
		},
		{
			ID: "unlock_overclock", Title: "NEW SPELL", Subtitle: "Overclock [R]",
			Description: "Motherboard surge: +70% firing rate & double Mana siphon output for 6s.",
			Color:       color.RGBA{R: 210, G: 80, B: 255, A: 255}, IsUnlock: true,
			Effects: []CardEffect{{Type: ModUnlockSpell, TargetID: "kernel_overclock"}},
		},
		{
			ID: "unlock_emp_wave", Title: "NEW SPELL", Subtitle: "EMP Disruption [F]",
			Description: "Sector EMP blast stripping enemy shields and freezing all targets.",
			Color:       color.RGBA{R: 0, G: 255, B: 220, A: 255}, IsUnlock: true,
			Effects: []CardEffect{{Type: ModUnlockSpell, TargetID: "emp_disruption"}},
		},
		{
			ID: "unlock_supernova", Title: "NEW SPELL", Subtitle: "Supernova [G]",
			Description: "Ultimate thermal overload wiping out 300 HP across a giant sector.",
			Color:       color.RGBA{R: 255, G: 40, B: 40, A: 255}, IsUnlock: true,
			Effects: []CardEffect{{Type: ModUnlockSpell, TargetID: "thermal_supernova"}},
		},
		// Passives
		{
			ID: "buff_damage", Title: "SYSTEM TWEAK", Subtitle: "Amplified Voltage",
			Description: "+20% Global damage for all defense nodes.",
			Color:       color.RGBA{R: 255, G: 90, B: 90, A: 255},
			Effects:     []CardEffect{{Type: ModTowerDamageMult, Value: 0.20}},
		},
		{
			ID: "buff_mana_gen", Title: "SYSTEM TWEAK", Subtitle: "Ether Superconductor",
			Description: "+30% Cyber-Mana generation speed from Siphon nodes.",
			Color:       color.RGBA{R: 180, G: 90, B: 255, A: 255},
			Effects:     []CardEffect{{Type: ModManaGenMult, Value: 0.30}},
		},
		{
			ID: "buff_speed", Title: "SYSTEM TWEAK", Subtitle: "Clock Boost",
			Description: "+15% Attack speed for all defense nodes.",
			Color:       color.RGBA{R: 255, G: 210, B: 60, A: 255},
			Effects:     []CardEffect{{Type: ModTowerSpeedMult, Value: 0.15}},
		},
		{
			ID: "buff_range", Title: "SYSTEM TWEAK", Subtitle: "Broadband Sensor",
			Description: "+20% Firing range for all defense nodes.",
			Color:       color.RGBA{R: 60, G: 220, B: 200, A: 255},
			Effects:     []CardEffect{{Type: ModTowerRangeMult, Value: 0.20}},
		},
		{
			ID: "buff_bounty", Title: "SYSTEM TWEAK", Subtitle: "Byte Harvest Engine",
			Description: "+25% Byte bounty rewarded from purged malware.",
			Color:       color.RGBA{R: 80, G: 255, B: 120, A: 255},
			Effects:     []CardEffect{{Type: ModByteBountyMult, Value: 0.25}},
		},
		{
			ID: "heal_kernel", Title: "SYSTEM TWEAK", Subtitle: "Kernel Patch",
			Description: "Instantly restore +35 Kernel Integrity and +50 Bytes.",
			Color:       color.RGBA{R: 80, G: 160, B: 255, A: 255},
			Effects:     []CardEffect{{Type: ModKernelHeal, Value: 35.0}},
		},
	}

	for _, t := range r.Towers {
		r.towerMap[t.ID] = t
	}
	for _, s := range r.Spells {
		r.spellMap[s.ID] = s
	}
	for _, e := range r.Enemies {
		r.enemyMap[e.ID] = e
	}
	for _, c := range r.Cards {
		r.cardMap[c.ID] = c
	}

	return r
}

func (r *Registry) GetTower(id string) *TowerDef {
	return r.towerMap[id]
}

func (r *Registry) GetSpell(id string) *SpellDef {
	return r.spellMap[id]
}

func (r *Registry) GetEnemy(id string) *EnemyDef {
	return r.enemyMap[id]
}

func (r *Registry) GetCard(id string) *CardDef {
	return r.cardMap[id]
}

func (r *Registry) CloneForNewRun() *Registry {
	return NewRegistry()
}
