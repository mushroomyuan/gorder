package integration

import (
	"context"

	_ "github.com/mushroomyuan/gorder/common/logging"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/product"
)

type StripeAPI struct {
	apiKey string
	stripe *stripe.Price
}

func NewStripeAPI() *StripeAPI {
	key := viper.GetString("stripe-key")
	//logrus.Infof("new stripe api key in stock: %s", key)
	if key == "" {
		logrus.Fatal("empty stripe-key")
	}
	return &StripeAPI{apiKey: key}
}

func (s *StripeAPI) GetPriceByProductID(ctx context.Context, pid string) (string, error) {
	stripe.Key = s.apiKey
	result, err := product.Get(pid, &stripe.ProductParams{})
	if err != nil {
		return "", err
	}
	return result.DefaultPrice.ID, err
}
