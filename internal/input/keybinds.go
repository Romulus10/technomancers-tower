package input

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const MaxToolbarSlots = 5

type KeybindsData struct {
	TowerKeys []string `json:"tower_keys"`
	SpellKeys []string `json:"spell_keys"`
	CodexKey  string   `json:"codex_key"`
	PauseKey  string   `json:"pause_key"`
}

type KeybindManager struct {
	SavePath     string
	TowerKeys    [MaxToolbarSlots]ebiten.Key
	SpellKeys    [MaxToolbarSlots]ebiten.Key
	CodexKey     ebiten.Key
	PauseKey     ebiten.Key
	IsRebinding  bool
	RebindTarget string // "tower_0", "spell_1", "codex", etc.
	RebindPrompt string
}

func NewKeybindManager() *KeybindManager {
	home, err := os.UserHomeDir()
	saveDir := "."
	if err == nil {
		saveDir = filepath.Join(home, ".technomancers-tower")
		_ = os.MkdirAll(saveDir, 0755)
	}

	km := &KeybindManager{
		SavePath: filepath.Join(saveDir, "keybinds.json"),
	}

	km.ResetDefaults()
	km.Load()
	return km
}

func (km *KeybindManager) ResetDefaults() {
	km.TowerKeys = [MaxToolbarSlots]ebiten.Key{
		ebiten.Key1, ebiten.Key2, ebiten.Key3, ebiten.Key4, ebiten.Key5,
	}
	km.SpellKeys = [MaxToolbarSlots]ebiten.Key{
		ebiten.KeyQ, ebiten.KeyW, ebiten.KeyE, ebiten.KeyR, ebiten.KeyF,
	}
	km.CodexKey = ebiten.KeyTab
	km.PauseKey = ebiten.KeyEscape
	km.IsRebinding = false
}

func (km *KeybindManager) StartRebind(target, prompt string) {
	km.IsRebinding = true
	km.RebindTarget = target
	km.RebindPrompt = prompt
}

func (km *KeybindManager) CancelRebind() {
	km.IsRebinding = false
	km.RebindTarget = ""
}

// Update listens for key presses during rebinding. Returns true if a key was bound.
func (km *KeybindManager) Update() bool {
	if !km.IsRebinding {
		return false
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		km.CancelRebind()
		return false
	}

	// Scan all possible keys
	for k := ebiten.Key(0); k <= ebiten.KeyMax; k++ {
		if inpututil.IsKeyJustPressed(k) {
			km.applyRebind(k)
			km.IsRebinding = false
			_ = km.Save()
			return true
		}
	}

	return false
}

func (km *KeybindManager) applyRebind(k ebiten.Key) {
	t := km.RebindTarget
	if strings.HasPrefix(t, "tower_") {
		idx := int(t[6] - '0')
		if idx >= 0 && idx < MaxToolbarSlots {
			km.TowerKeys[idx] = k
		}
	} else if strings.HasPrefix(t, "spell_") {
		idx := int(t[6] - '0')
		if idx >= 0 && idx < MaxToolbarSlots {
			km.SpellKeys[idx] = k
		}
	} else if t == "codex" {
		km.CodexKey = k
	} else if t == "pause" {
		km.PauseKey = k
	}
}

func (km *KeybindManager) KeyName(k ebiten.Key) string {
	s := k.String()
	s = strings.TrimPrefix(s, "Digit")
	s = strings.TrimPrefix(s, "Key")
	if len(s) == 0 {
		return "?"
	}
	return strings.ToUpper(s)
}

func (km *KeybindManager) Save() error {
	data := KeybindsData{
		TowerKeys: make([]string, MaxToolbarSlots),
		SpellKeys: make([]string, MaxToolbarSlots),
		CodexKey:  km.KeyName(km.CodexKey),
		PauseKey:  km.KeyName(km.PauseKey),
	}
	for i := 0; i < MaxToolbarSlots; i++ {
		data.TowerKeys[i] = km.KeyName(km.TowerKeys[i])
		data.SpellKeys[i] = km.KeyName(km.SpellKeys[i])
	}

	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(km.SavePath, bytes, 0644)
}

func (km *KeybindManager) Load() {
	bytes, err := os.ReadFile(km.SavePath)
	if err != nil {
		return
	}
	var data KeybindsData
	if err := json.Unmarshal(bytes, &data); err != nil {
		return
	}

	for i := 0; i < MaxToolbarSlots && i < len(data.TowerKeys); i++ {
		if k := parseKey(data.TowerKeys[i]); k >= 0 {
			km.TowerKeys[i] = k
		}
	}
	for i := 0; i < MaxToolbarSlots && i < len(data.SpellKeys); i++ {
		if k := parseKey(data.SpellKeys[i]); k >= 0 {
			km.SpellKeys[i] = k
		}
	}
	if k := parseKey(data.CodexKey); k >= 0 {
		km.CodexKey = k
	}
}

func parseKey(name string) ebiten.Key {
	name = strings.ToUpper(name)
	for k := ebiten.Key(0); k <= ebiten.KeyMax; k++ {
		str := strings.ToUpper(k.String())
		str = strings.TrimPrefix(str, "DIGIT")
		str = strings.TrimPrefix(str, "KEY")
		if str == name {
			return k
		}
	}
	return -1
}
