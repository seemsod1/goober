package entities

import "time"

type CarHistory struct {
	CarId      int      `gorm:"primaryKey; autoIncrement:false;not null"`
	Car        Car      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	RentInfoId int      `gorm:"primaryKey; autoIncrement:false;not null"`
	RentInfo   RentInfo `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
