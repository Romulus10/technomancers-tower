package state

import (
	"fmt"
	"image/color"
	"math"

	"technomancers-tower/internal/config"
	"technomancers-tower/internal/data"
	"technomancers-tower/internal/draft"
	"technomancers-tower/internal/economy"
	"technomancers-tower/internal/input"
	"technomancers-tower/internal/malware"
	"technomancers-tower/internal/meta"
	"technomancers-tower/internal/motherboard"
	"technomancers-tower/internal/nodes"
	"technomancers-tower/internal/save"
	"technomancers-tower/internal/spells"
	"technomancers-tower/internal/ui"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type GameMode int

const (
	ModeTitle GameMode = iota
	ModeSlotSelect
	ModePlaying
	ModeDraft
	ModeCodex
	ModePause
	ModeGameOver
	ModeMetaShop
)

type GameState struct {
	Mode              GameMode
	PreviousMode      GameMode
	SaveMgr           *save.SaveManager
	KeyMgr            *input.KeybindManager
	Settings          *config.SettingsManager
	Registry          *data.Registry
	Grid              *motherboard.Grid
	Run               *economy.RunState
	Malware           *malware.Spawner
	Towers            *nodes.TowerManager
	Spells            *spells.SpellManager
	Draft             *draft.DraftManager
	Meta              *meta.MetaManager
	UI                *ui.UI
	LoadoutTowers     [input.MaxToolbarSlots]string
	LoadoutSpells     [input.MaxToolbarSlots]string
	CodexTab          ui.CodexTab
	CodexScrollY      float64
	CoreHitFlash      float64
	ShopScrollY       float64
	ShopTargetScrollY float64
	TitleAnim         float64
	ShowHowToPlay     bool
	ScreenWidth       int
	ScreenHeight      int
}

func NewGameState(w, h int) *GameState {
	saveMgr := save.NewSaveManager()
	metaMgr := meta.NewMetaManager(saveMgr)
	keyMgr := input.NewKeybindManager()
	settingsMgr := config.NewSettingsManager()
	settingsMgr.Apply()
	registry := data.NewRegistry()

	gs := &GameState{
		Mode:         ModeTitle,
		PreviousMode: ModeTitle,
		SaveMgr:      saveMgr,
		KeyMgr:       keyMgr,
		Settings:     settingsMgr,
		Registry:     registry,
		Meta:         metaMgr,
		UI:           ui.NewUI(),
		ScreenWidth:  w,
		ScreenHeight: h,
		LoadoutTowers: [input.MaxToolbarSlots]string{
			"bit_driver", "mana_siphon", "crypto_miner", "tesla_bus", "cryo_cache",
		},
		LoadoutSpells: [input.MaxToolbarSlots]string{
			"chain_lightning", "chrono_freeze", "logic_bomb", "kernel_overclock", "emp_disruption",
		},
	}

	return gs
}

func (gs *GameState) LoadAndStartSlot(slotID int) {
	slotData, err := gs.SaveMgr.LoadSlot(slotID)
	if err != nil || slotData == nil {
		slotData = SlotDataDefault(slotID)
	}

	gs.Meta.ApplySlot(slotData)

	if slotData.HasActiveRun {
		gs.ResumeRun(slotData)
	} else {
		gs.StartNewRun()
	}
}

func SlotDataDefault(slotID int) *save.SlotData {
	return &save.SlotData{
		SlotID:       slotID,
		Talents:      make(map[string]int),
		HasActiveRun: false,
	}
}

func (gs *GameState) StartNewRun() {
	gs.Registry = data.NewRegistry()

	gs.Run = economy.NewRunState(
		gs.Meta.GetTalentLevel("kernel_shield"),
		gs.Meta.GetTalentLevel("boot_bytes"),
		gs.Meta.GetTalentLevel("mana_conductor"),
		gs.Meta.GetTalentLevel("overclock_nodes"),
		gs.Meta.GetTalentLevel("spell_efficiency"),
		gs.Meta.GetTalentLevel("scrap_leech"),
		gs.Meta.GetTalentLevel("tower_potency"),
		gs.Meta.GetTalentLevel("sensor_array"),
		gs.Meta.GetTalentLevel("nanite_resilience"),
	)

	gs.Grid = motherboard.NewGrid()
	gs.Malware = malware.NewSpawner(gs.Registry)
	gs.Towers = nodes.NewTowerManager(gs.Registry)
	gs.Spells = spells.NewSpellManager(gs.Registry)
	gs.Draft = draft.NewDraftManager(gs.Registry)
	gs.UI.HasTowerSelected = false
	gs.UI.HasSpellSelected = false
	gs.Mode = ModePlaying

	gs.SaveActiveRunState()
}

func (gs *GameState) ResumeRun(slotData *save.SlotData) {
	gs.Registry = data.NewRegistry()

	for _, tid := range slotData.ActiveRun.UnlockedTowers {
		if t := gs.Registry.GetTower(tid); t != nil {
			t.Unlocked = true
		}
	}
	for _, sid := range slotData.ActiveRun.UnlockedSpells {
		if s := gs.Registry.GetSpell(sid); s != nil {
			s.Unlocked = true
		}
	}

	gs.Run = &economy.RunState{
		Bytes:             slotData.ActiveRun.Bytes,
		TotalBytes:        slotData.ActiveRun.TotalBytes,
		Mana:              slotData.ActiveRun.Mana,
		MaxMana:           slotData.ActiveRun.MaxMana,
		KernelHP:          slotData.ActiveRun.KernelHP,
		MaxKernelHP:       slotData.ActiveRun.MaxKernelHP,
		Level:             slotData.ActiveRun.Level,
		CurrentXP:         slotData.ActiveRun.CurrentXP,
		TargetXP:          slotData.ActiveRun.TargetXP,
		PendingDrafts:     slotData.ActiveRun.PendingDrafts,
		RunTime:           slotData.ActiveRun.RunTime,
		Kills:             slotData.ActiveRun.Kills,
		BossKills:         slotData.ActiveRun.BossKills,
		TowerDamageMult:   slotData.ActiveRun.TowerDamageMult,
		TowerRangeMult:    slotData.ActiveRun.TowerRangeMult,
		TowerSpeedMult:    slotData.ActiveRun.TowerSpeedMult,
		ManaGenMult:       slotData.ActiveRun.ManaGenMult,
		ByteBountyMult:    slotData.ActiveRun.ByteBountyMult,
		SpellCooldownMult: math.Max(0.5, 1.0-float64(gs.Meta.GetTalentLevel("spell_efficiency"))*0.08),
		SpellCostMult:     math.Max(0.5, 1.0-float64(gs.Meta.GetTalentLevel("spell_efficiency"))*0.10),
	}

	gs.Grid = motherboard.NewGrid()
	gs.Malware = malware.NewSpawner(gs.Registry)
	gs.Malware.ElapsedRunTime = slotData.ActiveRun.RunTime
	gs.Malware.StartGraceTimer = 2.0

	gs.Towers = nodes.NewTowerManager(gs.Registry)
	for _, ts := range slotData.ActiveRun.Towers {
		if def := gs.Registry.GetTower(ts.TowerID); def != nil {
			gs.Grid.SetTower(ts.GridX, ts.GridY)
			wx, wy := gs.Grid.GridToScreenCenter(ts.GridX, ts.GridY)
			lvl := ts.Level
			if lvl < 1 {
				lvl = 1
			}
			tier := ts.Tier
			if tier < 1 {
				tier = 1
			}
			gs.Towers.Towers = append(gs.Towers.Towers, &nodes.Tower{
				ID:       ts.GridX*100 + ts.GridY,
				Def:      def,
				GridX:    ts.GridX,
				GridY:    ts.GridY,
				WorldX:   wx,
				WorldY:   wy,
				Tier:     tier,
				Level:    lvl,
				XP:       ts.XP,
				TargetXP: nodes.NextLevelXP(lvl),
			})
		}
	}
	gs.Towers.RecalculateBaseManaRate(gs.Run)

	gs.Spells = spells.NewSpellManager(gs.Registry)
	gs.Draft = draft.NewDraftManager(gs.Registry)
	gs.UI.HasTowerSelected = false
	gs.UI.HasSpellSelected = false
	gs.Mode = ModePlaying
}

func (gs *GameState) SaveActiveRunState() {
	if gs.Meta.CurrentSlotData == nil || gs.Run == nil {
		return
	}

	tSnapshots := make([]save.TowerSnapshot, 0, len(gs.Towers.Towers))
	for _, t := range gs.Towers.Towers {
		tSnapshots = append(tSnapshots, save.TowerSnapshot{
			TowerID: t.Def.ID,
			GridX:   t.GridX,
			GridY:   t.GridY,
			Tier:    t.Tier,
			Level:   t.Level,
			XP:      t.XP,
		})
	}

	unlockedTowers := make([]string, 0)
	for _, t := range gs.Registry.Towers {
		if t.Unlocked {
			unlockedTowers = append(unlockedTowers, t.ID)
		}
	}
	unlockedSpells := make([]string, 0)
	for _, s := range gs.Registry.Spells {
		if s.Unlocked {
			unlockedSpells = append(unlockedSpells, s.ID)
		}
	}

	gs.Meta.CurrentSlotData.HasActiveRun = true
	gs.Meta.CurrentSlotData.ActiveRun = save.RunSnapshot{
		Bytes:           gs.Run.Bytes,
		TotalBytes:      gs.Run.TotalBytes,
		Mana:            gs.Run.Mana,
		MaxMana:         gs.Run.MaxMana,
		KernelHP:        gs.Run.KernelHP,
		MaxKernelHP:     gs.Run.MaxKernelHP,
		Level:           gs.Run.Level,
		CurrentXP:       gs.Run.CurrentXP,
		TargetXP:        gs.Run.TargetXP,
		PendingDrafts:   gs.Run.PendingDrafts,
		RunTime:         gs.Run.RunTime,
		Kills:           gs.Run.Kills,
		BossKills:       gs.Run.BossKills,
		TowerDamageMult: gs.Run.TowerDamageMult,
		TowerRangeMult:  gs.Run.TowerRangeMult,
		TowerSpeedMult:  gs.Run.TowerSpeedMult,
		ManaGenMult:     gs.Run.ManaGenMult,
		ByteBountyMult:  gs.Run.ByteBountyMult,
		Towers:          tSnapshots,
		UnlockedTowers:  unlockedTowers,
		UnlockedSpells:  unlockedSpells,
	}

	_ = gs.Meta.SaveActiveSlot()
}

func (gs *GameState) ClearActiveRunSnapshot() {
	if gs.Meta.CurrentSlotData != nil {
		gs.Meta.CurrentSlotData.HasActiveRun = false
		_ = gs.Meta.SaveActiveSlot()
	}
}

func (gs *GameState) Update() error {
	cx, cy := ebiten.CursorPosition()
	justClicked := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	rightClicked := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight)

	// Intercept all inputs when actively rebinding a hotkey
	if gs.KeyMgr.IsRebinding {
		gs.KeyMgr.Update()
		return nil
	}

	switch gs.Mode {
	case ModeTitle:
		gs.TitleAnim += 1.0 / 60.0
		if gs.ShowHowToPlay {
			if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
				gs.ShowHowToPlay = false
			}
			if justClicked {
				if cx >= 280 && cx <= 520 && cy >= 490 && cy <= 530 {
					gs.ShowHowToPlay = false
				}
			}
		} else {
			if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
				gs.SaveMgr.RefreshSlots()
				gs.Mode = ModeSlotSelect
			}

			if justClicked {
				btnX := 260
				btnW := 280
				if cx >= btnX && cx <= btnX+btnW {
					if cy >= 255 && cy <= 299 {
						gs.SaveMgr.RefreshSlots()
						gs.Mode = ModeSlotSelect
					} else if cy >= 310 && cy <= 354 {
						gs.ShowHowToPlay = true
					} else if cy >= 365 && cy <= 409 {
						gs.PreviousMode = ModeTitle
						gs.CodexTab = ui.TabSettings
						gs.CodexScrollY = 0
						gs.Mode = ModeCodex
					} else if cy >= 420 && cy <= 464 {
						return ebiten.Termination
					}
				}
			}
		}

	case ModeSlotSelect:
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			gs.Mode = ModeTitle
		}
		gs.updateSlotSelect(cx, cy, justClicked)

	case ModePlaying:
		if rightClicked {
			gs.UI.HasTowerSelected = false
			gs.UI.HasSpellSelected = false
			gs.Spells.Deselect()
		} else if inpututil.IsKeyJustPressed(gs.KeyMgr.CodexKey) || inpututil.IsKeyJustPressed(ebiten.KeyL) {
			gs.PreviousMode = ModePlaying
			gs.Mode = ModeCodex
		} else if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			if gs.UI.HasTowerSelected || gs.UI.HasSpellSelected {
				gs.UI.HasTowerSelected = false
				gs.UI.HasSpellSelected = false
				gs.Spells.Deselect()
			} else {
				gs.SaveActiveRunState()
				gs.Mode = ModePause
			}
		} else {
			gs.updatePlaying(cx, cy, justClicked)
		}

	case ModeCodex:
		gs.updateCodex(cx, cy, justClicked)

	case ModePause:
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			gs.Mode = ModePlaying
		}

		if justClicked {
			if cx >= 260 && cx <= 540 {
				if cy >= 160 && cy <= 200 {
					gs.Mode = ModePlaying
				} else if cy >= 215 && cy <= 255 {
					gs.PreviousMode = ModePause
					gs.CodexTab = ui.TabSettings
					gs.CodexScrollY = 0
					gs.Mode = ModeCodex
				} else if cy >= 270 && cy <= 310 {
					gs.SaveActiveRunState()
					gs.SaveMgr.RefreshSlots()
					gs.Mode = ModeSlotSelect
				} else if cy >= 325 && cy <= 365 {
					gs.ClearActiveRunSnapshot()
					gs.StartNewRun()
				} else if cy >= 380 && cy <= 420 {
					gs.Mode = ModeMetaShop
				}
			}
		}

	case ModeDraft:
		if gs.Draft.Update(cx, cy, justClicked, gs.Run) {
			gs.SaveActiveRunState()
			if gs.Run.PendingDrafts > 0 {
				gs.Run.PendingDrafts--
				gs.Draft.GenerateDraft(gs.Run)
			} else {
				gs.Mode = ModePlaying
			}
		}

	case ModeGameOver:
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			gs.Mode = ModeMetaShop
		} else if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			gs.StartNewRun()
		}

		if justClicked {
			if cx >= 240 && cx <= 560 {
				if cy >= 390 && cy <= 430 {
					gs.Mode = ModeMetaShop
				} else if cy >= 445 && cy <= 485 {
					gs.StartNewRun()
				}
			}
		}

	case ModeMetaShop:
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			gs.StartNewRun()
		}

		cardH := 65
		gap := 10
		totalContentHeight := float64(len(gs.Meta.Talents)*(cardH+gap) + 20)
		viewportHeight := 445.0
		maxScroll := math.Max(0, totalContentHeight-viewportHeight)

		_, wy := ebiten.Wheel()
		if wy != 0 {
			gs.ShopTargetScrollY -= wy * 50.0
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
			gs.ShopTargetScrollY -= 12.0
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
			gs.ShopTargetScrollY += 12.0
		}

		if gs.ShopTargetScrollY < 0 {
			gs.ShopTargetScrollY = 0
		}
		if gs.ShopTargetScrollY > maxScroll {
			gs.ShopTargetScrollY = maxScroll
		}
		gs.ShopScrollY += (gs.ShopTargetScrollY - gs.ShopScrollY) * 0.25

		if justClicked {
			if cx >= 280 && cx <= 520 && cy >= 545 && cy <= 587 {
				gs.StartNewRun()
			} else if cy >= 75 && cy <= 530 {
				startY := 90.0
				for i, t := range gs.Meta.Talents {
					btnY := startY + float64(i*(cardH+gap)) + 12.0 - gs.ShopScrollY
					if cx >= 585 && cx <= 730 && float64(cy) >= btnY && float64(cy) <= btnY+40.0 {
						gs.Meta.BuyTalent(t.ID)
						break
					}
				}
			}
		}
	}

	return nil
}

func (gs *GameState) updateCodex(cx, cy int, justClicked bool) {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(gs.KeyMgr.CodexKey) {
		if gs.PreviousMode == ModeTitle || gs.PreviousMode == ModePause {
			gs.Mode = gs.PreviousMode
		} else {
			gs.Mode = ModePlaying
		}
		return
	}

	maxScroll := ui.GetCodexMaxScroll(gs.CodexTab, gs.Registry)

	// Wheel scrolling
	_, wy := ebiten.Wheel()
	if wy != 0 {
		gs.CodexScrollY -= wy * 40.0
	}

	// Keyboard scrolling (Arrow keys & W/S)
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		gs.CodexScrollY -= 6.0
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		gs.CodexScrollY += 6.0
	}

	// Clamp scroll strictly between 0 and maxScroll
	if gs.CodexScrollY < 0 {
		gs.CodexScrollY = 0
	}
	if gs.CodexScrollY > maxScroll {
		gs.CodexScrollY = maxScroll
	}

	if !justClicked {
		return
	}

	// 1. Close Button
	if cx >= 650 && cx <= 760 && cy >= 40 && cy <= 78 {
		if gs.PreviousMode == ModeTitle || gs.PreviousMode == ModePause {
			gs.Mode = gs.PreviousMode
		} else {
			gs.Mode = ModePlaying
		}
		return
	}

	// 2. Tabs Switcher (5 Tabs)
	for i := 0; i < 5; i++ {
		tx := 46 + i*(116+6)
		if cx >= tx && cx <= tx+116 && cy >= 40 && cy <= 78 {
			gs.CodexTab = ui.CodexTab(i)
			gs.CodexScrollY = 0
			return
		}
	}

	// 3. Tab Specific Interactions
	switch gs.CodexTab {
	case ui.TabTowers:
		startY := 100.0 - gs.CodexScrollY
		cardH := 72.0
		gap := 10.0
		for i, def := range gs.Registry.Towers {
			if !def.Unlocked {
				continue
			}
			cardY := startY + float64(i)*(cardH+gap)
			// Check Hotbar Slot buttons S1..S5
			for s := 0; s < input.MaxToolbarSlots; s++ {
				sbX := 515.0 + float64(s)*42.0
				sbY := cardY + 20.0
				if float64(cx) >= sbX && float64(cx) <= sbX+36.0 && float64(cy) >= sbY && float64(cy) <= sbY+32.0 {
					gs.LoadoutTowers[s] = def.ID
					return
				}
			}
		}

	case ui.TabSpells:
		startY := 100.0 - gs.CodexScrollY
		cardH := 72.0
		gap := 10.0
		for i, def := range gs.Registry.Spells {
			if !def.Unlocked {
				continue
			}
			cardY := startY + float64(i)*(cardH+gap)
			for s := 0; s < input.MaxToolbarSlots; s++ {
				sbX := 515.0 + float64(s)*42.0
				sbY := cardY + 20.0
				if float64(cx) >= sbX && float64(cx) <= sbX+36.0 && float64(cy) >= sbY && float64(cy) <= sbY+32.0 {
					gs.LoadoutSpells[s] = def.ID
					return
				}
			}
		}

	case ui.TabKeybinds:
		// Column 1: Tower slots
		for i := 0; i < input.MaxToolbarSlots; i++ {
			btnY := 195 + i*48
			if cx >= 290 && cx <= 365 && cy >= btnY && cy <= btnY+28 {
				gs.KeyMgr.StartRebind(fmt.Sprintf("tower_%d", i), fmt.Sprintf("Tower Slot %d", i+1))
				return
			}
		}

		// Column 2: Spell slots
		for i := 0; i < input.MaxToolbarSlots; i++ {
			btnY := 195 + i*48
			if cx >= 640 && cx <= 715 && cy >= btnY && cy <= btnY+28 {
				gs.KeyMgr.StartRebind(fmt.Sprintf("spell_%d", i), fmt.Sprintf("Spell Slot %d", i+1))
				return
			}
		}

		// Reset Defaults
		if cx >= 290 && cx <= 510 && cy >= 505 && cy <= 545 {
			gs.KeyMgr.ResetDefaults()
			_ = gs.KeyMgr.Save()
			return
		}

	case ui.TabSettings:
		startY := 160.0
		btnH := 48.0
		gap := 12.0
		// Preset Resolutions
		for i, res := range config.AvailableResolutions {
			ry := startY + float64(i)*(btnH+gap)
			if cx >= 90 && cx <= 710 && float64(cy) >= ry && float64(cy) <= ry+btnH {
				gs.Settings.SetResolution(res.Width, res.Height)
				return
			}
		}

		// Fullscreen Toggle
		fsY := startY + float64(len(config.AvailableResolutions))*(btnH+gap) + 10.0
		if cx >= 90 && cx <= 710 && float64(cy) >= fsY && float64(cy) <= fsY+btnH {
			gs.Settings.ToggleFullscreen()
			return
		}
	}
}

func (gs *GameState) updateSlotSelect(cx, cy int, justClicked bool) {
	if !justClicked {
		return
	}

	if cx >= 25 && cx <= 185 && cy >= 18 && cy <= 56 {
		gs.Mode = ModeTitle
		return
	}

	startY := 95
	cardH := 132
	gap := 16

	for i := 0; i < save.MaxSlots; i++ {
		slotID := i + 1
		cardY := startY + i*(cardH+gap)

		btnPlayY := cardY + 18
		if cx >= 560 && cx <= 740 && cy >= btnPlayY && cy <= btnPlayY+40 {
			gs.LoadAndStartSlot(slotID)
			return
		}

		if gs.SaveMgr.Slots[i].Exists {
			btnDelY := cardY + 70
			if cx >= 560 && cx <= 740 && cy >= btnDelY && cy <= btnDelY+38 {
				_ = gs.SaveMgr.DeleteSlot(slotID)
				return
			}
		}
	}
}

func (gs *GameState) updatePlaying(cx, cy int, justClicked bool) {
	dt := 1.0 / 60.0

	if gs.Malware.StartGraceTimer > 0 && inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		gs.Malware.StartGraceTimer = 0
	}

	// Dynamic Hotkeys for Tower Slots 1..5
	for i := 0; i < input.MaxToolbarSlots; i++ {
		if inpututil.IsKeyJustPressed(gs.KeyMgr.TowerKeys[i]) {
			if tID := gs.LoadoutTowers[i]; tID != "" {
				gs.selectTower(tID)
			}
		}
	}

	// Dynamic Hotkeys for Spell Slots 1..5
	for i := 0; i < input.MaxToolbarSlots; i++ {
		if inpututil.IsKeyJustPressed(gs.KeyMgr.SpellKeys[i]) {
			if sID := gs.LoadoutSpells[i]; sID != "" {
				gs.selectSpell(sID)
			}
		}
	}

	gs.Run.Update(dt)
	gs.Grid.Update(dt, cx, cy)
	if gs.CoreHitFlash > 0 {
		gs.CoreHitFlash -= dt
	}

	gameOver := false
	gs.Malware.ThreatLevel = gs.Towers.GetTotalOverclockTiers()
	gs.Malware.Update(dt, gs.Grid, func(e *malware.Enemy) {
		gs.CoreHitFlash = 0.25
		if gs.Run.DamageKernel(e.CoreDamage) {
			gameOver = true
		}
	})

	if gameOver {
		shards := gs.Run.CalculateFinalShards()
		gs.Meta.AddShards(shards)
		gs.Meta.TotalKills += gs.Run.Kills
		if gs.Run.RunTime > gs.Meta.HighScore {
			gs.Meta.HighScore = gs.Run.RunTime
		}
		gs.ClearActiveRunSnapshot()
		gs.Mode = ModeGameOver
		return
	}

	gs.Towers.Update(dt, gs.Malware.Enemies, gs.Run, gs.Malware)
	gs.Spells.Update(dt, gs.Run)

	if gs.Run.PendingDrafts > 0 {
		gs.Run.PendingDrafts--
		gs.Draft.GenerateDraft(gs.Run)
		gs.Mode = ModeDraft
		return
	}

	// KeyU hotkey for Upgrade All Promotable Towers
	if inpututil.IsKeyJustPressed(ebiten.KeyU) {
		if gs.Towers.PromoteAll(gs.Run) > 0 {
			gs.SaveActiveRunState()
		}
	}

	if justClicked {
		// Check Upgrade All button
		if count, _ := gs.Towers.GetPromotableInfo(); count > 0 {
			if cx >= 280 && cx <= 520 && cy >= 498 && cy <= 528 {
				if gs.Towers.PromoteAll(gs.Run) > 0 {
					gs.SaveActiveRunState()
				}
				return
			}
		}

		if cy >= 535 {
			// Check center Codex book button
			if cx >= 380 && cx <= 420 && cy >= 542 && cy <= 592 {
				gs.Mode = ModeCodex
				return
			}

			// Check Tower Slots [1..5]
			for i := 0; i < input.MaxToolbarSlots; i++ {
				bx := 14 + i*(68+6)
				if cx >= bx && cx <= bx+68 && cy >= 542 && cy <= 592 {
					if ebiten.IsKeyPressed(ebiten.KeyShift) {
						// Shift-Click opens instant rebind
						gs.KeyMgr.StartRebind(fmt.Sprintf("tower_%d", i), fmt.Sprintf("Tower Slot %d", i+1))
						return
					}
					if tID := gs.LoadoutTowers[i]; tID != "" {
						gs.selectTower(tID)
						return
					}
				}
			}

			// Check Spell Slots [1..5]
			for i := 0; i < input.MaxToolbarSlots; i++ {
				bx := 422 + i*(68+6)
				if cx >= bx && cx <= bx+68 && cy >= 542 && cy <= 592 {
					if ebiten.IsKeyPressed(ebiten.KeyShift) {
						gs.KeyMgr.StartRebind(fmt.Sprintf("spell_%d", i), fmt.Sprintf("Spell Slot %d", i+1))
						return
					}
					if sID := gs.LoadoutSpells[i]; sID != "" {
						gs.selectSpell(sID)
						return
					}
				}
			}
		} else {
			if gs.UI.HasTowerSelected {
				gx, gy, ok := gs.Grid.ScreenToGrid(cx, cy)
				if ok {
					if gs.Towers.BuildTower(gs.UI.SelectedTower, gx, gy, gs.Grid, gs.Run) {
						gs.SaveActiveRunState()
					}
				}
			} else if gs.UI.HasSpellSelected {
				if gs.Spells.CastSpell(gs.UI.SelectedSpell, float64(cx), float64(cy), gs.Malware.Enemies, gs.Run, gs.Malware) {
					gs.UI.HasSpellSelected = false
				}
			} else {
				gx, gy, ok := gs.Grid.ScreenToGrid(cx, cy)
				if ok {
					if gs.Towers.PromoteTowerAt(gx, gy, gs.Run) {
						gs.SaveActiveRunState()
					}
				}
			}
		}
	}
}

func (gs *GameState) selectTower(id string) {
	def := gs.Registry.GetTower(id)
	if def != nil && def.Unlocked {
		gs.UI.SelectedTower = id
		gs.UI.HasTowerSelected = true
		gs.UI.HasSpellSelected = false
	}
}

func (gs *GameState) selectSpell(id string) {
	def := gs.Registry.GetSpell(id)
	if def == nil || !def.Unlocked {
		gs.Spells.TriggerFlash(id)
		return
	}
	if gs.Spells.Cooldowns[id] > 0 || !gs.Run.CanAffordMana(def.ManaCost) {
		gs.Spells.TriggerFlash(id)
		return
	}

	gs.UI.SelectedSpell = id
	gs.UI.HasSpellSelected = true
	gs.UI.HasTowerSelected = false

	if def.CastType == data.CastGlobalBuff {
		gs.Spells.CastSpell(id, 400, 300, gs.Malware.Enemies, gs.Run, gs.Malware)
		gs.UI.HasSpellSelected = false
	}
}

func (gs *GameState) Draw(screen *ebiten.Image) {
	cx, cy := ebiten.CursorPosition()

	switch gs.Mode {
	case ModeTitle:
		gs.UI.DrawTitleScreen(screen, gs.TitleAnim, config.GetVersion(), cx, cy)
		if gs.ShowHowToPlay {
			gs.UI.DrawHowToPlay(screen)
		}

	case ModeSlotSelect:
		gs.UI.DrawSlotSelectScreen(screen, gs.SaveMgr, cx, cy)

	case ModePlaying, ModeDraft, ModePause:
		gs.Grid.Draw(screen)

		if gs.CoreHitFlash > 0 {
			coreX, coreY := gs.Grid.GridToScreenCenter(gs.Grid.CorePos.X, gs.Grid.CorePos.Y)
			flashR := float32(28.0 + (0.25-gs.CoreHitFlash)*40.0)
			vector.StrokeCircle(screen, float32(coreX), float32(coreY), flashR, 3, color.RGBA{R: 255, G: 50, B: 70, A: 240}, false)
			vector.FillCircle(screen, float32(coreX), float32(coreY), float32(motherboard.CellSize*0.6), color.RGBA{R: 255, G: 0, B: 50, A: 160}, false)
		}

		gs.Towers.Draw(screen)

		for _, e := range gs.Malware.Enemies {
			e.Draw(screen)
		}

		gs.Spells.Draw(screen)
		gs.UI.DrawPlacementPreview(screen, gs.Grid, gs.Towers, gs.Spells, gs.Run, cx, cy)
		gs.UI.DrawHUD(screen, gs.Run, gs.Malware.IsSurging, gs.Malware.SurgeTimeLeft, gs.Malware.StartGraceTimer, gs.Malware.WaveNumber, gs.Malware.CurrentPhase, gs.Malware.PhaseTimer, gs.Towers.GetTotalOverclockTiers())
		gs.UI.DrawToolbar(screen, gs.Towers, gs.Spells, gs.KeyMgr, gs.LoadoutTowers, gs.LoadoutSpells, gs.Run, cx, cy)

		if gs.Mode == ModeDraft {
			gs.Draft.Draw(screen)
		} else if gs.Mode == ModePause {
			gs.UI.DrawPauseMenu(screen)
		}

	case ModeCodex:
		if gs.Grid != nil && gs.Towers != nil && gs.Malware != nil && gs.Run != nil {
			gs.Grid.Draw(screen)
			gs.Towers.Draw(screen)
			for _, e := range gs.Malware.Enemies {
				e.Draw(screen)
			}
			gs.Spells.Draw(screen)
			gs.UI.DrawHUD(screen, gs.Run, gs.Malware.IsSurging, gs.Malware.SurgeTimeLeft, gs.Malware.StartGraceTimer, gs.Malware.WaveNumber, gs.Malware.CurrentPhase, gs.Malware.PhaseTimer, gs.Towers.GetTotalOverclockTiers())
			gs.UI.DrawToolbar(screen, gs.Towers, gs.Spells, gs.KeyMgr, gs.LoadoutTowers, gs.LoadoutSpells, gs.Run, cx, cy)
		} else {
			gs.UI.DrawTitleScreen(screen, gs.TitleAnim, config.GetVersion(), cx, cy)
		}
		gs.UI.DrawCodex(screen, gs.Registry, gs.KeyMgr, gs.Settings, gs.LoadoutTowers, gs.LoadoutSpells, gs.CodexTab, "", gs.CodexScrollY, cx, cy)

	case ModeGameOver:
		gs.Grid.Draw(screen)
		gs.Towers.Draw(screen)
		for _, e := range gs.Malware.Enemies {
			e.Draw(screen)
		}
		gs.UI.DrawGameOver(screen, gs.Run, gs.Meta)

	case ModeMetaShop:
		gs.UI.DrawMetaShop(screen, gs.Meta, gs.ShopScrollY, cx, cy)
	}

	// Always render rebind modal if active
	if gs.KeyMgr.IsRebinding {
		gs.UI.DrawRebindPrompt(screen, gs.KeyMgr)
	}
}

func (gs *GameState) Layout(outsideWidth, outsideHeight int) (int, int) {
	return gs.ScreenWidth, gs.ScreenHeight
}
