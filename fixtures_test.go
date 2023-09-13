package fixtures

import (
	"testing"

	"github.com/go-quicktest/qt"
)

func TestFixtutes_Register(t *testing.T) {
	repo := NewRepository()
	repo.Register(User{})
	repo.Register(Address{})
	repo.Register(Transaction{})
	// Duplicate registration should be ignored
	repo.Register(Transaction{})

	_, ok := repo.registry["users"]
	qt.Check(t, qt.IsTrue(ok), qt.Commentf("users table should be registered"))
	_, ok = repo.registry["addresses"]
	qt.Check(t, qt.IsTrue(ok), qt.Commentf("addresses table should be registered"))
	_, ok = repo.registry["transactions"]
	qt.Check(t, qt.IsTrue(ok), qt.Commentf("transactions table should be registered"))
}
