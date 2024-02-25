package entities

import "time"

type CarBrand struct {
	ID        int    `gorm:"primaryKey"`
	Name      string `gorm:"size:255;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
