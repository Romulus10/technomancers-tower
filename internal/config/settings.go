package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
)

type ResolutionPreset struct {
	Label  string
	Width  int
	Height int
}

var AvailableResolutions = []ResolutionPreset{
	{Label: "800 x 600 (1.0x Standard)", Width: 800, Height: 600},
	{Label: "1200 x 900 (1.5x HD)", Width: 1200, Height: 900},
	{Label: "1600 x 1200 (2.0x Retina)", Width: 1600, Height: 1200},
	{Label: "1920 x 1080 (1080p Fit)", Width: 1920, Height: 1080},
}

type SettingsData struct {
	WindowWidth  int  `json:"window_width"`
	WindowHeight int  `json:"window_height"`
	IsFullscreen bool `json:"is_fullscreen"`
	VSyncEnabled bool `json:"vsync_enabled"`
}

type SettingsManager struct {
	SavePath     string
	WindowWidth  int
	WindowHeight int
	IsFullscreen bool
	VSyncEnabled bool
}

func NewSettingsManager() *SettingsManager {
	home, err := os.UserHomeDir()
	saveDir := "."
	if err == nil {
		saveDir = filepath.Join(home, ".technomancers-tower")
		_ = os.MkdirAll(saveDir, 0755)
	}

	sm := &SettingsManager{
		SavePath:     filepath.Join(saveDir, "settings.json"),
		WindowWidth:  800,
		WindowHeight: 600,
		IsFullscreen: false,
		VSyncEnabled: true,
	}

	sm.Load()
	return sm
}

func (sm *SettingsManager) Apply() {
	ebiten.SetWindowSize(sm.WindowWidth, sm.WindowHeight)
	ebiten.SetFullscreen(sm.IsFullscreen)
	ebiten.SetVsyncEnabled(sm.VSyncEnabled)
}

func (sm *SettingsManager) SetResolution(width, height int) {
	sm.WindowWidth = width
	sm.WindowHeight = height
	if !sm.IsFullscreen {
		ebiten.SetWindowSize(width, height)
	}
	sm.Save()
}

func (sm *SettingsManager) ToggleFullscreen() {
	sm.IsFullscreen = !sm.IsFullscreen
	ebiten.SetFullscreen(sm.IsFullscreen)
	if !sm.IsFullscreen {
		ebiten.SetWindowSize(sm.WindowWidth, sm.WindowHeight)
	}
	sm.Save()
}

func (sm *SettingsManager) SetFullscreen(fs bool) {
	sm.IsFullscreen = fs
	ebiten.SetFullscreen(sm.IsFullscreen)
	if !sm.IsFullscreen {
		ebiten.SetWindowSize(sm.WindowWidth, sm.WindowHeight)
	}
	sm.Save()
}

func (sm *SettingsManager) Save() {
	data := SettingsData{
		WindowWidth:  sm.WindowWidth,
		WindowHeight: sm.WindowHeight,
		IsFullscreen: sm.IsFullscreen,
		VSyncEnabled: sm.VSyncEnabled,
	}

	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return
	}

	_ = os.WriteFile(sm.SavePath, bytes, 0644)
}

func (sm *SettingsManager) Load() {
	bytes, err := os.ReadFile(sm.SavePath)
	if err != nil {
		return
	}

	var data SettingsData
	if err := json.Unmarshal(bytes, &data); err != nil {
		return
	}

	if data.WindowWidth >= 800 {
		sm.WindowWidth = data.WindowWidth
	}
	if data.WindowHeight >= 600 {
		sm.WindowHeight = data.WindowHeight
	}
	sm.IsFullscreen = data.IsFullscreen
	sm.VSyncEnabled = data.VSyncEnabled
}
