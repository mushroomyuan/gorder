package query

import (
	"context"

	"github.com/mushroomyuan/gorder/common/decorator"
	domain "github.com/mushroomyuan/gorder/stock/domain/stock"
	"github.com/mushroomyuan/gorder/stock/entity"
	"github.com/sirupsen/logrus"
)

type GetItems struct {
	ItemIDs []string
}

type GetItemsResult struct {
	Items []*entity.Item
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
