package query

import (
	"context"

	"github.com/mushroomyuan/gorder/common/decorator"
	"github.com/mushroomyuan/gorder/common/genproto/orderpb"
	domain "github.com/mushroomyuan/gorder/stock/domain/stock"
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

// TODO: remove this
var stub = map[string]string{
	"1": "price_1Rp9RiFh5LY9IvRqBkfBWmYM",
	"2": "price_1RpC4hFh5LY9IvRqHyHoh1xO",
}

func (c checkIfItemsInStockHandler) Handle(ctx context.Context, query CheckIfItemsInStock) (*CheckIfItemsInStockResult, error) {
	var items []*orderpb.Item
	for _, item := range query.Items {
		// TODO: 改为从数据库或者stripe获取
		priceID, ok := stub[item.ID]
		if !ok {
			priceID = stub["1"]
		}
		items = append(items, &orderpb.Item{
			ID:       item.ID,
			Quantity: item.Quantity,
			PriceID:  priceID,
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
