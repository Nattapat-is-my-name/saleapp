package handler

import (
	"net/http"

	"saleapp/internal/dto/request"
	"saleapp/internal/service"
	"saleapp/pkg/response"

	"github.com/gin-gonic/gin"
	stripe "github.com/stripe/stripe-go/v76"
)

type PaymentHandler struct {
	paymentService *service.PaymentService
}

func NewPaymentHandler(paymentService *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService}
}

func (h *PaymentHandler) CreatePaymentIntent(c *gin.Context) {
	var req request.CreatePaymentIntentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request: "+err.Error())
		return
	}

	pi, payment, err := h.paymentService.CreatePaymentIntent(req.OrderID, req.Currency)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "PAYMENT_ERROR", "Failed to create payment: "+err.Error())
		return
	}

	response.Created(c, gin.H{
		"client_secret": pi.ClientSecret,
		"payment_id":    payment.ID,
		"amount":        pi.Amount,
		"currency":      pi.Currency,
	})
}

func (h *PaymentHandler) HandleWebhook(c *gin.Context) {
	var event stripe.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		response.Error(c, http.StatusBadRequest, "WEBHOOK_ERROR", "Webhook error: "+err.Error())
		return
	}

	if err := h.paymentService.HandleWebhook(&event); err != nil {
		response.Error(c, http.StatusInternalServerError, "WEBHOOK_PROCESSING_ERROR", "Webhook processing error")
		return
	}

	response.Success(c, gin.H{"received": true})
}

func (h *PaymentHandler) GetPaymentStatus(c *gin.Context) {
	orderID := c.Param("orderId")

	payment, err := h.paymentService.GetPaymentByOrderID(orderID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get payment")
		return
	}
	if payment == nil {
		response.Error(c, http.StatusNotFound, "NOT_FOUND", "Payment not found")
		return
	}

	response.Success(c, payment)
}
