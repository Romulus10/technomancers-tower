package ui

import (
	"fmt"
	"image/color"
	"math"

	"technomancers-tower/internal/data"
	"technomancers-tower/internal/economy"
	"technomancers-tower/internal/gfx"
	"technomancers-tower/internal/input"
	"technomancers-tower/internal/malware"
	"technomancers-tower/internal/meta"
	"technomancers-tower/internal/motherboard"
	"technomancers-tower/internal/nodes"
	"technomancers-tower/internal/save"
	"technomancers-tower/internal/spells"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// UI manages the HUD, toolbar, placement previews, shop, pause, title, and slot select screens.
type UI struct {
	SelectedTower     string
	HasTowerSelected  bool
	SelectedSpell     string
	HasSpellSelected  bool
	HoveredShopTalent string
}

func NewUI() *UI {
	return &UI{}
}

func (u *UI) DrawTitleScreen(screen *ebiten.Image, animTime float64, version string, cx, cy int) {
	screen.Fill(color.RGBA{R: 8, G: 12, B: 20, A: 255})

	// Animated circuit grid backdrop
	for x := 0; x < 800; x += 40 {
		for y := 0; y < 600; y += 40 {
			vector.StrokeRect(screen, float32(x), float32(y), 40, 40, 1, color.RGBA{R: 16, G: 24, B: 38, A: 255}, false)
		}
	}

	// Pulsing decorative holographic rings
	pulse := float32((math.Sin(animTime*2.0) + 1.0) * 0.5)
	vector.StrokeCircle(screen, 400, 150, 140+pulse*12, 1.5, color.RGBA{R: 0, G: 160, B: 220, A: 80}, false)
	vector.StrokeCircle(screen, 400, 150, 160+pulse*18, 1.5, color.RGBA{R: 160, G: 70, B: 255, A: 60}, false)

	// Outer framing holographic box
	gfx.DrawHolographicPanel(screen, 20, 20, 760, 560, gfx.ColorPanelBorder, color.RGBA{0, 0, 0, 0})

	// Glowing Title Banner
	titleBoxW := float32(560)
	titleBoxH := float32(95)
	titleBoxX := float32(120)
	titleBoxY := float32(95)

	gfx.DrawHolographicPanel(screen, titleBoxX, titleBoxY, titleBoxW, titleBoxH, gfx.ColorCyanNeon, gfx.ColorPanelDark)

	ebitenutil.DebugPrintAt(screen, "==================================================", 185, 110)
	ebitenutil.DebugPrintAt(screen, "        T E C H N O M A N C E R ' S   T O W E R   ", 185, 130)
	ebitenutil.DebugPrintAt(screen, "==================================================", 185, 150)
	ebitenutil.DebugPrintAt(screen, "CYBER-ARCANE ENDLESS ROGUELITE MAZE DEFENSE", 235, 172)

	// Menu Buttons
	btnX := float32(260)
	btnW := float32(280)
	btnH := float32(44)

	// 1. [ START GAME ]
	btn1Y := float32(255)
	hover1 := float32(cx) >= btnX && float32(cx) <= btnX+btnW && float32(cy) >= btn1Y && float32(cy) <= btn1Y+btnH
	bg1 := color.RGBA{R: 20, G: 65, B: 105, A: 255}
	border1 := gfx.ColorCyanNeon
	if hover1 {
		bg1 = color.RGBA{R: 30, G: 95, B: 155, A: 255}
		border1 = gfx.ColorPanelHighlight
	}
	gfx.DrawHolographicPanel(screen, btnX, btn1Y, btnW, btnH, border1, bg1)
	ebitenutil.DebugPrintAt(screen, "[ START GAME / PROFILES ]", int(btnX)+45, int(btn1Y)+15)

	// 2. [ HOW TO PLAY / CONTROLS ]
	btn2Y := float32(310)
	hover2 := float32(cx) >= btnX && float32(cx) <= btnX+btnW && float32(cy) >= btn2Y && float32(cy) <= btn2Y+btnH
	bg2 := color.RGBA{R: 28, G: 38, B: 58, A: 255}
	border2 := color.RGBA{R: 80, G: 160, B: 240, A: 255}
	if hover2 {
		bg2 = color.RGBA{R: 40, G: 55, B: 85, A: 255}
		border2 = gfx.ColorCyanNeon
	}
	gfx.DrawHolographicPanel(screen, btnX, btn2Y, btnW, btnH, border2, bg2)
	ebitenutil.DebugPrintAt(screen, "[ HOW TO PLAY / CONTROLS ]", int(btnX)+45, int(btn2Y)+15)

	// 3. [ SETTINGS & RESOLUTION ]
	btn3Y := float32(365)
	hover3 := float32(cx) >= btnX && float32(cx) <= btnX+btnW && float32(cy) >= btn3Y && float32(cy) <= btn3Y+btnH
	bg3 := color.RGBA{R: 35, G: 30, B: 52, A: 255}
	border3 := gfx.ColorPurpleArcane
	if hover3 {
		bg3 = color.RGBA{R: 50, G: 40, B: 75, A: 255}
		border3 = gfx.Brighten(gfx.ColorPurpleArcane, 0.4)
	}
	gfx.DrawHolographicPanel(screen, btnX, btn3Y, btnW, btnH, border3, bg3)
	ebitenutil.DebugPrintAt(screen, "[ SETTINGS & RESOLUTION ]", int(btnX)+45, int(btn3Y)+15)

	// 4. [ EXIT MAINFRAME ]
	btn4Y := float32(420)
	hover4 := float32(cx) >= btnX && float32(cx) <= btnX+btnW && float32(cy) >= btn4Y && float32(cy) <= btn4Y+btnH
	bg4 := color.RGBA{R: 40, G: 24, B: 32, A: 255}
	border4 := gfx.ColorCrimsonGlitch
	if hover4 {
		bg4 = color.RGBA{R: 65, G: 32, B: 45, A: 255}
		border4 = gfx.Brighten(gfx.ColorCrimsonGlitch, 0.3)
	}
	gfx.DrawHolographicPanel(screen, btnX, btn4Y, btnW, btnH, border4, bg4)
	ebitenutil.DebugPrintAt(screen, "[ EXIT MAINFRAME ]", int(btnX)+75, int(btn4Y)+15)

	// Dynamic Version footer
	textW := len(version) * 6
	posX := 400 - textW/2
	ebitenutil.DebugPrintAt(screen, version, posX, 545)
}

func (u *UI) DrawHowToPlay(screen *ebiten.Image) {
	// Dark backdrop
	vector.FillRect(screen, 0, 0, 800, 600, color.RGBA{R: 4, G: 6, B: 10, A: 240}, false)

	// Dialog Window
	boxX := float32(80)
	boxY := float32(50)
	boxW := float32(640)
	boxH := float32(500)

	gfx.DrawHolographicPanel(screen, boxX, boxY, boxW, boxH, gfx.ColorCyanNeon, gfx.ColorPanelDark)

	ebitenutil.DebugPrintAt(screen, "=== OPERATOR DIRECTIVE & SYSTEM CONTROLS ===", 240, 75)

	lines := []string{
		"1. THE MOTHERBOARD MATRIX & MAZE BUILDING:",
		"   - Malware packets stream continuously from perimeter ports toward the CPU Kernel.",
		"   - Place defense nodes (1-5) to construct mazes. Nodes cannot completely block all paths.",
		"   - Building nodes costs BYTES, earned from purging incoming malware packets.",
		"",
		"2. CYBER-MANA & ACTIVE SPELLCASTING:",
		"   - Build Mana Siphon Nodes [2] to harvest ambient ether into CYBER-MANA per second.",
		"   - Cast active player spells with [Q, W, E, R] (Chain Lightning, Freeze, Bomb, Overclock).",
		"   - Failed casts (insufficient mana / cooldown) flash the toolbar red.",
		"",
		"3. ROGUELITE DRAFTS & ENDLESS ESCALATION:",
		"   - Purging malware awards Runtime XP. Level-ups prompt a 3-Card Upgrade Draft.",
		"   - Unlock new defense blueprints, spells, and global motherboard multipliers.",
		"   - Repel periodic Rogue AI Bosses arriving during intense DDoS Surge alarms.",
		"",
		"4. ROOT ACCESS META-PROGRESSION & SAVING:",
		"   - Kernel Panic (death) extracts GLITCH SHARDS for permanent Motherboard Talents.",
		"   - All progress & mid-run snapshots are saved in Gzip-compressed binary profiles.",
	}

	startY := 115
	for i, l := range lines {
		ebitenutil.DebugPrintAt(screen, l, int(boxX)+30, startY+i*18)
	}

	// Close Button
	closeBtnX := float32(280)
	closeBtnY := float32(490)
	closeBtnW := float32(240)
	closeBtnH := float32(40)

	gfx.DrawHolographicPanel(screen, closeBtnX, closeBtnY, closeBtnW, closeBtnH, gfx.ColorCyanNeon, color.RGBA{R: 20, G: 70, B: 110, A: 255})
	ebitenutil.DebugPrintAt(screen, "[ CLOSE GUIDE ] (ESC / SPACE)", int(closeBtnX)+25, int(closeBtnY)+14)
}

func (u *UI) DrawHUD(screen *ebiten.Image, run *economy.RunState, spawnerActiveSurge bool, surgeTimeLeft float64, graceTimeLeft float64, waveNum int, phase malware.WavePhase, phaseTimeLeft float64, threatLevel int) {
	// Top Header Bar
	gfx.DrawHolographicPanel(screen, 0, 0, 800, 58, gfx.ColorPanelBorder, gfx.ColorPanelDark)

	// 1. Kernel Integrity (Health) Gradient Bar
	hpRatio := float32(run.KernelHP / run.MaxKernelHP)
	gfx.DrawGradientBar(screen, 16, 12, 130, 16, hpRatio, gfx.ColorCrimsonGlitch, gfx.Brighten(gfx.ColorCrimsonGlitch, 0.4))
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("KERNEL HP: %.0f/%.0f", run.KernelHP, run.MaxKernelHP), 20, 32)

	// 2. Cyber-Mana Bar & Siphon Rate
	manaRatio := float32(run.Mana / run.MaxMana)
	gfx.DrawGradientBar(screen, 155, 12, 130, 16, manaRatio, gfx.ColorPurpleArcane, gfx.Brighten(gfx.ColorPurpleArcane, 0.4))
	effectiveManaRate := run.BaseManaRate * run.ManaGenMult
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("MANA: %.0f (+%.1f/s)", run.Mana, effectiveManaRate), 160, 32)

	// 3. Bytes Currency & XP Gradient Bar
	effectiveByteRate := run.BaseByteRate * run.ByteBountyMult
	var byteStr string
	if effectiveByteRate > 0 {
		byteStr = fmt.Sprintf("BYTES: %s (+%.1f/s)", economy.FormatNumber(run.Bytes), effectiveByteRate)
	} else {
		byteStr = fmt.Sprintf("BYTES: %s", economy.FormatNumber(run.Bytes))
	}
	ebitenutil.DebugPrintAt(screen, byteStr, 295, 14)

	xpRatio := float32(run.CurrentXP / run.TargetXP)
	gfx.DrawGradientBar(screen, 295, 32, 100, 12, xpRatio, gfx.ColorEmeraldBio, gfx.Brighten(gfx.ColorEmeraldBio, 0.4))
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("LVL %d", run.Level), 402, 31)

	// 4. Wave & Phase Indicator
	var phaseText string
	phaseColor := gfx.ColorCyanNeon
	switch phase {
	case malware.PhaseBuild:
		phaseText = fmt.Sprintf("RECON (%.0fs)", phaseTimeLeft)
		phaseColor = gfx.ColorEmeraldBio
	case malware.PhaseSwarm:
		phaseText = fmt.Sprintf("SWARM (%.0fs)", phaseTimeLeft)
		phaseColor = gfx.ColorGoldMatrix
	case malware.PhaseSurge:
		phaseText = fmt.Sprintf("! DDoS SURGE (%.0fs) !", phaseTimeLeft)
		phaseColor = gfx.ColorCrimsonGlitch
	}
	if threatLevel > 0 {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("WAVE %d [T+%d]", waveNum, threatLevel), 475, 14)
	} else {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("WAVE %d", waveNum), 475, 14)
	}
	ebitenutil.DebugPrintAt(screen, phaseText, 475, 32)
	vector.FillRect(screen, 465, 10, 3, 38, phaseColor, false)

	// 5. Survival Time & Purged Stats
	mins := int(run.RunTime) / 60
	secs := int(run.RunTime) % 60
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TIME: %02d:%02d", mins, secs), 630, 14)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("PURGED: %d", run.Kills), 630, 32)

	// Pause prompt
	ebitenutil.DebugPrintAt(screen, "[ESC: PAUSE]", 710, 42)

	// Grace Period Banner
	if graceTimeLeft > 0 {
		gfx.DrawHolographicPanel(screen, 210, 75, 380, 52, gfx.ColorCyanNeon, color.RGBA{R: 18, G: 35, B: 55, A: 235})
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf(">>> PREPARING FIREWALLS: MALWARE IN %.1fs <<<", graceTimeLeft), 235, 84)
		ebitenutil.DebugPrintAt(screen, "[ Place defense nodes (1-5) | Press SPACE to Deploy ]", 225, 104)
	}
}

func (u *UI) DrawToolbar(screen *ebiten.Image, tm *nodes.TowerManager, sm *spells.SpellManager, km *input.KeybindManager, loadoutTowers [input.MaxToolbarSlots]string, loadoutSpells [input.MaxToolbarSlots]string, run *economy.RunState, cx, cy int) {
	// Bottom Toolbar Panel
	gfx.DrawHolographicPanel(screen, 0, 535, 800, 65, gfx.ColorPanelBorder, gfx.ColorPanelDark)

	cache := gfx.GetCache()

	// Defense Nodes (Slots 1..5)
	startX := float32(14)
	btnW := float32(68)
	btnH := float32(50)
	gap := float32(6)

	for i := 0; i < input.MaxToolbarSlots; i++ {
		tID := loadoutTowers[i]
		def := tm.Registry.GetTower(tID)
		bx := startX + float32(i)*(btnW+gap)
		by := float32(542)

		if def == nil {
			// Empty Slot
			gfx.DrawHolographicPanel(screen, bx, by, btnW, btnH, color.RGBA{R: 50, G: 60, B: 75, A: 180}, color.RGBA{R: 18, G: 20, B: 28, A: 200})
			keyLabel := km.KeyName(km.TowerKeys[i])
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("[%s]", keyLabel), int(bx)+4, int(by)+6)
			ebitenutil.DebugPrintAt(screen, "EMPTY", int(bx)+4, int(by)+24)
			continue
		}

		isSelected := u.HasTowerSelected && u.SelectedTower == def.ID
		canAfford := run.CanAffordBytes(def.BaseCost)

		bgColor := color.RGBA{R: 22, G: 28, B: 40, A: 255}
		if isSelected {
			bgColor = color.RGBA{R: 40, G: 65, B: 95, A: 255}
		}
		if !def.Unlocked {
			bgColor = color.RGBA{R: 18, G: 18, B: 22, A: 200}
		}

		borderColor := def.Color
		if !def.Unlocked || !canAfford {
			borderColor = color.RGBA{R: 70, G: 80, B: 95, A: 200}
		}
		gfx.DrawHolographicPanel(screen, bx, by, btnW, btnH, borderColor, bgColor)

		// Procedural Icon Badge
		if def.Unlocked {
			icon := cache.GetIconBadge(def.ID, true, def.Color)
			if icon != nil {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Scale(0.65, 0.65)
				op.GeoM.Translate(float64(bx+btnW-22), float64(by+6))
				screen.DrawImage(icon, op)
			}
		}

		keyLabel := km.KeyName(km.TowerKeys[i])
		if def.Unlocked {
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("[%s]%s", keyLabel, def.Name[:min(len(def.Name), 4)]), int(bx)+4, int(by)+6)
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%s B", economy.FormatNumber(def.BaseCost)), int(bx)+4, int(by)+28)
		} else {
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("[%s]LOCK", keyLabel), int(bx)+4, int(by)+18)
		}
	}

	// Center: Codex Button
	codexBtnX := float32(380)
	codexBtnY := float32(542)
	codexBtnW := float32(40)
	codexBtnH := float32(50)
	gfx.DrawHolographicPanel(screen, codexBtnX, codexBtnY, codexBtnW, codexBtnH, gfx.ColorCyanNeon, color.RGBA{R: 28, G: 38, B: 60, A: 255})
	ebitenutil.DebugPrintAt(screen, "TAB", int(codexBtnX)+9, int(codexBtnY)+10)
	ebitenutil.DebugPrintAt(screen, "BOOK", int(codexBtnX)+5, int(codexBtnY)+28)

	// Cyber-Spells (Slots 1..5)
	spellStartX := float32(422)

	for i := 0; i < input.MaxToolbarSlots; i++ {
		sID := loadoutSpells[i]
		def := sm.Registry.GetSpell(sID)
		bx := spellStartX + float32(i)*(btnW+gap)
		by := float32(542)

		if def == nil {
			// Empty Slot
			gfx.DrawHolographicPanel(screen, bx, by, btnW, btnH, color.RGBA{R: 60, G: 50, B: 75, A: 180}, color.RGBA{R: 20, G: 18, B: 28, A: 200})
			keyLabel := km.KeyName(km.SpellKeys[i])
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("[%s]", keyLabel), int(bx)+4, int(by)+6)
			ebitenutil.DebugPrintAt(screen, "EMPTY", int(bx)+4, int(by)+24)
			continue
		}

		isSelected := u.HasSpellSelected && u.SelectedSpell == def.ID
		canAfford := run.CanAffordMana(def.ManaCost)
		cd := sm.Cooldowns[def.ID]
		flash := sm.Flashes[def.ID]

		bgColor := color.RGBA{R: 26, G: 22, B: 38, A: 255}
		if isSelected {
			bgColor = color.RGBA{R: 60, G: 45, B: 85, A: 255}
		}
		if !def.Unlocked {
			bgColor = color.RGBA{R: 18, G: 18, B: 22, A: 200}
		}

		if flash > 0 {
			flashAlpha := uint8(255 * (flash / 0.35))
			bgColor = color.RGBA{R: 180, G: 20, B: 35, A: flashAlpha}
		}

		borderColor := def.Color
		if !def.Unlocked || !canAfford || cd > 0 {
			borderColor = color.RGBA{R: 70, G: 80, B: 95, A: 200}
		}
		if flash > 0 {
			borderColor = color.RGBA{R: 255, G: 60, B: 75, A: 255}
		}

		gfx.DrawHolographicPanel(screen, bx, by, btnW, btnH, borderColor, bgColor)

		if cd > 0 && def.Cooldown > 0 {
			cdRatio := float32(cd / (def.Cooldown * run.SpellCooldownMult))
			vector.FillRect(screen, bx, by, btnW, btnH*cdRatio, color.RGBA{R: 0, G: 0, B: 0, A: 160}, false)
		}

		// Procedural Icon Badge
		if def.Unlocked {
			icon := cache.GetIconBadge(def.ID, false, def.Color)
			if icon != nil {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Scale(0.65, 0.65)
				op.GeoM.Translate(float64(bx+btnW-22), float64(by+6))
				screen.DrawImage(icon, op)
			}
		}

		keyLabel := km.KeyName(km.SpellKeys[i])
		if def.Unlocked {
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("[%s]%s", keyLabel, def.Name[:min(len(def.Name), 4)]), int(bx)+4, int(by)+6)
			if cd > 0 {
				ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%.1fs", cd), int(bx)+4, int(by)+28)
			} else {
				actualMana := def.ManaCost * run.SpellCostMult
				ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%.0f MP", actualMana), int(bx)+4, int(by)+28)
			}
		} else {
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("[%s]LOCK", keyLabel), int(bx)+4, int(by)+18)
		}
	}

	// Upgrade All Button
	if count, totalCost := tm.GetPromotableInfo(); count > 0 {
		upgBtnX := float32(280)
		upgBtnY := float32(498)
		upgBtnW := float32(240)
		upgBtnH := float32(30)

		canAffordAll := run.CanAffordBytes(totalCost)
		bgColor := color.RGBA{R: 20, G: 48, B: 35, A: 245}
		borderColor := gfx.ColorEmeraldBio
		btnLabel := fmt.Sprintf("^ UPGRADE ALL (%d) | %sB [U]", count, economy.FormatNumber(totalCost))

		if !canAffordAll {
			bgColor = color.RGBA{R: 45, G: 38, B: 20, A: 240}
			borderColor = gfx.ColorGoldMatrix
			btnLabel = fmt.Sprintf("^ UPGRADE READY (%d) | %sB [U]", count, economy.FormatNumber(totalCost))
		}

		if float32(cx) >= upgBtnX && float32(cx) <= upgBtnX+upgBtnW && float32(cy) >= upgBtnY && float32(cy) <= upgBtnY+upgBtnH {
			bgColor.R = uint8(math.Min(255, float64(bgColor.R)+30))
			bgColor.G = uint8(math.Min(255, float64(bgColor.G)+30))
			bgColor.B = uint8(math.Min(255, float64(bgColor.B)+30))
		}

		gfx.DrawHolographicPanel(screen, upgBtnX, upgBtnY, upgBtnW, upgBtnH, borderColor, bgColor)
		ebitenutil.DebugPrintAt(screen, btnLabel, int(upgBtnX)+10, int(upgBtnY)+8)
	}
}

func (u *UI) DrawPlacementPreview(screen *ebiten.Image, grid *motherboard.Grid, tm *nodes.TowerManager, sm *spells.SpellManager, run *economy.RunState, cx, cy int) {
	if u.HasTowerSelected {
		def := tm.Registry.GetTower(u.SelectedTower)
		if def != nil {
			gx, gy, ok := grid.ScreenToGrid(cx, cy)
			if ok {
				canBuild := grid.CanBuildAt(gx, gy) && run.CanAffordBytes(def.BaseCost)
				sx := float32(motherboard.OffsetX + gx*motherboard.CellSize)
				sy := float32(motherboard.OffsetY + gy*motherboard.CellSize)
				sz := float32(motherboard.CellSize)

				boxColor := color.RGBA{R: 0, G: 255, B: 150, A: 160}
				if !canBuild {
					boxColor = color.RGBA{R: 255, G: 50, B: 50, A: 160}
				}
				vector.FillRect(screen, sx, sy, sz, sz, boxColor, false)

				if def.Range > 0 {
					rcx, rcy := grid.GridToScreenCenter(gx, gy)
					effectiveRange := float32(def.Range * run.TowerRangeMult)
					vector.StrokeCircle(screen, float32(rcx), float32(rcy), effectiveRange, 1.5, gfx.ColorCyanNeon, false)
				}
			}
		}
	}

	if u.HasSpellSelected {
		def := sm.Registry.GetSpell(u.SelectedSpell)
		if def != nil && def.TargetRadius > 0 {
			vector.StrokeCircle(screen, float32(cx), float32(cy), float32(def.TargetRadius), 2, def.Color, false)
			vector.FillCircle(screen, float32(cx), float32(cy), float32(def.TargetRadius), color.RGBA{
				R: def.Color.R, G: def.Color.G, B: def.Color.B, A: 40,
			}, false)
		}
	}

	// Hover inspection tooltip when hovering placed towers
	if !u.HasTowerSelected && !u.HasSpellSelected {
		gx, gy, ok := grid.ScreenToGrid(cx, cy)
		if ok {
			tower := tm.GetTowerAt(gx, gy)
			if tower != nil {
				// Draw Range indicator on hovered tower
				if tower.Def.Range > 0 {
					effectiveRange := float32(tower.Def.Range * run.TowerRangeMult * tower.GetRangeMult())
					vector.StrokeCircle(screen, float32(tower.WorldX), float32(tower.WorldY), effectiveRange, 1.2, gfx.ColorGoldMatrix, false)
				}

				// Tooltip Box near cursor
				tipX := float32(cx + 15)
				tipY := float32(cy - 20)
				tipW := float32(200)
				tipH := float32(100)
				if tower.Level >= nodes.MaxTowerLevel {
					tipH = 115
				}

				// Keep inside screen boundaries
				if tipX+tipW > 790 {
					tipX = float32(cx) - tipW - 10
				}
				if tipY+tipH > 530 {
					tipY = float32(530) - tipH
				}
				if tipY < 60 {
					tipY = 65
				}

				gfx.DrawHolographicPanel(screen, tipX, tipY, tipW, tipH, gfx.ColorCyanNeon, gfx.ColorPanelDark)

				// Header: Name & Tier / Level
				tierStr := ""
				if tower.Tier > 1 {
					tierStr = fmt.Sprintf(" [T%d]", tower.Tier)
				}
				ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%s%s [Lv.%d/%d]", tower.Def.Name, tierStr, tower.Level, nodes.MaxTowerLevel), int(tipX)+8, int(tipY)+6)

				// XP Bar
				ratio := float32(0.0)
				if tower.Level >= nodes.MaxTowerLevel {
					ratio = 1.0
					gfx.DrawGradientBar(screen, tipX+8, tipY+24, 184, 8, ratio, gfx.ColorGoldMatrix, gfx.Brighten(gfx.ColorGoldMatrix, 0.4))
					ebitenutil.DebugPrintAt(screen, "XP: MAX RANK (READY)", int(tipX)+8, int(tipY)+36)
				} else if tower.TargetXP > 0 {
					ratio = float32(tower.XP / tower.TargetXP)
					if ratio > 1 {
						ratio = 1
					}
					gfx.DrawGradientBar(screen, tipX+8, tipY+24, 184, 8, ratio, gfx.ColorCyanNeon, gfx.ColorFrostIce)
					ebitenutil.DebugPrintAt(screen, fmt.Sprintf("XP: %.0f / %.0f (%.0f%%)", tower.XP, tower.TargetXP, ratio*100), int(tipX)+8, int(tipY)+36)
				}

				// Scaled stats
				if tower.Def.Delivery.Type == data.DeliveryPassiveGenerator {
					if tower.Def.ByteRate > 0 {
						ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Byte Rate: +%.1f/s (x%.2f)", tower.GetByteRate()*run.ByteBountyMult, tower.GetByteRate()/tower.Def.ByteRate), int(tipX)+8, int(tipY)+54)
					} else {
						ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Mana Rate: +%.1f/s (x%.2f)", tower.GetManaRate()*run.ManaGenMult, tower.GetManaRate()/tower.Def.ManaRate), int(tipX)+8, int(tipY)+54)
					}
				} else {
					ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Dmg: x%.2f | Spd: x%.2f", tower.GetDamageMult(), tower.GetSpeedMult()), int(tipX)+8, int(tipY)+54)
					ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Range: %.0f (x%.2f)", tower.Def.Range*run.TowerRangeMult*tower.GetRangeMult(), tower.GetRangeMult()), int(tipX)+8, int(tipY)+72)
				}

				// Click to Upgrade Prompt at Max Level
				if tower.Level >= nodes.MaxTowerLevel {
					cost := tower.GetUpgradeCost()
					if run.CanAffordBytes(cost) {
						ebitenutil.DebugPrintAt(screen, fmt.Sprintf("[CLICK: UPGRADE %sB]", economy.FormatNumber(cost)), int(tipX)+8, int(tipY)+92)
					} else {
						ebitenutil.DebugPrintAt(screen, fmt.Sprintf("[NEED: %sB TO UPGRADE]", economy.FormatNumber(cost)), int(tipX)+8, int(tipY)+92)
					}
				}
			}
		}
	}
}

func (u *UI) DrawGameOver(screen *ebiten.Image, run *economy.RunState, metaMgr *meta.MetaManager) {
	vector.FillRect(screen, 0, 0, 800, 600, color.RGBA{R: 8, G: 10, B: 16, A: 235}, false)

	gfx.DrawHolographicPanel(screen, 180, 90, 440, 420, gfx.ColorCrimsonGlitch, gfx.ColorPanelDark)

	ebitenutil.DebugPrintAt(screen, "=== KERNEL PANIC: CORE INTEGRITY LOST ===", 240, 120)

	mins := int(run.RunTime) / 60
	secs := int(run.RunTime) % 60
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Survival Runtime: %02d:%02d", mins, secs), 240, 165)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Malware Packets Purged: %d", run.Kills), 240, 195)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Rogue AI Bosses Defeated: %d", run.BossKills), 240, 225)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Bytes Harvested: %s", economy.FormatNumber(run.TotalBytes)), 240, 255)

	vector.StrokeLine(screen, 220, 290, 580, 290, 1, color.RGBA{R: 60, G: 70, B: 90, A: 255}, false)

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("GLITCH SHARDS EXTRACTED: +%d", run.ShardsEarned), 240, 315)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Total Glitch Shards in Archive: %d", metaMgr.GlitchShards), 240, 345)

	gfx.DrawHolographicPanel(screen, 240, 390, 320, 40, gfx.ColorCyanNeon, color.RGBA{R: 35, G: 55, B: 85, A: 255})
	ebitenutil.DebugPrintAt(screen, "[ ENTER ROOT ACCESS / META-SHOP ] (SPACE)", 265, 404)

	gfx.DrawHolographicPanel(screen, 240, 445, 320, 40, gfx.ColorEmeraldBio, color.RGBA{R: 30, G: 45, B: 65, A: 255})
	ebitenutil.DebugPrintAt(screen, "[ QUICK REBOOT KERNEL ] (ENTER)", 295, 459)
}

func (u *UI) DrawPauseMenu(screen *ebiten.Image) {
	// Dark overlay
	vector.FillRect(screen, 0, 0, 800, 600, color.RGBA{R: 5, G: 8, B: 14, A: 230}, false)

	// Pause Box
	gfx.DrawHolographicPanel(screen, 220, 95, 360, 410, gfx.ColorCyanNeon, gfx.ColorPanelDark)

	ebitenutil.DebugPrintAt(screen, "=== MAINFRAME SUSPENDED (PAUSED) ===", 260, 120)

	// Buttons
	btnW := float32(280)
	btnH := float32(40)
	btnX := float32(260)

	// 1. Resume
	b1Y := float32(160)
	gfx.DrawHolographicPanel(screen, btnX, b1Y, btnW, btnH, gfx.ColorCyanNeon, color.RGBA{R: 30, G: 60, B: 90, A: 255})
	ebitenutil.DebugPrintAt(screen, "[ RESUME RUN ] (ESC)", 335, int(b1Y)+14)

	// 2. Settings
	b2Y := float32(215)
	gfx.DrawHolographicPanel(screen, btnX, b2Y, btnW, btnH, gfx.ColorPurpleArcane, color.RGBA{R: 35, G: 30, B: 52, A: 255})
	ebitenutil.DebugPrintAt(screen, "[ DISPLAY & SETTINGS ]", 325, int(b2Y)+14)

	// 3. Save & Exit
	b3Y := float32(270)
	gfx.DrawHolographicPanel(screen, btnX, b3Y, btnW, btnH, gfx.ColorCyanGlow, color.RGBA{R: 35, G: 50, B: 75, A: 255})
	ebitenutil.DebugPrintAt(screen, "[ SAVE & EXIT TO SLOTS ]", 315, int(b3Y)+14)

	// 4. Abandon
	b4Y := float32(325)
	gfx.DrawHolographicPanel(screen, btnX, b4Y, btnW, btnH, gfx.ColorCrimsonGlitch, color.RGBA{R: 45, G: 35, B: 55, A: 255})
	ebitenutil.DebugPrintAt(screen, "[ ABANDON & RESTART RUN ]", 310, int(b4Y)+14)

	// 5. Meta-Shop
	b5Y := float32(380)
	gfx.DrawHolographicPanel(screen, btnX, b5Y, btnW, btnH, gfx.ColorPurpleArcane, color.RGBA{R: 35, G: 45, B: 65, A: 255})
	ebitenutil.DebugPrintAt(screen, "[ ROOT ACCESS ARCHIVE ]", 320, int(b5Y)+14)
}

func (u *UI) DrawSlotSelectScreen(screen *ebiten.Image, sm *save.SaveManager, cx, cy int) {
	screen.Fill(color.RGBA{R: 10, G: 14, B: 22, A: 255})

	// Header
	gfx.DrawHolographicPanel(screen, 0, 0, 800, 75, gfx.ColorCyanNeon, gfx.ColorPanelDark)

	// Back button
	btnBackX := float32(25)
	btnBackY := float32(18)
	btnBackW := float32(160)
	btnBackH := float32(38)
	hoverBack := float32(cx) >= btnBackX && float32(cx) <= btnBackX+btnBackW && float32(cy) >= btnBackY && float32(cy) <= btnBackY+btnBackH
	bgBack := color.RGBA{R: 28, G: 35, B: 50, A: 255}
	if hoverBack {
		bgBack = color.RGBA{R: 40, G: 55, B: 80, A: 255}
	}
	gfx.DrawHolographicPanel(screen, btnBackX, btnBackY, btnBackW, btnBackH, gfx.ColorCyanNeon, bgBack)
	ebitenutil.DebugPrintAt(screen, "< MAIN MENU [ESC]", int(btnBackX)+15, int(btnBackY)+12)

	ebitenutil.DebugPrintAt(screen, "=== SELECT SAVE PROFILE ===", 320, 18)
	ebitenutil.DebugPrintAt(screen, "Gzip-Compressed Binary Storage | Multi-Slot Profiles", 270, 42)

	// Draw 3 Slots
	startY := float32(95.0)
	cardH := float32(132.0)
	gap := float32(16.0)
	cardW := float32(720.0)

	for i := 0; i < save.MaxSlots; i++ {
		slot := sm.Slots[i]
		cardY := startY + float32(i)*(cardH+gap)

		borderColor := color.RGBA{R: 50, G: 70, B: 100, A: 255}
		if slot.Exists {
			borderColor = gfx.ColorCyanNeon
		}
		gfx.DrawHolographicPanel(screen, 40, cardY, cardW, cardH, borderColor, gfx.ColorPanelDark)

		slotTitle := fmt.Sprintf("SLOT 0%d: %s", i+1, slot.ProfileName)
		if !slot.Exists {
			slotTitle = fmt.Sprintf("SLOT 0%d: [ EMPTY PROFILE - INITIALIZE ]", i+1)
		}
		ebitenutil.DebugPrintAt(screen, slotTitle, 60, int(cardY)+15)

		if slot.Exists {
			mins := int(slot.HighScore) / 60
			secs := int(slot.HighScore) % 60
			statsStr := fmt.Sprintf("Glitch Shards: %d  |  High Score: %02d:%02d  |  Lifetime Kills: %d", slot.GlitchShards, mins, secs, slot.TotalKills)
			ebitenutil.DebugPrintAt(screen, statsStr, 60, int(cardY)+40)

			lastSavedStr := fmt.Sprintf("Last Saved: %s", slot.LastSaved.Format("2006-01-02 15:04:05"))
			ebitenutil.DebugPrintAt(screen, lastSavedStr, 60, int(cardY)+62)

			if slot.HasActiveRun {
				runMins := int(slot.ActiveRunDuration) / 60
				runSecs := int(slot.ActiveRunDuration) % 60
				runStatus := fmt.Sprintf(">>> ACTIVE RUN IN PROGRESS: %02d:%02d (Kernel HP: %.0f) <<<", runMins, runSecs, slot.ActiveRunHP)
				ebitenutil.DebugPrintAt(screen, runStatus, 60, int(cardY)+85)
			} else {
				ebitenutil.DebugPrintAt(screen, "Status: Ready for deployment in Motherboard Matrix", 60, int(cardY)+85)
			}
		} else {
			ebitenutil.DebugPrintAt(screen, "No profile data stored. Click [INITIALIZE] to start fresh operator profile.", 60, int(cardY)+55)
		}

		// Buttons:
		btnPlayX := float32(560)
		btnPlayY := cardY + 18
		btnPlayW := float32(180)
		btnPlayH := float32(40)

		playText := "[ RESUME RUN ]"
		if !slot.Exists {
			playText = "[ INITIALIZE ]"
		} else if !slot.HasActiveRun {
			playText = "[ LAUNCH RUN ]"
		}

		gfx.DrawHolographicPanel(screen, btnPlayX, btnPlayY, btnPlayW, btnPlayH, gfx.ColorCyanNeon, color.RGBA{R: 20, G: 70, B: 110, A: 255})
		ebitenutil.DebugPrintAt(screen, playText, int(btnPlayX)+28, int(btnPlayY)+14)

		if slot.Exists {
			btnDelX := float32(560)
			btnDelY := cardY + 70
			btnDelW := float32(180)
			btnDelH := float32(38)

			gfx.DrawHolographicPanel(screen, btnDelX, btnDelY, btnDelW, btnDelH, gfx.ColorCrimsonGlitch, color.RGBA{R: 45, G: 25, B: 35, A: 255})
			ebitenutil.DebugPrintAt(screen, "[ DELETE PROFILE ]", int(btnDelX)+20, int(btnDelY)+12)
		}
	}

	// Bottom Bar
	gfx.DrawHolographicPanel(screen, 0, 540, 800, 60, gfx.ColorPanelBorder, gfx.ColorPanelDark)
	ebitenutil.DebugPrintAt(screen, "[ Profiles are automatically preserved in compressed binary storage ]", 185, 562)
}

func (u *UI) DrawMetaShop(screen *ebiten.Image, metaMgr *meta.MetaManager, scrollY float64, cx, cy int) {
	screen.Fill(color.RGBA{R: 12, G: 16, B: 24, A: 255})

	// List Talents with Scroll Offset
	startY := float32(90.0)
	cardH := float32(65.0)
	gap := float32(10.0)
	cardW := float32(710.0)

	for i, t := range metaMgr.Talents {
		cardY := startY + float32(i)*(cardH+gap) - float32(scrollY)

		if cardY+cardH < 75 || cardY > 530 {
			continue
		}

		cost := t.CurrentCost()
		canBuy := cost > 0 && metaMgr.GlitchShards >= cost

		gfx.DrawHolographicPanel(screen, 35, cardY, cardW, cardH, color.RGBA{R: 50, G: 70, B: 100, A: 255}, gfx.ColorPanelDark)

		levelStr := fmt.Sprintf("Rank %d/%d", t.Level, t.MaxLevel)
		if t.Level >= t.MaxLevel {
			levelStr = "MAXED"
		}
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%s (%s)", t.Name, levelStr), 55, int(cardY)+12)
		ebitenutil.DebugPrintAt(screen, t.Description, 55, int(cardY)+34)

		btnX := float32(585)
		btnY := cardY + 12
		btnW := float32(145)
		btnH := float32(40)

		btnBg := color.RGBA{R: 35, G: 55, B: 85, A: 255}
		if !canBuy {
			btnBg = color.RGBA{R: 25, G: 28, B: 35, A: 255}
		}

		btnBorder := gfx.ColorCyanNeon
		if !canBuy {
			btnBorder = color.RGBA{R: 60, G: 70, B: 85, A: 200}
		}
		gfx.DrawHolographicPanel(screen, btnX, btnY, btnW, btnH, btnBorder, btnBg)

		if cost > 0 {
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Upgrade (%d Shards)", cost), int(btnX)+8, int(btnY)+14)
		} else {
			ebitenutil.DebugPrintAt(screen, "MAX RANK", int(btnX)+38, int(btnY)+14)
		}
	}

	// Scrollbar Track & Thumb (Right side)
	totalContent := float32(len(metaMgr.Talents)*(int(cardH)+int(gap)) + 20)
	viewportH := float32(440.0)
	if totalContent > viewportH {
		trackX := float32(760)
		trackY := float32(85)
		trackW := float32(6)
		trackH := viewportH

		vector.FillRect(screen, trackX, trackY, trackW, trackH, color.RGBA{R: 20, G: 28, B: 40, A: 200}, false)

		maxScroll := totalContent - viewportH
		thumbH := float32(math.Max(40, float64(viewportH*(viewportH/totalContent))))
		scrollRatio := float32(scrollY) / maxScroll
		if scrollRatio > 1.0 {
			scrollRatio = 1.0
		}
		thumbY := trackY + (trackH-thumbH)*scrollRatio

		vector.FillRect(screen, trackX, thumbY, trackW, thumbH, gfx.ColorCyanNeon, false)
	}

	// Fixed Header
	gfx.DrawHolographicPanel(screen, 0, 0, 800, 75, gfx.ColorPurpleArcane, gfx.ColorPanelDark)

	ebitenutil.DebugPrintAt(screen, "=== ROOT ACCESS: MOTHERBOARD ARCHITECTURE ARCHIVE ===", 210, 16)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Archived Glitch Shards: %d  (Scroll with Wheel / Arrow Keys / W & S)", metaMgr.GlitchShards), 150, 42)

	// Fixed Bottom Bar: Start Run
	gfx.DrawHolographicPanel(screen, 0, 530, 800, 70, gfx.ColorCyanNeon, gfx.ColorPanelDark)

	startBtnX := float32(280)
	startBtnY := float32(545)
	startBtnW := float32(240)
	startBtnH := float32(42)

	gfx.DrawHolographicPanel(screen, startBtnX, startBtnY, startBtnW, startBtnH, gfx.ColorCyanNeon, color.RGBA{R: 20, G: 80, B: 120, A: 255})
	ebitenutil.DebugPrintAt(screen, "[ BOOT RUN ] (SPACE)", int(startBtnX)+50, int(startBtnY)+15)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
