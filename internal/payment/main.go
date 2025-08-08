package main

import (
	"context"

	"github.com/mushroomyuan/gorder/common/broker"
	"github.com/mushroomyuan/gorder/common/config"
	"github.com/mushroomyuan/gorder/common/logging"
	"github.com/mushroomyuan/gorder/common/server"
	"github.com/mushroomyuan/gorder/common/tracing"
	"github.com/mushroomyuan/gorder/payment/infrastructure/consumer"
	"github.com/mushroomyuan/gorder/payment/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	logging.Init()
	if err := config.NewViperConfig(); err != nil {
		logrus.Fatal(err)
	}
}

func main() {

	serverType := viper.GetString("payment.server-to-run")
	serviceName := viper.GetString("payment.service-name")
	logrus.Infof("serverType: %s, serviceName: %s", serverType, serviceName)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shutdown, err := tracing.InitJaegerProvider(viper.GetString("jaeger.url"), serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer shutdown(ctx)

	application, cleanup := service.NewApplication(ctx)
	defer cleanup()

	ch, closeCh := broker.Connect(
		viper.GetString("rabbitmq.user"),
		viper.GetString("rabbitmq.password"),
		viper.GetString("rabbitmq.host"),
		viper.GetString("rabbitmq.port"),
	)

	defer func() {
		_ = ch.Close()
		_ = closeCh()
	}()

	go consumer.NewConsumer(application).Listen(ch)

	paymentHandler := NewPaymentHandler(ch)
	switch serverType {
	case "http":
		server.RunHTTPServer(serviceName, paymentHandler.RegisterRoutes)
	case "grpc":
		//TODO: implement me
		logrus.Panic("work in progress")
	default:
		logrus.Panicf("unknow server type: %s", serverType)
	}

}
