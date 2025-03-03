package main

import (
	"context"
	"github.com/mushroomyuan/gorder/common/config"
	"github.com/mushroomyuan/gorder/common/genproto/stockpb"
	"github.com/mushroomyuan/gorder/common/server"
	"github.com/mushroomyuan/gorder/stock/ports"
	"github.com/mushroomyuan/gorder/stock/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func init() {
	if err := config.NewViperConfig(); err != nil {
		logrus.Fatal(err)
	}
}

func main() {
	serviceName := viper.GetString("stock.service-name")
	serviceType := viper.GetString("stock.server-to-run")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	application := service.NewApplication(ctx)

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
