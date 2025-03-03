// Package ports provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package ports

// CreateOrderRequest defines model for CreateOrderRequest.
type CreateOrderRequest struct {
	CustomerID string             `json:"customerID"`
	Items      []ItemWithQuantity `json:"items"`
}

// Error defines model for Error.
type Error struct {
	Message *string `json:"message,omitempty"`
}

// Item defines model for Item.
type Item struct {
	Id       *string `json:"id,omitempty"`
	Name     *string `json:"name,omitempty"`
	PriceID  *string `json:"priceID,omitempty"`
	Quantity *int32  `json:"quantity,omitempty"`
}

// ItemWithQuantity defines model for ItemWithQuantity.
type ItemWithQuantity struct {
	Id       *string `json:"id,omitempty"`
	Quantity *int32  `json:"quantity,omitempty"`
}

// Order defines model for Order.
type Order struct {
	CustomerID   *string `json:"customerID,omitempty"`
	Id           *string `json:"id,omitempty"`
	Items        *[]Item `json:"items,omitempty"`
	PaymmentLink *string `json:"paymmentLink,omitempty"`
}

// PostCustomerCustomerIDOrdersJSONRequestBody defines body for PostCustomerCustomerIDOrders for application/json ContentType.
type PostCustomerCustomerIDOrdersJSONRequestBody = CreateOrderRequest
