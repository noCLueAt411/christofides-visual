package visual

import (
	"christofides-algo/algorithm"
	buildgraph "christofides-algo/build_graph"
	"christofides-algo/model"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	Graph         model.Graph
	Width, Height int
	Step          Step
	MST           []model.Edge
	Matching      []model.Edge
	Tour          []int
	KeyPressed    bool
}

func NewGame(n int, windowWidth, windowHeight float32) *Game {
	g := buildgraph.BuildNewGraph(n, windowWidth, windowHeight)

	return &Game{
		Graph:  g,
		Width:  int(windowWidth),
		Height: int(windowHeight),
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
			odds := algorithm.FindOddDegreeNodes(g.MST, g.Graph.Points)
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
	// Alle Kanten (nur im Idle-Zustand)
	if g.Step == Idle {
		for i := 0; i < g.Graph.Nodes; i++ {
			for j := i + 1; j < g.Graph.Nodes; j++ {
				p1 := g.Graph.Points[i]
				p2 := g.Graph.Points[j]
				vector.StrokeLine(screen, p1.X, p1.Y, p2.X, p2.Y, 1, color.RGBA{80, 80, 80, 255}, false)
			}
		}
	}

	// MST-Kanten
	if g.Step == MST || g.Step == Matching {
		for _, e := range g.MST {
			vector.StrokeLine(screen, e.From.X, e.From.Y, e.To.X, e.To.Y, 1, color.RGBA{0, 200, 0, 255}, false)
		}
	}

	// Matching-Kanten
	if g.Step == Matching {
		for _, e := range g.Matching {
			vector.StrokeLine(screen, e.From.X, e.From.Y, e.To.X, e.To.Y, 1, color.RGBA{200, 0, 200, 255}, false)
		}
	}

	// Tour
	if g.Step >= Tour {
		for i := 0; i < len(g.Tour)-1; i++ {
			from := g.Graph.Points[g.Tour[i]]
			to := g.Graph.Points[g.Tour[i+1]]
			vector.StrokeLine(screen, from.X, from.Y, to.X, to.Y, 1, color.RGBA{0, 100, 255, 255}, false)
		}
	}

	// Knoten
	for i, p := range g.Graph.Points {
		vector.DrawFilledCircle(screen, p.X, p.Y, 3, color.RGBA{0, 100, 255, 255}, false)
		ebitenutil.DebugPrintAt(screen, string(rune('A'+i)), int(p.X+5), int(p.Y+5))
	}

	// Schrittanzeige
	steps := []string{"MST", "Matching", "Tour"}
	for i, s := range steps {
		prefix := "   "
		if Step(i+1) == g.Step {
			prefix = "→ "
		}
		ebitenutil.DebugPrintAt(screen, prefix+s, 10, 20+i*15)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.Width, g.Height
}
