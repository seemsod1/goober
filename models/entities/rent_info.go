package entities

import "time"

type RentInfo struct {
	ID        int          `gorm:"primaryKey"`
	StartDate time.Time    `gorm:"not null"`
	EndDate   time.Time    `gorm:"not null"`
	Price     float64      `gorm:"not null"`
	FromId    int          `gorm:"not null"`
	From      RentLocation `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ReturnId  int          `gorm:"not null"`
	Return    RentLocation `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
