package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mushroomyuan/gorder/common/broker"
	"github.com/mushroomyuan/gorder/common/entity"
	"github.com/mushroomyuan/gorder/common/logging"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	stripe "github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/webhook"
	"go.opentelemetry.io/otel"
)

type PaymentHandler struct {
	channel *amqp.Channel
}

func NewPaymentHandler(ch *amqp.Channel) *PaymentHandler {
	return &PaymentHandler{channel: ch}
}

// stripe listen --forward-to localhost:8284/api/webhook
func (h *PaymentHandler) RegisterRoutes(c *gin.Engine) {
	c.POST("/api/webhook", h.HandleWebhook)
}

func (h *PaymentHandler) HandleWebhook(c *gin.Context) {
	logrus.WithContext(c.Request.Context()).Info("received webhook from stripe")
	var err error
	defer func() {
		if err != nil {
			logging.Warnf(c.Request.Context(), nil, "handleWebhook error: %v", err)
		} else {
			logging.Infof(c.Request.Context(), nil, "%s", "handleWebhook success")
		}
	}()

	const MaxBodyBytes = int64(1024 * 1024)
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		err = errors.Wrapf(err, "failed to read request body: %v\n", err)
		c.Writer.WriteHeader(http.StatusServiceUnavailable)
		c.JSON(http.StatusServiceUnavailable, err)
		return
	}
	// Pass the request body and Stripe-Signature header to ConstructEvent, along
	// with the webhook signing key.
	event, err := webhook.ConstructEvent(payload, c.Request.Header.Get("Stripe-Signature"),
		viper.GetString("ENDPOINT_STRIPE_SECRET"))

	if err != nil {
		err = errors.Wrapf(err, "failed to construct event: %v\n", err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	switch event.Type {
	case stripe.EventTypeCheckoutSessionCompleted:
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		if session.PaymentStatus == stripe.CheckoutSessionPaymentStatusPaid {
			var items []*entity.Item
			_ = json.Unmarshal([]byte(session.Metadata["items"]), &items)

			t := otel.Tracer("rabbitmq")
			ctx, span := t.Start(c.Request.Context(), fmt.Sprintf("rabbitmq.%s.publish", broker.EventOrderPaid))
			defer span.End()

			err = broker.PublishEvent(ctx, broker.PublishEventReq{
				Channel:  h.channel,
				Routing:  broker.FanOut,
				Queue:    "",
				Exchange: broker.EventOrderPaid,
				Body: entity.NewValidOrder(
					items, session.Metadata["paymentLink"],
					string(stripe.CheckoutSessionPaymentStatusPaid),
					session.Metadata["customerID"],
					session.Metadata["orderID"],
				),
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "failed to publish event",
					"err":     err.Error(),
				})
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
