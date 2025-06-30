package algorithm

import (
	"christofides-algo/model"
	"math"
)

func FindOddDegreeNodes(mst []model.Edge, points []model.Point) []int {
	degree := make(map[model.Point]int)
	for _, e := range mst {
		degree[e.From]++
		degree[e.To]++
	}

	odds := []int{}
	for i, p := range points {
		if degree[p]%2 == 1 {
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

			newPair := model.Edge{
				From:   g.Points[a],
				To:     g.Points[b],
				Weight: w,
			}

			newRemaining := append([]int{}, remaining[1:i]...)
			newRemaining = append(newRemaining, remaining[i+1:]...)
			helper(newRemaining, append(current, newPair), cost+w)
		}
	}
	helper(nodes, []model.Edge{}, 0)
	return best
}

// Multigraph als Adjazenzliste (mit Indices)
func MergeEdges(mst, matching []model.Edge, g model.Graph) map[int][]int {
	graph := make(map[int][]int)
	pointToIndex := make(map[model.Point]int)
	for i, p := range g.Points {
		pointToIndex[p] = i
	}

	add := func(e model.Edge) {
		from := pointToIndex[e.From]
		to := pointToIndex[e.To]
		graph[from] = append(graph[from], to)
		graph[to] = append(graph[to], from)
	}

	for _, e := range mst {
		add(e)
	}
	for _, e := range matching {
		add(e)
	}
	return graph
}

func EulerTour(graph map[int][]int, start int) []int {
	var tour []int
	var dfs func(int)
	visited := make(map[[2]int]bool)

	dfs = func(u int) {
		for len(graph[u]) > 0 {
			v := graph[u][0]
			graph[u] = graph[u][1:]

			// Entferne Rückkante
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

func Shortcut(tour []int) []int {
	seen := make(map[int]bool)
	ham := []int{}
	for _, node := range tour {
		if !seen[node] {
			ham = append(ham, node)
			seen[node] = true
		}
	}
	if len(ham) > 0 {
		ham = append(ham, ham[0])
	}
	return ham
}

func Christofides(g model.Graph) ([]int, float32) {
	mst := BuildMST(g) // muss angepasst sein: gibt []model.Edge mit Point zurück
	odd := FindOddDegreeNodes(mst, g.Points)
	matching := MinPerfectMatching(odd, g)
	multigraph := MergeEdges(mst, matching, g)
	euler := EulerTour(multigraph, 0)
	tour := Shortcut(euler)
	cost := computeTourCost(tour, g)
	return tour, cost
}

func computeTourCost(tour []int, g model.Graph) float32 {
	var total float32
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

func reverse(s []int) []int {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}
