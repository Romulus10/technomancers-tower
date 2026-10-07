package motherboard_test

import (
	"testing"

	"technomancers-tower/internal/motherboard"
)

func TestGridFlowfieldAndBlocking(t *testing.T) {
	grid := motherboard.NewGrid()

	// Initial pathfinding validation from spawners to core
	if len(grid.Spawners) == 0 {
		t.Fatalf("expected spawners on grid")
	}

	for _, sp := range grid.Spawners {
		if grid.DistMatrix[sp.X][sp.Y] >= 9999 {
			t.Errorf("spawner at (%d, %d) is unreachable from core", sp.X, sp.Y)
		}
	}

	// Verify valid tower placement
	canPlace := grid.CanBuildAt(2, 2)
	if !canPlace {
		t.Errorf("expected cell (2, 2) to be buildable")
	}

	placed := grid.SetTower(2, 2)
	if !placed || grid.Cells[2][2] != motherboard.CellTower {
		t.Errorf("failed to set tower at (2, 2)")
	}

	// Cannot place on Core
	if grid.CanBuildAt(grid.CorePos.X, grid.CorePos.Y) {
		t.Errorf("should not be able to place tower on Core")
	}

	// Cannot place on Spawners
	sp := grid.Spawners[0]
	if grid.CanBuildAt(sp.X, sp.Y) {
		t.Errorf("should not be able to place tower on Spawner")
	}
}

func BenchmarkFlowfieldCalculation(b *testing.B) {
	grid := motherboard.NewGrid()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		grid.RecomputeFlowfield()
	}
}
