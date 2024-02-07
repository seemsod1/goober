package models

import "time"

type CarType struct {
	ID        int    `gorm:"primaryKey"`
	Name      string `gorm:"size:255;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
