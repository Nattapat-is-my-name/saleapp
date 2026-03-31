package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type PaymentStatus string

const (
	PaymentPending   PaymentStatus = "pending"
	PaymentSucceeded PaymentStatus = "succeeded"
	PaymentFailed    PaymentStatus = "failed"
	PaymentRefunded  PaymentStatus = "refunded"
)

type Payment struct {
	ID              uuid.UUID       `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrderID         uuid.UUID       `gorm:"type:uuid;not null;uniqueIndex" json:"order_id"`
	Order           Order           `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	StripePaymentID string          `gorm:"size:100;uniqueIndex" json:"stripe_payment_id"`
	Amount          decimal.Decimal `gorm:"type:decimal(10,2);not null" json:"amount"`
	Currency        string          `gorm:"size:3;default:'usd'" json:"currency"`
	Status          PaymentStatus   `gorm:"type:varchar(20);default:'pending'" json:"status"`
	PaymentMethod   string          `gorm:"size:50" json:"payment_method"`
	ErrorMessage    string          `gorm:"type:text" json:"error_message,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

func (p *Payment) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}
