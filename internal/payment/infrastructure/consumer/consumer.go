package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mushroomyuan/gorder/common/broker"
	"github.com/mushroomyuan/gorder/common/genproto/orderpb"
	"github.com/mushroomyuan/gorder/payment/app"
	"github.com/mushroomyuan/gorder/payment/app/command"
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
			c.handleMessge(msg, q)
		}
	}()
	select {}
}

func (c *Consumer) handleMessge(msg amqp.Delivery, q amqp.Queue) {
	logrus.Infof("Payment receive a message from %s,msg = %s", q.Name, string(msg.Body))
	ctx := broker.ExtractRabbitMQHeaders(context.Background(), msg.Headers)
	t := otel.Tracer("rabbitmq")
	_, span := t.Start(ctx, fmt.Sprintf("rabbitmq.%s.consume", q.Name))
	defer span.End()

	o := &orderpb.Order{}
	if err := json.Unmarshal(msg.Body, o); err != nil {
		logrus.Infof("failed to unmarshal message from %s,err=%v", q.Name, err)
		_ = msg.Nack(false, false)
		return
	}
	if _, err := c.app.Commands.CreatePayment.Handle(ctx, command.CreatePayment{
		Order: o,
	}); err != nil {
		// TODO:retry
		logrus.Infof("failed to create payment for order %s,err=%v", o.ID, err)
		_ = msg.Nack(false, false)
		return
	}
	span.AddEvent("payment.created")
	err := msg.Ack(false)
	if err != nil {
		logrus.Warnf("failed to ack message from %s,err=%v", q.Name, err)
	}
	logrus.Infof("create payment for order %s success", o.ID)

}
