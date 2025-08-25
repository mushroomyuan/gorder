package grpc

import (
	"context"

	"github.com/mushroomyuan/gorder/common/genproto/orderpb"
	"github.com/mushroomyuan/gorder/common/tracing"
	"google.golang.org/grpc/status"
)

type OrderGRPC struct {
	client orderpb.OrderServiceClient
}

func NewOrderGRPC(client orderpb.OrderServiceClient) *OrderGRPC {
	return &OrderGRPC{client: client}
}

func (o *OrderGRPC) UpdateOrder(ctx context.Context, order *orderpb.Order) (err error) {
	ctx, span := tracing.Start(ctx, "order_grpc.update_order")
	defer span.End()
	_, err = o.client.UpdateOrder(ctx, order)
	return status.Convert(err).Err()
}
