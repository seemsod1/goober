package entities

import "time"

type UserHistory struct {
	UserID     int      `gorm:"primaryKey; autoIncrement:false;not null"`
	User       User     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	RentInfoID int      `gorm:"primaryKey; autoIncrement:false;not null"`
	RentInfo   RentInfo `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
