package adapters

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/mushroomyuan/gorder/common/logging"
	domain "github.com/mushroomyuan/gorder/order/domain/order"
	"github.com/sirupsen/logrus"
)

type MemoryOrderRepository struct {
	lock  *sync.RWMutex
	store []*domain.Order
}

func NewMemoryOrderRepository() *MemoryOrderRepository {
	//s := make([]*domain.Order, 0)
	//s = append(s, &domain.Order{
	//	ID:          "fake-ID",
	//	CustomerID:  "fake-CustomerID",
	//	Status:      "fake-Status",
	//	PaymentLink: "fake-PaymentLink",
	//	Items:       nil,
	//})
	return &MemoryOrderRepository{
		lock:  &sync.RWMutex{},
		store: make([]*domain.Order, 0),
		//store: s,
	}
}

func (m *MemoryOrderRepository) Create(ctx context.Context, order *domain.Order) (created *domain.Order, err error) {
	_, deferLog := logging.WhenRequest(ctx, "MemoryOrderRepository.Create", map[string]any{"order": order})
	defer deferLog(created, &err)
	m.lock.Lock()
	defer m.lock.Unlock()
	newOrder := &domain.Order{
		ID:          strconv.FormatInt(time.Now().Unix(), 10),
		CustomerID:  order.CustomerID,
		Status:      order.Status,
		PaymentLink: order.PaymentLink,
		Items:       order.Items,
	}
	m.store = append(m.store, newOrder)
	return newOrder, nil
}

func (m *MemoryOrderRepository) Get(ctx context.Context, id, customerID string) (got *domain.Order, err error) {
	_, deferLog := logging.WhenRequest(ctx, "MemoryOrderRepository.Get", map[string]any{
		"id":          id,
		"customer_id": customerID,
	})
	defer deferLog(got, &err)
	m.lock.RLock()
	defer m.lock.RUnlock()
	for _, order := range m.store {
		if order.CustomerID == customerID && order.ID == id {
			return order, nil
		}
	}
	return nil, domain.NotFoundError{OrderID: id}
}

func (m *MemoryOrderRepository) Update(
	ctx context.Context,
	o *domain.Order,
	updateFn func(context.Context, *domain.Order) (*domain.Order, error),
) (err error) {
	_, deferLog := logging.WhenRequest(ctx, "MemoryOrderRepository.Update", map[string]any{"order": o})
	defer deferLog(nil, &err)

	m.lock.Lock()
	defer m.lock.Unlock()
	found := false
	defer func() {
		logrus.Infof("memory_order_repo_update||found=%t||orderID=%s||customerID=%s", found, o.ID, o.CustomerID)
	}()
	for i, order := range m.store {
		if order.ID == o.ID && order.CustomerID == o.CustomerID {
			found = true
			updateOrder, err := updateFn(ctx, o)
			if err != nil {
				return err
			}
			m.store[i] = updateOrder
		}
	}
	if !found {
		return domain.NotFoundError{OrderID: o.ID}
	}
	return nil
}
