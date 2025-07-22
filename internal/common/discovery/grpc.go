package discovery

import (
	"context"
	"time"

	"github.com/mushroomyuan/gorder/common/discovery/consul"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func RegistryToConsul(ctx context.Context, serviceName string) (func() error, error) {
	registry, err := consul.New(viper.GetString("consul.addr"))
	if err != nil {
		return func() error {
			return nil
		}, err
	}
	instanceId := generateInstanceId(serviceName)
	grpcAddr := viper.Sub(serviceName).GetString("grpc-addr")
	if err := registry.Register(ctx, instanceId, serviceName, grpcAddr); err != nil {
		return func() error {
			return nil
		}, err
	}
	go func() {
		for {
			if err := registry.HealthCheck(instanceId, serviceName); err != nil {
				logrus.Panicf("no heart beat from %s to registry,err=%v", serviceName, err)
			}
			time.Sleep(time.Second * 1)
		}
	}()

	logrus.WithFields(logrus.Fields{
		"serviceName": serviceName,
		"addr":        grpcAddr,
	}).Info("Registry to consul success")
	return func() error {
		return registry.DeRegister(ctx, instanceId, serviceName)
	}, nil
}
