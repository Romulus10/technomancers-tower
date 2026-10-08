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

	if len(metaMgr.Talents) != 25 {
		t.Errorf("expected 25 talents across 5 categories, got %d", len(metaMgr.Talents))
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

	// Next rank cost should scale (10 * 1.5 = 15)
	nextCost := talent.CurrentCost()
	if nextCost != 15 {
		t.Errorf("expected rank 2 cost 15, got %d", nextCost)
	}

	// Buy Rank 2
	bought = metaMgr.BuyTalent("kernel_shield")
	if !bought || metaMgr.GetTalentLevel("kernel_shield") != 2 || metaMgr.GlitchShards != 75 {
		t.Errorf("failed to buy rank 2: level=%d, shards=%d", metaMgr.GetTalentLevel("kernel_shield"), metaMgr.GlitchShards)
	}

	// Refund Rank 2 (should refund 15 shards)
	refunded := metaMgr.RefundTalent("kernel_shield")
	if !refunded || metaMgr.GetTalentLevel("kernel_shield") != 1 || metaMgr.GlitchShards != 90 {
		t.Errorf("failed to refund rank 2: level=%d, shards=%d", metaMgr.GetTalentLevel("kernel_shield"), metaMgr.GlitchShards)
	}

	// Refund Rank 1 (should refund 10 shards)
	refunded = metaMgr.RefundTalent("kernel_shield")
	if !refunded || metaMgr.GetTalentLevel("kernel_shield") != 0 || metaMgr.GlitchShards != 100 {
		t.Errorf("failed to refund rank 1: level=%d, shards=%d", metaMgr.GetTalentLevel("kernel_shield"), metaMgr.GlitchShards)
	}

	// Cannot refund rank 0
	if metaMgr.RefundTalent("kernel_shield") {
		t.Errorf("should not be able to refund rank 0")
	}

	// Buy multiple talents and test full ResetAllTalents
	metaMgr.BuyTalent("kernel_shield")   // 10
	metaMgr.BuyTalent("kernel_shield")   // 15
	metaMgr.BuyTalent("overclock_nodes") // 25
	metaMgr.BuyTalent("boot_bytes")      // 15
	if metaMgr.GlitchShards != 35 {
		t.Errorf("expected 35 shards left, got %d", metaMgr.GlitchShards)
	}

	refundedAmount := metaMgr.ResetAllTalents()
	if refundedAmount != 65 || metaMgr.GlitchShards != 100 {
		t.Errorf("expected 65 refunded to reach 100 shards, got refunded=%d, shards=%d", refundedAmount, metaMgr.GlitchShards)
	}

	for _, tName := range []string{"kernel_shield", "overclock_nodes", "boot_bytes"} {
		if metaMgr.GetTalentLevel(tName) != 0 {
			t.Errorf("expected talent %s to be reset to 0", tName)
		}
	}
}
