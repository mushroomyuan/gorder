package service

import (
	"context"

	"github.com/mushroomyuan/gorder/common/broker"
	grpcClient "github.com/mushroomyuan/gorder/common/client"
	"github.com/mushroomyuan/gorder/common/metrics"
	"github.com/mushroomyuan/gorder/order/adapters"
	"github.com/mushroomyuan/gorder/order/adapters/grpc"
	"github.com/mushroomyuan/gorder/order/app"
	"github.com/mushroomyuan/gorder/order/app/command"
	"github.com/mushroomyuan/gorder/order/app/query"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewApplication(ctx context.Context) (app app.Application, cleanup func()) {

	stockClient, closeStockClient, err := grpcClient.NewStockGRPCClient(ctx)
	if err != nil {
		panic(err)
	}
	stockGRPC := grpc.NewStockGRPC(stockClient)
	ch, closeCh := broker.Connect(
		viper.GetString("rabbitmq.user"),
		viper.GetString("rabbitmq.password"),
		viper.GetString("rabbitmq.host"),
		viper.GetString("rabbitmq.port"),
	)

	return newApplication(ctx, stockGRPC, ch), func() {
		_ = closeStockClient()
		_ = ch.Close()
		_ = closeCh()
	}
}

func newApplication(_ context.Context, stockGRPC query.StockService, ch *amqp.Channel) app.Application {
	orderInmemRepo := adapters.NewMemoryOrderRepository()
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricsClient := metrics.TodoMetrics{}
	return app.Application{
		Commands: app.Commands{
			CreateOrder: command.NewCreateOrderHandler(orderInmemRepo, stockGRPC, ch, logger, metricsClient),
			UpdateOrder: command.NewUpdateOrderHandler(orderInmemRepo, logger, metricsClient),
		},
		Queries: app.Queries{
			GetCustomerOrder: query.NewGetCustomerOrderHandler(orderInmemRepo, logger, metricsClient),
		},
	}
}
