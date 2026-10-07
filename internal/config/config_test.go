package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"technomancers-tower/internal/config"
)

func TestGetVersion(t *testing.T) {
	// 1. Test Default fallback
	config.Version = ""
	_ = os.Unsetenv("GAME_VERSION")
	_ = os.Unsetenv("APP_VERSION")

	if ver := config.GetVersion(); ver != config.DefaultVersion {
		t.Errorf("expected default version %s, got %s", config.DefaultVersion, ver)
	}

	// 2. Test Environment variable override
	_ = os.Setenv("GAME_VERSION", "v1.5.0-env")
	if ver := config.GetVersion(); ver != "v1.5.0-env" {
		t.Errorf("expected env version v1.5.0-env, got %s", ver)
	}
	_ = os.Unsetenv("GAME_VERSION")

	// 3. Test Compile-time injection override
	config.Version = "v2.0.0-injected"
	if ver := config.GetVersion(); ver != "v2.0.0-injected" {
		t.Errorf("expected injected version v2.0.0-injected, got %s", ver)
	}
	config.Version = ""
}

func TestSettingsManager(t *testing.T) {
	tmpDir := t.TempDir()
	savePath := filepath.Join(tmpDir, "settings.json")

	sm := &config.SettingsManager{
		SavePath:     savePath,
		WindowWidth:  800,
		WindowHeight: 600,
		IsFullscreen: false,
		VSyncEnabled: true,
	}

	sm.SetResolution(1200, 900)
	if sm.WindowWidth != 1200 || sm.WindowHeight != 900 {
		t.Errorf("expected 1200x900, got %dx%d", sm.WindowWidth, sm.WindowHeight)
	}

	// Verify persistence file exists and loads
	sm2 := &config.SettingsManager{
		SavePath: savePath,
	}
	sm2.Load()

	if sm2.WindowWidth != 1200 || sm2.WindowHeight != 900 {
		t.Errorf("expected loaded resolution 1200x900, got %dx%d", sm2.WindowWidth, sm2.WindowHeight)
	}
}
