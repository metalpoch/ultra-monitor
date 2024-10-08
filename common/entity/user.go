package entity

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID             uint `gorm:"primaryKey"`
	Fullname       string
	Email          string
	Password       string
	ChangePassword bool
	IsAdmin        bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeleteAt       gorm.DeletedAt
}

type UserResponse struct {
	Fullname string `bson:"fullname"`
	P00      uint   `bson:"p00"`
	Email    string `bson:"email"`
	Password string `bson:"password"`
}
