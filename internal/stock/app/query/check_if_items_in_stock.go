package query

import (
	"context"
	"strings"
	"time"

	"github.com/mushroomyuan/gorder/common/decorator"
	"github.com/mushroomyuan/gorder/common/handler/redis"
	domain "github.com/mushroomyuan/gorder/stock/domain/stock"
	"github.com/mushroomyuan/gorder/stock/entity"
	"github.com/mushroomyuan/gorder/stock/infrastructure/integration"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	redisLockPrefix = "check_stock_"
)

type CheckIfItemsInStock struct {
	Items []*entity.ItemWithQuantity
}

type CheckIfItemsInStockResult struct {
	Items []*entity.Item
}

type CheckIfItemsInStockHandler decorator.QueryHandler[CheckIfItemsInStock, []*entity.Item]

type checkIfItemsInStockHandler struct {
	stockRepo domain.Repository
	stripeAPI *integration.StripeAPI
}

// Deprecated
// var stub = map[string]string{
// 	"1": "price_1Rp9RiFh5LY9IvRqBkfBWmYM",
// 	"2": "price_1RpC4hFh5LY9IvRqHyHoh1xO",
// }

func NewCheckIfItemsInStockHandler(
	stockRepo domain.Repository,
	stripeAPI *integration.StripeAPI,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) CheckIfItemsInStockHandler {
	if stockRepo == nil {
		panic("stockRepo is nil")
	}
	if stripeAPI == nil {
		panic("stripeAPI is nil")
	}
	return decorator.ApplyQueryDecorators[CheckIfItemsInStock, []*entity.Item](
		checkIfItemsInStockHandler{
			stockRepo: stockRepo,
			stripeAPI: stripeAPI,
		},
		logger,
		metricsClient,
	)
}

func (c checkIfItemsInStockHandler) Handle(ctx context.Context, query CheckIfItemsInStock) ([]*entity.Item, error) {
	if err := lock(ctx, getLockKey(query)); err != nil {
		return nil, errors.Wrapf(err, "redis lock error,key=%s", getLockKey(query))
	}
	defer func() {
		if err := unlock(ctx, getLockKey(query)); err != nil {
			logrus.Warnf("redis unlock failed: %v", err)
		}
	}()

	var res []*entity.Item
	for _, item := range query.Items {
		priceID, err := c.stripeAPI.GetPriceByProductID(ctx, item.ID)
		if err != nil || priceID == "" {
			logrus.Warningf("GetPriceByProductID error:itemID=%s,err= %v", item.ID, err)
			return nil, err
		}
		res = append(res, &entity.Item{
			ID:       item.ID,
			Quantity: item.Quantity,
			PriceID:  priceID,
		})
	}

	if err := c.checkStock(ctx, query.Items); err != nil {
		return nil, err
	}
	//TODO:扣库存
	return res, nil
}

func getLockKey(query CheckIfItemsInStock) string {
	var ids []string
	for _, item := range query.Items {
		ids = append(ids, item.ID)
	}
	return redisLockPrefix + strings.Join(ids, ",")
}

func unlock(ctx context.Context, key string) error {
	return redis.Del(ctx, redis.LocalClient(), key)
}

func lock(ctx context.Context, key string) error {
	return redis.SetNX(ctx, redis.LocalClient(), key, "1", 5*time.Minute)
}

func (c checkIfItemsInStockHandler) checkStock(ctx context.Context, query []*entity.ItemWithQuantity) error {
	ids := make([]string, len(query))
	for i, item := range query {
		ids[i] = item.ID
	}
	records, err := c.stockRepo.GetStock(ctx, ids)
	if err != nil {
		return err
	}
	//var idQuantityMap map[string]int32
	idQuantityMap := make(map[string]int32, 100)
	for _, r := range records {
		idQuantityMap[r.ID] += r.Quantity
	}
	var (
		ok       = true
		failedOn []struct {
			ID   string
			Want int32
			Have int32
		}
	)

	for _, item := range query {
		if item.Quantity > idQuantityMap[item.ID] {
			ok = false
			failedOn = append(failedOn, struct {
				ID   string
				Want int32
				Have int32
			}{ID: item.ID, Want: item.Quantity, Have: idQuantityMap[item.ID]})
		}
	}
	if ok {
		return c.stockRepo.UpdateStock(ctx, query, func(
			ctx context.Context,
			existing []*entity.ItemWithQuantity,
			query []*entity.ItemWithQuantity,
		) ([]*entity.ItemWithQuantity, error) {
			var newItems []*entity.ItemWithQuantity
			for _, e := range existing {
				for _, q := range query {
					if e.ID == q.ID {
						newItems = append(newItems, &entity.ItemWithQuantity{
							ID:       e.ID,
							Quantity: e.Quantity - q.Quantity,
						})
					}
				}
			}
			return newItems, nil
		})
	}
	return domain.ExceedStockError{FailedOn: failedOn}
}
