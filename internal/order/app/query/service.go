package query

import (
	"context"
	"github.com/mushroomyuan/gorder/common/genproto/orderpb"
)

type StockService interface {
	CheckIfItemsInStock(ctx context.Context, items []*orderpb.ItemWithQuantity) error
	GetItems(ctx context.Context, itemsIDs []string) ([]*orderpb.Item, error)
}
