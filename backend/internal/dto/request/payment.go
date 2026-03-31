package request

type CreatePaymentIntentRequest struct {
	OrderID string `json:"order_id" binding:"required,uuid"`
	Currency string `json:"currency" binding:"omitempty"`
}
