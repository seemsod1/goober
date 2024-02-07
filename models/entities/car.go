package models

import (
	"time"
)

type Car struct {
	ID         int          `gorm:"primaryKey"`
	TypeID     int          `gorm:"not null"`
	Type       CarType      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ModelId    int          `gorm:"not null"`
	Model      CarModel     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Year       int          `gorm:"not null"`
	Plate      string       `gorm:"size:255;unique;index;not null"`
	Price      float64      `gorm:"not null"`
	LocationId int          `gorm:"not null"`
	Location   RentLocation `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Available  bool         `gorm:"default:true"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
