package entities

import "time"

type User struct {
	ID         int       `gorm:"primaryKey"`
	Name       string    `gorm:"size:255;not null"`
	BirthDate  time.Time `gorm:"not null"`
	Email      string    `gorm:"size:255;unique;index;not null"`
	Password   string    `gorm:"size:255;not null"`
	Phone      string    `gorm:"size:255;unique;index;not null"`
	IsVerified bool      `gorm:"default:false"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	RoleId     int      `gorm:"not null"`
	Role       UserRole `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type UserRole struct {
	ID        int    `gorm:"primaryKey"`
	Name      string `gorm:"size:255;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
