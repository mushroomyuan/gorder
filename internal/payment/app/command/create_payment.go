package command

import (
	"context"

	"github.com/mushroomyuan/gorder/common/decorator"
	"github.com/mushroomyuan/gorder/common/genproto/orderpb"
	"github.com/mushroomyuan/gorder/common/tracing"
	domain "github.com/mushroomyuan/gorder/payment/domain"
	"github.com/sirupsen/logrus"
)

type CreatePayment struct {
	Order *orderpb.Order
}

type CreatePaymentHandler decorator.CommandHandler[CreatePayment, string]

type createPaymentHandler struct {
	processor domain.Processor
	orderGRPC OrderService
}

func NewCreatePaymentHandler(
	processor domain.Processor,
	orderGRPC OrderService,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) CreatePaymentHandler {
	if processor == nil {
		panic("NewCreatePaymentHandler paramter processor is nil")
	}
	if orderGRPC == nil {
		panic("NewCreatePaymentHandler paramter orderGRPC is nil")
	}
	return decorator.ApplyCommandDecorators[CreatePayment, string](
		createPaymentHandler{
			processor: processor,
			orderGRPC: orderGRPC,
		},
		logger,
		metricsClient,
	)
}

func (c createPaymentHandler) Handle(ctx context.Context, cmd CreatePayment) (string, error) {
	_, span := tracing.Start(ctx, "create_payment")
	defer span.End()

	paymentLink, err := c.processor.CreatePaymentLink(ctx, cmd.Order)
	if err != nil {
		return "", err
	}
	logrus.Infof("create payment link for order:%s succese,paymentLink:%s", cmd.Order.ID, paymentLink)
	newOrder := &orderpb.Order{
		ID:          cmd.Order.ID,
		CustomerID:  cmd.Order.CustomerID,
		Status:      "waiting for payment",
		PaymentLink: paymentLink,
		Items:       cmd.Order.Items,
	}

	err = c.orderGRPC.UpdateOrder(ctx, newOrder)
	return paymentLink, err
}
