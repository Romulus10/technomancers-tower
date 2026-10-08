package gfx

import (
	"image/color"
	"testing"
)

func TestSpriteCacheAndGenerators(t *testing.T) {
	cache := GetCache()
	if cache == nil {
		t.Fatalf("expected non-nil SpriteCache")
	}

	// 1. Test CPU Core & Spawner Port sprites
	core := cache.GetCPUCore()
	if core == nil {
		t.Fatalf("expected non-nil CPU Core sprite")
	}
	port := cache.GetSpawnerPort()
	if port == nil {
		t.Fatalf("expected non-nil Spawner Port sprite")
	}

	// 2. Test Procedural Tower Sprite generation across tiers
	towerIDs := []string{
		"bit_driver", "mana_siphon", "crypto_miner", "quantum_miner",
		"tesla_bus", "cryo_cache", "logic_mortar", "emp_railgun",
		"plasma_flak", "overclock_beacon", "singularity_vortex", "bit_leech",
	}

	for _, id := range towerIDs {
		for tier := 1; tier <= 5; tier++ {
			sprite := cache.GetTowerSprite(id, tier, color.RGBA{R: 0, G: 220, B: 255, A: 255})
			if sprite == nil {
				t.Errorf("failed to generate sprite for tower %s at tier %d", id, tier)
			}
		}
	}

	// 3. Test Procedural Enemy Sprite generation
	enemyIDs := []string{
		"swarmer", "sprite", "golem", "trojan_shard",
		"healer_worm", "phantom_glitch", "buffer_overflower",
		"boss_daemon", "boss_zeroday",
	}

	for _, id := range enemyIDs {
		isBoss := id == "boss_daemon" || id == "boss_zeroday"
		sprite := cache.GetEnemySprite(id, isBoss, color.RGBA{R: 255, G: 50, B: 80, A: 255}, 10.0)
		if sprite == nil {
			t.Errorf("failed to generate sprite for enemy %s", id)
		}
	}

	// 4. Test Icon Badges
	badge := cache.GetIconBadge("bit_driver", true, color.RGBA{R: 0, G: 220, B: 255, A: 255})
	if badge == nil {
		t.Fatalf("expected non-nil icon badge")
	}
}

func TestColorMath(t *testing.T) {
	c1 := color.RGBA{R: 0, G: 0, B: 0, A: 255}
	c2 := color.RGBA{R: 100, G: 100, B: 100, A: 255}

	lerped := LerpColor(c1, c2, 0.5)
	if lerped.R != 50 || lerped.G != 50 || lerped.B != 50 {
		t.Errorf("unexpected LerpColor result: %+v", lerped)
	}

	faded := FadeAlpha(c2, 0.5)
	if faded.A != 127 {
		t.Errorf("unexpected FadeAlpha result: %d", faded.A)
	}

	brightened := Brighten(c1, 0.5)
	if brightened.R != 127 {
		t.Errorf("unexpected Brighten result: %+v", brightened)
	}

	darkened := Darken(c2, 0.5)
	if darkened.R != 50 {
		t.Errorf("unexpected Darken result: %+v", darkened)
	}
}
