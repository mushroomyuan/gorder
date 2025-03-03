package ports

import (
	context "context"
	"github.com/mushroomyuan/gorder/common/genproto/stockpb"
	"github.com/mushroomyuan/gorder/stock/app"
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
	//TODO implement me
	panic("implement me")
}

func (G GRPCServer) CheckIfItemsInStock(ctx context.Context, requset *stockpb.CheckIfItemsInStockRequset) (*stockpb.CheckIfItemsInStockResponse, error) {
	//TODO implement me
	panic("implement me")
}
