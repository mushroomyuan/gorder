package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mushroomyuan/gorder/order/app"
	"github.com/mushroomyuan/gorder/order/app/query"
	"net/http"
)

type HTTPServer struct {
	app app.Application
}

func (s *HTTPServer) PostCustomerCustomerIDOrders(c *gin.Context, customerID string) {

}

func (s *HTTPServer) GetCustomerCustomerIDOrdersOrderID(c *gin.Context, customerID string, orderID string) {
	o, err := s.app.Queries.GetCustomerOrder.Handle(c, query.GetCustomerOrder{
		OrderID:    "fake-ID",
		CustomerID: "fake-CustomerID",
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, o)
}
