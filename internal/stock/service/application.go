package service

import (
	"context"

	"github.com/mushroomyuan/gorder/common/metrics"
	"github.com/mushroomyuan/gorder/stock/adapters"
	"github.com/mushroomyuan/gorder/stock/app"
	"github.com/mushroomyuan/gorder/stock/app/query"
	"github.com/mushroomyuan/gorder/stock/infrastructure/integration"
	"github.com/mushroomyuan/gorder/stock/infrastructure/persistent"
	"github.com/spf13/viper"
)

func NewApplication(_ context.Context) app.Application {
	//stockRepo := adapters.NewMemoryStockRepository()
	db := persistent.NewMySQL()
	stockRepo := adapters.NewMySQLStockRepository(db)
	stripeAPI := integration.NewStripeAPI()

	metricsClient := metrics.NewPrometheusMetricsClient(&metrics.PrometheusMetricsClientConfig{
		Host:        viper.GetString("stock.metrics-export-addr"),
		ServiceName: viper.GetString("stock.service-name"),
	})

	return app.Application{
		Commands: app.Commands{},
		Queries: app.Queries{
			CheckIfItemsInStock: query.NewCheckIfItemsInStockHandler(stockRepo, stripeAPI, metricsClient),
			GetItems:            query.NewGetItemsHandler(stockRepo, metricsClient),
		},
	}
}
