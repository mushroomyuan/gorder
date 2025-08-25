package grpc

import (
	"context"
	"errors"

	"github.com/mushroomyuan/gorder/common/genproto/orderpb"
	"github.com/mushroomyuan/gorder/common/genproto/stockpb"
	"github.com/mushroomyuan/gorder/common/logging"
)

type StockGRPC struct {
	client stockpb.StockServiceClient
}

func NewStockGRPC(client stockpb.StockServiceClient) *StockGRPC {
	return &StockGRPC{client: client}
}

func (s StockGRPC) CheckIfItemsInStock(ctx context.Context, items []*orderpb.ItemWithQuantity) (resp *stockpb.CheckIfItemsInStockResponse, err error) {
	_, dLog := logging.WhenRequest(ctx, "StockGRPC.CheckIfItemsInStock", items)
	defer dLog(resp, &err)
	if items == nil {
		return nil, errors.New("grpc items cannot be nil")
	}
	return s.client.CheckIfItemsInStock(ctx, &stockpb.CheckIfItemsInStockRequest{Items: items})

}

func (s StockGRPC) GetItems(ctx context.Context, itemsIDs []string) (items []*orderpb.Item, err error) {
	_, dLog := logging.WhenRequest(ctx, "StockGRPC.GetItems", items)
	defer dLog(items, &err)

	var getItemsRequest = &stockpb.GetItemsRequest{
		ItemsIDs: itemsIDs,
	}
	response, err := s.client.GetItems(ctx, getItemsRequest)
	if err != nil {
		return nil, err
	}
	return response.Items, nil
}
