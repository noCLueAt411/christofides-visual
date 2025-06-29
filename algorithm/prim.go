package algorithm

import (
	"christofides-algo/model"
	. "christofides-algo/model"
	"container/heap"
)

// Priority Queue für Kanten (min-heap)
type edgeHeap []Edge

func (h edgeHeap) Len() int            { return len(h) }
func (h edgeHeap) Less(i, j int) bool  { return h[i].Weight < h[j].Weight }
func (h edgeHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *edgeHeap) Push(x interface{}) { *h = append(*h, x.(Edge)) }
func (h *edgeHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

// BuildMST erstellt einen MST mit dem Prim-Algorithmus und gibt die Kantenliste zurück
func BuildMST(g model.Graph) []Edge {
	n := g.Nodes
	visited := make([]bool, n)
	mst := []Edge{}
	h := &edgeHeap{}
	heap.Init(h)

	// Start bei Knoten 0
	visited[0] = true
	for j := 1; j < n; j++ {
		heap.Push(h, Edge{From: 0, To: j, Weight: g.Edges[0][j]})
	}

	for len(mst) < n-1 && h.Len() > 0 {
		e := heap.Pop(h).(Edge)
		if visited[e.To] {
			continue
		}

		visited[e.To] = true
		mst = append(mst, e)

		for j := 0; j < n; j++ {
			if !visited[j] {
				heap.Push(h, Edge{From: e.To, To: j, Weight: g.Edges[e.To][j]})
			}
		}
	}

	return mst
}
