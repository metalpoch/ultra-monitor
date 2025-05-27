package entity

import (
	"time"
)

type User struct {
	ID             uint32    `db:"id"`
	Fullname       string    `db:"fullname"`
	Email          string    `db:"email"`
	Password       string    `db:"password"`
	ChangePassword bool      `db:"change_password"`
	IsAdmin        bool      `db:"is_admin"`
	CreatedAt      time.Time `db:"created_at"`
}

type UserResponse struct {
	P00      uint32 `db:"p00"`
	Fullname string `db:"fullname"`
	Email    string `db:"email"`
	Password string `db:"created_at"`
}
