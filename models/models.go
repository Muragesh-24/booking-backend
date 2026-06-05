package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	gorm.Model

	Name       string `json:"name" gorm:"type:varchar(100);not null"`
	Email      string `json:"email" gorm:"type:varchar(100);uniqueIndex;not null;index"`
	RollNumber string `json:"roll" gorm:"type:varchar(50);not null;index"`
	IsVerified bool   `json:"is_verified" gorm:"default:false;index"`
}

// Booking represents a booking/order in the system
type Booking struct {
	ID uint `gorm:"primaryKey" json:"id"`

	Name  string `json:"name" gorm:"type:varchar(100);not null;index"`
	Email string `json:"email" gorm:"type:varchar(100);not null;index"`
	Phone string `json:"phone" gorm:"type:varchar(20);not null;index"`

	UTR string `json:"utr" gorm:"type:varchar(100);uniqueIndex;not null"`

	Kannadigas    int `json:"kannadigas" gorm:"default:0"`
	NonKannadigas int `json:"nonKannadigas" gorm:"default:0"`
	Total         int `json:"total" gorm:"default:0"`

	Status      string `json:"status" gorm:"type:varchar(30);default:'not_entered';index"`
	StatusCount int    `json:"statuscount" gorm:"default:0"`
	IsVerified  bool   `json:"is_verified" gorm:"default:false;index"`

	CreatedAt time.Time `json:"createdAt" gorm:"index"`
	UpdatedAt time.Time `json:"updatedAt"`
}