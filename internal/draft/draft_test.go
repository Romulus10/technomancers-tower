package draft_test

import (
	"testing"

	"technomancers-tower/internal/data"
	"technomancers-tower/internal/draft"
	"technomancers-tower/internal/economy"
)

func TestDraftManagerGenerationAndFiltering(t *testing.T) {
	reg := data.NewRegistry()
	dm := draft.NewDraftManager(reg)
	run := economy.NewRunState(nil)

	dm.GenerateDraft(run)
	if len(dm.OfferedCards) != 3 {
		t.Fatalf("expected 3 offered draft cards, got %d", len(dm.OfferedCards))
	}

	// Select first card and apply modifier
	selected := dm.OfferedCards[0]
	dm.ApplyCard(selected, run)

	// Verify modifier applied
	for _, eff := range selected.Effects {
		switch eff.Type {
		case data.ModUnlockTower:
			tower := reg.GetTower(eff.TargetID)
			if tower == nil || !tower.Unlocked {
				t.Errorf("expected tower %s to be unlocked", eff.TargetID)
			}
		case data.ModUnlockSpell:
			spell := reg.GetSpell(eff.TargetID)
			if spell == nil || !spell.Unlocked {
				t.Errorf("expected spell %s to be unlocked", eff.TargetID)
			}
		case data.ModTowerDamageMult:
			if run.TowerDamageMult <= 1.0 {
				t.Errorf("expected TowerDamageMult to increase")
			}
		}
	}
}

func TestDraftExtendedChoicesAndReroll(t *testing.T) {
	reg := data.NewRegistry()
	dm := draft.NewDraftManager(reg)
	talents := economy.MapTalents{
		"extra_draft_slot": 1, // 4 choices
		"draft_rerolls":    2, // 2 rerolls
	}
	run := economy.NewRunState(talents)

	dm.GenerateDraft(run)
	if len(dm.OfferedCards) != 4 {
		t.Fatalf("expected 4 offered draft cards, got %d", len(dm.OfferedCards))
	}

	// Reroll draft
	rerolled := dm.Reroll(run)
	if !rerolled || run.DraftRerolls != 1 || len(dm.OfferedCards) != 4 {
		t.Errorf("expected successful reroll with 1 reroll left and 4 cards, got rerolled=%v, rerolls=%d, cards=%d", rerolled, run.DraftRerolls, len(dm.OfferedCards))
	}
}

func TestWrapText(t *testing.T) {
	longText := "Unlocks long-range artillery node dealing high explosive splash AoE."
	lines := draft.WrapText(longText, 26)

	if len(lines) <= 1 {
		t.Fatalf("expected multiple wrapped lines, got %d", len(lines))
	}

	for i, l := range lines {
		if len(l) > 26 {
			t.Errorf("line %d exceeds max length 26: %q (len=%d)", i, l, len(l))
		}
	}
}

func TestQuantumMinerCardUnlock(t *testing.T) {
	reg := data.NewRegistry()
	dm := draft.NewDraftManager(reg)
	run := economy.NewRunState(nil)

	qMiner := reg.GetTower("quantum_miner")
	if qMiner == nil || qMiner.Unlocked {
		t.Fatalf("expected quantum_miner to be locked by default")
	}

	card := reg.GetCard("unlock_quantum_miner")
	if card == nil {
		t.Fatalf("expected unlock_quantum_miner card to exist in registry")
	}

	dm.ApplyCard(card, run)
	if !qMiner.Unlocked {
		t.Errorf("expected quantum_miner to be unlocked after applying card")
	}
}
