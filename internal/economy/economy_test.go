package economy_test

import (
	"testing"

	"technomancers-tower/internal/economy"
)

func TestNewRunStateWithMetaModifiers(t *testing.T) {
	// metaShieldLvl=2, metaBootLvl=3, metaManaLvl=1, metaSpeedLvl=2, metaSpellLvl=1, metaScrapLvl=2, metaDamageLvl=3, metaRangeLvl=1, metaNaniteLvl=0
	run := economy.NewRunState(2, 3, 1, 2, 1, 2, 3, 1, 0)

	// MaxHP: 100 + 2*25 = 150
	if run.MaxKernelHP != 150.0 || run.KernelHP != 150.0 {
		t.Errorf("expected HP 150, got %.1f/%.1f", run.KernelHP, run.MaxKernelHP)
	}

	// Starting Bytes: 120 + 3*40 = 240
	if run.Bytes != 240.0 || run.TotalBytes != 240.0 {
		t.Errorf("expected Bytes 240, got %.1f", run.Bytes)
	}

	// Modifiers
	if run.ManaGenMult != 1.20 {
		t.Errorf("expected ManaGenMult 1.20, got %.2f", run.ManaGenMult)
	}
	if run.TowerDamageMult != 1.18 {
		t.Errorf("expected TowerDamageMult 1.18, got %.2f", run.TowerDamageMult)
	}
}

func TestRunStateTransactionsAndEconomy(t *testing.T) {
	run := economy.NewRunState(0, 0, 0, 0, 0, 0, 0, 0, 0)

	// Bytes transactions
	if !run.CanAffordBytes(50) {
		t.Errorf("expected to afford 50 bytes with 120 starting")
	}
	run.SpendBytes(50)
	if run.Bytes != 70 {
		t.Errorf("expected 70 bytes left, got %.1f", run.Bytes)
	}

	// Mana transactions
	run.Mana = 40
	if !run.CanAffordMana(30) {
		t.Errorf("expected to afford 30 mana with 40 available")
	}
	run.SpendMana(30)
	if run.Mana != 10 {
		t.Errorf("expected 10 mana left, got %.1f", run.Mana)
	}

	// Passive mana generation over time
	run.BaseManaRate = 10.0 // 10 mana/sec
	run.ManaGenMult = 1.5
	run.Update(1.0) // 1 second elapsed
	if run.Mana != 25.0 {
		t.Errorf("expected 25.0 mana after 1s generation, got %.1f", run.Mana)
	}

	// Kernel damage and game over trigger
	gameOver := run.DamageKernel(40)
	if gameOver || run.KernelHP != 60 {
		t.Errorf("expected 60 HP remaining without game over, got %.1f", run.KernelHP)
	}
	gameOver = run.DamageKernel(70)
	if !gameOver || run.KernelHP != 0 {
		t.Errorf("expected game over on lethal damage, got gameOver=%v, HP=%.1f", gameOver, run.KernelHP)
	}
}

func TestGlitchShardCalculation(t *testing.T) {
	run := economy.NewRunState(0, 0, 0, 0, 0, 0, 0, 0, 0)
	run.RunTime = 120.0 // 2 minutes (120/10 = 12 shards)
	run.Kills = 100     // 100 kills (100/20 = 5 shards)
	run.BossKills = 2   // 2 bosses (2*15 = 30 shards)

	shards := run.CalculateFinalShards()
	expected := 12 + 5 + 30
	if shards != expected {
		t.Errorf("expected %d shards, got %d", expected, shards)
	}
}

func TestRunStateLevelUpAndPendingDrafts(t *testing.T) {
	run := economy.NewRunState(0, 0, 0, 0, 0, 0, 0, 0, 0)
	if run.Level != 1 || run.PendingDrafts != 0 || run.TargetXP != 40.0 {
		t.Fatalf("unexpected initial run state: level=%d, drafts=%d, targetXP=%.1f", run.Level, run.PendingDrafts, run.TargetXP)
	}

	// Add 40 XP -> triggers 1 level up & queues 1 draft
	leveledUp := run.AddKill(10, 40, false)
	if !leveledUp || run.Level != 2 || run.PendingDrafts != 1 || run.CurrentXP != 0 {
		t.Errorf("expected level 2 with 1 pending draft, got level=%d, drafts=%d, xp=%.1f", run.Level, run.PendingDrafts, run.CurrentXP)
	}

	// Add enough XP for 2 consecutive levels (TargetXP is ~54, then ~73)
	// 54 + 73 + 10 = 137 XP -> should level up to 4 with 2 additional pending drafts
	leveledUp = run.AddXP(137.0)
	if !leveledUp || run.Level != 4 || run.PendingDrafts != 3 {
		t.Errorf("expected level 4 with 3 pending drafts, got level=%d, drafts=%d", run.Level, run.PendingDrafts)
	}
}
