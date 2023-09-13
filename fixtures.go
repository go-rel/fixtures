// Package fixtures importing data from YAML files in database for REL.
//
// Usage:
//
//	repo := fixtures.New()
//	// Register all needed types
//	repo.Register(&MyTableType{})
//
//	// Import data from YAML content
//	err := repo.Import(ctx, db, content)
package fixtures

import (
	"context"
	"log"
	"reflect"

	"github.com/go-rel/rel"
)

// Logger to be used by the repository to notify about warnings.
type Logger interface {
	Warn(msg string)
}

type defaultLogger struct{}

func (l defaultLogger) Warn(msg string) {
	log.Println(msg)
}

// BeforeSave interface can be implemented by a type to allow changing type data before saving data to database.
type BeforeSave interface {
	BeforeSave(context.Context) error
}

// Repsitory of fixtures that can be loaded and imported.
type Repository struct {
	log      Logger
	registry map[string]any
}

// New creates a new fixtures repository.
func New() *Repository {
	return &Repository{
		log:      defaultLogger{},
		registry: make(map[string]any, 10),
	}
}

// Register a type that can be loaded as fixture.
func (r *Repository) Register(v any) {
	if reflect.ValueOf(v).Kind() != reflect.Ptr {
		v = reflect.New(reflect.TypeOf(v)).Interface()
	}

	name := rel.NewDocument(v).Meta().Table()
	if _, ok := r.registry[name]; ok {
		return
	}

	r.registry[name] = v
}

// SetLogger sets the logger to be used by the repository to notify about warnings.
func (r *Repository) SetLogger(l Logger) {
	r.log = l
}
