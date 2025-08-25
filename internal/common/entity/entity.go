package entity

import (
	"errors"
	"fmt"
	"strings"
)

type Item struct {
	ID       string
	Name     string
	Quantity int32
	PriceID  string
}

func (it Item) validate() error {
	//if err := util.AssertNotEmpty(it.ID, it.Name, it.PriceID); err != nil {
	//	return err
	//}
	var invalidFields []string
	if it.ID == "" {
		invalidFields = append(invalidFields, "ID")
	}
	if it.Name == "" {
		invalidFields = append(invalidFields, "Name")
	}
	if it.PriceID == "" {
		invalidFields = append(invalidFields, "PriceID")
	}
	return fmt.Errorf("item %v is invalid,empty fields = [%s]", it, strings.Join(invalidFields, ", "))
}

func NewItem(ID string, name string, quantity int32, priceID string) *Item {
	return &Item{ID: ID, Name: name, Quantity: quantity, PriceID: priceID}
}

func NewValidItem(ID string, name string, quantity int32, priceID string) (*Item, error) {
	item := NewItem(ID, name, quantity, priceID)
	if err := item.validate(); err != nil {
		return nil, err
	}
	return item, nil

}

type ItemWithQuantity struct {
	ID       string
	Quantity int32
}

func (iq ItemWithQuantity) validate() error {
	//if err := util.AssertNotEmpty(it.ID, it.Name, it.PriceID); err != nil {
	//	return err
	//}
	var invalidFields []string
	if iq.ID == "" {
		invalidFields = append(invalidFields, "ID")
	}

	return errors.New(strings.Join(invalidFields, ","))
}

func NewItemWithQuantity(ID string, quantity int32) *ItemWithQuantity {
	return &ItemWithQuantity{ID: ID, Quantity: quantity}
}

func NewValidItemWithQuantity(ID string, quantity int32) (*ItemWithQuantity, error) {
	item := NewItemWithQuantity(ID, quantity)
	if err := item.validate(); err != nil {
		return nil, err
	}
	return item, nil
}

type Order struct {
	ID          string
	CustomerID  string
	Status      string
	PaymentLink string
	Items       []*Item
}

func (o *Order) validate() error {
	var invalidFields []string
	if o.ID == "" {
		invalidFields = append(invalidFields, "ID")
	}
	if o.CustomerID == "" {
		invalidFields = append(invalidFields, "CustomerID")
	}
	if o.Status == "" {
		invalidFields = append(invalidFields, "Status")
	}

	for _, item := range o.Items {
		if err := item.validate(); err != nil {
			invalidFields = append(invalidFields, "Items")
			break
		}
	}
	return errors.New(strings.Join(invalidFields, ","))
}

func NewOrder(items []*Item, paymentLink string, status string, customerID string, ID string) *Order {
	return &Order{Items: items, PaymentLink: paymentLink, Status: status, CustomerID: customerID, ID: ID}
}

func NewValidOrder(items []*Item, paymentLink string, status string, customerID string, ID string) (*Order, error) {
	order := NewOrder(items, paymentLink, status, customerID, ID)
	if err := order.validate(); err != nil {
		return nil, err
	}
	return order, nil
}
