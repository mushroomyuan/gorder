package service

import (
	"context"

	"github.com/mushroomyuan/gorder/common/broker"
	grpcClient "github.com/mushroomyuan/gorder/common/client"
	"github.com/mushroomyuan/gorder/common/metrics"
	"github.com/mushroomyuan/gorder/payment/adapters/grpc"
	"github.com/mushroomyuan/gorder/payment/app"
	"github.com/mushroomyuan/gorder/payment/app/command"
	"github.com/mushroomyuan/gorder/payment/domain"
	"github.com/mushroomyuan/gorder/payment/infrastructure/processor"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewApplication(ctx context.Context) (app app.Application, cleanup func()) {
	orderClient, closeOrderClient, err := grpcClient.NewOrderGRPCClient(ctx)
	if err != nil {
		panic(err)
	}

	orderGRPC := grpc.NewOrderGRPC(orderClient)

	ch, closeCh := broker.Connect(
		viper.GetString("rabbitmq.user"),
		viper.GetString("rabbitmq.password"),
		viper.GetString("rabbitmq.host"),
		viper.GetString("rabbitmq.port"),
	)
	// memProcessor := processor.NewInmemProcessor()
	stripeProcessor := processor.NewStripeProcessor(viper.GetString("stripe-key"))
	logrus.Infof("stripe-key: %s", viper.GetString("stripe-key"))

	return newApplication(ctx, orderGRPC, stripeProcessor), func() {
		_ = closeOrderClient()
		_ = ch.Close()
		_ = closeCh()
	}

}

func newApplication(_ context.Context, orderGRPC command.OrderService, processor domain.Processor) app.Application {

	metricsClient := metrics.NewPrometheusMetricsClient(&metrics.PrometheusMetricsClientConfig{
		Host:        viper.GetString("payment.metrics-export-addr"),
		ServiceName: viper.GetString("payment.service-name"),
	})

	return app.Application{
		Commands: app.Commands{
			CreatePayment: command.NewCreatePaymentHandler(processor, orderGRPC, metricsClient),
		},
	}
}
