package models

import "time"

type CarAssignment struct {
	CarId     int        `gorm:"not null"`
	Car       Car        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	PurposeId int        `gorm:"not null"`
	Purpose   CarPurpose `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
