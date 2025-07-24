package main

import (
	"github.com/mushroomyuan/gorder/common/broker"
	"github.com/mushroomyuan/gorder/common/config"
	"github.com/mushroomyuan/gorder/common/logging"
	"github.com/mushroomyuan/gorder/common/server"
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
	serverName := viper.GetString("payment.service-name")
	logrus.Infof("serverType: %s, serverName: %s", serverType, serverName)

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

	paymentHandler := NewPaymentHandler()
	switch serverType {
	case "http":
		server.RunHTTPServer(serverName, paymentHandler.RegisterRoutes)
	case "grpc":
		//TODO: implement me
		logrus.Panic("work in progress")
	default:
		logrus.Panicf("unknow server type: %s", serverType)
	}

}
