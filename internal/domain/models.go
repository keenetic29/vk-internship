package domain

import (
	"time"
)

type User struct {
	ID       	uint   	`gorm:"primaryKey"`
	Username 	string 	`gorm:"unique;not null"`
	Password 	string 	`gorm:"not null"`
	CreatedAt 	time.Time
}

type Advertisement struct {
	ID          uint   	`gorm:"primaryKey"`
	Title       string 	`gorm:"not null;size:100"`
	Description string 	`gorm:"not null;size:1000"`
	ImageURL    string 	`gorm:"not null"`
	Price       float64 `gorm:"not null"`
	UserID      uint    `gorm:"not null"`
	User        User    `gorm:"foreignKey:UserID"`
	IsOwner     bool    `gorm:"-" json:"is_owner"`
	CreatedAt   time.Time
}

type AuthToken struct {
	UserID 		uint   	`gorm:"not null"`
	Token  		string 	`gorm:"primaryKey"`
	ExpiresAt 	time.Time
}