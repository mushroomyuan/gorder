package main

import (
	"context"

	_ "github.com/mushroomyuan/gorder/common/config"
	"github.com/mushroomyuan/gorder/common/discovery"
	"github.com/mushroomyuan/gorder/common/genproto/stockpb"
	"github.com/mushroomyuan/gorder/common/logging"
	"github.com/mushroomyuan/gorder/common/server"
	"github.com/mushroomyuan/gorder/common/tracing"
	"github.com/mushroomyuan/gorder/stock/ports"
	"github.com/mushroomyuan/gorder/stock/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func init() {
	logging.Init()
}

func main() {
	serviceName := viper.GetString("stock.service-name")
	serviceType := viper.GetString("stock.server-to-run")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	shutdown, err := tracing.InitJaegerProvider(viper.GetString("jaeger.url"), serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer shutdown(ctx)

	application := service.NewApplication(ctx)
	deregisterFunc, err := discovery.RegistryToConsul(ctx, serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer func() {
		_ = deregisterFunc()
	}()

	switch serviceType {
	case "grpc":
		server.RunGRPCServer(serviceName, func(server *grpc.Server) {
			svc := ports.NewGRPCServer(application)
			stockpb.RegisterStockServiceServer(server, svc)
		})
	case "http":
	// TODO
	default:
		panic("unknown service type: " + serviceType)

	}

}
