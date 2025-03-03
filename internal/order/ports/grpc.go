package ports

import (
	context "context"
	"github.com/mushroomyuan/gorder/common/genproto/orderpb"
	"github.com/mushroomyuan/gorder/order/app"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GRPCServer struct {
	app app.Application
}

func NewGRPCServer(app app.Application) *GRPCServer {
	return &GRPCServer{app: app}
}

func (G GRPCServer) CreateOrder(ctx context.Context, request *orderpb.CreateOrderRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (G GRPCServer) GetOrder(ctx context.Context, request *orderpb.GetOrderRequest) (*orderpb.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (G GRPCServer) UpdataOrder(ctx context.Context, order *orderpb.Order) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}
