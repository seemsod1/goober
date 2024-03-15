package entities

import "time"

type VerificationData struct {
	ID        int    `gorm:"primaryKey"`
	Email     string `gorm:"size:255;unique;not null"`
	Code      string `gorm:"size:255;not null"`
	ExpiresAt time.Time
}
