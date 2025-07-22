package query

import (
	"github.com/mushroomyuan/gorder/common/decorator"
	"github.com/mushroomyuan/gorder/common/genproto/orderpb"
	domain "github.com/mushroomyuan/gorder/stock/domain/stock"
	"context"
	"github.com/sirupsen/logrus"
)

type CheckIfItemsInStock struct {
	Items []*orderpb.ItemWithQuantity
}

type CheckIfItemsInStockResult struct {
	Items []*orderpb.Item
}

type CheckIfItemsInStockHandler decorator.QueryHandler[CheckIfItemsInStock, *CheckIfItemsInStockResult]


type checkIfItemsInStockHandler struct {
	stockRepo domain.Repository
}

func (c checkIfItemsInStockHandler) Handle(ctx context.Context, query CheckIfItemsInStock) (*CheckIfItemsInStockResult, error) {
	var items []*orderpb.Item
	for _,item := range query.Items {	
		items = append(items, &orderpb.Item{
			ID:       item.ID,
			Quantity: item.Quantity,
		})
	}
	
	return &CheckIfItemsInStockResult{Items: items}, nil
}


func NewCheckIfItemsInStockHandler(
	stockRepo domain.Repository,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) CheckIfItemsInStockHandler {
	return decorator.ApplyQueryDecorators[CheckIfItemsInStock, *CheckIfItemsInStockResult](
		checkIfItemsInStockHandler{stockRepo: stockRepo},
		logger,
		metricsClient,
	)
}