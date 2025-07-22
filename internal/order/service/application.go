package service

import (
	"context"
	grpcClient "github.com/mushroomyuan/gorder/common/client"
	"github.com/mushroomyuan/gorder/common/metrics"
	"github.com/mushroomyuan/gorder/order/adapters"
	"github.com/mushroomyuan/gorder/order/adapters/grpc"
	"github.com/mushroomyuan/gorder/order/app"
	"github.com/mushroomyuan/gorder/order/app/command"
	"github.com/mushroomyuan/gorder/order/app/query"
	"github.com/sirupsen/logrus"
)

func NewApplication(ctx context.Context) (app app.Application, cleanup func()) {

	stockClient, closeStockClient, err := grpcClient.NewStockGRPCClient(ctx)
	if err != nil {
		panic(err)
	}
	stockGRPC := grpc.NewStockGRPC(stockClient)
	return newApplication(ctx, stockGRPC), func() {
		_ = closeStockClient()
	}
}

func newApplication(_ context.Context, stockGRPC query.StockService) app.Application {
	orderInmemRepo := adapters.NewMemoryOrderRepository()
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricsClient := metrics.TodoMetrics{}
	return app.Application{
		Commands: app.Commands{
			CreateOrder: command.NewCreateOrderHandler(orderInmemRepo, stockGRPC, logger, metricsClient),
			UpdateOrder: command.NewUpdateOrderHandler(orderInmemRepo, logger, metricsClient),
		},
		Queries: app.Queries{
			GetCustomerOrder: query.NewGetCustomerOrderHandler(orderInmemRepo, logger, metricsClient),
		},
	}
}
