package config

import (
	"os"
	"strings"
)

// Injected via CI compile flags: -ldflags "-X technomancers-tower/internal/config.Version=vX.Y.Z"
var Version = ""

const DefaultVersion = "v0.1.0-dev"

// GetVersion returns the effective version, checking ldflags injection first,
// then environment variables (GAME_VERSION / APP_VERSION), and lastly the default fallback.
func GetVersion() string {
	if strings.TrimSpace(Version) != "" {
		return strings.TrimSpace(Version)
	}

	if envVer := os.Getenv("GAME_VERSION"); strings.TrimSpace(envVer) != "" {
		return strings.TrimSpace(envVer)
	}

	if envVer := os.Getenv("APP_VERSION"); strings.TrimSpace(envVer) != "" {
		return strings.TrimSpace(envVer)
	}

	return DefaultVersion
}
