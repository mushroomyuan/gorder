package adapters

import (
	"context"
	"github.com/mushroomyuan/gorder/common/genproto/orderpb"
	domain "github.com/mushroomyuan/gorder/stock/domain/stock"
	"sync"
)

type MemoryStockRepository struct {
	lock  *sync.RWMutex
	store map[string]*orderpb.Item
}

var stub = map[string]*orderpb.Item{
	"item_id": {
		ID:       "foo_item",
		Name:     "stub_item",
		Quantity: 10000,
		PriceID:  "stub_item_price_id",
	},
}

func NewMemoryStockRepository() *MemoryStockRepository {
	return &MemoryStockRepository{
		lock:  &sync.RWMutex{},
		store: stub,
	}
}

func (m MemoryStockRepository) GetItems(ctx context.Context, ids []string) ([]*orderpb.Item, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	var (
		res     = make([]*orderpb.Item, len(ids))
		missing []string
	)
	for _, id := range ids {
		if item, exit := m.store[id]; exit {
			res = append(res, item)
		} else {
			missing = append(missing, id)
		}

	}
	if len(res) == len(ids) {
		return res, nil
	}
	return res, domain.NotFoundError{Missing: missing}
}
