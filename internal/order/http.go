package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mushroomyuan/gorder/common"
	client "github.com/mushroomyuan/gorder/common/client/order"
	"github.com/mushroomyuan/gorder/order/app"
	"github.com/mushroomyuan/gorder/order/app/command"
	"github.com/mushroomyuan/gorder/order/app/dto"
	"github.com/mushroomyuan/gorder/order/app/query"
	"github.com/mushroomyuan/gorder/order/convertor"
)

type HTTPServer struct {
	common.BaseResponse
	app app.Application
}

func (H HTTPServer) PostCustomerCustomerIdOrders(c *gin.Context, customerID string) {
	var (
		req  client.CreateOrderRequest
		resp dto.CreateOrderResponse
		err  error
	)
	defer func() {
		H.Response(c, err, &resp)
	}()

	if err = c.ShouldBindJSON(&req); err != nil {
		return
	}
	r, err := H.app.Commands.CreateOrder.Handle(c.Request.Context(), command.CreateOrder{
		CustomerID: req.CustomerId,
		Items:      convertor.NewItemWithQuantityConvertor().ClientsToEntities(req.Items),
	})
	if err != nil {
		return
	}
	resp = dto.CreateOrderResponse{
		OrderID:     r.OrderID,
		CustomerID:  req.CustomerId,
		RedirectURL: fmt.Sprintf("http://localhost:8282/success?customerID=%s&orderID=%s", req.CustomerId, r.OrderID),
	}
}

func (s *HTTPServer) GetCustomerCustomerIdOrdersOrderId(c *gin.Context, customerId string, orderId string) {
	// ctx, span := tracing.Start(c, "GetCustomerCustomerIDOrdersOrderID")
	// defer span.End()

	var (
		err  error
		resp struct {
			Order *client.Order
		}
	)

	defer func() {
		s.Response(c, err, resp)
	}()

	o, err := s.app.Queries.GetCustomerOrder.Handle(c.Request.Context(), query.GetCustomerOrder{
		OrderID:    orderId,
		CustomerID: customerId,
	})
	if err != nil {
		return
	}
	resp.Order = convertor.NewOrderConvertor().EntityToClient(o)
}
