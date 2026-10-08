package economy_test

import (
	"testing"

	"technomancers-tower/internal/economy"
)

func TestNewRunStateWithMetaModifiers(t *testing.T) {
	talents := economy.MapTalents{
		"kernel_shield":        2, // +50 HP -> 150
		"boot_bytes":           3, // +120 Bytes -> 240
		"mana_conductor":       1, // +20% Mana gen -> 1.20
		"overclock_nodes":      2, // +10% Speed -> 1.10
		"spell_efficiency":     1, // -8% CD, -10% cost
		"scrap_leech":          2, // +20% Byte bounty -> 1.20
		"tower_potency":        3, // +18% Dmg -> 1.18
		"sensor_array":         1, // +6% Range -> 1.06
		"hardened_firewall":    2, // -2 dmg reduction
		"natural_mana_regen":   2, // +1.0 mana/sec
		"expanded_mana_pool":   2, // +50 max mana -> 150
		"shards_harvest":       2, // +30% shards
		"starting_level":       1, // start at lvl 2 + 1 bonus draft
		"draft_rerolls":        2, // 2 rerolls
		"extra_draft_slot":     1, // 4 draft choices
	}
	run := economy.NewRunState(talents)

	// MaxHP: 100 + 2*25 = 150
	if run.MaxKernelHP != 150.0 || run.KernelHP != 150.0 {
		t.Errorf("expected HP 150, got %.1f/%.1f", run.KernelHP, run.MaxKernelHP)
	}

	// Starting Bytes: 120 + 3*40 = 240
	if run.Bytes != 240.0 || run.TotalBytes != 240.0 {
		t.Errorf("expected Bytes 240, got %.1f", run.Bytes)
	}

	// Max Mana: 100 + 2*25 = 150
	if run.MaxMana != 150.0 {
		t.Errorf("expected MaxMana 150, got %.1f", run.MaxMana)
	}

	// Modifiers
	if run.ManaGenMult != 1.20 {
		t.Errorf("expected ManaGenMult 1.20, got %.2f", run.ManaGenMult)
	}
	if run.TowerDamageMult != 1.18 {
		t.Errorf("expected TowerDamageMult 1.18, got %.2f", run.TowerDamageMult)
	}
	if run.Level != 2 || run.PendingDrafts != 1 {
		t.Errorf("expected level 2 with 1 pending draft from starting_level, got lvl=%d, drafts=%d", run.Level, run.PendingDrafts)
	}
	if run.DraftCardChoices != 4 {
		t.Errorf("expected 4 draft card choices, got %d", run.DraftCardChoices)
	}
	if run.DraftRerolls != 2 {
		t.Errorf("expected 2 draft rerolls, got %d", run.DraftRerolls)
	}
}

func TestRunStateTransactionsAndEconomy(t *testing.T) {
	run := economy.NewRunState(nil)

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

	// Kernel damage with no armor
	gameOver := run.DamageKernel(40)
	if gameOver || run.KernelHP != 60 {
		t.Errorf("expected 60 HP remaining without game over, got %.1f", run.KernelHP)
	}
	gameOver = run.DamageKernel(70)
	if !gameOver || run.KernelHP != 0 {
		t.Errorf("expected game over on lethal damage, got gameOver=%v, HP=%.1f", gameOver, run.KernelHP)
	}
}

func TestArmorAndNaniteRegen(t *testing.T) {
	talents := economy.MapTalents{
		"hardened_firewall": 5, // -5 damage reduction
		"nanite_regen":      2, // +2 HP every 15s
	}
	run := economy.NewRunState(talents)

	// 10 damage reduced by 5 -> 5 damage taken (100 -> 95)
	run.DamageKernel(10)
	if run.KernelHP != 95.0 {
		t.Errorf("expected 95 HP after 10 - 5 armor damage, got %.1f", run.KernelHP)
	}

	// 15 seconds elapsed -> Nanite regen triggers +2 HP (95 -> 97)
	run.Update(15.0)
	if run.KernelHP != 97.0 {
		t.Errorf("expected 97 HP after nanite regen, got %.1f", run.KernelHP)
	}
}

func TestGlitchShardCalculation(t *testing.T) {
	talents := economy.MapTalents{
		"shards_harvest": 2, // +30% shards
	}
	run := economy.NewRunState(talents)
	run.RunTime = 120.0 // 2 minutes (120/10 = 12 shards)
	run.Kills = 100     // 100 kills (100/20 = 5 shards)
	run.BossKills = 2   // 2 bosses (2*15 = 30 shards)

	shards := run.CalculateFinalShards()
	// Base = 12 + 5 + 30 = 47. 47 * 1.30 = 61.1 -> 61
	if shards != 61 {
		t.Errorf("expected 61 shards with +30%% harvest talent, got %d", shards)
	}
}

func TestRunStateLevelUpAndPendingDrafts(t *testing.T) {
	run := economy.NewRunState(nil)
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
