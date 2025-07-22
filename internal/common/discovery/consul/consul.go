package consul

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	consul "github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
)

var (
	consulClient *Registry
	once         sync.Once
	initErr      error
)

func New(consulAddr string) (*Registry, error) {
	once.Do(func() {
		config := consul.DefaultConfig()
		config.Address = consulAddr
		client, err := consul.NewClient(config)
		if err != nil {
			initErr = err
			return
		}
		consulClient = &Registry{client: client}
	})
	if initErr != nil {
		return nil, initErr
	}
	return consulClient, nil
}

type Registry struct {
	client *consul.Client
}

func (r *Registry) Register(_ context.Context, instanceId, serviceName, hostPort string) error {
	parts := strings.Split(hostPort, ":")
	if len(parts) != 2 {
		return errors.New("invalid host:port format")
	}
	host := parts[0]
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("invalid port format: %w", err)
	}
	r.client.Agent().ServiceRegister(&consul.AgentServiceRegistration{
		ID:      instanceId,
		Name:    serviceName,
		Port:    port,
		Address: host,
		Check: &consul.AgentServiceCheck{
			CheckID:                        instanceId,
			TLSSkipVerify:                  true,
			TTL:                            "5s",
			Timeout:                        "5s",
			DeregisterCriticalServiceAfter: "10s",
			HTTP:                           fmt.Sprintf("http://%s:%d/health", host, port),
		},
	})
	return nil
}

func (r *Registry) DeRegister(_ context.Context, instanceId, serviceName string) error {
	logrus.WithFields(logrus.Fields{
		"instanceId":  instanceId,
		"serviceName": serviceName,
	}).Info("DeRegister from consul")
	return r.client.Agent().CheckDeregister(instanceId)
}

func (r *Registry) Discover(ctx context.Context, serviceName string) ([]string, error) {
	entries, _, err := r.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, err
	}
	var ips []string
	for _, e := range entries {
		ips = append(ips, fmt.Sprintf("%s:%d", e.Service.Address, e.Service.Port))
	}
	return ips, nil
}

func (r *Registry) HealthCheck(instanceId, serviceName string) error {
	return r.client.Agent().UpdateTTL(instanceId, "online", consul.HealthPassing)
}
