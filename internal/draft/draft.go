package draft

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"strings"

	"technomancers-tower/internal/data"
	"technomancers-tower/internal/economy"
	"technomancers-tower/internal/gfx"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// WrapText splits text into multiple lines so that no line exceeds maxLineLen characters.
func WrapText(text string, maxLineLen int) []string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return nil
	}

	var lines []string
	currentLine := words[0]

	for _, w := range words[1:] {
		if len(currentLine)+1+len(w) <= maxLineLen {
			currentLine += " " + w
		} else {
			lines = append(lines, currentLine)
			currentLine = w
		}
	}
	lines = append(lines, currentLine)
	return lines
}

type DraftManager struct {
	OfferedCards []*data.CardDef
	Active       bool
	HoverIndex   int
	Registry     *data.Registry
}

func NewDraftManager(reg *data.Registry) *DraftManager {
	return &DraftManager{
		OfferedCards: make([]*data.CardDef, 0),
		HoverIndex:   -1,
		Registry:     reg,
	}
}

func (dm *DraftManager) GenerateDraft(run *economy.RunState) {
	pool := make([]*data.CardDef, 0)

	for _, card := range dm.Registry.Cards {
		if card.IsUnlock {
			// Check if already unlocked
			alreadyUnlocked := false
			for _, eff := range card.Effects {
				if eff.Type == data.ModUnlockTower {
					if t := dm.Registry.GetTower(eff.TargetID); t != nil && t.Unlocked {
						alreadyUnlocked = true
					}
				} else if eff.Type == data.ModUnlockSpell {
					if s := dm.Registry.GetSpell(eff.TargetID); s != nil && s.Unlocked {
						alreadyUnlocked = true
					}
				}
			}
			if !alreadyUnlocked {
				pool = append(pool, card)
			}
		} else {
			// Standard passive / system tweak
			pool = append(pool, card)
		}
	}

	choicesCount := 3
	if run != nil && run.DraftCardChoices > 0 {
		choicesCount = run.DraftCardChoices
	}

	if len(pool) <= choicesCount {
		dm.OfferedCards = pool
	} else {
		rand.Shuffle(len(pool), func(i, j int) { pool[i], pool[j] = pool[j], pool[i] })
		dm.OfferedCards = pool[:choicesCount]
	}

	dm.Active = true
	dm.HoverIndex = -1
}

func (dm *DraftManager) Reroll(run *economy.RunState) bool {
	if run == nil || run.DraftRerolls <= 0 {
		return false
	}
	run.DraftRerolls--
	dm.GenerateDraft(run)
	return true
}

func (dm *DraftManager) ApplyCard(card *data.CardDef, run *economy.RunState) {
	for _, eff := range card.Effects {
		switch eff.Type {
		case data.ModUnlockTower:
			if t := dm.Registry.GetTower(eff.TargetID); t != nil {
				t.Unlocked = true
			}
		case data.ModUnlockSpell:
			if s := dm.Registry.GetSpell(eff.TargetID); s != nil {
				s.Unlocked = true
			}
		case data.ModTowerDamageMult:
			run.TowerDamageMult += eff.Value
		case data.ModTowerRangeMult:
			run.TowerRangeMult += eff.Value
		case data.ModTowerSpeedMult:
			run.TowerSpeedMult += eff.Value
		case data.ModManaGenMult:
			run.ManaGenMult += eff.Value
		case data.ModByteBountyMult:
			run.ByteBountyMult += eff.Value
		case data.ModKernelHeal:
			run.KernelHP = math.Min(run.MaxKernelHP, run.KernelHP+eff.Value)
			run.Bytes += 50
		}
	}
}

func (dm *DraftManager) getCardLayout() (cardW, cardH, gap, startX, startY float32) {
	numCards := len(dm.OfferedCards)
	startY = float32(145)
	if numCards <= 3 {
		cardW = 210
		cardH = 310
		gap = 25
	} else if numCards == 4 {
		cardW = 175
		cardH = 310
		gap = 18
	} else {
		cardW = 142
		cardH = 310
		gap = 12
	}
	totalW := float32(numCards)*cardW + float32(math.Max(0, float64(numCards-1)))*gap
	startX = (float32(800) - totalW) / 2
	return
}

func (dm *DraftManager) Update(cursorX, cursorY int, justClicked bool, run *economy.RunState) bool {
	if !dm.Active {
		return false
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyR) && run != nil && run.DraftRerolls > 0 {
		dm.Reroll(run)
		return false
	}

	// Check Reroll button click
	if run != nil && run.DraftRerolls > 0 {
		btnX, btnY, btnW, btnH := float32(310), float32(480), float32(180), float32(34)
		if float32(cursorX) >= btnX && float32(cursorX) <= btnX+btnW && float32(cursorY) >= btnY && float32(cursorY) <= btnY+btnH {
			if justClicked {
				dm.Reroll(run)
				return false
			}
		}
	}

	dm.HoverIndex = -1
	cardW, cardH, gap, startX, startY := dm.getCardLayout()

	for i, card := range dm.OfferedCards {
		cx := startX + float32(i)*(cardW+gap)
		cy := startY

		if float32(cursorX) >= cx && float32(cursorX) <= cx+cardW && float32(cursorY) >= cy && float32(cursorY) <= cy+cardH {
			dm.HoverIndex = i
			if justClicked {
				dm.ApplyCard(card, run)
				dm.Active = false
				return true
			}
		}
	}

	return false
}

func (dm *DraftManager) Draw(screen *ebiten.Image, run *economy.RunState) {
	if !dm.Active {
		return
	}

	// Dim holographic background
	vector.FillRect(screen, 0, 0, 800, 600, color.RGBA{R: 5, G: 8, B: 15, A: 235}, false)

	// Header Banner Box
	gfx.DrawHolographicPanel(screen, 180, 75, 440, 42, gfx.ColorCyanNeon, gfx.ColorPanelDark)
	ebitenutil.DebugPrintAt(screen, "=== RUNTIME LEVEL UP: SELECT AN UPGRADE ===", 245, 88)

	cardW, cardH, gap, startX, startY := dm.getCardLayout()
	numCards := len(dm.OfferedCards)
	maxCharLen := 25
	if numCards >= 4 {
		maxCharLen = 20
	}

	for i, card := range dm.OfferedCards {
		cx := startX + float32(i)*(cardW+gap)
		cy := startY

		isHovered := dm.HoverIndex == i
		if isHovered {
			cy -= 8 // slight lift on hover
		}

		// Card background & holographic frame
		bgColor := color.RGBA{R: 18, G: 24, B: 36, A: 255}
		borderColor := gfx.FadeAlpha(card.Color, 0.7)
		if isHovered {
			bgColor = color.RGBA{R: 28, G: 38, B: 58, A: 255}
			borderColor = gfx.Brighten(card.Color, 0.3)
		}

		gfx.DrawHolographicPanel(screen, cx, cy, cardW, cardH, borderColor, bgColor)

		// Top Banner
		vector.FillRect(screen, cx+2, cy+2, cardW-4, 30, color.RGBA{R: card.Color.R / 3, G: card.Color.G / 3, B: card.Color.B / 3, A: 255}, false)
		ebitenutil.DebugPrintAt(screen, card.Title, int(cx)+10, int(cy)+10)

		// Procedural Icon Badge if card has TargetID
		if len(card.Effects) > 0 && card.Effects[0].TargetID != "" {
			targetID := card.Effects[0].TargetID
			isTower := card.Effects[0].Type == data.ModUnlockTower
			icon := gfx.GetCache().GetIconBadge(targetID, isTower, card.Color)
			if icon != nil {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(cx+cardW-34), float64(cy+3))
				screen.DrawImage(icon, op)
			}
		}

		// Card Subtitle
		subtitleLines := WrapText(card.Subtitle, maxCharLen)
		subY := int(cy) + 48
		for sIdx, sLine := range subtitleLines {
			ebitenutil.DebugPrintAt(screen, sLine, int(cx)+10, subY+sIdx*15)
		}

		// Card Description (word-wrapped)
		descLines := WrapText(card.Description, maxCharLen)
		descY := int(cy) + 95
		for dIdx, dLine := range descLines {
			ebitenutil.DebugPrintAt(screen, dLine, int(cx)+10, descY+dIdx*16)
		}

		// Click prompt
		if isHovered {
			promptBg := color.RGBA{R: 20, G: 70, B: 110, A: 240}
			btnY := cy + cardH - 45
			vector.FillRect(screen, cx+10, btnY, cardW-20, 32, promptBg, false)
			vector.StrokeRect(screen, cx+10, btnY, cardW-20, 32, 1.5, gfx.ColorCyanNeon, false)
			label := "[ INSTALL ]"
			if numCards <= 3 {
				label = "[ INSTALL UPGRADE ]"
			}
			ebitenutil.DebugPrintAt(screen, label, int(cx)+int(cardW/2)-len(label)*3-10, int(btnY)+10)
		}
	}

	// Draw Reroll Button if charges available
	if run != nil && run.DraftRerolls > 0 {
		btnX, btnY, btnW, btnH := float32(310), float32(480), float32(180), float32(34)
		gfx.DrawHolographicPanel(screen, btnX, btnY, btnW, btnH, gfx.ColorGoldMatrix, gfx.ColorPanelDark)
		rerollText := fmt.Sprintf("[R] REROLL (%d Left)", run.DraftRerolls)
		ebitenutil.DebugPrintAt(screen, rerollText, int(btnX)+20, int(btnY)+10)
	}
}
