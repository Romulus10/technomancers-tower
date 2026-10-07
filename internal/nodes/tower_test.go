package nodes_test

import (
	"image/color"
	"math"
	"testing"

	"technomancers-tower/internal/data"
	"technomancers-tower/internal/economy"
	"technomancers-tower/internal/malware"
	"technomancers-tower/internal/motherboard"
	"technomancers-tower/internal/nodes"
)

func TestTowerBuildingAndManaRecalculation(t *testing.T) {
	reg := data.NewRegistry()
	tm := nodes.NewTowerManager(reg)
	grid := motherboard.NewGrid()
	run := economy.NewRunState(0, 0, 0, 0, 0, 0, 0, 0, 0) // 120 starting bytes

	// 1. Build Bit Driver (Cost 30)
	success := tm.BuildTower("bit_driver", 2, 2, grid, run)
	if !success || len(tm.Towers) != 1 || run.Bytes != 90 {
		t.Errorf("failed to build Bit Driver: success=%v, bytes=%.1f", success, run.Bytes)
	}

	// 2. Build Mana Siphon (Cost 45, ManaRate 4.0)
	success = tm.BuildTower("mana_siphon", 3, 3, grid, run)
	if !success || len(tm.Towers) != 2 || run.Bytes != 45 {
		t.Errorf("failed to build Mana Siphon: success=%v, bytes=%.1f", success, run.Bytes)
	}

	if run.BaseManaRate != 4.0 {
		t.Errorf("expected BaseManaRate 4.0, got %.1f", run.BaseManaRate)
	}

	// 3. Build Tesla Bus (Cost 75) - Cannot afford with 45 bytes remaining
	success = tm.BuildTower("tesla_bus", 4, 4, grid, run)
	if success {
		t.Errorf("should not afford tesla_bus with 45 bytes")
	}
}

func TestTowerAttacksAndOnHitPipeline(t *testing.T) {
	reg := data.NewRegistry()
	tm := nodes.NewTowerManager(reg)
	spawner := malware.NewSpawner(reg)
	grid := motherboard.NewGrid()
	run := economy.NewRunState(0, 0, 0, 0, 0, 0, 0, 0, 0)
	run.Bytes = 1000

	built := tm.BuildTower("bit_driver", 1, 1, grid, run)
	if !built || len(tm.Towers) != 1 {
		t.Fatalf("failed to build bit_driver at (1, 1)")
	}

	// Enemy located within range (range is 100)
	twx, twy := grid.GridToScreenCenter(1, 1)
	enemy := &malware.Enemy{
		ID:        1,
		Def:       &data.EnemyDef{ID: "swarmer", Color: color.RGBA{R: 255, G: 0, B: 0, A: 255}},
		X:         twx + 60,
		Y:         twy,
		HP:        50,
		MaxHP:     50,
		BaseSpeed: 80,
	}
	enemies := []*malware.Enemy{enemy}

	// Step tower update to trigger attack (small dt so projectile stays in flight)
	tm.Towers[0].Timer = 0
	tm.Update(0.01, enemies, run, spawner)

	if len(tm.Projectiles) != 1 {
		t.Fatalf("expected 1 projectile spawned, got %d", len(tm.Projectiles))
	}

	// Step update with fast projectile to hit enemy
	tm.Projectiles[0].Speed = 10000
	tm.Update(0.1, enemies, run, spawner)

	if enemy.HP >= 50 {
		t.Errorf("expected enemy to take damage from projectile hit, hp=%.1f", enemy.HP)
	}
	if len(tm.Projectiles) != 0 {
		t.Errorf("expected projectile to be consumed on hit, remaining=%d", len(tm.Projectiles))
	}

	// Verify tower gained XP from combat
	if tm.Towers[0].XP <= 0 {
		t.Errorf("expected tower to gain combat XP, got %.1f", tm.Towers[0].XP)
	}
}

func TestTowerLevelProgressionAndStatScaling(t *testing.T) {
	reg := data.NewRegistry()
	tm := nodes.NewTowerManager(reg)
	grid := motherboard.NewGrid()
	run := economy.NewRunState(0, 0, 0, 0, 0, 0, 0, 0, 0)
	run.Bytes = 1000

	// 1. Test Combat Tower Leveling
	tm.BuildTower("bit_driver", 1, 1, grid, run)
	tower := tm.Towers[0]

	if tower.Level != 1 || tower.GetDamageMult() != 1.0 || tower.GetSpeedMult() != 1.0 || tower.GetRangeMult() != 1.0 {
		t.Fatalf("unexpected initial stats at Level 1")
	}

	// Add 100 XP -> Level 2
	leveledUp := tower.AddXP(100.0)
	if !leveledUp || tower.Level != 2 {
		t.Errorf("expected level up to 2, got leveledUp=%v, level=%d", leveledUp, tower.Level)
	}
	if tower.GetDamageMult() != 1.10 || tower.GetSpeedMult() != 1.05 || tower.GetRangeMult() != 1.03 {
		t.Errorf("unexpected stats at Level 2: dmg=%.2f, spd=%.2f, range=%.2f", tower.GetDamageMult(), tower.GetSpeedMult(), tower.GetRangeMult())
	}

	// Add XP to reach Max Level 5 (250 + 550 + 1100 = 1900 XP)
	tower.AddXP(1900.0)
	if tower.Level != nodes.MaxTowerLevel {
		t.Errorf("expected Max Level %d, got %d", nodes.MaxTowerLevel, tower.Level)
	}
	if tower.GetDamageMult() != 1.40 || tower.GetSpeedMult() != 1.20 || tower.GetRangeMult() != 1.12 {
		t.Errorf("unexpected stats at Level 5: dmg=%.2f, spd=%.2f, range=%.2f", tower.GetDamageMult(), tower.GetSpeedMult(), tower.GetRangeMult())
	}

	// Further XP should be ignored at max level
	tower.AddXP(500.0)
	if tower.Level != nodes.MaxTowerLevel || tower.XP != 0 {
		t.Errorf("expected level to stay capped at %d with 0 excess XP, got level=%d, xp=%.1f", nodes.MaxTowerLevel, tower.Level, tower.XP)
	}

	// 2. Test Passive Mana Siphon Leveling & Mana Rate Recalculation
	tm.BuildTower("mana_siphon", 3, 3, grid, run)
	siphon := tm.Towers[1]
	if siphon.GetManaRate() != 4.0 || run.BaseManaRate != 4.0 {
		t.Errorf("expected initial BaseManaRate 4.0, got siphon=%.1f, run=%.1f", siphon.GetManaRate(), run.BaseManaRate)
	}

	// Level up Mana Siphon to Level 2
	siphon.AddXP(100.0)
	tm.RecalculateBaseManaRate(run)

	if siphon.Level != 2 || siphon.GetManaRate() != 4.4 || run.BaseManaRate != 4.4 {
		t.Errorf("expected Level 2 siphon with ManaRate 4.4, got level=%d, siphon=%.1f, run=%.1f", siphon.Level, siphon.GetManaRate(), run.BaseManaRate)
	}

	// 3. Test Passive Crypto Miner Building, Leveling & Byte Rate Recalculation
	tm.BuildTower("crypto_miner", 4, 4, grid, run)
	miner := tm.Towers[2]
	if miner.GetByteRate() != 2.0 || run.BaseByteRate != 2.0 {
		t.Errorf("expected initial BaseByteRate 2.0, got miner=%.1f, run=%.1f", miner.GetByteRate(), run.BaseByteRate)
	}

	// Level up Crypto Miner to Level 2
	miner.AddXP(100.0)
	tm.RecalculateBaseManaRate(run)

	if miner.Level != 2 || miner.GetByteRate() != 2.2 || run.BaseByteRate != 2.2 {
		t.Errorf("expected Level 2 miner with ByteRate 2.2, got level=%d, miner=%.1f, run=%.1f", miner.Level, miner.GetByteRate(), run.BaseByteRate)
	}
}

func TestTowerOverclockPromotionAndTierScaling(t *testing.T) {
	reg := data.NewRegistry()
	tm := nodes.NewTowerManager(reg)
	grid := motherboard.NewGrid()
	run := economy.NewRunState(0, 0, 0, 0, 0, 0, 0, 0, 0)
	run.Bytes = 500

	tm.BuildTower("bit_driver", 2, 2, grid, run)
	tower := tm.Towers[0]

	// 1. Initial State (Tier 1, Level 1)
	if tower.Tier != 1 || tower.Level != 1 {
		t.Fatalf("expected Tier 1 Level 1, got T%d L%d", tower.Tier, tower.Level)
	}

	// Cannot promote at Level 1
	if tower.CanPromote(run) || tm.PromoteTowerAt(2, 2, run) {
		t.Errorf("should not be able to promote tower before reaching Max Level")
	}

	// 2. Level tower to Max Level 5
	tower.AddXP(100.0 + 250.0 + 550.0 + 1100.0)
	if tower.Level != nodes.MaxTowerLevel {
		t.Fatalf("expected tower at Max Level %d, got %d", nodes.MaxTowerLevel, tower.Level)
	}

	// Upgrade cost for Tier 1 -> Tier 2 on Bit Driver (BaseCost 30): round(30 * 0.50 * 1) = 15
	costT1 := tower.GetUpgradeCost()
	if costT1 != 15 {
		t.Errorf("expected upgrade cost 15, got %.1f", costT1)
	}

	// Check cannot promote if player doesn't have enough bytes
	run.Bytes = 10
	if tower.CanPromote(run) || tm.PromoteTowerAt(2, 2, run) {
		t.Errorf("should not promote if player cannot afford cost")
	}

	// 3. Promote to Tier 2
	run.Bytes = 100
	promoted := tm.PromoteTowerAt(2, 2, run)
	if !promoted {
		t.Fatalf("expected successful promotion to Tier 2")
	}
	if run.Bytes != 85 {
		t.Errorf("expected 85 bytes remaining after 15 cost, got %.1f", run.Bytes)
	}
	if tower.Tier != 2 || tower.Level != 1 || tower.XP != 0 || tower.TargetXP != nodes.NextLevelXP(1) {
		t.Errorf("invalid state after promotion: Tier=%d, Level=%d, XP=%.1f, TargetXP=%.1f", tower.Tier, tower.Level, tower.XP, tower.TargetXP)
	}

	// Tier 2 Level 1 Multipliers:
	// Damage: (1 + 1*0.35) * (1 + 0) = 1.35
	// Speed:  (1 + 1*0.20) * (1 + 0) = 1.20
	// Range:  1 + 1*0.05 + 0 = 1.05
	if math.Abs(tower.GetDamageMult()-1.35) > 0.001 || math.Abs(tower.GetSpeedMult()-1.20) > 0.001 || math.Abs(tower.GetRangeMult()-1.05) > 0.001 {
		t.Errorf("unexpected Tier 2 Level 1 stats: dmg=%.3f, spd=%.3f, range=%.3f", tower.GetDamageMult(), tower.GetSpeedMult(), tower.GetRangeMult())
	}

	// 4. Level up in Tier 2 to Max Level 5
	tower.AddXP(100.0 + 250.0 + 550.0 + 1100.0)
	if tower.Level != nodes.MaxTowerLevel {
		t.Fatalf("expected Level 5 in Tier 2, got %d", tower.Level)
	}

	// Tier 2 Level 5 Multipliers:
	// Damage: (1 + 0.35) * (1 + 0.40) = 1.35 * 1.40 = 1.89
	// Speed:  (1 + 0.20) * (1 + 0.20) = 1.20 * 1.20 = 1.44
	// Range:  1 + 0.05 + 4*0.03 = 1.17
	if math.Abs(tower.GetDamageMult()-1.89) > 0.001 || math.Abs(tower.GetSpeedMult()-1.44) > 0.001 || math.Abs(tower.GetRangeMult()-1.17) > 0.001 {
		t.Errorf("unexpected Tier 2 Level 5 stats: dmg=%.3f, spd=%.3f, range=%.3f", tower.GetDamageMult(), tower.GetSpeedMult(), tower.GetRangeMult())
	}

	// 5. Tier 2 -> Tier 3 Upgrade Cost: round(30 * 0.50 * 2) = 30
	costT2 := tower.GetUpgradeCost()
	if costT2 != 30 {
		t.Errorf("expected Tier 2 upgrade cost 30, got %.1f", costT2)
	}

	promoted = tm.PromoteTowerAt(2, 2, run)
	if !promoted || tower.Tier != 3 || tower.Level != 1 {
		t.Fatalf("expected promotion to Tier 3, got promoted=%v, Tier=%d, Level=%d", promoted, tower.Tier, tower.Level)
	}
}

func TestTowerPromoteAllBatchUpgrade(t *testing.T) {
	reg := data.NewRegistry()
	tm := nodes.NewTowerManager(reg)
	grid := motherboard.NewGrid()
	run := economy.NewRunState(0, 0, 0, 0, 0, 0, 0, 0, 0)
	run.Bytes = 1000

	// Build 3 towers: 2 Bit Drivers (BaseCost 30) + 1 Tesla Bus (BaseCost 75)
	tm.BuildTower("bit_driver", 1, 1, grid, run)
	tm.BuildTower("bit_driver", 2, 2, grid, run)
	tm.BuildTower("tesla_bus", 3, 3, grid, run)

	// Level up towers 0 and 1 to Max Level 5; leave tower 2 at Level 1
	tm.Towers[0].AddXP(2000.0)
	tm.Towers[1].AddXP(2000.0)

	count, totalCost := tm.GetPromotableInfo()
	if count != 2 {
		t.Errorf("expected 2 promotable towers, got %d", count)
	}
	// Total cost = 15 + 15 = 30
	if totalCost != 30 {
		t.Errorf("expected total upgrade cost 30, got %.1f", totalCost)
	}

	run.Bytes = 50
	upgradedCount := tm.PromoteAll(run)
	if upgradedCount != 2 {
		t.Errorf("expected 2 towers batch upgraded, got %d", upgradedCount)
	}
	if run.Bytes != 20 {
		t.Errorf("expected 20 bytes remaining (50 - 30), got %.1f", run.Bytes)
	}
	if tm.Towers[0].Tier != 2 || tm.Towers[0].Level != 1 || tm.Towers[1].Tier != 2 || tm.Towers[1].Level != 1 {
		t.Errorf("unexpected tier states after PromoteAll: T0=(T%d, L%d), T1=(T%d, L%d)",
			tm.Towers[0].Tier, tm.Towers[0].Level, tm.Towers[1].Tier, tm.Towers[1].Level)
	}

	// Now none are promotable
	count, _ = tm.GetPromotableInfo()
	if count != 0 {
		t.Errorf("expected 0 promotable towers after upgrade, got %d", count)
	}

	// Threat level: T0 is Tier 2 (1 point), T1 is Tier 2 (1 point), T2 is Tier 1 (0 points) -> Total 2
	threat := tm.GetTotalOverclockTiers()
	if threat != 2 {
		t.Errorf("expected total overclock threat level 2, got %d", threat)
	}
}
