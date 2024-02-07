package models

import "time"

type UserHistory struct {
	ID         int      `gorm:"primaryKey"`
	UserID     int      `gorm:"not null"`
	User       User     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	RentInfoID int      `gorm:"not null"`
	RentInfo   RentInfo `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
