package fixtures

import (
	"bytes"
	"testing"

	"github.com/go-quicktest/qt"
)

func TestRepository_ImportOrder(t *testing.T) {
	repo := NewRepository()
	repo.Register(User{})
	repo.Register(Transaction{})
	repo.Register(Address{})

	order, err := repo.importOrder()
	qt.Assert(t, qt.IsNil(err))
	qt.Assert(t, qt.HasLen(order, 3))

	qt.Check(t, qt.DeepEquals(order, []string{"addresses", "users", "transactions"}))
}

func TestRepository_DrawDependencies(t *testing.T) {
	repo := NewRepository()
	repo.Register(User{})
	repo.Register(Transaction{})
	repo.Register(Address{})

	var buf bytes.Buffer
	err := repo.DrawDependencies(&buf)
	qt.Assert(t, qt.IsNil(err))

	// fmt.Println(buf.String())

	qt.Check(t, qt.StringContains(buf.String(), `"transactions" -> "addresses"`))
	qt.Check(t, qt.StringContains(buf.String(), `"transactions" -> "users"`))
	qt.Check(t, qt.StringContains(buf.String(), `"users" -> "addresses"`))
}
