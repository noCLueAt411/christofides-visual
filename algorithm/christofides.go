package algorithm

import (
	"christofides-algo/model"
	"math"
)

func FindOddDegreeNodes(mst []model.Edge, nodeCount int) []int {
	degree := make([]int, nodeCount)
	for _, e := range mst {
		degree[e.From]++
		degree[e.To]++
	}
	odds := []int{}
	for i, d := range degree {
		if d%2 == 1 {
			odds = append(odds, i)
		}
	}
	return odds
}

// Brute-Force Minimum Perfect Matching (nur für ≤8 Knoten!)
func MinPerfectMatching(odds []int, g model.Graph) []model.Edge {
	if len(odds)%2 != 0 {
		panic("Number of odd-degree nodes must be even")
	}
	return bestMatching(odds, g)
}

func bestMatching(nodes []int, g model.Graph) []model.Edge {
	best := []model.Edge{}
	minCost := float32(math.MaxFloat32)

	var helper func([]int, []model.Edge, float32)
	helper = func(remaining []int, current []model.Edge, cost float32) {
		if len(remaining) == 0 {
			if cost < minCost {
				best = make([]model.Edge, len(current))
				copy(best, current)
				minCost = cost
			}
			return
		}
		a := remaining[0]
		for i := 1; i < len(remaining); i++ {
			b := remaining[i]
			w := g.Edges[a][b]
			newPair := model.Edge{From: a, To: b, Weight: w}

			// remove a & b from remaining
			newRemaining := append([]int{}, remaining[1:i]...)
			newRemaining = append(newRemaining, remaining[i+1:]...)
			helper(newRemaining, append(current, newPair), cost+w)
		}
	}
	helper(nodes, []model.Edge{}, 0)
	return best
}

// Merge MST + Matching → Multigraph (als Adjazenzliste)
func MergeEdges(mst, matching []model.Edge, nodeCount int) map[int][]int {
	graph := make(map[int][]int)
	addEdge := func(from, to int) {
		graph[from] = append(graph[from], to)
		graph[to] = append(graph[to], from)
	}
	for _, e := range mst {
		addEdge(e.From, e.To)
	}
	for _, e := range matching {
		addEdge(e.From, e.To)
	}
	return graph
}

// Hierholzer-Algorithmus für Eulerkreis
func EulerTour(graph map[int][]int, start int) []int {
	var tour []int
	var dfs func(int)
	visited := make(map[[2]int]bool)

	dfs = func(u int) {
		for len(graph[u]) > 0 {
			v := graph[u][0]
			graph[u] = graph[u][1:]

			// remove reverse edge
			for i, w := range graph[v] {
				if w == u {
					graph[v] = append(graph[v][:i], graph[v][i+1:]...)
					break
				}
			}

			e := [2]int{min(u, v), max(u, v)}
			if visited[e] {
				continue
			}
			visited[e] = true
			dfs(v)
		}
		tour = append(tour, u)
	}
	dfs(start)
	return reverse(tour)
}

func reverse(s []int) []int {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

// Shortcut Euler-Tour → Hamiltonkreis
func Shortcut(tour []int) []int {
	seen := make(map[int]bool)
	ham := []int{}
	for _, node := range tour {
		if !seen[node] {
			ham = append(ham, node)
			seen[node] = true
		}
	}
	// Rückkehr zum Start
	if len(ham) > 0 {
		ham = append(ham, ham[0])
	}
	return ham
}

// Christofides kombiniert alle Schritte
func Christofides(g model.Graph) ([]int, float32) {
	mst := BuildMST(g)
	odd := FindOddDegreeNodes(mst, g.Nodes)
	matching := MinPerfectMatching(odd, g)
	multigraph := MergeEdges(mst, matching, g.Nodes)
	euler := EulerTour(multigraph, 0)
	tour := Shortcut(euler)
	cost := computeTourCost(tour, g)
	return tour, cost
}

func computeTourCost(tour []int, g model.Graph) float32 {
	var total float32 = 0
	for i := 0; i < len(tour)-1; i++ {
		total += g.Edges[tour[i]][tour[i+1]]
	}
	return total
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
