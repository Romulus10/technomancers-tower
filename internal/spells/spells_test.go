package spells_test

import (
	"image/color"
	"testing"

	"technomancers-tower/internal/data"
	"technomancers-tower/internal/economy"
	"technomancers-tower/internal/malware"
	"technomancers-tower/internal/spells"
)

func TestSpellCastingAndCooldowns(t *testing.T) {
	reg := data.NewRegistry()
	sm := spells.NewSpellManager(reg)
	spawner := malware.NewSpawner(reg)
	run := economy.NewRunState(0, 0, 0, 0, 0, 0, 0, 0, 0)
	run.Mana = 100

	// Chain lightning is unlocked by default
	enemy := &malware.Enemy{
		ID:        1,
		Def:       &data.EnemyDef{ID: "swarmer", Color: color.RGBA{R: 255, G: 0, B: 0, A: 255}},
		X:         200,
		Y:         200,
		HP:        100,
		MaxHP:     100,
		BaseSpeed: 80,
	}
	enemies := []*malware.Enemy{enemy}

	// 1. Cast Chain Lightning (Cost 25, Dmg 60)
	success := sm.CastSpell("chain_lightning", 200, 200, enemies, run, spawner)
	if !success || run.Mana != 75 || enemy.HP != 40 {
		t.Errorf("failed cast: success=%v, mana=%.1f, enemyHP=%.1f", success, run.Mana, enemy.HP)
	}

	// 2. Cannot cast again immediately due to cooldown
	success = sm.CastSpell("chain_lightning", 200, 200, enemies, run, spawner)
	if success {
		t.Errorf("should not be able to cast while on cooldown")
	}

	// 3. Fast forward cooldown
	sm.Update(10.0, run)
	if sm.Cooldowns["chain_lightning"] > 0 {
		t.Errorf("expected cooldown to expire after 10s")
	}

	// 4. Insufficient Mana check
	run.Mana = 10
	success = sm.CastSpell("chain_lightning", 200, 200, enemies, run, spawner)
	if success {
		t.Errorf("should not cast with insufficient mana")
	}
}
