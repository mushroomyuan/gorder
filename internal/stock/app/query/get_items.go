package query

import (
	"github.com/mushroomyuan/gorder/common/decorator"
	"github.com/mushroomyuan/gorder/common/genproto/orderpb"
	domain "github.com/mushroomyuan/gorder/stock/domain/stock"
	"context"
	"github.com/sirupsen/logrus"
)

type GetItems struct {
	ItemIDs []string
}

type GetItemsResult struct {
	Items []*orderpb.Item
}

type GetItemsHandler decorator.QueryHandler[GetItems, *GetItemsResult]

type getItemsHandler struct {
	stockRepo domain.Repository
}

func (c getItemsHandler) Handle(ctx context.Context, query GetItems) (*GetItemsResult, error) {
	items, err := c.stockRepo.GetItems(ctx, query.ItemIDs)
	if err != nil {
		return nil, err
	}
	return &GetItemsResult{Items: items}, nil
}

func NewGetItemsHandler(
	stockRepo domain.Repository,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) GetItemsHandler {
	if stockRepo == nil {
		panic("stockRepo is nil")
	}
	return decorator.ApplyQueryDecorators[GetItems, *GetItemsResult](
		getItemsHandler{stockRepo: stockRepo},
		logger,
		metricsClient,
	)
}
