package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/mushroomyuan/gorder/common/broker"
	grpcClient "github.com/mushroomyuan/gorder/common/client"
	_ "github.com/mushroomyuan/gorder/common/config"
	"github.com/mushroomyuan/gorder/common/logging"
	"github.com/mushroomyuan/gorder/common/tracing"
	"github.com/mushroomyuan/gorder/kitchen/adapters"
	"github.com/mushroomyuan/gorder/kitchen/infrastructure/consumer"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	logging.Init()
}

func main() {
	serviceName := viper.GetString("kitchen.service-name")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shutdown, err := tracing.InitJaegerProvider(viper.GetString("jaeger.url"), serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer shutdown(ctx)

	//deregisterFunc, err := discovery.RegistryToConsul(ctx, serviceName)
	//if err != nil {
	//	logrus.Fatal(err)
	//}
	//defer func() {
	//	_ = deregisterFunc()
	//}()

	orderClient, closeFunc, err := grpcClient.NewOrderGRPCClient(ctx)
	if err != nil {
		logrus.Fatal(err)
	}
	defer closeFunc()

	orderGRPC := adapters.NewOrderGRPC(orderClient)

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
	go consumer.NewConsumer(orderGRPC).Listen(ch)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		logrus.Infof("received shutdown signal, shutting down gracefully")
		os.Exit(0)
	}()
	logrus.Printf("to exit,press Ctrl+C")
	select {}
}
