package adapters

import (
	"context"
	"strconv"
	"sync"
	"time"

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

func (m *MemoryOrderRepository) Create(_ context.Context, order *domain.Order) (*domain.Order, error) {
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
	logrus.WithFields(logrus.Fields{
		"input_order":        order,
		"store_after_create": m.store,
	}).Debug("memory order created")
	return newOrder, nil
}

func (m *MemoryOrderRepository) Get(_ context.Context, id, customerID string) (*domain.Order, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	for _, order := range m.store {
		if order.CustomerID == customerID && order.ID == id {
			logrus.Debugf("memory_order_repo_get||found||id=%s||CustomerID=%s||res=%+v", order.ID, order.CustomerID, *order)
			return order, nil
		}
	}
	return nil, domain.NotFoundError{OrderID: id}
}

func (m *MemoryOrderRepository) Update(ctx context.Context, o *domain.Order, updateFn func(context.Context, *domain.Order) (*domain.Order, error)) error {
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
