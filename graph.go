package fixtures

import (
	"errors"
	"fmt"
	"io"

	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"
	"github.com/go-rel/rel"
)

func (r Repository) importOrder() ([]string, error) {
	g := graph.New(graph.StringHash, graph.Directed(), graph.PreventCycles())

	// Register all tables
	for name := range r.registry {
		_ = g.AddVertex(name)
	}

	// Register all relations between tables
	for name, v := range r.registry {
		meta := rel.NewDocument(v).Meta()
		for _, field := range meta.BelongsTo() {
			fk := meta.Association(field)
			if err := g.AddEdge(fk.DocumentMeta().Table(), name); errors.Is(err, graph.ErrEdgeCreatesCycle) {
				r.log.Warn(fmt.Sprintf("foreign key cycle detected, skipping table %s and %s relation", name, fk.DocumentMeta().Table()))
			}
		}
		for _, field := range meta.HasOne() {
			fk := meta.Association(field)
			if err := g.AddEdge(name, fk.DocumentMeta().Table()); errors.Is(err, graph.ErrEdgeCreatesCycle) {
				r.log.Warn(fmt.Sprintf("foreign key cycle detected, skipping table %s and %s relation", fk.DocumentMeta().Table(), name))
			}
		}
		for _, field := range meta.HasMany() {
			fk := meta.Association(field)
			if err := g.AddEdge(name, fk.DocumentMeta().Table()); errors.Is(err, graph.ErrEdgeCreatesCycle) {
				r.log.Warn(fmt.Sprintf("foreign key cycle detected, skipping table %s and %s relation", fk.DocumentMeta().Table(), name))
			}
		}
	}

	return graph.TopologicalSort(g)
}

// DrawDependencies draws a graph of all tables and their relations in DOT format.
func (r Repository) DrawDependencies(w io.Writer) error {
	g := graph.New(graph.StringHash, graph.Directed())

	// Register all tables
	for name := range r.registry {
		_ = g.AddVertex(name)
	}

	// Register all relations between tables
	for name, v := range r.registry {
		meta := rel.NewDocument(v).Meta()
		for _, field := range meta.BelongsTo() {
			fk := meta.Association(field)
			_ = g.AddEdge(name, fk.DocumentMeta().Table(), graph.EdgeAttribute("label", fk.ReferenceField()))
		}
	}
	for name, v := range r.registry {
		meta := rel.NewDocument(v).Meta()
		for _, field := range meta.HasOne() {
			fk := meta.Association(field)
			_ = g.AddEdge(fk.DocumentMeta().Table(), name, graph.EdgeAttribute("label", fk.ForeignField()))
		}
		for _, field := range meta.HasMany() {
			fk := meta.Association(field)
			_ = g.AddEdge(fk.DocumentMeta().Table(), name, graph.EdgeAttribute("label", fk.ForeignField()))
		}
	}

	return draw.DOT(g, w)
}
