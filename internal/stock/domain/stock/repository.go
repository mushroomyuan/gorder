package stock

import (
	"context"
	"fmt"
	"strings"

	"github.com/mushroomyuan/gorder/common/entity"
)

type Repository interface {
	GetItems(ctx context.Context, ids []string) ([]*entity.Item, error)
	GetStock(ctx context.Context, ids []string) ([]*entity.ItemWithQuantity, error)
	UpdateStock(
		ctx context.Context,
		data []*entity.ItemWithQuantity,
		updateFn func(
			ctx context.Context,
			existing []*entity.ItemWithQuantity,
			query []*entity.ItemWithQuantity,
		) ([]*entity.ItemWithQuantity, error),
	) error
}

type NotFoundError struct {
	Missing []string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("these items not found in stock: %s", strings.Join(e.Missing, ","))
}

type ExceedStockError struct {
	FailedOn []struct {
		ID   string
		Want int32
		Have int32
	}
}

func (e ExceedStockError) Error() string {
	var info []string
	for _, item := range e.FailedOn {
		info = append(info, fmt.Sprintf("product_id=%s,want=%d,have=%d", item.ID, item.Want, item.Have))
	}
	return fmt.Sprintf("these items is not sufficient for product: %s", strings.Join(info, ","))
}
