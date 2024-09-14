package entity

import "time"

type User struct {
	ID             uint      `sql:"id"`
	Fullname       string    `sql:"fullname"`
	Email          string    `sql:"email"`
	Password       string    `sql:"password"`
	ChangePassword bool      `sql:"change_password"`
	IsAdmin        bool      `sql:"is_admin"`
	IsDisabled     bool      `sql:"is_disabled"`
	CreatedAt      time.Time `sql:"created_at"`
	UpdatedAt      time.Time `sql:"updated_at"`
}

type UserResponse struct {
	Fullname string `bson:"fullname"`
	P00      uint   `bson:"p00"`
	Email    string `bson:"email"`
	Password string `bson:"password"`
}
