package convertor

import (
	client "github.com/mushroomyuan/gorder/common/client/order"
	"github.com/mushroomyuan/gorder/common/genproto/orderpb"
	domain "github.com/mushroomyuan/gorder/order/domain/order"
	"github.com/mushroomyuan/gorder/order/entity"
)

type OrderConvertor struct{}
type ItemConvertor struct{}
type ItemWithQuantityConvertor struct{}

func (c *OrderConvertor) EntityToProto(o *domain.Order) *orderpb.Order {
	c.check(o)
	return &orderpb.Order{
		ID:          o.ID,
		CustomerID:  o.CustomerID,
		Status:      o.Status,
		Items:       NewItemConvertor().EntitiesToProtos(o.Items),
		PaymentLink: o.PaymentLink,
	}
}

func (c *OrderConvertor) ProtoToEntity(o *orderpb.Order) *domain.Order {
	c.check(o)
	return &domain.Order{
		ID:          o.ID,
		CustomerID:  o.CustomerID,
		Status:      o.Status,
		PaymentLink: o.PaymentLink,
		Items:       NewItemConvertor().ProtosToEntities(o.Items),
	}
}

func (c *OrderConvertor) EntityToClient(o *domain.Order) *client.Order {
	c.check(o)
	return &client.Order{
		CustomerId:  o.CustomerID,
		Id:          o.ID,
		Status:      o.Status,
		PaymentLink: o.PaymentLink,
		Items:       NewItemConvertor().EntitiesToClients(o.Items),
	}
}

func (c *OrderConvertor) ClientToEntity(o *client.Order) *domain.Order {
	c.check(o)
	return &domain.Order{
		CustomerID:  o.CustomerId,
		ID:          o.Id,
		Status:      o.Status,
		PaymentLink: o.PaymentLink,
		Items:       NewItemConvertor().ClientsToEntities(o.Items),
	}
}

func (c *OrderConvertor) check(o any) {
	if o == nil {
		panic("cannot convert nil order")
	}
}

func (c *ItemConvertor) EntitiesToProtos(items []*entity.Item) (res []*orderpb.Item) {
	for _, i := range items {
		res = append(res, c.EntityToProto(i))
	}
	return
}

func (c *ItemConvertor) EntityToProto(i *entity.Item) *orderpb.Item {
	return &orderpb.Item{
		ID:       i.ID,
		Name:     i.Name,
		Quantity: i.Quantity,
		PriceID:  i.PriceID,
	}
}

func (c *ItemConvertor) ProtosToEntities(items []*orderpb.Item) (res []*entity.Item) {
	for _, i := range items {
		res = append(res, c.ProtToEntity(i))
	}
	return
}

func (c *ItemConvertor) ProtToEntity(o *orderpb.Item) *entity.Item {
	return &entity.Item{
		ID:       o.ID,
		Name:     o.Name,
		Quantity: o.Quantity,
		PriceID:  o.PriceID,
	}
}

func (c *ItemConvertor) EntitiesToClients(items []*entity.Item) (res []client.Item) {
	for _, i := range items {
		res = append(res, c.EntityToClient(i))
	}
	return
}

func (c *ItemConvertor) EntityToClient(e *entity.Item) client.Item {
	return client.Item{
		Id:       e.ID,
		Name:     e.Name,
		Quantity: e.Quantity,
		PriceId:  e.PriceID,
	}
}

func (c *ItemConvertor) ClientsToEntities(items []client.Item) (res []*entity.Item) {
	for _, i := range items {
		res = append(res, c.ClientToEntity(i))
	}
	return
}

func (c *ItemConvertor) ClientToEntity(e client.Item) *entity.Item {
	return &entity.Item{
		ID:       e.Id,
		Name:     e.Name,
		Quantity: e.Quantity,
		PriceID:  e.PriceId,
	}
}

func (c *ItemWithQuantityConvertor) EntitiesToProtos(items []*entity.ItemWithQuantity) (res []*orderpb.ItemWithQuantity) {
	for _, i := range items {
		res = append(res, c.EntityToProto(i))
	}
	return
}

func (c *ItemWithQuantityConvertor) EntityToProto(i *entity.ItemWithQuantity) *orderpb.ItemWithQuantity {
	return &orderpb.ItemWithQuantity{
		ID:       i.ID,
		Quantity: i.Quantity,
	}
}

func (c *ItemWithQuantityConvertor) ProtosToEntities(items []*orderpb.ItemWithQuantity) (res []*entity.ItemWithQuantity) {
	for _, i := range items {
		res = append(res, c.ProtoToEntity(i))
	}
	return
}

func (c *ItemWithQuantityConvertor) ProtoToEntity(i *orderpb.ItemWithQuantity) *entity.ItemWithQuantity {
	return &entity.ItemWithQuantity{
		ID:       i.ID,
		Quantity: i.Quantity,
	}
}

func (c *ItemWithQuantityConvertor) ClientsToEntities(items []client.ItemWithQuantity) (res []*entity.ItemWithQuantity) {
	for _, i := range items {
		res = append(res, c.ClientToEntity(i))
	}
	return
}

func (c *ItemWithQuantityConvertor) ClientToEntity(i client.ItemWithQuantity) *entity.ItemWithQuantity {
	return &entity.ItemWithQuantity{
		ID:       i.Id,
		Quantity: i.Quantity,
	}
}

func (c *ItemWithQuantityConvertor) EntitiesToClients(items []entity.ItemWithQuantity) (res []*client.ItemWithQuantity) {
	for _, i := range items {
		res = append(res, c.EntityToClient(i))
	}
	return
}

func (c *ItemWithQuantityConvertor) EntityToClient(i entity.ItemWithQuantity) *client.ItemWithQuantity {
	return &client.ItemWithQuantity{
		Id:       i.ID,
		Quantity: i.Quantity,
	}
}
