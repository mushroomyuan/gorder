package consumer

import (
	"github.com/mushroomyuan/gorder/common/broker"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type Consumer struct{}

func NewConsumer() *Consumer {
	return &Consumer{}
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
			c.handleMessge(msg, q, ch)
		}
	}()
	select {}
}

func (c *Consumer) handleMessge(msg amqp.Delivery, q amqp.Queue, ch *amqp.Channel) {
	logrus.Infof("Payment receive a message from %s,msg = %s", q.Name, string(msg.Body))
	err := msg.Ack(false)
	if err != nil {
		logrus.Warnf("failed to ack message from %s,err=%v", q.Name, err)
	}
}
