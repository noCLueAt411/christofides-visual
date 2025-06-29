package buildgraph

import (
	"christofides-algo/model"
	"fmt"
	"math"
	"math/rand"
	"os"

	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"
)

type point struct {
	x, y float64
}

// BuildNewGraph erzeugt einen vollständigen, symmetrischen, metrischen Graph
// und gibt ihn als model.Graph zurück
func BuildNewGraph(numberOfNodes int) model.Graph {
	points := make([]point, numberOfNodes)
	for i := 0; i < numberOfNodes; i++ {
		points[i] = point{
			x: rand.Float64() * 100,
			y: rand.Float64() * 100,
		}
	}

	edges64 := make([][]float64, numberOfNodes)
	for i := 0; i < numberOfNodes; i++ {
		edges64[i] = make([]float64, numberOfNodes)
		for j := 0; j < numberOfNodes; j++ {
			if i == j {
				edges64[i][j] = 0
			} else {
				dist := euclideanDistance(points[i], points[j])
				edges64[i][j] = dist
			}
		}
	}

	// Symmetrie
	for i := 0; i < numberOfNodes; i++ {
		for j := i + 1; j < numberOfNodes; j++ {
			avg := (edges64[i][j] + edges64[j][i]) / 2
			edges64[i][j] = avg
			edges64[j][i] = avg
		}
	}

	// Jetzt zu float32 konvertieren
	edges := make([][]float32, numberOfNodes)
	for i := 0; i < numberOfNodes; i++ {
		edges[i] = make([]float32, numberOfNodes)
		for j := 0; j < numberOfNodes; j++ {
			edges[i][j] = float32(edges64[i][j])
		}
	}

	return model.Graph{
		Nodes: numberOfNodes,
		Edges: edges,
	}
}

// SaveModelGraphAsDOT konvertiert einen model.Graph in einen DOT-Graph
// und speichert ihn in eine Datei
func SaveModelGraphAsDOT(g model.Graph, filename string) error {
	// Neuen ungerichteten Graph erzeugen
	dotGraph := graph.New(graph.IntHash, graph.Weighted())

	// Knoten hinzufügen
	for i := 0; i < g.Nodes; i++ {
		_ = dotGraph.AddVertex(i)
	}

	// Kanten mit Gewichten hinzufügen
	for i := 0; i < g.Nodes; i++ {
		for j := i + 1; j < g.Nodes; j++ {
			weight := g.Edges[i][j]
			if weight > 0 {
				_ = dotGraph.AddEdge(i, j, graph.EdgeAttribute("label", fmt.Sprintf("%.1f", weight)))
			}
		}
	}

	// DOT-Datei erzeugen
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("could not create DOT file: %v", err)
	}
	defer file.Close()

	err = draw.DOT(dotGraph, file)
	if err != nil {
		return fmt.Errorf("could not write DOT file: %v", err)
	}

	fmt.Printf("Graph saved as DOT: %s\n", filename)
	return nil
}

func euclideanDistance(p1, p2 point) float64 {
	dx := p1.x - p2.x
	dy := p1.y - p2.y
	return math.Sqrt(dx*dx + dy*dy)
}
