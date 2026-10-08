package meta

import (
	"technomancers-tower/internal/save"
)

type TalentCategory string

const (
	CategoryDefense TalentCategory = "DEFENSE"
	CategoryOffense TalentCategory = "OFFENSE"
	CategoryArcane  TalentCategory = "ARCANE"
	CategoryEconomy TalentCategory = "ECONOMY"
	CategoryUtility TalentCategory = "UTILITY"
)

var AllCategories = []TalentCategory{
	CategoryDefense,
	CategoryOffense,
	CategoryArcane,
	CategoryEconomy,
	CategoryUtility,
}

// Talent represents a permanent meta-progression talent.
type Talent struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Category    TalentCategory `json:"category"`
	Level       int            `json:"level"`
	MaxLevel    int            `json:"max_level"`
	BaseCost    int            `json:"base_cost"`
	CostMult    float64        `json:"cost_mult"`
}

func (t *Talent) CostAtLevel(level int) int {
	if level < 0 || level >= t.MaxLevel {
		return -1
	}
	cost := float64(t.BaseCost)
	for i := 0; i < level; i++ {
		cost *= t.CostMult
	}
	return int(cost)
}

func (t *Talent) CurrentCost() int {
	if t.Level >= t.MaxLevel {
		return -1
	}
	return t.CostAtLevel(t.Level)
}

func (t *Talent) PreviousRankCost() int {
	if t.Level <= 0 {
		return 0
	}
	return t.CostAtLevel(t.Level - 1)
}

func (t *Talent) TotalInvested() int {
	total := 0
	for i := 0; i < t.Level; i++ {
		c := t.CostAtLevel(i)
		if c > 0 {
			total += c
		}
	}
	return total
}

// MetaManager handles permanent talents connected to the multi-slot save system.
type MetaManager struct {
	SaveMgr         *save.SaveManager
	CurrentSlotData *save.SlotData
	GlitchShards    int
	HighScore       float64
	TotalKills      int
	Talents         []*Talent
	talentMap       map[string]*Talent
}

func NewMetaManager(sm *save.SaveManager) *MetaManager {
	m := &MetaManager{
		SaveMgr:   sm,
		talentMap: make(map[string]*Talent),
		Talents: []*Talent{
			// ==================== 1. DEFENSE & CORE ====================
			{
				ID:          "kernel_shield",
				Name:        "Kernel Shielding",
				Description: "+25 Max Kernel Integrity per rank.",
				Category:    CategoryDefense,
				Level:       0,
				MaxLevel:    10,
				BaseCost:    10,
				CostMult:    1.5,
			},
			{
				ID:          "nanite_regen",
				Name:        "Nanite Regeneration",
				Description: "+1.0 Kernel HP restored every 15s per rank.",
				Category:    CategoryDefense,
				Level:       0,
				MaxLevel:    5,
				BaseCost:    25,
				CostMult:    1.6,
			},
			{
				ID:          "kernel_emp",
				Name:        "Kernel EMP Pulse",
				Description: "When Kernel is hit, release shockwave dealing 30 dmg/rank and slowing malware by 50%.",
				Category:    CategoryDefense,
				Level:       0,
				MaxLevel:    5,
				BaseCost:    30,
				CostMult:    1.6,
			},
			{
				ID:          "hardened_firewall",
				Name:        "Hardened Firewall",
				Description: "Reduces all incoming Kernel damage by -1 (min 1) per rank.",
				Category:    CategoryDefense,
				Level:       0,
				MaxLevel:    5,
				BaseCost:    35,
				CostMult:    1.65,
			},
			{
				ID:          "nanite_resilience",
				Name:        "Self-Repair Subroutine",
				Description: "Automatically restore +15 Kernel Integrity on Boss kill per rank.",
				Category:    CategoryDefense,
				Level:       0,
				MaxLevel:    5,
				BaseCost:    40,
				CostMult:    1.65,
			},

			// ==================== 2. OFFENSE & NODES ====================
			{
				ID:          "overclock_nodes",
				Name:        "Global Overclock",
				Description: "+5% Attack rate to all defense nodes per rank.",
				Category:    CategoryOffense,
				Level:       0,
				MaxLevel:    10,
				BaseCost:    25,
				CostMult:    1.6,
			},
			{
				ID:          "tower_potency",
				Name:        "Amplified Core Logic",
				Description: "+6% Global damage to all defense nodes per rank.",
				Category:    CategoryOffense,
				Level:       0,
				MaxLevel:    10,
				BaseCost:    30,
				CostMult:    1.55,
			},
			{
				ID:          "sensor_array",
				Name:        "Broadband Sensor Bus",
				Description: "+6% Firing range to all defense nodes per rank.",
				Category:    CategoryOffense,
				Level:       0,
				MaxLevel:    8,
				BaseCost:    25,
				CostMult:    1.5,
			},
			{
				ID:          "critical_subroutines",
				Name:        "Critical Subroutines",
				Description: "+4% Tower Critical Hit chance (dealing 200% damage) per rank.",
				Category:    CategoryOffense,
				Level:       0,
				MaxLevel:    5,
				BaseCost:    35,
				CostMult:    1.6,
			},
			{
				ID:          "status_amplification",
				Name:        "Corrosive Payloads",
				Description: "+12% Status effect duration (burn, slow, freeze) per rank.",
				Category:    CategoryOffense,
				Level:       0,
				MaxLevel:    5,
				BaseCost:    25,
				CostMult:    1.5,
			},

			// ==================== 3. ARCANE & SPELLS ====================
			{
				ID:          "spell_efficiency",
				Name:        "Subroutine Optimization",
				Description: "-8% Spell cooldowns and -10% Mana cost per rank.",
				Category:    CategoryArcane,
				Level:       0,
				MaxLevel:    5,
				BaseCost:    35,
				CostMult:    1.7,
			},
			{
				ID:          "mana_conductor",
				Name:        "Mana Bus Conductor",
				Description: "+20% Mana Siphon generation speed per rank.",
				Category:    CategoryArcane,
				Level:       0,
				MaxLevel:    8,
				BaseCost:    20,
				CostMult:    1.5,
			},
			{
				ID:          "expanded_mana_pool",
				Name:        "Expanded Buffer Pool",
				Description: "+25 Max Mana capacity per rank.",
				Category:    CategoryArcane,
				Level:       0,
				MaxLevel:    6,
				BaseCost:    20,
				CostMult:    1.5,
			},
			{
				ID:          "natural_mana_regen",
				Name:        "Ambient Ether Siphon",
				Description: "Passively generate +0.5 Mana per second per rank.",
				Category:    CategoryArcane,
				Level:       0,
				MaxLevel:    6,
				BaseCost:    25,
				CostMult:    1.55,
			},
			{
				ID:          "spell_critical",
				Name:        "Arcane Resonance",
				Description: "+5% Spell Crit chance (2.0x dmg) and +8% Spell AOE radius per rank.",
				Category:    CategoryArcane,
				Level:       0,
				MaxLevel:    5,
				BaseCost:    30,
				CostMult:    1.6,
			},

			// ==================== 4. ECONOMY & MINING ====================
			{
				ID:          "boot_bytes",
				Name:        "Bootloader Cache",
				Description: "+40 Starting Bytes per rank.",
				Category:    CategoryEconomy,
				Level:       0,
				MaxLevel:    8,
				BaseCost:    15,
				CostMult:    1.45,
			},
			{
				ID:          "scrap_leech",
				Name:        "Byte Harvester",
				Description: "+10% Byte reward from purged malware per rank.",
				Category:    CategoryEconomy,
				Level:       0,
				MaxLevel:    8,
				BaseCost:    20,
				CostMult:    1.5,
			},
			{
				ID:          "byte_interest",
				Name:        "Algorithmic Interest",
				Description: "Gain +4% bonus bytes on unspent balance every 30s interval per rank.",
				Category:    CategoryEconomy,
				Level:       0,
				MaxLevel:    5,
				BaseCost:    30,
				CostMult:    1.6,
			},
			{
				ID:          "salvage_efficiency",
				Name:        "Recycle Protocol",
				Description: "+10% Tower sell/demolish byte refund rate per rank (base 50%).",
				Category:    CategoryEconomy,
				Level:       0,
				MaxLevel:    5,
				BaseCost:    20,
				CostMult:    1.5,
			},
			{
				ID:          "miner_boost",
				Name:        "Quantum Miner Firmware",
				Description: "Passive Byte Miner nodes produce +15% more bytes/sec per rank.",
				Category:    CategoryEconomy,
				Level:       0,
				MaxLevel:    6,
				BaseCost:    25,
				CostMult:    1.55,
			},

			// ==================== 5. UTILITY & META ====================
			{
				ID:          "shards_harvest",
				Name:        "Glitch Shard Siphon",
				Description: "+15% Glitch Shards earned on run completion per rank.",
				Category:    CategoryUtility,
				Level:       0,
				MaxLevel:    8,
				BaseCost:    30,
				CostMult:    1.5,
			},
			{
				ID:          "draft_rerolls",
				Name:        "Draft Reroll Protocol",
				Description: "Start run with +1 Draft Reroll charge per rank.",
				Category:    CategoryUtility,
				Level:       0,
				MaxLevel:    5,
				BaseCost:    35,
				CostMult:    1.6,
			},
			{
				ID:          "extra_draft_slot",
				Name:        "Extended Draft Buffer",
				Description: "+1 Additional card choice offered during draft selections per rank.",
				Category:    CategoryUtility,
				Level:       0,
				MaxLevel:    2,
				BaseCost:    75,
				CostMult:    2.0,
			},
			{
				ID:          "starting_level",
				Name:        "Accelerated Firmware",
				Description: "Start run at Level 1 + 1 bonus draft selection per rank.",
				Category:    CategoryUtility,
				Level:       0,
				MaxLevel:    3,
				BaseCost:    50,
				CostMult:    1.8,
			},
			{
				ID:          "xp_multiplier",
				Name:        "XP Overclocking",
				Description: "+8% XP gained from purged malware per rank.",
				Category:    CategoryUtility,
				Level:       0,
				MaxLevel:    6,
				BaseCost:    25,
				CostMult:    1.5,
			},
		},
	}

	for _, t := range m.Talents {
		m.talentMap[t.ID] = t
	}

	return m
}

func (m *MetaManager) GetTalentsByCategory(cat TalentCategory) []*Talent {
	res := make([]*Talent, 0)
	for _, t := range m.Talents {
		if t.Category == cat {
			res = append(res, t)
		}
	}
	return res
}

func (m *MetaManager) TotalInvestedShards() int {
	total := 0
	for _, t := range m.Talents {
		total += t.TotalInvested()
	}
	return total
}

func (m *MetaManager) ApplySlot(slotData *save.SlotData) {
	m.CurrentSlotData = slotData
	m.GlitchShards = slotData.GlitchShards
	m.HighScore = slotData.HighScoreTime
	m.TotalKills = slotData.TotalKills

	// Reset all talent levels first
	for _, t := range m.Talents {
		t.Level = 0
	}

	// Apply loaded ranks
	for id, lvl := range slotData.Talents {
		if t, ok := m.talentMap[id]; ok {
			if lvl > t.MaxLevel {
				lvl = t.MaxLevel
			}
			t.Level = lvl
		}
	}
}

func (m *MetaManager) GetTalentLevel(id string) int {
	if t, ok := m.talentMap[id]; ok {
		return t.Level
	}
	return 0
}

func (m *MetaManager) BuyTalent(id string) bool {
	t, ok := m.talentMap[id]
	if !ok || t.Level >= t.MaxLevel {
		return false
	}
	cost := t.CurrentCost()
	if m.GlitchShards >= cost {
		m.GlitchShards -= cost
		t.Level++
		_ = m.SaveActiveSlot()
		return true
	}
	return false
}

func (m *MetaManager) RefundTalent(id string) bool {
	t, ok := m.talentMap[id]
	if !ok || t.Level <= 0 {
		return false
	}
	refund := t.PreviousRankCost()
	m.GlitchShards += refund
	t.Level--
	_ = m.SaveActiveSlot()
	return true
}

func (m *MetaManager) ResetAllTalents() int {
	totalRefunded := 0
	for _, t := range m.Talents {
		totalRefunded += t.TotalInvested()
		t.Level = 0
	}
	m.GlitchShards += totalRefunded
	_ = m.SaveActiveSlot()
	return totalRefunded
}

func (m *MetaManager) AddShards(amount int) {
	m.GlitchShards += amount
	_ = m.SaveActiveSlot()
}

func (m *MetaManager) SaveActiveSlot() error {
	if m.CurrentSlotData == nil {
		return nil
	}
	m.CurrentSlotData.GlitchShards = m.GlitchShards
	m.CurrentSlotData.HighScoreTime = m.HighScore
	m.CurrentSlotData.TotalKills = m.TotalKills
	if m.CurrentSlotData.Talents == nil {
		m.CurrentSlotData.Talents = make(map[string]int)
	}
	for _, t := range m.Talents {
		m.CurrentSlotData.Talents[t.ID] = t.Level
	}

	return m.SaveMgr.SaveSlot(m.CurrentSlotData)
}
