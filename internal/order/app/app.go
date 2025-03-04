package app

import (
	"github.com/mushroomyuan/gorder/order/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct{}

type Queries struct {
	GetCustomerOrder query.GetCustomerOrderHandler
}
