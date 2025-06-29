package visual

import (
	"christofides-algo/algorithm"
	buildgraph "christofides-algo/build_graph"
	"christofides-algo/model"
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	Graph      model.Graph
	Points     []point
	Step       Step
	MST        []model.Edge
	Matching   []model.Edge
	Tour       []int
	KeyPressed bool
}

type point struct {
	x, y float32
}

func NewGame(n int) *Game {
	g := buildgraph.BuildNewGraph(n)

	// random coords for display (same size as node count)
	points := make([]point, n)
	for i := 0; i < n; i++ {
		points[i] = point{
			x: rand.Float32()*600 + 50,
			y: rand.Float32()*400 + 50,
		}
	}

	return &Game{
		Graph:  g,
		Points: points,
		Step:   Idle,
	}
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyRight) && !g.KeyPressed {
		g.Step = (g.Step + 1) % Done
		g.KeyPressed = true

		switch g.Step {
		case MST:
			g.MST = algorithm.BuildMST(g.Graph)
		case Matching:
			odds := algorithm.FindOddDegreeNodes(g.MST, g.Graph.Nodes)
			g.Matching = algorithm.MinPerfectMatching(odds, g.Graph)
		case Tour:
			t, _ := algorithm.Christofides(g.Graph)
			g.Tour = t
		}
	}
	if !ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.KeyPressed = false
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Kanten
	for i := 0; i < g.Graph.Nodes; i++ {
		for j := i + 1; j < g.Graph.Nodes; j++ {
			vector.StrokeLine(screen, g.Points[i].x, g.Points[i].y, g.Points[j].x, g.Points[j].y, 1, color.RGBA{80, 80, 80, 255}, false)
		}
	}

	// MST-Kanten
	if g.Step >= MST {
		for _, e := range g.MST {
			vector.StrokeLine(screen, g.Points[e.From].x, g.Points[e.From].y, g.Points[e.To].x, g.Points[e.To].y, 1, color.RGBA{0, 200, 0, 255}, false)
		}
	}

	// Matching-Kanten
	if g.Step >= Matching {
		for _, e := range g.Matching {
			vector.StrokeLine(screen, g.Points[e.From].x, g.Points[e.From].y, g.Points[e.To].x, g.Points[e.To].y, 1, color.RGBA{200, 0, 200, 255}, false)
		}
	}

	// Tour
	if g.Step >= Tour {
		for i := 0; i < len(g.Tour)-1; i++ {
			from := g.Tour[i]
			to := g.Tour[i+1]
			vector.StrokeLine(screen, g.Points[from].x, g.Points[from].y, g.Points[to].x, g.Points[to].y, 1, color.RGBA{0, 100, 255, 255}, false)
		}
	}

	// Knoten
	for i, p := range g.Points {
		vector.StrokeLine(screen, p.x-3, p.y-3, 6, 6, 1, color.RGBA{0, 100, 255, 255}, false)
		ebitenutil.DebugPrintAt(screen, string(rune('A'+i)), int(p.x+5), int(p.y+5))
	}

	// Text
	ebitenutil.DebugPrint(screen, "→: Nächster Schritt | Schritte: Kanten, MST, Matching, Tour")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 800, 600
}
