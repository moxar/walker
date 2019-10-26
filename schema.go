package walker

import (
	"strings"

	"gopkg.in/yaml.v2"
)

// Schema of the database, represented as graph.
type Schema struct {
	Verticies []Vertex
	Arcs      []Arc
	aliases   map[Vertex]Vertex
}

// String returns the yaml representation of the schema.
func (s Schema) String() string {
	raw, _ := yaml.Marshal(s)
	return string(raw)
}

// NewSchema prepares the schema for aliases.
func NewSchema() *Schema {
	return &Schema{aliases: make(map[Vertex]Vertex)}
}

// Vertex of the schema.
type Vertex = string

// Arc of the schema. Arcs are two-ways.
type Arc struct {
	Source, Destination Vertex

	// Link expresses the relation between the source and destination.
	Link string
}

// Alias the canonical table into the alias on the given links. If no link is provided, all relations are aliased.
func (s *Schema) Alias(canon, alias Vertex, links ...Vertex) {
	s.Verticies = append(s.Verticies, alias)
	s.aliases[alias] = canon

	// Look for the alias in each arc.
	for i, a := range s.Arcs {
		replaceArc := func(dst string) {
			chunks := strings.Fields(a.Link)
			for i, c := range chunks {
				if strings.HasPrefix(c, canon+".") {
					chunks[i] = alias + strings.TrimLeft(c, canon)
				}
			}
			link := strings.Join(chunks, " ")
			s.Arcs[i] = Arc{
				Source:      alias,
				Destination: dst,
				Link:        link,
			}
		}

		// Replace only the alias link that matches the arc.
		for _, link := range links {
			if (a.Source == canon && a.Destination == link) || (a.Destination == canon && a.Source == link) {
				replaceArc(link)
				break
			}
		}

		// No link is declared, meaning every link of the arc must be replaced.
		if len(links) == 0 {
			var link string
			if a.Source == canon {
				link = a.Destination
			}
			if a.Destination == canon {
				link = a.Source
			}
			if link != "" {
				replaceArc(link)
			}
		}
	}
}
