package query

import (
	"github.com/mushroomyuan/gorder/common/decorator"
	"github.com/mushroomyuan/gorder/common/genproto/orderpb"
	domian "github.com/mushroomyuan/gorder/internal/stock/domain/stock"
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