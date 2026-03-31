package service

import (
	"fmt"
	"log"
	"os"

	"saleapp/internal/models"
	"saleapp/internal/repository"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	stripe "github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/paymentintent"
)

type PaymentService struct {
	paymentRepo repository.PaymentRepository
	orderRepo   repository.OrderRepository
}

func NewPaymentService(paymentRepo repository.PaymentRepository, orderRepo repository.OrderRepository) *PaymentService {
	return &PaymentService{
		paymentRepo: paymentRepo,
		orderRepo:   orderRepo,
	}
}

func (s *PaymentService) InitStripe() {
	stripeKey := os.Getenv("STRIPE_SECRET_KEY")
	if stripeKey == "" {
		stripeKey = "sk_test_placeholder"
	}
	stripe.Key = stripeKey
}

func (s *PaymentService) CreatePaymentIntent(orderID string, currency string) (*stripe.PaymentIntent, *models.Payment, error) {
	if currency == "" {
		currency = "usd"
	}

	orderUUID, err := uuid.Parse(orderID)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid order ID: %w", err)
	}

	order, err := s.orderRepo.GetByID(orderUUID)
	if err != nil {
		return nil, nil, fmt.Errorf("order not found: %w", err)
	}
	if order == nil {
		return nil, nil, fmt.Errorf("order not found")
	}

	amountCents := int64(order.Total.Mul(decimal.NewFromInt(100)).Round(0).IntPart())

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amountCents),
		Currency: stripe.String(currency),
		Metadata: map[string]string{
			"order_id": orderID,
		},
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create payment intent: %w", err)
	}

	payment := &models.Payment{
		OrderID:         order.ID,
		StripePaymentID: pi.ID,
		Amount:          order.Total,
		Currency:        currency,
		Status:          models.PaymentPending,
	}

	if err := s.paymentRepo.Create(payment); err != nil {
		return nil, nil, fmt.Errorf("failed to save payment: %w", err)
	}

	return pi, payment, nil
}

// ProcessWebhookEvent handles a verified Stripe event synchronously.
// This is called after signature verification and idempotency check.
func (s *PaymentService) ProcessWebhookEvent(event *stripe.Event) error {
	switch event.Type {
	case "payment_intent.succeeded":
		pi, ok := event.Data.Object["id"].(string)
		if !ok {
			return fmt.Errorf("failed to extract payment intent ID from event")
		}
		return s.paymentRepo.UpdateStatusByStripeID(pi, models.PaymentSucceeded)
	case "payment_intent.payment_failed":
		pi, ok := event.Data.Object["id"].(string)
		if !ok {
			return fmt.Errorf("failed to extract payment intent ID from event")
		}
		return s.paymentRepo.UpdateStatusByStripeID(pi, models.PaymentFailed)
	case "charge.refunded":
		pi, ok := event.Data.Object["payment_intent"].(string)
		if !ok {
			return fmt.Errorf("failed to extract payment intent ID from charge.refunded event")
		}
		return s.paymentRepo.UpdateStatusByStripeID(pi, models.PaymentRefunded)
	}
	return nil
}

// HandleWebhookAsync processes a Stripe webhook event asynchronously.
// It marks the event as processed first (idempotency), then processes it.
// Use this for returning 200 fast to Stripe while processing in background.
func (s *PaymentService) HandleWebhookAsync(event stripe.Event) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered in webhook processing: %v", r)
		}
	}()

	// Mark event processed FIRST to handle retries idempotently
	if err := s.paymentRepo.MarkEventProcessed(event.ID); err != nil {
		log.Printf("Failed to mark event %s as processed: %v", event.ID, err)
	}

	// Process the event
	if err := s.ProcessWebhookEvent(&event); err != nil {
		log.Printf("Failed to process webhook event %s: %v", event.ID, err)
	}
}

// IsEventProcessed checks if a Stripe event has already been processed.
func (s *PaymentService) IsEventProcessed(eventID string) (bool, error) {
	return s.paymentRepo.IsEventProcessed(eventID)
}

func (s *PaymentService) GetPaymentByOrderID(orderID string) (*models.Payment, error) {
	return s.paymentRepo.GetByOrderID(orderID)
}

func (s *PaymentService) GetPaymentByStripeID(stripeID string) (*models.Payment, error) {
	return s.paymentRepo.GetByStripeID(stripeID)
}
