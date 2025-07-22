package ports

import (
	context "context"
	"github.com/mushroomyuan/gorder/common/genproto/orderpb"
	"github.com/mushroomyuan/gorder/common/genproto/stockpb"
	"github.com/mushroomyuan/gorder/stock/app"
	"github.com/sirupsen/logrus"
)

//func NewGRPCServer() *GRPCServer {
//	return &GRPCServer{}
//}

type GRPCServer struct {
	app app.Application
}

func NewGRPCServer(app app.Application) *GRPCServer {
	return &GRPCServer{app: app}
}

func (G GRPCServer) GetItems(ctx context.Context, request *stockpb.GetItemsRequest) (*stockpb.GetItemsResponse, error) {
	logrus.Info("rpc_request_in,stock.GetItems")
	defer func() {
		logrus.Info("rpc_request_out,stock.GetItems")
	}()
	fake := []*orderpb.Item{
		{ID: "fake-item-from-stock-GetItems"},
	}
	return &stockpb.GetItemsResponse{Items: fake}, nil
}

func (G GRPCServer) CheckIfItemsInStock(ctx context.Context, requset *stockpb.CheckIfItemsInStockRequset) (*stockpb.CheckIfItemsInStockResponse, error) {
	logrus.Info("rpc_request_in,stock.CheckIfItemsInStock")
	defer func() {
		logrus.Info("rpc_request_out,stock.CheckIfItemsInStock")
	}()
	return nil, nil
}
