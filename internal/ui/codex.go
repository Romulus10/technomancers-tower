package ui

import (
	"fmt"
	"image/color"
	"math"

	"technomancers-tower/internal/config"
	"technomancers-tower/internal/data"
	"technomancers-tower/internal/economy"
	"technomancers-tower/internal/input"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type CodexTab int

const (
	TabTowers CodexTab = iota
	TabSpells
	TabEnemies
	TabKeybinds
	TabSettings
)

// GetCodexMaxScroll calculates the maximum downward scroll distance for the given tab.
func GetCodexMaxScroll(tab CodexTab, reg *data.Registry) float64 {
	switch tab {
	case TabTowers:
		cardH := 72.0
		gap := 10.0
		totalH := float64(len(reg.Towers))*(cardH+gap) + 15.0
		viewportH := 550.0 - 85.0
		return math.Max(0, totalH-viewportH)
	case TabSpells:
		cardH := 72.0
		gap := 10.0
		totalH := float64(len(reg.Spells))*(cardH+gap) + 15.0
		viewportH := 550.0 - 85.0
		return math.Max(0, totalH-viewportH)
	case TabEnemies:
		cardH := 76.0
		gap := 10.0
		totalH := float64(len(reg.Enemies))*(cardH+gap) + 15.0
		viewportH := 550.0 - 85.0
		return math.Max(0, totalH-viewportH)
	default:
		return 0
	}
}

func (u *UI) DrawCodex(screen *ebiten.Image, reg *data.Registry, km *input.KeybindManager, settings *config.SettingsManager, loadoutTowers [input.MaxToolbarSlots]string, loadoutSpells [input.MaxToolbarSlots]string, activeTab CodexTab, selectedItemID string, scrollY float64, cx, cy int) {
	// Dark semi-transparent background
	vector.FillRect(screen, 0, 0, 800, 600, color.RGBA{R: 6, G: 9, B: 15, A: 245}, false)

	// Main Codex Container
	boxX := float32(30)
	boxY := float32(25)
	boxW := float32(740)
	boxH := float32(550)

	vector.FillRect(screen, boxX, boxY, boxW, boxH, color.RGBA{R: 14, G: 18, B: 28, A: 255}, false)

	// 1. Render Tab Content First (so header masks it on top)
	switch activeTab {
	case TabTowers:
		u.drawTowersCodex(screen, reg, km, loadoutTowers, boxX, boxY, boxW, boxH, scrollY, cx, cy)
	case TabSpells:
		u.drawSpellsCodex(screen, reg, km, loadoutSpells, boxX, boxY, boxW, boxH, scrollY, cx, cy)
	case TabEnemies:
		u.drawEnemiesCodex(screen, reg, boxX, boxY, boxW, boxH, scrollY, cx, cy)
	case TabKeybinds:
		u.drawKeybindsTab(screen, km, boxX, boxY, boxW, boxH, cx, cy)
	case TabSettings:
		u.drawSettingsTab(screen, settings, boxX, boxY, boxW, boxH, cx, cy)
	}

	// 2. Top Header Masking Panel (prevents scrolled cards from overlapping tabs)
	vector.FillRect(screen, boxX+2, boxY+2, boxW-4, 62, color.RGBA{R: 14, G: 18, B: 28, A: 255}, false)

	// Top Navigation Tabs (5 Tabs)
	tabW := float32(116)
	tabH := float32(38)
	tabs := []string{"[ 1. TOWERS ]", "[ 2. SPELLS ]", "[ 3. MALWARE ]", "[ 4. KEYS ]", "[ 5. SETTINGS ]"}

	for i, tTitle := range tabs {
		tx := boxX + 16 + float32(i)*(tabW+6)
		ty := boxY + 15
		isSelected := int(activeTab) == i

		bg := color.RGBA{R: 22, G: 28, B: 42, A: 255}
		border := color.RGBA{R: 60, G: 80, B: 110, A: 255}
		if isSelected {
			bg = color.RGBA{R: 20, G: 70, B: 115, A: 255}
			border = color.RGBA{R: 0, G: 220, B: 255, A: 255}
		}

		vector.FillRect(screen, tx, ty, tabW, tabH, bg, false)
		vector.StrokeRect(screen, tx, ty, tabW, tabH, 1.5, border, false)
		ebitenutil.DebugPrintAt(screen, tTitle, int(tx)+8, int(ty)+12)
	}

	// Close Button on top right
	closeX := boxX + boxW - 95
	closeY := boxY + 15
	vector.FillRect(screen, closeX, closeY, 80, tabH, color.RGBA{R: 35, G: 22, B: 30, A: 255}, false)
	vector.StrokeRect(screen, closeX, closeY, 80, tabH, 1.5, color.RGBA{R: 255, G: 80, B: 100, A: 255}, false)
	ebitenutil.DebugPrintAt(screen, "[ CLOSE ]", int(closeX)+10, int(closeY)+12)

	// Divider
	vector.StrokeLine(screen, boxX+15, boxY+62, boxX+boxW-15, boxY+62, 1.5, color.RGBA{R: 40, G: 55, B: 80, A: 255}, false)

	// 3. Bottom Masking Bar & Container Border
	vector.FillRect(screen, boxX+2, boxY+boxH-8, boxW-4, 8, color.RGBA{R: 14, G: 18, B: 28, A: 255}, false)
	vector.StrokeRect(screen, boxX, boxY, boxW, boxH, 2, color.RGBA{R: 0, G: 200, B: 240, A: 255}, false)

	// 4. Scrollbar Track & Indicator Thumb
	maxScroll := GetCodexMaxScroll(activeTab, reg)
	if maxScroll > 0 {
		trackX := boxX + boxW - 12
		trackY := boxY + 70
		trackH := boxH - 85
		vector.FillRect(screen, trackX, trackY, 5, trackH, color.RGBA{R: 20, G: 28, B: 40, A: 220}, false)

		thumbH := float32(math.Max(32, float64(trackH)*(float64(trackH)/(float64(trackH)+maxScroll))))
		scrollRatio := float32(scrollY / maxScroll)
		if scrollRatio > 1 {
			scrollRatio = 1
		}
		thumbY := trackY + scrollRatio*(trackH-thumbH)
		vector.FillRect(screen, trackX, thumbY, 5, thumbH, color.RGBA{R: 0, G: 200, B: 240, A: 255}, false)
	}

	// Active Rebinding Overlay Prompt
	if km.IsRebinding {
		u.DrawRebindPrompt(screen, km)
	}
}

func (u *UI) drawTowersCodex(screen *ebiten.Image, reg *data.Registry, km *input.KeybindManager, loadout [input.MaxToolbarSlots]string, boxX, boxY, boxW, boxH float32, scrollY float64, cx, cy int) {
	cardW := float32(700)
	cardH := float32(72)
	gap := float32(10)
	startY := boxY + 75 - float32(scrollY)

	for i, def := range reg.Towers {
		cardY := startY + float32(i)*(cardH+gap)
		if cardY+cardH < boxY+65 || cardY > boxY+boxH-15 {
			continue
		}

		bg := color.RGBA{R: 20, G: 25, B: 38, A: 255}
		border := def.Color
		if !def.Unlocked {
			bg = color.RGBA{R: 16, G: 18, B: 24, A: 200}
			border = color.RGBA{R: 60, G: 70, B: 85, A: 180}
		}

		vector.FillRect(screen, boxX+20, cardY, cardW, cardH, bg, false)
		vector.StrokeRect(screen, boxX+20, cardY, cardW, cardH, 1.5, border, false)

		// Icon / Color Indicator
		vector.FillRect(screen, boxX+30, cardY+12, 48, 48, def.Color, false)
		vector.StrokeRect(screen, boxX+30, cardY+12, 48, 48, 1.5, color.RGBA{R: 255, G: 255, B: 255, A: 200}, false)

		if def.Unlocked {
			ebitenutil.DebugPrintAt(screen, def.Name, int(boxX)+90, int(cardY)+10)
			ebitenutil.DebugPrintAt(screen, def.Description, int(boxX)+90, int(cardY)+30)

			statsStr := fmt.Sprintf("Cost: %s B | Range: %.0f | Cooldown: %.2fs", economy.FormatNumber(def.BaseCost), def.Range, def.Cooldown)
			if def.ManaRate > 0 {
				statsStr += fmt.Sprintf(" | Mana: +%.1f/s", def.ManaRate)
			}
			if def.ByteRate > 0 {
				statsStr += fmt.Sprintf(" | Income: +%.1f B/s", def.ByteRate)
			}
			ebitenutil.DebugPrintAt(screen, statsStr, int(boxX)+90, int(cardY)+50)

			// Hotbar Slot Assignment Buttons [S1..S5]
			for s := 0; s < input.MaxToolbarSlots; s++ {
				sbX := boxX + 485 + float32(s)*42
				sbY := cardY + 20
				sbW := float32(36)
				sbH := float32(32)

				isAssigned := loadout[s] == def.ID
				sBg := color.RGBA{R: 25, G: 35, B: 55, A: 255}
				sBorder := color.RGBA{R: 70, G: 90, B: 120, A: 255}
				if isAssigned {
					sBg = color.RGBA{R: 0, G: 160, B: 220, A: 255}
					sBorder = color.RGBA{R: 160, G: 240, B: 255, A: 255}
				}

				vector.FillRect(screen, sbX, sbY, sbW, sbH, sBg, false)
				vector.StrokeRect(screen, sbX, sbY, sbW, sbH, 1.2, sBorder, false)

				keyLabel := km.KeyName(km.TowerKeys[s])
				ebitenutil.DebugPrintAt(screen, keyLabel, int(sbX)+12, int(sbY)+8)
			}
		} else {
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%s [LOCKED BLUEPRINT]", def.Name), int(boxX)+90, int(cardY)+15)
			ebitenutil.DebugPrintAt(screen, "Unlock this defense blueprint via Runtime Level-Up Draft or Root Access Archive.", int(boxX)+90, int(cardY)+40)
		}
	}
}

func (u *UI) drawSpellsCodex(screen *ebiten.Image, reg *data.Registry, km *input.KeybindManager, loadout [input.MaxToolbarSlots]string, boxX, boxY, boxW, boxH float32, scrollY float64, cx, cy int) {
	cardW := float32(700)
	cardH := float32(72)
	gap := float32(10)
	startY := boxY + 75 - float32(scrollY)

	for i, def := range reg.Spells {
		cardY := startY + float32(i)*(cardH+gap)
		if cardY+cardH < boxY+65 || cardY > boxY+boxH-15 {
			continue
		}

		bg := color.RGBA{R: 24, G: 20, B: 35, A: 255}
		border := def.Color
		if !def.Unlocked {
			bg = color.RGBA{R: 16, G: 18, B: 24, A: 200}
			border = color.RGBA{R: 60, G: 70, B: 85, A: 180}
		}

		vector.FillRect(screen, boxX+20, cardY, cardW, cardH, bg, false)
		vector.StrokeRect(screen, boxX+20, cardY, cardW, cardH, 1.5, border, false)

		vector.FillCircle(screen, boxX+54, cardY+36, 22, def.Color, false)
		vector.StrokeCircle(screen, boxX+54, cardY+36, 22, 1.5, color.RGBA{R: 255, G: 255, B: 255, A: 200}, false)

		if def.Unlocked {
			ebitenutil.DebugPrintAt(screen, def.Name, int(boxX)+90, int(cardY)+10)
			ebitenutil.DebugPrintAt(screen, def.Description, int(boxX)+90, int(cardY)+30)

			statsStr := fmt.Sprintf("Mana Cost: %.0f MP | Cooldown: %.1fs | Radius: %.0f", def.ManaCost, def.Cooldown, def.TargetRadius)
			ebitenutil.DebugPrintAt(screen, statsStr, int(boxX)+90, int(cardY)+50)

			for s := 0; s < input.MaxToolbarSlots; s++ {
				sbX := boxX + 485 + float32(s)*42
				sbY := cardY + 20
				sbW := float32(36)
				sbH := float32(32)

				isAssigned := loadout[s] == def.ID
				sBg := color.RGBA{R: 35, G: 25, B: 55, A: 255}
				sBorder := color.RGBA{R: 100, G: 70, B: 140, A: 255}
				if isAssigned {
					sBg = color.RGBA{R: 160, G: 70, B: 255, A: 255}
					sBorder = color.RGBA{R: 220, G: 180, B: 255, A: 255}
				}

				vector.FillRect(screen, sbX, sbY, sbW, sbH, sBg, false)
				vector.StrokeRect(screen, sbX, sbY, sbW, sbH, 1.2, sBorder, false)

				keyLabel := km.KeyName(km.SpellKeys[s])
				ebitenutil.DebugPrintAt(screen, keyLabel, int(sbX)+12, int(sbY)+8)
			}
		} else {
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%s [LOCKED SPELL]", def.Name), int(boxX)+90, int(boxX)+15)
			ebitenutil.DebugPrintAt(screen, "Unlock this Cyber-Spell via Runtime Level-Up Draft or Root Access Archive.", int(boxX)+90, int(cardY)+40)
		}
	}
}

func (u *UI) drawEnemiesCodex(screen *ebiten.Image, reg *data.Registry, boxX, boxY, boxW, boxH float32, scrollY float64, cx, cy int) {
	cardW := float32(700)
	cardH := float32(76)
	gap := float32(10)
	startY := boxY + 75 - float32(scrollY)

	for i, def := range reg.Enemies {
		cardY := startY + float32(i)*(cardH+gap)
		if cardY+cardH < boxY+65 || cardY > boxY+boxH-15 {
			continue
		}

		bg := color.RGBA{R: 24, G: 20, B: 28, A: 255}
		border := def.Color
		if def.IsBoss {
			bg = color.RGBA{R: 35, G: 16, B: 24, A: 255}
		}

		vector.FillRect(screen, boxX+20, cardY, cardW, cardH, bg, false)
		vector.StrokeRect(screen, boxX+20, cardY, cardW, cardH, 1.5, border, false)

		// Icon / Unit Representation
		vector.FillCircle(screen, boxX+54, cardY+38, def.Radius*1.2, def.Color, false)
		if def.BaseShield > 0 {
			vector.StrokeCircle(screen, boxX+54, cardY+38, def.Radius*1.2+3, 1.5, color.RGBA{R: 80, G: 190, B: 255, A: 240}, false)
		}

		// Text info
		nameLabel := def.Name
		if def.IsBoss {
			nameLabel = fmt.Sprintf("! BOSS: %s !", def.Name)
		}
		ebitenutil.DebugPrintAt(screen, nameLabel, int(boxX)+90, int(cardY)+10)
		ebitenutil.DebugPrintAt(screen, def.Description, int(boxX)+90, int(cardY)+28)

		statLine := fmt.Sprintf("HP: %.0f | Shield: %.0f | Speed: %.0f | Bounty: +%.0f Bytes | XP: +%.0f | Breach Dmg: %.0f",
			def.BaseHP, def.BaseShield, def.BaseSpeed, def.BaseBounty, def.BaseXP, def.CoreDamage)
		ebitenutil.DebugPrintAt(screen, statLine, int(boxX)+90, int(cardY)+50)
	}
}

func (u *UI) drawKeybindsTab(screen *ebiten.Image, km *input.KeybindManager, boxX, boxY, boxW, boxH float32, cx, cy int) {
	ebitenutil.DebugPrintAt(screen, "=== CUSTOMIZE OPERATOR KEYBINDINGS ===", int(boxX)+250, int(boxY)+80)
	ebitenutil.DebugPrintAt(screen, "Click any key badge below to rebind. Preferences auto-save globally.", int(boxX)+180, int(boxY)+102)

	// Column 1: Defense Node Hotkeys
	col1X := boxX + 40
	col1Y := boxY + 135
	ebitenutil.DebugPrintAt(screen, "DEFENSE NODE SLOTS [1-5]:", int(col1X), int(col1Y))

	for i := 0; i < input.MaxToolbarSlots; i++ {
		rowY := col1Y + 30 + float32(i)*48
		vector.FillRect(screen, col1X, rowY, 310, 40, color.RGBA{R: 20, G: 26, B: 38, A: 255}, false)
		vector.StrokeRect(screen, col1X, rowY, 310, 40, 1.2, color.RGBA{R: 50, G: 70, B: 100, A: 255}, false)

		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Tower Slot %d:", i+1), int(col1X)+15, int(rowY)+12)

		btnX := col1X + 220
		btnY := rowY + 6
		vector.FillRect(screen, btnX, btnY, 75, 28, color.RGBA{R: 25, G: 60, B: 95, A: 255}, false)
		vector.StrokeRect(screen, btnX, btnY, 75, 28, 1.2, color.RGBA{R: 0, G: 200, B: 255, A: 255}, false)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("[ %s ]", km.KeyName(km.TowerKeys[i])), int(btnX)+15, int(btnY)+6)
	}

	// Column 2: Cyber-Spell Hotkeys
	col2X := boxX + 390
	col2Y := boxY + 135
	ebitenutil.DebugPrintAt(screen, "CYBER-SPELL SLOTS [Q-W-E-R-F]:", int(col2X), int(col2Y))

	for i := 0; i < input.MaxToolbarSlots; i++ {
		rowY := col2Y + 30 + float32(i)*48
		vector.FillRect(screen, col2X, rowY, 310, 40, color.RGBA{R: 24, G: 20, B: 36, A: 255}, false)
		vector.StrokeRect(screen, col2X, rowY, 310, 40, 1.2, color.RGBA{R: 80, G: 60, B: 110, A: 255}, false)

		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Spell Slot %d:", i+1), int(col2X)+15, int(rowY)+12)

		btnX := col2X + 220
		btnY := rowY + 6
		vector.FillRect(screen, btnX, btnY, 75, 28, color.RGBA{R: 60, G: 35, B: 95, A: 255}, false)
		vector.StrokeRect(screen, btnX, btnY, 75, 28, 1.2, color.RGBA{R: 180, G: 120, B: 255, A: 255}, false)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("[ %s ]", km.KeyName(km.SpellKeys[i])), int(btnX)+15, int(btnY)+6)
	}

	// Reset Defaults Button
	resetBtnX := boxX + 260
	resetBtnY := boxY + 480
	vector.FillRect(screen, resetBtnX, resetBtnY, 220, 40, color.RGBA{R: 45, G: 25, B: 35, A: 255}, false)
	vector.StrokeRect(screen, resetBtnX, resetBtnY, 220, 40, 1.5, color.RGBA{R: 255, G: 90, B: 110, A: 255}, false)
	ebitenutil.DebugPrintAt(screen, "[ RESET DEFAULT KEYS ]", int(resetBtnX)+30, int(resetBtnY)+14)
}

func (u *UI) drawSettingsTab(screen *ebiten.Image, settings *config.SettingsManager, boxX, boxY, boxW, boxH float32, cx, cy int) {
	ebitenutil.DebugPrintAt(screen, "=== DISPLAY & RESOLUTION CONFIGURATION ===", int(boxX)+230, int(boxY)+80)
	ebitenutil.DebugPrintAt(screen, "Select display resolution or toggle fullscreen mode. Settings auto-save.", int(boxX)+165, int(boxY)+102)

	// Resolution Presets
	startX := boxX + 60
	startY := boxY + 135
	btnW := float32(620)
	btnH := float32(48)
	gap := float32(12)

	for i, res := range config.AvailableResolutions {
		ry := startY + float32(i)*(btnH+gap)
		isActive := settings.WindowWidth == res.Width && settings.WindowHeight == res.Height

		bg := color.RGBA{R: 20, G: 26, B: 38, A: 255}
		border := color.RGBA{R: 50, G: 70, B: 95, A: 255}
		if isActive {
			bg = color.RGBA{R: 20, G: 65, B: 105, A: 255}
			border = color.RGBA{R: 0, G: 220, B: 255, A: 255}
		}

		vector.FillRect(screen, startX, ry, btnW, btnH, bg, false)
		vector.StrokeRect(screen, startX, ry, btnW, btnH, 1.5, border, false)

		status := "[ SELECT ]"
		if isActive {
			status = ">>> ACTIVE <<<"
		}

		ebitenutil.DebugPrintAt(screen, res.Label, int(startX)+25, int(ry)+16)
		ebitenutil.DebugPrintAt(screen, status, int(startX)+480, int(ry)+16)
	}

	// Fullscreen Toggle Row
	fsY := startY + float32(len(config.AvailableResolutions))*(btnH+gap) + 10
	fsBg := color.RGBA{R: 28, G: 22, B: 38, A: 255}
	fsBorder := color.RGBA{R: 120, G: 70, B: 200, A: 255}
	if settings.IsFullscreen {
		fsBg = color.RGBA{R: 55, G: 25, B: 90, A: 255}
		fsBorder = color.RGBA{R: 200, G: 120, B: 255, A: 255}
	}

	vector.FillRect(screen, startX, fsY, btnW, btnH, fsBg, false)
	vector.StrokeRect(screen, startX, fsY, btnW, btnH, 1.5, fsBorder, false)

	fsStatus := "[ OFF ] (Click to Enable)"
	if settings.IsFullscreen {
		fsStatus = "[ FULLSCREEN ON ] (Active)"
	}
	ebitenutil.DebugPrintAt(screen, "FULLSCREEN DISPLAY MODE", int(startX)+25, int(fsY)+16)
	ebitenutil.DebugPrintAt(screen, fsStatus, int(startX)+420, int(fsY)+16)
}

func (u *UI) DrawRebindPrompt(screen *ebiten.Image, km *input.KeybindManager) {
	vector.FillRect(screen, 0, 0, 800, 600, color.RGBA{R: 2, G: 4, B: 8, A: 220}, false)

	pW := float32(440)
	pH := float32(180)
	pX := float32(180)
	pY := float32(210)

	vector.FillRect(screen, pX, pY, pW, pH, color.RGBA{R: 18, G: 25, B: 40, A: 255}, false)
	vector.StrokeRect(screen, pX, pY, pW, pH, 2.5, color.RGBA{R: 0, G: 220, B: 255, A: 255}, false)

	ebitenutil.DebugPrintAt(screen, "=== REBIND HOTKEY ===", int(pX)+145, int(pY)+25)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Binding: %s", km.RebindPrompt), int(pX)+40, int(pY)+65)
	ebitenutil.DebugPrintAt(screen, ">>> PRESS ANY KEY ON YOUR KEYBOARD <<<", int(pX)+65, int(pY)+100)
	ebitenutil.DebugPrintAt(screen, "[ Press ESC to Cancel ]", int(pX)+140, int(pY)+135)
}
