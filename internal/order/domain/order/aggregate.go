package order

import "errors"

type Identity struct {
	CustomerID string
	OrderID    string
}

type AggregateRoot struct {
	Identity    Identity
	OrderEntity *Order
}

func NewAggregateRoot(identity Identity, orderEntity *Order) *AggregateRoot {
	return &AggregateRoot{Identity: identity, OrderEntity: orderEntity}
}

func (a *AggregateRoot) BusinessIdentity() Identity {
	return Identity{
		CustomerID: a.Identity.CustomerID,
		OrderID:    a.OrderEntity.ID,
	}
}

func (a *AggregateRoot) Validate() error {
	if a.Identity.CustomerID == "" || a.Identity.OrderID == "" {
		return errors.New("invalid identity")
	}
	if a.OrderEntity == nil {
		return errors.New("empty order entity")
	}
	return nil
}
