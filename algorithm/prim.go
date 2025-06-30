package algorithm

import (
	"christofides-algo/model"
	"container/heap"
)

type edgeHeapItem struct {
	from   int
	to     int
	weight float32
}

type edgeHeap []edgeHeapItem

func (h edgeHeap) Len() int            { return len(h) }
func (h edgeHeap) Less(i, j int) bool  { return h[i].weight < h[j].weight }
func (h edgeHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *edgeHeap) Push(x interface{}) { *h = append(*h, x.(edgeHeapItem)) }
func (h *edgeHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

// BuildMST erstellt einen MST mit dem Prim-Algorithmus und gibt []model.Edge mit Point-Koordinaten zurück
func BuildMST(g model.Graph) []model.Edge {
	n := g.Nodes
	visited := make([]bool, n)
	mst := []model.Edge{}
	h := &edgeHeap{}
	heap.Init(h)

	// Start bei Knoten 0
	visited[0] = true
	for j := 1; j < n; j++ {
		heap.Push(h, edgeHeapItem{from: 0, to: j, weight: g.Edges[0][j]})
	}

	for len(mst) < n-1 && h.Len() > 0 {
		e := heap.Pop(h).(edgeHeapItem)
		if visited[e.to] {
			continue
		}

		visited[e.to] = true
		mst = append(mst, model.Edge{
			From:   g.Points[e.from],
			To:     g.Points[e.to],
			Weight: e.weight,
		})

		for j := 0; j < n; j++ {
			if !visited[j] {
				heap.Push(h, edgeHeapItem{from: e.to, to: j, weight: g.Edges[e.to][j]})
			}
		}
	}

	return mst
}
