package main

import "github.com/gin-gonic/gin"

type PaymentHandler struct {
}

func NewPaymentHandler() *PaymentHandler {
	return &PaymentHandler{}
}

func (h *PaymentHandler) RegisterRoutes(c *gin.Engine) {
	c.POST("/api/webhook",h.HandleWebhook)
}

func (h *PaymentHandler) HandleWebhook(c *gin.Context) {
	//TODO: implement me
	panic("implement me")
}
