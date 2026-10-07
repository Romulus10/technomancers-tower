package input_test

import (
	"os"
	"path/filepath"
	"testing"

	"technomancers-tower/internal/input"

	"github.com/hajimehoshi/ebiten/v2"
)

func TestKeybindManagerDefaultsAndRebinds(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "keybinds_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	km := input.NewKeybindManager()
	km.SavePath = filepath.Join(tempDir, "keybinds.json")
	km.ResetDefaults()

	if len(km.TowerKeys) != 5 {
		t.Errorf("expected 5 tower keys, got %d", len(km.TowerKeys))
	}
	if len(km.SpellKeys) != 5 {
		t.Errorf("expected 5 spell keys, got %d", len(km.SpellKeys))
	}

	// Verify default bindings
	expectedTowers := [5]ebiten.Key{ebiten.Key1, ebiten.Key2, ebiten.Key3, ebiten.Key4, ebiten.Key5}
	for i, k := range expectedTowers {
		if km.TowerKeys[i] != k {
			t.Errorf("expected tower slot %d to be %v, got %v", i+1, k, km.TowerKeys[i])
		}
	}

	expectedSpells := [5]ebiten.Key{ebiten.KeyQ, ebiten.KeyW, ebiten.KeyE, ebiten.KeyR, ebiten.KeyF}
	for i, k := range expectedSpells {
		if km.SpellKeys[i] != k {
			t.Errorf("expected spell slot %d to be %v, got %v", i+1, k, km.SpellKeys[i])
		}
	}

	// Save and reload
	km.TowerKeys[0] = ebiten.KeyZ
	km.SpellKeys[4] = ebiten.KeySpace
	if err := km.Save(); err != nil {
		t.Fatalf("failed to save keybinds: %v", err)
	}

	km2 := input.NewKeybindManager()
	km2.SavePath = filepath.Join(tempDir, "keybinds.json")
	km2.Load()

	if km2.TowerKeys[0] != ebiten.KeyZ {
		t.Errorf("expected reloaded tower slot 0 to be KeyZ, got %v", km2.TowerKeys[0])
	}
	if km2.SpellKeys[4] != ebiten.KeySpace {
		t.Errorf("expected reloaded spell slot 4 to be KeySpace, got %v", km2.SpellKeys[4])
	}
}
