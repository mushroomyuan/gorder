package adapters

import (
	"context"

	"github.com/mushroomyuan/gorder/common/genproto/orderpb"
)

type OrderGRPC struct {
	client orderpb.OrderServiceClient
}

func NewOrderGRPC(client orderpb.OrderServiceClient) *OrderGRPC {
	return &OrderGRPC{client: client}
}

func (g *OrderGRPC) UpdateOrder(ctx context.Context, req *orderpb.Order) error {
	_, err := g.client.UpdataOrder(ctx, req)
	return err
}
