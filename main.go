package main

import (
	"fmt"
	"log"

	"technomancers-tower/internal/config"
	"technomancers-tower/internal/state"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	WindowWidth  = 800
	WindowHeight = 600
)

func main() {
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ver := config.GetVersion()
	ebiten.SetWindowTitle(fmt.Sprintf("Technomancer's Tower [%s] - Cyber-Arcane Roguelite TD", ver))

	gameState := state.NewGameState(WindowWidth, WindowHeight)
	if err := ebiten.RunGame(gameState); err != nil {
		log.Fatal(err)
	}
}
