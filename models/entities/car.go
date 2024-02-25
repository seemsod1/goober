package entities

import (
	"time"
)

type Car struct {
	ID         int          `gorm:"primaryKey"`
	TypeId     int          `gorm:"not null"`
	Type       CarType      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ModelId    int          `gorm:"not null"`
	Model      CarModel     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Bags       int          `gorm:"not null"`
	Passengers int          `gorm:"not null"`
	Year       int          `gorm:"not null"`
	Plate      string       `gorm:"size:255;unique;index;not null"`
	Price      float64      `gorm:"not null"`
	Color      string       `gorm:"size:255;not null"`
	LocationId int          `gorm:"not null"`
	Location   RentLocation `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Available  bool         `gorm:"default:true"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
