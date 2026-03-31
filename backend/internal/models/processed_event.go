package models

import (
	"time"
)

type ProcessedEvent struct {
	EventID     string    `gorm:"primaryKey;size:255"`
	ProcessedAt time.Time `gorm:"autoCreateTime"`
}
