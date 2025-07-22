package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/mushroomyuan/gorder/common/config"
	"github.com/mushroomyuan/gorder/common/genproto/orderpb"
	"github.com/mushroomyuan/gorder/common/server"
	"github.com/mushroomyuan/gorder/order/ports"
	"github.com/mushroomyuan/gorder/order/service"
	"github.com/mushroomyuan/gorder/common/discovery"
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
	// 测试启动一个http实例
	//log.Println("Listening:8082")
	//mux := http.NewServeMux()
	//mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//	log.Println(r.RequestURI)
	//	_, _ = w.Write([]byte("<h1>Welcome to the home page</h1>"))
	//})
	//mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
	//	log.Println(r.RequestURI)
	//	_, _ = io.WriteString(w, "pong")
	//})
	//http.Handle("/", mux)
	//if err := http.ListenAndServe(":8082", mux); err != nil {
	//	log.Fatal(err)
	//}
	serviceName := viper.Sub("order").GetString("service-name")
	//serviceType := viper.GetString("stock.service-to-run")


	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	application, cleanup := service.NewApplication(ctx)
	defer cleanup()

	deregisterFunc, err := discovery.RegistryToConsul(ctx, serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer func() {
		_ = deregisterFunc()
	}()

	go server.RunGRPCServer(serviceName, func(server *grpc.Server) {
		svc := ports.NewGRPCServer(application)
		orderpb.RegisterOrderServiceServer(server, svc)
	})

	// 启动http服务
	server.RunHTTPServer(serviceName, func(router *gin.Engine) {
		ports.RegisterHandlersWithOptions(router, &HTTPServer{
			app: application,
		}, ports.GinServerOptions{
			BaseURL:      "/api",
			Middlewares:  nil,
			ErrorHandler: nil,
		})
	})

}
