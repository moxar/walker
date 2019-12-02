package walker

import "github.com/beefsack/go-astar"

type pather struct {
	Vertex
	Neighbors []astar.Pather
}

// PathNeighbors returns the list of neighbors of a Vertex.
func (p *pather) PathNeighbors() []astar.Pather {
	return p.Neighbors
}

// PathNeighborCost returns the constant cost of 1: every path has the same weight.
func (p *pather) PathNeighborCost(to astar.Pather) float64 {
	return 1
}

// PathEstimatedCost returns the heuristic distance between p and the target.
// We return 0 because because we cannot estimate the distance using Manhattan approach.
// Using 0 fallsback to Dijkstra's algorythm. See http://theory.stanford.edu/~amitp/GameProgramming/Heuristics.html
func (p *pather) PathEstimatedCost(to astar.Pather) float64 {
	return 0
}
