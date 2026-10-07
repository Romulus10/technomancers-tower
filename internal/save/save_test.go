package save_test

import (
	"testing"

	"technomancers-tower/internal/save"
)

func TestSaveLoadManager(t *testing.T) {
	tmpDir := t.TempDir()
	sm := &save.SaveManager{
		BaseDir: tmpDir,
	}

	data := &save.SlotData{
		SlotID: 1,
		Talents: map[string]int{
			"kernel_shield": 3,
			"boot_bytes":    2,
		},
		GlitchShards:  150,
		HighScoreTime: 420.5,
		TotalKills:    95,
		HasActiveRun:  true,
		ActiveRun: save.RunSnapshot{
			Bytes:          500,
			TotalBytes:     1200,
			Level:          4,
			CurrentXP:      80,
			TargetXP:       150,
			RunTime:        125.0,
			Kills:          45,
			UnlockedTowers: []string{"bit_driver", "tesla_bus"},
			Towers: []save.TowerSnapshot{
				{TowerID: "bit_driver", GridX: 2, GridY: 3, Tier: 2, Level: 3, XP: 175.5},
			},
		},
	}

	err := sm.SaveSlot(data)
	if err != nil {
		t.Fatalf("failed to save slot: %v", err)
	}

	// Verify compressed slot file exists
	loaded, err := sm.LoadSlot(1)
	if err != nil {
		t.Fatalf("failed to load slot: %v", err)
	}

	if loaded.SlotID != 1 || loaded.GlitchShards != 150 || loaded.Talents["kernel_shield"] != 3 {
		t.Errorf("loaded data mismatch: %+v", loaded)
	}

	if !loaded.HasActiveRun || loaded.ActiveRun.Bytes != 500 || len(loaded.ActiveRun.UnlockedTowers) != 2 {
		t.Errorf("loaded active run mismatch: %+v", loaded.ActiveRun)
	}
	if len(loaded.ActiveRun.Towers) != 1 || loaded.ActiveRun.Towers[0].Tier != 2 || loaded.ActiveRun.Towers[0].Level != 3 || loaded.ActiveRun.Towers[0].XP != 175.5 {
		t.Errorf("loaded tower snapshot mismatch: %+v", loaded.ActiveRun.Towers)
	}

	// Test Slot refresh and metadata
	sm.RefreshSlots()
	if !sm.Slots[0].Exists || sm.Slots[0].GlitchShards != 150 {
		t.Errorf("slot metadata mismatch: %+v", sm.Slots[0])
	}

	// Test Deletion
	if err := sm.DeleteSlot(1); err != nil {
		t.Fatalf("failed to delete slot: %v", err)
	}

	sm.RefreshSlots()
	if sm.Slots[0].Exists {
		t.Errorf("expected slot 1 to be deleted")
	}
}
