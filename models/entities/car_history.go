package models

import "time"

type CarHistory struct {
	ID         int      `gorm:"primaryKey"`
	CarID      int      `gorm:"not null"`
	Car        Car      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	RentInfoID int      `gorm:"not null"`
	RentInfo   RentInfo `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
