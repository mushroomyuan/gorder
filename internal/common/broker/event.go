package broker

import (
	"context"
	"encoding/json"

	"github.com/mushroomyuan/gorder/common/logging"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	EventOrderCreated = "order.created"
	EventOrderPaid    = "order.paid"
)

type RoutingType string

const (
	FanOut RoutingType = "fanout"
	Direct RoutingType = "direct"
)

type PublishEventReq struct {
	Channel  *amqp.Channel
	Routing  RoutingType
	Queue    string
	Exchange string
	Body     any
}

func PublishEvent(ctx context.Context, p PublishEventReq) (err error) {
	_, deferLog := logging.WhenEventPublish(ctx, p)
	defer deferLog(nil, &err)
	if err = checkParam(p); err != nil {
		return err
	}
	switch p.Routing {
	default:
		logging.Panicf(ctx, nil, "unknown routing type: %s", string(p.Routing))
	case FanOut:
		return fanOut(ctx, p)
	case Direct:
		return directQueue(ctx, p)
	}
	return nil
}

func checkParam(p PublishEventReq) error {
	if p.Channel == nil {
		return errors.New("channel is nil")
	}
	return nil
}

func directQueue(ctx context.Context, p PublishEventReq) (err error) {
	_, err = p.Channel.QueueDeclare(p.Queue, true, false, false, false, nil)
	if err != nil {
		return err
	}
	jsonBody, err := json.Marshal(p.Body)
	if err != nil {
		return err
	}
	return doPublish(ctx, p.Channel, p.Exchange, p.Queue, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Body:         jsonBody,
		Headers:      InjectRabbitMQHeaders(ctx),
	})
}

func doPublish(ctx context.Context, ch *amqp.Channel, exchange, key string, mandatory bool, immediate bool, msg amqp.Publishing) error {
	if err := ch.PublishWithContext(ctx, exchange, key, mandatory, immediate, msg); err != nil {
		logging.Warnf(ctx, nil, "_publish_event_failed||exchange=%s||key=%s||msg=%v", exchange, key, msg)
		return errors.Wrapf(err, "_publish_event_failed||exchange=%s||key=%s||msg=%v", exchange, key, msg)
	}
	return nil
}

func fanOut(ctx context.Context, p PublishEventReq) (err error) {
	jsonBody, err := json.Marshal(p.Body)
	if err != nil {
		return err
	}
	return doPublish(ctx, p.Channel, p.Exchange, "", false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Body:         jsonBody,
		Headers:      InjectRabbitMQHeaders(ctx),
	})
}
