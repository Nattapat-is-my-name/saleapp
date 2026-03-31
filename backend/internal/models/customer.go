package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Customer struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email     string    `gorm:"uniqueIndex;size:255" json:"email"`
	Phone     string    `gorm:"size:20" json:"phone"`
	FirstName string    `gorm:"size:100" json:"first_name"`
	LastName  string    `gorm:"size:100" json:"last_name"`
	Address   string    `gorm:"type:text" json:"address"`
	Notes     string    `gorm:"type:text" json:"notes"`
	Orders    []Order   `gorm:"foreignKey:CustomerID" json:"orders,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (c *Customer) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

func (c *Customer) FullName() string {
	return c.FirstName + " " + c.LastName
}
