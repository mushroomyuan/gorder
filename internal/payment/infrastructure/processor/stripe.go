package processor

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mushroomyuan/gorder/common/entity"
	"github.com/mushroomyuan/gorder/common/tracing"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
)

type StripeProcessor struct {
	apiKey string
}

func NewStripeProcessor(apiKey string) *StripeProcessor {
	if apiKey == "" {
		panic("apiKey is required")
	}

	stripe.Key = apiKey
	return &StripeProcessor{apiKey: apiKey}
}

const (
	SuccessURL = "http://localhost:8282/success"
)

func (s *StripeProcessor) CreatePaymentLink(ctx context.Context, order *entity.Order) (string, error) {
	_, span := tracing.Start(ctx, "stripe_processor.create_payment_link")
	defer span.End()
	var items []*stripe.CheckoutSessionLineItemParams
	for _, item := range order.Items {
		items = append(items, &stripe.CheckoutSessionLineItemParams{
			Price:    stripe.String(item.PriceID),
			Quantity: stripe.Int64(int64(item.Quantity)),
		})
	}
	marshallItems, err := json.Marshal(order.Items)
	if err != nil {
		return "", err
	}
	metadata := map[string]string{
		"orderID":     order.ID,
		"customerID":  order.CustomerID,
		"status":      order.Status,
		"items":       string(marshallItems),
		"paymentLink": order.PaymentLink,
	}
	params := &stripe.CheckoutSessionParams{
		LineItems:  items,
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(fmt.Sprintf("%s?customerID=%s&orderID=%s", SuccessURL, order.CustomerID, order.ID)),
		Metadata:   metadata,
	}
	result, err := session.New(params)
	if err != nil {
		return "", err
	}

	return result.URL, nil
}
