package command

import (
	"context"

	"github.com/mushroomyuan/gorder/common/decorator"
	"github.com/mushroomyuan/gorder/common/logging"
	domain "github.com/mushroomyuan/gorder/order/domain/order"
	"github.com/sirupsen/logrus"
)

type UpdateOrder struct {
	Order    *domain.Order
	UpdateFn func(context.Context, *domain.Order) (*domain.Order, error)
}

type UpdateOrderHandler decorator.QueryHandler[UpdateOrder, any]

type updateOrderHandler struct {
	orderRepo domain.Repository
	// stockGRPC
}

func (u updateOrderHandler) Handle(ctx context.Context, cmd UpdateOrder) (any, error) {
	var err error
	defer logging.WhenCommandExecuted(ctx, "UpdateOrderHandler", cmd, err)
	if cmd.UpdateFn == nil {
		logrus.Panicf("UpdateOrder handler must have UpdateOrder function,cmd=%+v", cmd)
	}
	if err = u.orderRepo.Update(ctx, cmd.Order, cmd.UpdateFn); err != nil {
		return nil, err
	}
	return nil, nil
}

func NewUpdateOrderHandler(
	orderRepo domain.Repository,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) UpdateOrderHandler {
	if orderRepo == nil {
		panic("orderRepo is nil")
	}
	return decorator.ApplyCommandDecorators[UpdateOrder, any](
		updateOrderHandler{orderRepo: orderRepo},
		logger,
		metricsClient,
	)
}
