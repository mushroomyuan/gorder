package client

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/mushroomyuan/gorder/common/discovery"
	"github.com/mushroomyuan/gorder/common/genproto/orderpb"
	"github.com/mushroomyuan/gorder/common/genproto/stockpb"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewStockGRPCClient(ctx context.Context) (client stockpb.StockServiceClient, close func() error, err error) {
	if !WaitForStockGRPCClient(viper.GetDuration("dial-grpc-timeout") * time.Second) {
		return nil, nil, errors.New("stock grpc not available")
	}
	grpcAddr, err := discovery.GetServiceAddr(ctx, viper.GetString("stock.service-name"))
	if err != nil {
		return nil, func() error { return err }, err
	}
	if grpcAddr == "" {
		logrus.Warn("empty grpc addr for stock grpc")
	}

	opts := grpcDialOptions(grpcAddr)
	if err != nil {
		return nil, func() error { return err }, err
	}
	conn, err := grpc.NewClient(grpcAddr, opts...)
	if err != nil {
		return nil, func() error { return err }, err
	}
	return stockpb.NewStockServiceClient(conn), conn.Close, nil
}

func NewOrderGRPCClient(ctx context.Context) (client orderpb.OrderServiceClient, close func() error, err error) {
	if !WaitForOrderGRPCClient(viper.GetDuration("dial-grpc-timeout") * time.Second) {
		return nil, nil, errors.New("order grpc not available")
	}
	grpcAddr, err := discovery.GetServiceAddr(ctx, viper.GetString("order.service-name"))
	if err != nil {
		return nil, func() error { return err }, err
	}
	if grpcAddr == "" {
		logrus.Warn("empty grpc addr for order grpc")
	}

	opts := grpcDialOptions(grpcAddr)
	if err != nil {
		return nil, func() error { return err }, err
	}
	conn, err := grpc.NewClient(grpcAddr, opts...)
	if err != nil {
		return nil, func() error { return err }, err
	}
	return orderpb.NewOrderServiceClient(conn), conn.Close, nil
}

func grpcDialOptions(_ string) []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
}

func WaitForOrderGRPCClient(timeout time.Duration) bool {
	logrus.Infof("configured timeout = %v", timeout)
	logrus.Infof("waiting for order grpc client timeout:%v seconds", timeout.Seconds())
	return waitFor(viper.GetString("order.grpc-addr"), timeout)
}

func WaitForStockGRPCClient(timeout time.Duration) bool {
	logrus.Infof("configured timeout = %v", timeout)
	logrus.Infof("waiting for stock grpc client timeout:%v seconds", timeout.Seconds())
	return waitFor(viper.GetString("stock.grpc-addr"), timeout)
}

func waitFor(addr string, timeout time.Duration) bool {
	portAvailable := make(chan struct{})
	timeoutChan := time.After(timeout)
	go func() {
		for {
			select {
			case <-timeoutChan:
				return
			default:
				// continue
			}
			_, err := net.Dial("tcp", addr)
			if err == nil {
				close(portAvailable)
				return
			}
			time.Sleep(200 * time.Millisecond)

		}
	}()

	select {
	case <-portAvailable:
		return true
	case <-timeoutChan:
		return false
	}
}
