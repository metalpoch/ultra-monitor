package entity

import "time"

type User struct {
	ID             uint      `gorm:"id"`
	Fullname       string    `gorm:"fullname"`
	Email          string    `gorm:"email"`
	Password       string    `gorm:"password"`
	ChangePassword bool      `gorm:"change_password"`
	IsAdmin        bool      `gorm:"is_admin"`
	IsDisabled     bool      `gorm:"is_disabled"`
	CreatedAt      time.Time `gorm:"created_at"`
	UpdatedAt      time.Time `gorm:"updated_at"`
}

type UserResponse struct {
	Fullname string `bson:"fullname"`
	P00      uint   `bson:"p00"`
	Email    string `bson:"email"`
	Password string `bson:"password"`
}
