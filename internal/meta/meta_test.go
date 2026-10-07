package meta_test

import (
	"testing"

	"technomancers-tower/internal/meta"
	"technomancers-tower/internal/save"
)

func TestMetaManagerTalentProgression(t *testing.T) {
	tmpDir := t.TempDir()
	saveMgr := &save.SaveManager{BaseDir: tmpDir}
	metaMgr := meta.NewMetaManager(saveMgr)

	slotData := &save.SlotData{
		SlotID:       1,
		GlitchShards: 100,
		Talents:      make(map[string]int),
	}
	metaMgr.ApplySlot(slotData)

	if metaMgr.GlitchShards != 100 {
		t.Errorf("expected 100 shards, got %d", metaMgr.GlitchShards)
	}

	talent := metaMgr.Talents[0] // kernel_shield (base_cost=10, mult=1.5)
	initialCost := talent.CurrentCost()
	if initialCost != 10 {
		t.Errorf("expected initial cost 10, got %d", initialCost)
	}

	// Buy Rank 1
	bought := metaMgr.BuyTalent("kernel_shield")
	if !bought || metaMgr.GetTalentLevel("kernel_shield") != 1 || metaMgr.GlitchShards != 90 {
		t.Errorf("failed to buy talent: level=%d, shards=%d", metaMgr.GetTalentLevel("kernel_shield"), metaMgr.GlitchShards)
	}

	// Next rank cost should scale
	nextCost := talent.CurrentCost()
	if nextCost != 15 {
		t.Errorf("expected rank 2 cost 15, got %d", nextCost)
	}

	// Cannot buy with insufficient shards
	metaMgr.GlitchShards = 5
	bought = metaMgr.BuyTalent("kernel_shield")
	if bought {
		t.Errorf("should not be able to buy with insufficient shards")
	}
}
