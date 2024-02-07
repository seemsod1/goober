package models

import "time"

type CarModel struct {
	ID        int      `gorm:"primaryKey"`
	Name      string   `gorm:"size:255;not null"`
	BrandId   int      `gorm:"not null"`
	Brand     CarBrand `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
