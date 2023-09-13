package fixtures

import (
	"context"
	"strings"
	"time"
)

// User schema.
type User struct {
	ID   int
	Name string
	Age  int16

	Address   Address `ref:"address_id" fk:"id" autoload:"true"`
	AddressID int

	CreatedAt time.Time

	Deleted bool

	Transactions []Transaction `ref:"id" fk:"buyer_id"`
}

// Transaction schema.
type Transaction struct {
	ID     int
	Item   string
	Status string

	Buyer   User `ref:"buyer_id" fk:"id"`
	BuyerID int

	DeliveryAddress   *Address `ref:"delivery_address_id" fk:"id"`
	DeliveryAddressID int
}

// Address schema.
type Address struct {
	ID   int
	City string
}

func (a *Address) BeforeSave(_ context.Context) error {
	a.City = strings.TrimSpace(a.City)

	return nil
}
