package walker

import (
	"fmt"
	"strings"

	"github.com/RyanCarrier/dijkstra"
)

// Graph extends a dijkstra's graph based on the schema.
type Graph struct {
	graph  *dijkstra.Graph
	schema *Schema
}

// NewGraph prepares a graph from a schema.
func NewGraph(s *Schema) (*Graph, error) {
	g := dijkstra.NewGraph()
	for _, v := range s.Verticies {
		g.AddMappedVertex(v)
	}
	for _, a := range s.Arcs {
		if err := g.AddMappedArc(a.Source, a.Destination, 1); err != nil {
			return nil, err
		}
		if err := g.AddMappedArc(a.Destination, a.Source, 1); err != nil {
			return nil, err
		}
	}
	return &Graph{graph: g, schema: s}, nil
}

// shortest returns the shortest path between from and to.
func (g *Graph) shortest(from, to string) ([]string, error) {
	src, _ := g.graph.GetMapping(from)
	dst, _ := g.graph.GetMapping(to)
	p, err := g.graph.Shortest(src, dst)
	if err != nil {
		return nil, fmt.Errorf("no path between %s and %s", from, to)
	}
	path := make([]string, 0, len(p.Path))
	for _, arc := range p.Path {
		name, _ := g.graph.GetMapped(arc)
		path = append(path, name)
	}

	table := path[0]
	if alias, ok := g.schema.aliases[table]; ok {
		table = alias + " AS " + table
	}
	out := make([]string, 0, len(p.Path))
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
	return out, nil
}

// From builds the FROM ... JOIN clause of the query to link the src with all dst.
func (g *Graph) From(src string, dst ...string) (string, error) {
	out := make([]string, 0, 2*len(g.schema.Verticies)-1)
	indexes := make(map[string]struct{})
	for _, to := range dst {
		if to == src {
			continue
		}
		path, err := g.shortest(src, to)
		if err != nil {
			return "", err
		}
		for _, p := range path {
			if _, ok := indexes[p]; !ok {
				indexes[p] = struct{}{}
				out = append(out, p)
			}
		}
	}
	return strings.Join(out, " "), nil
}
