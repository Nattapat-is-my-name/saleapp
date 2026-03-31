package service

import (
	"fmt"
	"os"

	"saleapp/internal/models"
	"saleapp/internal/repository"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stripe/stripe-go/v76"
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

func (s *PaymentService) HandleWebhook(event *stripe.Event) error {
	switch event.Type {
	case "payment_intent.succeeded":
		pi := event.Data.Object["id"].(string)
		return s.paymentRepo.UpdateStatusByStripeID(pi, models.PaymentSucceeded)
	case "payment_intent.payment_failed":
		pi := event.Data.Object["id"].(string)
		return s.paymentRepo.UpdateStatusByStripeID(pi, models.PaymentFailed)
	case "charge.refunded":
		pi := event.Data.Object["payment_intent"].(string)
		return s.paymentRepo.UpdateStatusByStripeID(pi, models.PaymentRefunded)
	}
	return nil
}

func (s *PaymentService) GetPaymentByOrderID(orderID string) (*models.Payment, error) {
	return s.paymentRepo.GetByOrderID(orderID)
}

func (s *PaymentService) GetPaymentByStripeID(stripeID string) (*models.Payment, error) {
	return s.paymentRepo.GetByStripeID(stripeID)
}
