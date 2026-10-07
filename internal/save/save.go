package save

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const MaxSlots = 3

// TowerSnapshot stores placed defense node info.
type TowerSnapshot struct {
	TowerID string
	GridX   int
	GridY   int
	Tier    int
	Level   int
	XP      float64
}

// RunSnapshot stores active in-progress run state.
type RunSnapshot struct {
	Bytes           float64
	TotalBytes      float64
	Mana            float64
	MaxMana         float64
	KernelHP        float64
	MaxKernelHP     float64
	Level           int
	CurrentXP       float64
	TargetXP        float64
	PendingDrafts   int
	RunTime         float64
	Kills           int
	BossKills       int
	TowerDamageMult float64
	TowerRangeMult  float64
	TowerSpeedMult  float64
	ManaGenMult     float64
	ByteBountyMult  float64
	Towers          []TowerSnapshot
	UnlockedTowers  []string
	UnlockedSpells  []string
}

// SlotData stores all data for a single save slot.
type SlotData struct {
	SlotID        int
	ProfileName   string
	GlitchShards  int
	HighScoreTime float64
	TotalKills    int
	Talents       map[string]int
	LastSaved     time.Time
	HasActiveRun  bool
	ActiveRun     RunSnapshot
}

// SlotSummary contains lightweight summary for menu rendering.
type SlotSummary struct {
	SlotID            int
	Exists            bool
	ProfileName       string
	GlitchShards      int
	HighScore         float64
	TotalKills        int
	LastSaved         time.Time
	HasActiveRun      bool
	ActiveRunDuration float64
	ActiveRunHP       float64
}

// SaveManager handles multi-file compressed binary save operations.
type SaveManager struct {
	BaseDir    string
	ActiveSlot int
	Slots      [MaxSlots]SlotSummary
}

func NewSaveManager() *SaveManager {
	home, err := os.UserHomeDir()
	var saveDir string
	if err == nil {
		saveDir = filepath.Join(home, ".technomancers-tower", "saves")
	} else {
		saveDir = filepath.Join(".", "saves")
	}
	_ = os.MkdirAll(saveDir, 0755)

	sm := &SaveManager{
		BaseDir:    saveDir,
		ActiveSlot: 1,
	}

	sm.RefreshSlots()
	return sm
}

func (sm *SaveManager) getSlotPath(slotID int) string {
	return filepath.Join(sm.BaseDir, fmt.Sprintf("slot_%d.dat", slotID))
}

func (sm *SaveManager) RefreshSlots() {
	for i := 1; i <= MaxSlots; i++ {
		path := sm.getSlotPath(i)
		data, err := sm.readCompressedSlot(path)
		if err != nil || data == nil {
			sm.Slots[i-1] = SlotSummary{
				SlotID: i,
				Exists: false,
			}
		} else {
			sm.Slots[i-1] = SlotSummary{
				SlotID:            i,
				Exists:            true,
				ProfileName:       data.ProfileName,
				GlitchShards:      data.GlitchShards,
				HighScore:         data.HighScoreTime,
				TotalKills:        data.TotalKills,
				LastSaved:         data.LastSaved,
				HasActiveRun:      data.HasActiveRun,
				ActiveRunDuration: data.ActiveRun.RunTime,
				ActiveRunHP:       data.ActiveRun.KernelHP,
			}
		}
	}
}

func (sm *SaveManager) LoadSlot(slotID int) (*SlotData, error) {
	path := sm.getSlotPath(slotID)
	data, err := sm.readCompressedSlot(path)
	if err != nil {
		return nil, err
	}
	if data == nil {
		// New slot default
		data = &SlotData{
			SlotID:       slotID,
			ProfileName:  fmt.Sprintf("OPERATOR 0%d", slotID),
			Talents:      make(map[string]int),
			LastSaved:    time.Now(),
			HasActiveRun: false,
		}
	}
	sm.ActiveSlot = slotID
	return data, nil
}

func (sm *SaveManager) SaveSlot(data *SlotData) error {
	data.LastSaved = time.Now()
	if data.ProfileName == "" {
		data.ProfileName = fmt.Sprintf("OPERATOR 0%d", data.SlotID)
	}

	path := sm.getSlotPath(data.SlotID)
	tmpPath := path + ".tmp"

	// 1. Encode with Gob
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	enc := gob.NewEncoder(gz)
	if err := enc.Encode(data); err != nil {
		gz.Close()
		return err
	}
	if err := gz.Close(); err != nil {
		return err
	}

	// 2. Write to temp file
	if err := os.WriteFile(tmpPath, buf.Bytes(), 0644); err != nil {
		return err
	}

	// 3. Atomic rename to prevent corruption
	if err := os.Rename(tmpPath, path); err != nil {
		return err
	}

	sm.RefreshSlots()
	return nil
}

func (sm *SaveManager) DeleteSlot(slotID int) error {
	path := sm.getSlotPath(slotID)
	_ = os.Remove(path)
	_ = os.Remove(path + ".tmp")
	sm.RefreshSlots()
	return nil
}

func (sm *SaveManager) readCompressedSlot(path string) (*SlotData, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	gz, err := gzip.NewReader(bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}
	defer gz.Close()

	var data SlotData
	dec := gob.NewDecoder(gz)
	if err := dec.Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}
