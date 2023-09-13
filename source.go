package fixtures

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/go-rel/rel"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
)

type fixtureLoader struct {
	repo  *Repository
	table string
	data  map[string][]any
}

func (fixtureLoader) setValue(doc *rel.Document, field string, val ast.Node) {
	switch v := val.(type) {
	case *ast.IntegerNode:
		doc.SetValue(field, v.Value)
	case *ast.StringNode:
		typ, ok := doc.Type(field)
		if !ok {
			doc.SetValue(field, v.Value)
		}
		switch typ.String() {
		case "time.Time":
			if fv, err := time.Parse(time.RFC3339, v.Value); err == nil {
				doc.SetValue(field, fv)
			} else {
				doc.SetValue(field, v.Value)
			}
		default:
			doc.SetValue(field, v.Value)
		}
	case *ast.FloatNode:
		doc.SetValue(field, v.Value)
	case *ast.BoolNode:
		doc.SetValue(field, v.Value)
	case *ast.NullNode:
		doc.SetValue(field, nil)
	}
}

func (w *fixtureLoader) readRecord(table string, arr ast.Node) {
	datav := w.repo.registry[table]
	inst := reflect.New(reflect.TypeOf(datav).Elem()).Interface()
	doc := rel.NewDocument(inst)
	fields := doc.Meta().Fields()

	for _, val := range arr.(*ast.MappingNode).Values {
		field := val.Key.String()
		if !slices.Contains(fields, field) {
			continue
		}

		w.setValue(doc, field, val.Value)
	}

	w.data[table] = append(w.data[table], inst)
}

func (w *fixtureLoader) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}
	if node.Type() == ast.SequenceType {
		// Single table YAML file
		if w.table != "" {
			for _, arr := range node.(*ast.SequenceNode).Values {
				w.readRecord(w.table, arr)
			}
			return w
		}
		// Multiple table YAML file
		for table := range w.repo.registry {
			if node.GetPath() == "$."+table {
				for _, arr := range node.(*ast.SequenceNode).Values {
					w.readRecord(table, arr)
				}
			}
		}
	}

	return w
}

// ImportFromYAML imports data from YAML file content.
func (r *Repository) ImportFromYAML(ctx context.Context, db rel.Repository, source []byte) error {
	p, err := parser.ParseBytes(source, 0)
	if err != nil {
		return err
	}
	if len(p.Docs) == 0 {
		return nil
	}

	w := &fixtureLoader{
		repo: r,
		data: make(map[string][]any, len(r.registry)),
	}
	ast.Walk(w, p.Docs[0])

	return r.importData(ctx, db, w.data)
}

// ImportFromDir imports data from YAML files in a directory.
func (r *Repository) ImportFromDir(ctx context.Context, db rel.Repository, path string) error {
	dir, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	w := &fixtureLoader{
		repo: r,
		data: make(map[string][]any, len(r.registry)),
	}

	for _, file := range dir {
		if file.IsDir() || filepath.Ext(file.Name()) != ".yaml" {
			continue
		}

		p, err := parser.ParseFile(filepath.Join(path, file.Name()), 0)
		if err != nil {
			return err
		}

		if len(p.Docs) == 0 {
			continue
		}

		w.table = strings.TrimSuffix(file.Name(), ".yaml")
		ast.Walk(w, p.Docs[0])
	}

	return r.importData(ctx, db, w.data)
}

func (r *Repository) importData(ctx context.Context, db rel.Repository, data map[string][]any) error {
	tables, err := r.importOrder()
	if err != nil {
		r.log.Warn(fmt.Sprintf("failed to get table import order: %v", err))

		// Fallback to list of map keys
		for table := range r.registry {
			tables = append(tables, table)
		}
	}

	return db.Transaction(ctx, func(ctx context.Context) error {
		// TODO: Disable foreign key checks

		for _, table := range tables {
			if _, ok := data[table]; !ok {
				continue
			}

			for _, v := range data[table] {
				if err := db.Insert(ctx, v); err != nil {
					return err
				}
			}
		}

		return nil
	})
}
