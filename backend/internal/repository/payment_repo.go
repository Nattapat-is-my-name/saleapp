package repository

import (
	"saleapp/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentRepository interface {
	Create(payment *models.Payment) error
	GetByOrderID(orderID string) (*models.Payment, error)
	GetByStripeID(stripeID string) (*models.Payment, error)
	UpdateStatusByStripeID(stripeID string, status models.PaymentStatus) error
	UpdateStatus(orderID string, status models.PaymentStatus) error
	IsEventProcessed(eventID string) (bool, error)
	MarkEventProcessed(eventID string) error
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(payment *models.Payment) error {
	return r.db.Create(payment).Error
}

func (r *paymentRepository) GetByOrderID(orderID string) (*models.Payment, error) {
	var payment models.Payment
	orderUUID, err := uuid.Parse(orderID)
	if err != nil {
		return nil, err
	}
	err = r.db.Where("order_id = ?", orderUUID).First(&payment).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &payment, err
}

func (r *paymentRepository) GetByStripeID(stripeID string) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.Where("stripe_payment_id = ?", stripeID).First(&payment).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &payment, err
}

func (r *paymentRepository) UpdateStatusByStripeID(stripeID string, status models.PaymentStatus) error {
	return r.db.Model(&models.Payment{}).
		Where("stripe_payment_id = ?", stripeID).
		Update("status", status).Error
}

func (r *paymentRepository) UpdateStatus(orderID string, status models.PaymentStatus) error {
	orderUUID, err := uuid.Parse(orderID)
	if err != nil {
		return err
	}
	return r.db.Model(&models.Payment{}).
		Where("order_id = ?", orderUUID).
		Update("status", status).Error
}

func (r *paymentRepository) IsEventProcessed(eventID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.ProcessedEvent{}).Where("event_id = ?", eventID).Count(&count).Error
	return count > 0, err
}

func (r *paymentRepository) MarkEventProcessed(eventID string) error {
	return r.db.Create(&models.ProcessedEvent{EventID: eventID}).Error
}
