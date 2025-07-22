package command

import (
	"context"
	"github.com/mushroomyuan/gorder/common/decorator"
	"github.com/mushroomyuan/gorder/common/genproto/orderpb"
	"github.com/mushroomyuan/gorder/order/app/query"
	domain "github.com/mushroomyuan/gorder/order/domain/order"
	"github.com/sirupsen/logrus"
	"errors"
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
}

func (c createOrderHandler) Handle(ctx context.Context, cmd CreateOrder) (*CreateOrderResult, error) {
	validateItems,err:= c.validate(ctx,cmd.Items)
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
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) CreateOrderHandler {
	if orderRepo == nil {
		panic("orderRepo is nil")
	}
	return decorator.ApplyCommandDecorators[CreateOrder, *CreateOrderResult](
		createOrderHandler{orderRepo: orderRepo, stockGRPC: stockGRPC},
		logger,
		metricsClient,
	)
}
