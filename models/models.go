package models

//user

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	gorm.Model
	Name       string     `json:"name" gorm:"type:varchar(100);not null"`
	Email      string     `json:"email" gorm:"uniqueIndex;type:varchar(100);not null"`
	// Password   string     `json:"password" gorm:"type:varchar(255);not null"`
	RollNumber string        `json:"roll" gorm:"not null"`
	IsVerified bool       `json:"is_verified" gorm:"default:false"`
	
}

//booking
type Booking struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Phone         string    `json:"phone"`
	UTR           string    `json:"utr"`
	Kannadigas    int       `json:"kannadigas"`
	NonKannadigas int       `json:"nonKannadigas"`
	Total         int       `json:"total"`
	CreatedAt     time.Time `json:"createdAt"`


    Status      string    `json:"status" gorm:"default:not_entered"` 
	StatusCount int  		`json:"statuscount" gorm:"default:0"`
    IsVerified  bool      `json:"is_verified" gorm:"default:false"`    
}