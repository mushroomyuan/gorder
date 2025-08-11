package ports

import (
	"context"

	"github.com/mushroomyuan/gorder/common/genproto/orderpb"
	"github.com/mushroomyuan/gorder/order/app"
	"github.com/mushroomyuan/gorder/order/app/command"
	"github.com/mushroomyuan/gorder/order/app/query"
	"github.com/mushroomyuan/gorder/order/convertor"
	domain "github.com/mushroomyuan/gorder/order/domain/order"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GRPCServer struct {
	app app.Application
}

func NewGRPCServer(app app.Application) *GRPCServer {
	return &GRPCServer{app: app}
}

func (G GRPCServer) CreateOrder(ctx context.Context, request *orderpb.CreateOrderRequest) (_ *emptypb.Empty, err error) {
	logrus.Infof("order_grpc||create_order||request_in||request=%v", request)

	_, err = G.app.Commands.CreateOrder.Handle(ctx, command.CreateOrder{
		CustomerID: request.CustomerID,
		Items:      convertor.NewItemWithQuantityConvertor().ProtosToEntities(request.Items),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create order: %v", err.Error())
	}

	return nil, err
}

func (G GRPCServer) GetOrder(ctx context.Context, request *orderpb.GetOrderRequest) (*orderpb.Order, error) {
	logrus.Infof("order_grpc||get_order||request_in||request=%v", request)
	order, err := G.app.Queries.GetCustomerOrder.Handle(ctx, query.GetCustomerOrder{
		OrderID:    request.OrderID,
		CustomerID: request.CustomerID,
	})
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "order not found: %v", err.Error())
	}
	return convertor.NewOrderConvertor().EntityToProto(order), nil
}

func (G GRPCServer) UpdataOrder(ctx context.Context, request *orderpb.Order) (_ *emptypb.Empty, err error) {
	logrus.Infof("order_grpc||update_order||request_in||request=%v", request)
	order, err := domain.NewOrder(
		request.ID,
		request.CustomerID,
		request.Status,
		request.PaymentLink,
		convertor.NewItemConvertor().ProtosToEntities(request.Items))
	if err != nil {
		err = status.Errorf(codes.Internal, "failed to update order: %v", err.Error())
		return nil, err
	}
	_, err = G.app.Commands.UpdateOrder.Handle(ctx, command.UpdateOrder{
		Order: order,
		UpdateFn: func(ctx context.Context, order *domain.Order) (*domain.Order, error) {
			return order, nil
		},
	})
	return nil, err
}
