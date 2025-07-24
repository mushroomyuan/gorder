package command

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/mushroomyuan/gorder/common/broker"
	"github.com/mushroomyuan/gorder/common/decorator"
	"github.com/mushroomyuan/gorder/common/genproto/orderpb"
	"github.com/mushroomyuan/gorder/order/app/query"
	domain "github.com/mushroomyuan/gorder/order/domain/order"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type CreateOrder struct {
	CustomerID string
	Items      []*orderpb.ItemWithQuantity
}

type CreateOrderResult struct {
	OrderID string
}

type CreateOrderHandler decorator.QueryHandler[CreateOrder, *CreateOrderResult]

type createOrderHandler struct {
	orderRepo domain.Repository
	stockGRPC query.StockService
	channel   *amqp.Channel
}

func (c createOrderHandler) Handle(ctx context.Context, cmd CreateOrder) (*CreateOrderResult, error) {
	validateItems, err := c.validate(ctx, cmd.Items)
	if err != nil {
		return nil, err
	}
	o, err := c.orderRepo.Create(ctx, &domain.Order{
		CustomerID: cmd.CustomerID,
		Items:      validateItems,
	})
	if err != nil {
		return nil, err
	}
	q, err := c.channel.QueueDeclare(broker.EventOrderCreated, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	marshalledOrder, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}
	err = c.channel.PublishWithContext(ctx, "", q.Name, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Body:         marshalledOrder,
	})
	if err != nil {
		return nil, err
	}
	return &CreateOrderResult{OrderID: o.ID}, nil
}

func (c createOrderHandler) validate(ctx context.Context, itemsWithQuantity []*orderpb.ItemWithQuantity) ([]*orderpb.Item, error) {
	if len(itemsWithQuantity) == 0 {
		return nil, errors.New("items is empty")
	}
	itemsWithQuantity = packItems(itemsWithQuantity)
	response, err := c.stockGRPC.CheckIfItemsInStock(ctx, itemsWithQuantity)
	if err != nil {
		return nil, err
	}
	return response.Items, nil
}

func packItems(items []*orderpb.ItemWithQuantity) []*orderpb.ItemWithQuantity {
	mergedItems := make(map[string]int32)
	for _, item := range items {
		mergedItems[item.ID] += item.Quantity
	}
	var result []*orderpb.ItemWithQuantity
	for id, quantity := range mergedItems {
		result = append(result, &orderpb.ItemWithQuantity{
			ID:       id,
			Quantity: quantity,
		})
	}
	return result
}

func NewCreateOrderHandler(
	orderRepo domain.Repository,
	stockGRPC query.StockService,
	channel *amqp.Channel,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) CreateOrderHandler {
	if orderRepo == nil {
		panic("orderRepo is nil")
	}
	if channel == nil {
		panic("channel is nil")
	}
	return decorator.ApplyCommandDecorators[CreateOrder, *CreateOrderResult](
		createOrderHandler{
			orderRepo: orderRepo,
			stockGRPC: stockGRPC,
			channel:   channel,
		},
		logger,
		metricsClient,
	)
}
