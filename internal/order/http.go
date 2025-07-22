package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mushroomyuan/gorder/common/genproto/orderpb"
	"github.com/mushroomyuan/gorder/order/app"
	"github.com/mushroomyuan/gorder/order/app/command"
	"github.com/mushroomyuan/gorder/order/app/query"
	"net/http"
)

type HTTPServer struct {
	app app.Application
}

func (s *HTTPServer) PostCustomerCustomerIDOrders(c *gin.Context, customerID string) {
	var req orderpb.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	r, err := s.app.Commands.CreateOrder.Handle(c, command.CreateOrder{
		CustomerID: customerID,
		Items:      req.Items,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":     "success",
		"customer_id": customerID,
		"order_id":    r.OrderID,
	})
}

func (s *HTTPServer) GetCustomerCustomerIDOrdersOrderID(c *gin.Context, customerID string, orderID string) {
	o, err := s.app.Queries.GetCustomerOrder.Handle(c, query.GetCustomerOrder{
		OrderID:    orderID,
		CustomerID: customerID,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, o)
}
