package meta

import (
	"technomancers-tower/internal/save"
)

// Talent represents a permanent meta-progression talent.
type Talent struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Level       int     `json:"level"`
	MaxLevel    int     `json:"max_level"`
	BaseCost    int     `json:"base_cost"`
	CostMult    float64 `json:"cost_mult"`
}

func (t *Talent) CurrentCost() int {
	if t.Level >= t.MaxLevel {
		return -1
	}
	cost := float64(t.BaseCost)
	for i := 0; i < t.Level; i++ {
		cost *= t.CostMult
	}
	return int(cost)
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
			{
				ID:          "kernel_shield",
				Name:        "Kernel Shielding",
				Description: "+25 Max Kernel Integrity per rank.",
				Level:       0,
				MaxLevel:    10,
				BaseCost:    10,
				CostMult:    1.5,
			},
			{
				ID:          "boot_bytes",
				Name:        "Bootloader Cache",
				Description: "+40 Starting Bytes per rank.",
				Level:       0,
				MaxLevel:    8,
				BaseCost:    15,
				CostMult:    1.45,
			},
			{
				ID:          "mana_conductor",
				Name:        "Mana Bus Conductor",
				Description: "+20% Mana Siphon generation speed per rank.",
				Level:       0,
				MaxLevel:    8,
				BaseCost:    20,
				CostMult:    1.5,
			},
			{
				ID:          "overclock_nodes",
				Name:        "Global Overclock",
				Description: "+5% Attack rate to all defense nodes per rank.",
				Level:       0,
				MaxLevel:    10,
				BaseCost:    25,
				CostMult:    1.6,
			},
			{
				ID:          "spell_efficiency",
				Name:        "Subroutine Optimization",
				Description: "-8% Spell cooldowns and -10% Mana cost per rank.",
				Level:       0,
				MaxLevel:    5,
				BaseCost:    35,
				CostMult:    1.7,
			},
			{
				ID:          "scrap_leech",
				Name:        "Byte Harvester",
				Description: "+10% Byte reward from purged malware per rank.",
				Level:       0,
				MaxLevel:    8,
				BaseCost:    20,
				CostMult:    1.5,
			},
			{
				ID:          "tower_potency",
				Name:        "Amplified Core Logic",
				Description: "+6% Global damage to all defense nodes per rank.",
				Level:       0,
				MaxLevel:    10,
				BaseCost:    30,
				CostMult:    1.55,
			},
			{
				ID:          "sensor_array",
				Name:        "Broadband Sensor Bus",
				Description: "+6% Firing range to all defense nodes per rank.",
				Level:       0,
				MaxLevel:    8,
				BaseCost:    25,
				CostMult:    1.5,
			},
			{
				ID:          "nanite_resilience",
				Name:        "Self-Repair Subroutine",
				Description: "Automatically restore +15 Kernel Integrity on Boss kill.",
				Level:       0,
				MaxLevel:    5,
				BaseCost:    40,
				CostMult:    1.65,
			},
		},
	}

	for _, t := range m.Talents {
		m.talentMap[t.ID] = t
	}

	return m
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
		m.SaveActiveSlot()
		return true
	}
	return false
}

func (m *MetaManager) AddShards(amount int) {
	m.GlitchShards += amount
	m.SaveActiveSlot()
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
