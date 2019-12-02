package walker

import (
	"fmt"
	"strings"

	"github.com/beefsack/go-astar"
)

// Graph extends a dijkstra's graph based on the schema.
type Graph struct {
	schema  *Schema
	pathers map[Vertex]*pather
}

// NewGraph prepares a graph from a schema.
func NewGraph(s *Schema) *Graph {

	var g Graph
	g.pathers = make(map[Vertex]*pather)
	g.schema = s

	for _, v := range s.Verticies {
		g.pathers[v] = &pather{Vertex: v}
	}

	for _, p := range g.pathers {
		p.Neighbors = make([]astar.Pather, 0, len(s.Arcs))
		for _, a := range s.Arcs {
			var v Vertex
			switch p.Vertex {
			case a.Source:
				v = a.Destination
			case a.Destination:
				v = a.Source
			default:
				continue
			}
			p.Neighbors = append(p.Neighbors, g.pathers[v])
		}
	}

	return &g
}

// shortest returns the shortest path between from and to.
func (g *Graph) shortest(from, to string) ([]string, error) {
	if from == to || to == "" {
		return []string{from}, nil
	}

	f, t := g.pathers[from], g.pathers[to]
	if f == nil || t == nil {
		return nil, &PathError{Src: from, Dst: to}
	}

	pathers, _, ok := astar.Path(f, t)
	if !ok {
		return nil, &PathError{Src: from, Dst: to}
	}
	path := make([]string, 0, len(pathers))
	for _, pa := range pathers {
		p := pa.(*pather)
		path = append(path, p.Vertex)
	}

	// astar sometimes messes up with the vertcies order. We don't realy care about them, except
	// the first one that must be the "main" table.
	if path[0] != from {
		path[0], path[len(path)-1] = path[len(path)-1], path[0]
	}
	return path, nil
}

// drawPath links and aliases the elements of the path.
func (g *Graph) drawPath(path []string) []string {

	table := path[0]
	if alias, ok := g.schema.aliases[table]; ok {
		table = alias + " AS " + table
	}
	out := make([]string, 0, len(path))
	out = append(out, table)
	for i := 0; i < len(path)-1; i++ {
		v1, v2 := path[i], path[i+1]
		for _, a := range g.schema.Arcs {
			var table string
			if a.Source == v1 && a.Destination == v2 {
				table = a.Destination
			}
			if a.Destination == v1 && a.Source == v2 {
				table = a.Source
			}
			if table != "" {
				if alias, ok := g.schema.aliases[v2]; ok {
					table = alias + " AS " + table
				}
				out = append(out, fmt.Sprintf("JOIN %s ON %s", table, a.Link))
			}
		}
	}
	return out
}

// From builds the FROM ... JOIN clause of the query to link the src with all dst.
func (g *Graph) From(src string, dst ...string) (string, error) {
	out := make([]string, 0, 2*len(dst)+1)
	indexes := make(map[string]struct{})
	if len(dst) == 0 {
		dst = append(dst, src)
	}
	for _, to := range dst {
		path, err := g.shortest(src, to)
		if err != nil {
			return "", err
		}
		path = g.drawPath(path)
		for _, p := range path {
			if _, ok := indexes[p]; !ok {
				indexes[p] = struct{}{}
				out = append(out, p)
			}
		}
	}
	return strings.Join(out, " "), nil
}
