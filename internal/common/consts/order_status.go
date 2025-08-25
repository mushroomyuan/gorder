package consts

type OrderStatus string

const (
	OrderStatusPending           string = "pending"
	OrderStatusWaitingForPayment string = "waiting for payment"
	OrderStatusPaid              string = "paid"
	OrderStatusReady             string = "ready"
)
