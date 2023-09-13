package fixtures

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-quicktest/qt"
	"github.com/go-rel/rel"
	"github.com/go-rel/rel/migrator"
	"github.com/go-rel/sqlite3"

	_ "github.com/mattn/go-sqlite3"
)

func createTestDB(t *testing.T) rel.Repository {
	temp, err := os.MkdirTemp(os.TempDir(), "fixtures-test")
	qt.Assert(t, qt.IsNil(err))

	adapter, err := sqlite3.Open(filepath.Join(temp, "test.db"))
	qt.Assert(t, qt.IsNil(err), qt.Commentf("failed to create SQLite3 database"))

	db := rel.New(adapter)

	m := migrator.New(db)
	m.Register(1, func(schema *rel.Schema) {
		schema.CreateTable("addresses", func(t *rel.Table) {
			t.ID("id", rel.Primary(true))
			t.String("city", rel.Limit(200))
		})

		schema.CreateTable("users", func(t *rel.Table) {
			t.ID("id", rel.Primary(true))
			t.String("name", rel.Limit(100))
			t.SmallInt("age")
			t.Int("address_id", rel.Required(true))
			t.DateTime("created_at", rel.Required(true))
			t.Bool("deleted", rel.Required(true), rel.Default(false))

			t.ForeignKey("address_id", "addresses", "id")
		})

		schema.CreateTable("transactions", func(t *rel.Table) {
			t.ID("id", rel.Primary(true))
			t.String("item", rel.Limit(100))
			t.String("status", rel.Default("pending"), rel.Limit(20))
			t.Int("buyer_id", rel.Required(true))
			t.Int("delivery_address_id")

			t.ForeignKey("buyer_id", "users", "id")
			t.ForeignKey("delivery_address_id", "addresses", "id")
		})
	}, func(_ *rel.Schema) {})

	m.Migrate(context.TODO())

	t.Cleanup(func() {
		adapter.Close()

		_ = os.RemoveAll(temp)
	})

	return db
}

func TestFixtures_Import(t *testing.T) {
	repo := New()
	repo.Register(&User{})
	repo.Register(&Address{})
	repo.Register(&Transaction{})

	db := createTestDB(t)

	err := repo.Import(context.TODO(), db,
		[]byte(`---
users:
- id: 1
  name: John Doe
  age: 20
  created_at: 2019-01-01T00:00:00Z
  address_id: 1
addresses:
- id: 1
  city: New York
transactions:
`))
	qt.Assert(t, qt.IsNil(err))

	user := User{}
	err = db.Find(context.TODO(), &user, rel.Eq("id", 1))
	qt.Assert(t, qt.IsNil(err))

	qt.Check(t, qt.Equals(user.ID, 1))
	qt.Check(t, qt.Equals(user.Name, "John Doe"))
	qt.Check(t, qt.Equals(user.Age, 20))
	qt.Check(t, qt.Equals(user.Address.ID, 1))
	qt.Check(t, qt.Equals(user.Address.City, "New York"))
}

func TestFixtures_ImportDir(t *testing.T) {
	repo := New()
	repo.Register(&User{})
	repo.Register(&Address{})
	repo.Register(&Transaction{})

	db := createTestDB(t)

	err := repo.ImportDir(context.TODO(), db, "testdata/sample/")
	qt.Assert(t, qt.IsNil(err))

	user := User{}
	err = db.Find(context.TODO(), &user, rel.Eq("id", 1))
	qt.Assert(t, qt.IsNil(err))

	qt.Check(t, qt.Equals(user.ID, 1))
	qt.Check(t, qt.Equals(user.Name, "John Doe"))
	qt.Check(t, qt.Equals(user.Age, 20))
	qt.Check(t, qt.Equals(user.Address.ID, 1))
	qt.Check(t, qt.Equals(user.Address.City, "New York"))
}
