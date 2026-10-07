package motherboard

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	Cols     = 20
	Rows     = 14
	CellSize = 34
	OffsetX  = 35
	OffsetY  = 65
)

// CellType indicates what occupies a grid cell.
type CellType int

const (
	CellEmpty CellType = iota
	CellCore
	CellSpawner
	CellTower
)

// GridPoint is a 2D integer coordinate.
type GridPoint struct {
	X, Y int
}

// Direction vector for flowfield.
type Direction struct {
	DX, DY float64
}

// Grid manages the Motherboard Matrix, pathfinding, and placement validation.
type Grid struct {
	Cells       [Cols][Rows]CellType
	CorePos     GridPoint
	Spawners    []GridPoint
	DistMatrix  [Cols][Rows]int
	FlowMatrix  [Cols][Rows]Direction
	HoverCell   GridPoint
	HasHover    bool
	TraceOffset float64
}

func NewGrid() *Grid {
	g := &Grid{
		CorePos: GridPoint{X: Cols / 2, Y: Rows / 2},
		Spawners: []GridPoint{
			{X: Cols / 2, Y: 0},        // North Port
			{X: Cols / 2, Y: Rows - 1}, // South Port
			{X: 0, Y: Rows / 2},        // West Port
			{X: Cols - 1, Y: Rows / 2}, // East Port
		},
	}

	// Mark Core
	g.Cells[g.CorePos.X][g.CorePos.Y] = CellCore

	// Mark Spawners
	for _, sp := range g.Spawners {
		g.Cells[sp.X][sp.Y] = CellSpawner
	}

	g.RecomputeFlowfield()
	return g
}

func (g *Grid) InBounds(x, y int) bool {
	return x >= 0 && x < Cols && y >= 0 && y < Rows
}

func (g *Grid) ScreenToGrid(sx, sy int) (int, int, bool) {
	if sx < OffsetX || sy < OffsetY {
		return -1, -1, false
	}
	gx := (sx - OffsetX) / CellSize
	gy := (sy - OffsetY) / CellSize
	if g.InBounds(gx, gy) {
		return gx, gy, true
	}
	return -1, -1, false
}

func (g *Grid) GridToScreenCenter(gx, gy int) (float64, float64) {
	sx := float64(OffsetX + gx*CellSize + CellSize/2)
	sy := float64(OffsetY + gy*CellSize + CellSize/2)
	return sx, sy
}

// CanBuildAt checks if a cell is valid for building.
func (g *Grid) CanBuildAt(gx, gy int) bool {
	if !g.InBounds(gx, gy) {
		return false
	}
	if g.Cells[gx][gy] != CellEmpty {
		return false
	}

	// Temporarily place tower and check if flowfield connects all spawners to the core
	g.Cells[gx][gy] = CellTower
	valid := g.validateConnectivity()
	g.Cells[gx][gy] = CellEmpty

	return valid
}

func (g *Grid) validateConnectivity() bool {
	dist := g.calculateBFS()
	for _, sp := range g.Spawners {
		if dist[sp.X][sp.Y] >= 9999 {
			return false
		}
	}
	return true
}

func (g *Grid) calculateBFS() [Cols][Rows]int {
	var dist [Cols][Rows]int
	for x := 0; x < Cols; x++ {
		for y := 0; y < Rows; y++ {
			dist[x][y] = 9999
		}
	}

	queue := []GridPoint{g.CorePos}
	dist[g.CorePos.X][g.CorePos.Y] = 0

	dirs := []GridPoint{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		for _, d := range dirs {
			nx, ny := curr.X+d.X, curr.Y+d.Y
			if g.InBounds(nx, ny) && g.Cells[nx][ny] != CellTower {
				if dist[nx][ny] > dist[curr.X][curr.Y]+1 {
					dist[nx][ny] = dist[curr.X][curr.Y] + 1
					queue = append(queue, GridPoint{nx, ny})
				}
			}
		}
	}

	return dist
}

// RecomputeFlowfield updates distances and path vector direction fields for enemies.
func (g *Grid) RecomputeFlowfield() {
	g.DistMatrix = g.calculateBFS()

	dirs := []GridPoint{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}

	for x := 0; x < Cols; x++ {
		for y := 0; y < Rows; y++ {
			if g.Cells[x][y] == CellCore {
				g.FlowMatrix[x][y] = Direction{0, 0}
				continue
			}

			bestDist := g.DistMatrix[x][y]
			bestDir := Direction{0, 0}

			for _, d := range dirs {
				nx, ny := x+d.X, y+d.Y
				if g.InBounds(nx, ny) && g.Cells[nx][ny] != CellTower {
					if g.DistMatrix[nx][ny] < bestDist {
						bestDist = g.DistMatrix[nx][ny]
						bestDir = Direction{DX: float64(d.X), DY: float64(d.Y)}
					}
				}
			}

			g.FlowMatrix[x][y] = bestDir
		}
	}
}

// SetTower places a tower at (gx, gy) and recomputes the flowfield.
func (g *Grid) SetTower(gx, gy int) bool {
	if !g.CanBuildAt(gx, gy) {
		return false
	}
	g.Cells[gx][gy] = CellTower
	g.RecomputeFlowfield()
	return true
}

// RemoveTower clears a tower at (gx, gy) and recomputes the flowfield.
func (g *Grid) RemoveTower(gx, gy int) {
	if g.InBounds(gx, gy) && g.Cells[gx][gy] == CellTower {
		g.Cells[gx][gy] = CellEmpty
		g.RecomputeFlowfield()
	}
}

func (g *Grid) Update(dt float64, cursorX, cursorY int) {
	g.TraceOffset += dt * 40.0
	gx, gy, ok := g.ScreenToGrid(cursorX, cursorY)
	g.HasHover = ok
	if ok {
		g.HoverCell = GridPoint{gx, gy}
	}
}

func (g *Grid) Draw(screen *ebiten.Image) {
	// Motherboard Background
	boardW := float32(Cols * CellSize)
	boardH := float32(Rows * CellSize)
	vector.FillRect(screen, OffsetX-4, OffsetY-4, boardW+8, boardH+8, color.RGBA{R: 12, G: 16, B: 24, A: 255}, false)
	vector.StrokeRect(screen, OffsetX-4, OffsetY-4, boardW+8, boardH+8, 2, color.RGBA{R: 0, G: 180, B: 200, A: 255}, false)

	// Draw Grid Lines & Traces
	for x := 0; x < Cols; x++ {
		for y := 0; y < Rows; y++ {
			sx := float32(OffsetX + x*CellSize)
			sy := float32(OffsetY + y*CellSize)
			sz := float32(CellSize)

			// Subtle grid cell
			vector.StrokeRect(screen, sx, sy, sz, sz, 1, color.RGBA{R: 20, G: 32, B: 48, A: 255}, false)
		}
	}

	// Draw Flowfield Guide Arrows / Circuit Traces
	for x := 0; x < Cols; x++ {
		for y := 0; y < Rows; y++ {
			if g.Cells[x][y] == CellEmpty {
				dir := g.FlowMatrix[x][y]
				if dir.DX != 0 || dir.DY != 0 {
					cx := float32(OffsetX + x*CellSize + CellSize/2)
					cy := float32(OffsetY + y*CellSize + CellSize/2)
					tx := cx + float32(dir.DX*8)
					ty := cy + float32(dir.DY*8)
					vector.StrokeLine(screen, cx, cy, tx, ty, 1.2, color.RGBA{R: 28, G: 65, B: 85, A: 160}, false)
				}
			}
		}
	}

	// Draw Perimeter Ports (Spawners)
	for _, sp := range g.Spawners {
		sx := float32(OffsetX + sp.X*CellSize)
		sy := float32(OffsetY + sp.Y*CellSize)
		sz := float32(CellSize)
		vector.FillRect(screen, sx+2, sy+2, sz-4, sz-4, color.RGBA{R: 200, G: 40, B: 60, A: 200}, false)
		vector.StrokeRect(screen, sx+2, sy+2, sz-4, sz-4, 1.5, color.RGBA{R: 255, G: 100, B: 120, A: 255}, false)
	}

	// Draw Kernel CPU Core
	csx := float32(OffsetX + g.CorePos.X*CellSize)
	csy := float32(OffsetY + g.CorePos.Y*CellSize)
	csz := float32(CellSize)
	vector.FillRect(screen, csx+2, csy+2, csz-4, csz-4, color.RGBA{R: 0, G: 160, B: 240, A: 255}, false)
	vector.StrokeRect(screen, csx+1, csy+1, csz-2, csz-2, 2, color.RGBA{R: 160, G: 240, B: 255, A: 255}, false)

	// Draw Core Pulse
	pulse := float32((math.Sin(g.TraceOffset*0.1) + 1.0) * 0.5)
	vector.FillCircle(screen, csx+csz/2, csy+csz/2, 6+pulse*3, color.RGBA{R: 255, G: 255, B: 255, A: 230}, false)
}
