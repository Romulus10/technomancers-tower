package particles

import (
	"image/color"
	"testing"
)

func TestParticleManager(t *testing.T) {
	pm := NewParticleManager()
	if pm == nil {
		t.Fatalf("expected non-nil ParticleManager")
	}

	// 1. Emit hit sparks
	pm.EmitHitSparks(100, 100, 10, color.RGBA{R: 255, G: 200, B: 50, A: 255})
	if len(pm.Particles) != 10 {
		t.Errorf("expected 10 particles, got %d", len(pm.Particles))
	}

	// 2. Emit explosion
	pm.EmitExplosion(200, 200, 30, color.RGBA{R: 255, G: 50, B: 50, A: 255})
	if len(pm.Particles) <= 10 {
		t.Errorf("expected additional explosion particles")
	}

	// 3. Combat Text
	pm.AddDamageText(150, 200, 200, true)
	if len(pm.CombatTexts) != 1 {
		t.Errorf("expected 1 combat text, got %d", len(pm.CombatTexts))
	}

	// 4. Update simulation & lifecycle
	pm.Update(0.1)
	if pm.CombatTexts[0].Y >= 200 {
		t.Errorf("expected combat text to float upward, Y: %f", pm.CombatTexts[0].Y)
	}

	// Expire all particles and texts
	pm.Update(2.0)
	if len(pm.Particles) != 0 {
		t.Errorf("expected particles to expire after 2s, remaining: %d", len(pm.Particles))
	}
	if len(pm.CombatTexts) != 0 {
		t.Errorf("expected combat texts to expire after 2s, remaining: %d", len(pm.CombatTexts))
	}
}

func TestScreenShake(t *testing.T) {
	shake := NewScreenShake()
	shake.AddTrauma(0.8)

	if shake.Trauma != 0.8 {
		t.Errorf("expected trauma 0.8, got %f", shake.Trauma)
	}

	ox, oy := shake.GetOffset()
	if ox == 0 && oy == 0 {
		t.Errorf("expected non-zero offset with active trauma")
	}

	shake.Update(1.0)
	if shake.Trauma >= 0.8 {
		t.Errorf("expected trauma decay over time")
	}
}
