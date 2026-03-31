package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Category struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name      string    `gorm:"size:100;not null" json:"name"`
	Slug      string    `gorm:"size:100;uniqueIndex" json:"slug"`
	ParentID  *uuid.UUID `gorm:"type:uuid" json:"parent_id,omitempty"`
	Parent    *Category `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Products  []Product `gorm:"foreignKey:CategoryID" json:"products,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

type Product struct {
	ID          uuid.UUID       `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	SKU         string          `gorm:"uniqueIndex;size:50;not null" json:"sku"`
	Name        string          `gorm:"size:255;not null" json:"name"`
	Description string          `gorm:"type:text" json:"description"`
	Price       decimal.Decimal `gorm:"type:decimal(10,2);not null" json:"price"`
	Cost        decimal.Decimal `gorm:"type:decimal(10,2)" json:"cost"`
	Stock       int             `gorm:"default:0" json:"stock"`
	CategoryID *uuid.UUID      `gorm:"type:uuid" json:"category_id,omitempty"`
	Category   *Category       `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	IsActive    bool            `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

func (p *Product) Margin() decimal.Decimal {
	if p.Cost.IsZero() {
		return decimal.Zero
	}
	return p.Price.Sub(p.Cost).Div(p.Cost).Mul(decimal.NewFromInt(100))
}

func (p *Product) IsInStock() bool {
	return p.Stock > 0
}

func (p *Product) HasLowStock(threshold int) bool {
	return p.Stock > 0 && p.Stock <= threshold
}
