package request

import "github.com/shopspring/decimal"

type CreateOrderItemRequest struct {
	ProductID string          `json:"product_id" binding:"required"`
	Quantity  int             `json:"quantity" binding:"required,min=1"`
	Discount  decimal.Decimal `json:"discount"`
}

type CreateOrderRequest struct {
	CustomerID    *string                  `json:"customer_id"`
	PaymentMethod string                   `json:"payment_method" binding:"required"`
	Notes         string                   `json:"notes"`
	Items         []CreateOrderItemRequest `json:"items" binding:"required,min=1,dive"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending completed cancelled refunded"`
}

type ListOrdersRequest struct {
	Page       int    `form:"page" binding:"min=0"`
	Limit      int    `form:"limit" binding:"min=0,max=100"`
	CustomerID string `form:"customer_id"`
	Status     string `form:"status"`
	StartDate  string `form:"start_date"`
	EndDate    string `form:"end_date"`
}
