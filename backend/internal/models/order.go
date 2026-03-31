package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type OrderStatus string

const (
	StatusPending   OrderStatus = "pending"
	StatusCompleted OrderStatus = "completed"
	StatusCancelled OrderStatus = "cancelled"
	StatusRefunded  OrderStatus = "refunded"
)

type Order struct {
	ID            uuid.UUID       `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrderNumber   string          `gorm:"uniqueIndex;size:50;not null" json:"order_number"`
	CustomerID    *uuid.UUID      `gorm:"type:uuid" json:"customer_id,omitempty"`
	Customer      *Customer       `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	UserID        uuid.UUID       `gorm:"type:uuid;not null" json:"user_id"`
	User          User            `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Status        OrderStatus     `gorm:"type:varchar(20);default:'pending'" json:"status"`
	Subtotal      decimal.Decimal `gorm:"type:decimal(10,2);not null" json:"subtotal"`
	Tax           decimal.Decimal `gorm:"type:decimal(10,2);default:0" json:"tax"`
	Discount      decimal.Decimal `gorm:"type:decimal(10,2);default:0" json:"discount"`
	Total         decimal.Decimal `gorm:"type:decimal(10,2);not null" json:"total"`
	PaymentMethod string          `gorm:"size:50" json:"payment_method"`
	Notes         string          `gorm:"type:text" json:"notes"`
	Items         []OrderItem     `gorm:"foreignKey:OrderID" json:"items,omitempty"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

func (o *Order) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}

type OrderItem struct {
	ID        uuid.UUID       `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrderID   uuid.UUID       `gorm:"type:uuid;not null" json:"order_id"`
	ProductID uuid.UUID       `gorm:"type:uuid;not null" json:"product_id"`
	Product   Product         `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Quantity  int             `gorm:"not null" json:"quantity"`
	UnitPrice decimal.Decimal `gorm:"type:decimal(10,2);not null" json:"unit_price"`
	Discount  decimal.Decimal `gorm:"type:decimal(10,2);default:0" json:"discount"`
	Total     decimal.Decimal `gorm:"type:decimal(10,2);not null" json:"total"`
	CreatedAt time.Time       `json:"created_at"`
}

func (oi *OrderItem) BeforeCreate(tx *gorm.DB) error {
	if oi.ID == uuid.Nil {
		oi.ID = uuid.New()
	}
	return nil
}

func (oi *OrderItem) CalculateTotal() decimal.Decimal {
	return oi.UnitPrice.Mul(decimal.NewFromInt(int64(oi.Quantity))).Sub(oi.Discount)
}
