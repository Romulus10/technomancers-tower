package draft

import (
	"image/color"
	"math"
	"math/rand"
	"strings"

	"technomancers-tower/internal/data"
	"technomancers-tower/internal/economy"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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

	if len(pool) < 3 {
		dm.OfferedCards = pool
	} else {
		rand.Shuffle(len(pool), func(i, j int) { pool[i], pool[j] = pool[j], pool[i] })
		dm.OfferedCards = pool[:3]
	}

	dm.Active = true
	dm.HoverIndex = -1
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

func (dm *DraftManager) Update(cursorX, cursorY int, justClicked bool, run *economy.RunState) bool {
	if !dm.Active {
		return false
	}

	dm.HoverIndex = -1
	cardW := float32(210)
	cardH := float32(310)
	startY := float32(145)
	gap := float32(30)
	totalW := float32(len(dm.OfferedCards))*cardW + float32(len(dm.OfferedCards)-1)*gap
	startX := (float32(800) - totalW) / 2

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

func (dm *DraftManager) Draw(screen *ebiten.Image) {
	if !dm.Active {
		return
	}

	// Dim background
	vector.FillRect(screen, 0, 0, 800, 600, color.RGBA{R: 5, G: 8, B: 15, A: 220}, false)

	// Header
	ebitenutil.DebugPrintAt(screen, "=== RUNTIME LEVEL UP: SELECT AN UPGRADE ===", 245, 105)

	cardW := float32(210)
	cardH := float32(310)
	startY := float32(145)
	gap := float32(30)
	totalW := float32(len(dm.OfferedCards))*cardW + float32(len(dm.OfferedCards)-1)*gap
	startX := (float32(800) - totalW) / 2

	for i, card := range dm.OfferedCards {
		cx := startX + float32(i)*(cardW+gap)
		cy := startY

		isHovered := dm.HoverIndex == i
		if isHovered {
			cy -= 8 // slight lift on hover
		}

		// Card background
		bgColor := color.RGBA{R: 18, G: 24, B: 36, A: 255}
		if isHovered {
			bgColor = color.RGBA{R: 26, G: 36, B: 55, A: 255}
		}
		vector.FillRect(screen, cx, cy, cardW, cardH, bgColor, false)

		// Border
		borderCol := card.Color
		if !isHovered {
			borderCol = color.RGBA{R: card.Color.R / 2, G: card.Color.G / 2, B: card.Color.B / 2, A: 200}
		}
		vector.StrokeRect(screen, cx, cy, cardW, cardH, 2, borderCol, false)

		// Top Banner
		vector.FillRect(screen, cx+2, cy+2, cardW-4, 30, color.RGBA{R: card.Color.R / 3, G: card.Color.G / 3, B: card.Color.B / 3, A: 255}, false)
		ebitenutil.DebugPrintAt(screen, card.Title, int(cx)+15, int(cy)+10)

		// Card Subtitle (wrapped if needed)
		subtitleLines := WrapText(card.Subtitle, 25)
		subY := int(cy) + 48
		for sIdx, sLine := range subtitleLines {
			ebitenutil.DebugPrintAt(screen, sLine, int(cx)+15, subY+sIdx*15)
		}

		// Card Description (word-wrapped to fit card width)
		descLines := WrapText(card.Description, 26)
		descY := int(cy) + 95
		for dIdx, dLine := range descLines {
			ebitenutil.DebugPrintAt(screen, dLine, int(cx)+15, descY+dIdx*16)
		}

		// Click prompt
		if isHovered {
			ebitenutil.DebugPrintAt(screen, "[ CLICK TO INSTALL ]", int(cx)+35, int(cy)+270)
		}
	}
}
