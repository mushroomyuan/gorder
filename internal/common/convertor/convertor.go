package convertor

import (
	client "github.com/mushroomyuan/gorder/common/client/order"
	"github.com/mushroomyuan/gorder/common/entity"
	"github.com/mushroomyuan/gorder/common/genproto/orderpb"
)

type OrderConvertor struct{}
type ItemConvertor struct{}
type ItemWithQuantityConvertor struct{}

func (c *OrderConvertor) EntityToProto(o *entity.Order) *orderpb.Order {
	c.check(o)
	return &orderpb.Order{
		ID:          o.ID,
		CustomerID:  o.CustomerID,
		Status:      o.Status,
		Items:       NewItemConvertor().EntitiesToProtos(o.Items),
		PaymentLink: o.PaymentLink,
	}
}

func (c *OrderConvertor) ProtoToEntity(o *orderpb.Order) *entity.Order {
	c.check(o)
	return entity.NewOrder(
		NewItemConvertor().ProtosToEntities(o.Items),
		o.PaymentLink,
		o.Status,
		o.CustomerID,
		o.ID,
	)
}

func (c *OrderConvertor) EntityToClient(o *entity.Order) *client.Order {
	c.check(o)
	return &client.Order{
		CustomerId:  o.CustomerID,
		Id:          o.ID,
		Status:      o.Status,
		PaymentLink: o.PaymentLink,
		Items:       NewItemConvertor().EntitiesToClients(o.Items),
	}
}

func (c *OrderConvertor) ClientToEntity(o *client.Order) *entity.Order {
	c.check(o)
	return entity.NewOrder(
		NewItemConvertor().ClientsToEntities(o.Items),
		o.PaymentLink,
		o.Status,
		o.CustomerId,
		o.Id,
	)

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
	return entity.NewItem(o.ID, o.Name, o.Quantity, o.PriceID)
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
	return entity.NewItem(e.Id, e.Name, e.Quantity, e.PriceId)
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
	return entity.NewItemWithQuantity(i.ID, i.Quantity)
}

func (c *ItemWithQuantityConvertor) ClientsToEntities(items []client.ItemWithQuantity) (res []*entity.ItemWithQuantity) {
	for _, i := range items {
		res = append(res, c.ClientToEntity(i))
	}
	return
}

func (c *ItemWithQuantityConvertor) ClientToEntity(i client.ItemWithQuantity) *entity.ItemWithQuantity {
	return entity.NewItemWithQuantity(i.Id, i.Quantity)
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
