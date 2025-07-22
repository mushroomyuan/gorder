package command

import (
	"context"
	"github.com/mushroomyuan/gorder/common/decorator"
	"github.com/mushroomyuan/gorder/common/genproto/orderpb"
	"github.com/mushroomyuan/gorder/order/app/query"
	domain "github.com/mushroomyuan/gorder/order/domain/order"
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
}

func (c createOrderHandler) Handle(ctx context.Context, cmd CreateOrder) (*CreateOrderResult, error) {
	// TODO: call stock grpc to get items.
	err := c.stockGRPC.CheckIfItemsInStock(ctx, cmd.Items)
	if err != nil {
		return nil, err
	}
	resp, err := c.stockGRPC.GetItems(ctx, []string{"123"})
	if err != nil {
		return nil, err
	}
	logrus.Info("createOrderHandler||err from stockGRPC GetItems", resp)
	var stockResponse []*orderpb.Item
	for _, item := range cmd.Items {
		stockResponse = append(stockResponse, &orderpb.Item{
			ID:       item.ID,
			Quantity: item.Quantity,
		})
	}
	o, err := c.orderRepo.Create(ctx, &domain.Order{
		CustomerID: cmd.CustomerID,
		Items:      stockResponse,
	})
	if err != nil {
		return nil, err
	}
	return &CreateOrderResult{OrderID: o.ID}, nil
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
