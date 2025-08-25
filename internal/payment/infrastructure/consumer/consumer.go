package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mushroomyuan/gorder/common/broker"
	"github.com/mushroomyuan/gorder/common/entity"
	"github.com/mushroomyuan/gorder/common/logging"
	"github.com/mushroomyuan/gorder/payment/app"
	"github.com/mushroomyuan/gorder/payment/app/command"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
)

type Consumer struct {
	app app.Application
}

func NewConsumer(app app.Application) *Consumer {
	return &Consumer{
		app: app,
	}
}

func (c *Consumer) Listen(ch *amqp.Channel) {
	q, err := ch.QueueDeclare(broker.EventOrderCreated, true, false, false, false, nil)
	if err != nil {
		logrus.Fatal(err)
	}
	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		logrus.Warnf("failed to consume queue=%s,err=%v", q.Name, err)
	}

	go func() {
		for msg := range msgs {
			c.handleMessage(ch, msg, q)
		}
	}()
	select {}
}

func (c *Consumer) handleMessage(ch *amqp.Channel, msg amqp.Delivery, q amqp.Queue) {
	t := otel.Tracer("rabbitmq")
	ctx, span := t.Start(broker.ExtractRabbitMQHeaders(context.Background(), msg.Headers), fmt.Sprintf("rabbitmq.%s.consume", q.Name))
	defer span.End()
	logging.Infof(ctx, nil, "Payment receive a message from %s,msg = %s", q.Name, string(msg.Body))
	var err error
	defer func() {
		if err != nil {
			logging.Warnf(ctx, nil, "consume failed|| from=%s ||msg=%+v ||err=%v", q.Name, msg, err)
			_ = msg.Nack(false, false)
		} else {
			logging.Infof(ctx, nil, "%s", "consume success")
			_ = msg.Ack(false)
		}
	}()

	o := &entity.Order{}
	if err := json.Unmarshal(msg.Body, o); err != nil {
		err = errors.Wrapf(err, "failed to unmarshal order")
		return
	}
	if _, err := c.app.Commands.CreatePayment.Handle(ctx, command.CreatePayment{
		Order: o,
	}); err != nil {
		err = errors.Wrapf(err, "failed to create payment")
		if err = broker.HandleRetry(ctx, ch, &msg); err != nil {
			err = errors.Wrapf(err, "retry_error,error handling retry,messageID=%s,error=%v", msg.MessageId, err)
		}
		return
	}
	span.AddEvent("payment.created")
	logrus.Infof("create payment for order %s success", o.ID)
}
