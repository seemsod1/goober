package entities

import "time"

type RentLocation struct {
	ID          int    `gorm:"primaryKey"`
	FullAddress string `gorm:"size:255;not null"`
	UserId      int    `gorm:"not null"`
	User        User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CityId      int    `gorm:"not null"`
	City        City   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
