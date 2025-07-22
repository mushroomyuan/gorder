package grpc

import (
	"context"

	"github.com/mushroomyuan/gorder/common/genproto/orderpb"
	"github.com/mushroomyuan/gorder/common/genproto/stockpb"
	"github.com/sirupsen/logrus"
)

type StockGRPC struct {
	client stockpb.StockServiceClient
}

func NewStockGRPC(client stockpb.StockServiceClient) *StockGRPC {
	return &StockGRPC{client: client}
}

func (s StockGRPC) CheckIfItemsInStock(ctx context.Context, items []*orderpb.ItemWithQuantity) error {
	response, err := s.client.CheckIfItemsInStock(ctx, &stockpb.CheckIfItemsInStockRequset{Items: items})
	logrus.Info("stock_grpc response", response)
	return err
}

func (s StockGRPC) GetItems(ctx context.Context, itemsIDs []string) ([]*orderpb.Item, error) {
	var getItemsRequest = &stockpb.GetItemsRequest{
		ItemsIDs: itemsIDs,
	}
	response, err := s.client.GetItems(ctx, getItemsRequest)
	logrus.Info("get_items response", response)
	if err != nil {
		return nil, err
	}
	return response.Items, nil
}
