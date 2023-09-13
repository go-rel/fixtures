# Fixtures

[![Go Reference](https://pkg.go.dev/badge/github.com/go-rel/fixtures.svg)](https://pkg.go.dev/github.com/go-rel/fixtures)
[![Tests](https://github.com/go-rel/fixtures/actions/workflows/test.yaml/badge.svg?branch=main)](https://github.com/go-rel/fixtures/actions/workflows/test.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-rel/fixtures)](https://goreportcard.com/report/github.com/go-rel/fixtures)
[![codecov](https://codecov.io/gh/go-rel/fixtures/branch/main/graph/badge.svg?token=yxBdKVPXip)](https://codecov.io/gh/go-rel/fixtures)
[![Gitter chat](https://badges.gitter.im/go-rel/rel.png)](https://gitter.im/go-rel/rel)

Fixture importing in database for REL.

## Example

### Using single YAML file

YAML file must contain properties named after database table names with array of objects with fields as column names.

```go
package main

import (
	"context"

	"github.com/go-rel/fixtures"
	"github.com/go-rel/rel"
)

func main() {
	repo := fixtures.NewRepository()
	// Register all needed types
	repo.Register(&User{})
	repo.Register(&Address{})
	repo.Register(&Transaction{})

	// TODO db := rel.New(adapter)

	if err := repo.ImportFromYAML(context.Background(), db,
		[]byte(`---
users:
- id: 1
  name: John Doe
  age: 20
  created_at: 2019-01-01T06:10:00Z
  address_id: 1
addresses:
- id: 1
  city: New York
`)); err != nil {
		panic(err)
	}
}
```

### Using directory with YAML files

Directory must contain YAML files named as table names with extension `.yaml` containing just array of objects
with fields as column names.

```go
package main

import (
	"context"

	"github.com/go-rel/fixtures"
	"github.com/go-rel/rel"
)

func main() {
	repo := fixtures.NewRepository()
	// Register all needed types
	repo.Register(&User{})
	repo.Register(&Address{})
	repo.Register(&Transaction{})

	// TODO db := rel.New(adapter)

	if err := repo.ImportFromDir(context.Background(), db, "path/to/dir/")); err != nil {
		panic(err)
	}
}
```
