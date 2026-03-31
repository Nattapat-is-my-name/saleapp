package handler

import (
	"io"
	"log"
	"net/http"
	"os"

	"saleapp/internal/service"
	"saleapp/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v76/webhook"
)

type PaymentHandler struct {
	paymentService *service.PaymentService
}

func NewPaymentHandler(paymentService *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService}
}

func (h *PaymentHandler) CreatePaymentIntent(c *gin.Context) {
	var req struct {
		OrderID  string `json:"order_id" binding:"required"`
		Currency string `json:"currency"`
	}
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

// HandleWebhook handles incoming Stripe webhook events.
// It verifies the signature, checks for duplicates, and returns 200 immediately
// before processing the event asynchronously.
func (h *PaymentHandler) HandleWebhook(c *gin.Context) {
	// 1. Read raw body — required for Stripe signature verification
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Failed to read webhook body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Can't read body"})
		return
	}

	// 2. Verify Stripe signature
	sigHeader := c.GetHeader("Stripe-Signature")
	webhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	if webhookSecret == "" {
		log.Println("WARNING: STRIPE_WEBHOOK_SECRET not set, using test mode")
		webhookSecret = "whsec_test_placeholder"
	}

	event, err := webhook.ConstructEvent(payload, sigHeader, webhookSecret)
	if err != nil {
		// Signature verification failed — log and return 400
		// Stripe will NOT retry events with invalid signatures
		log.Printf("Webhook signature verification failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid signature"})
		return
	}

	// 3. Idempotency check — skip already-processed events
	processed, err := h.paymentService.IsEventProcessed(event.ID)
	if err != nil {
		// DB error — log but still return 200 to avoid Stripe retry storms
		// The event may or may not have been processed; we'll find out on retry
		log.Printf("DB error checking event %s: %v", event.ID, err)
	}
	if processed {
		log.Printf("Event %s already processed, skipping", event.ID)
		c.JSON(http.StatusOK, gin.H{"received": true})
		return
	}

	// 4. Return 200 immediately, process async
	// Stripe expects a 200 response within ~30 seconds
	go h.paymentService.HandleWebhookAsync(event)

	c.JSON(http.StatusOK, gin.H{"received": true})
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
