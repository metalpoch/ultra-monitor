package entity

import (
	"time"
)

type User struct {
	ID             int32     `db:"id"`
	ChangePassword bool      `db:"change_password"`
	IsAdmin        bool      `db:"is_admin"`
	IsDisabled     bool      `db:"is_disabled"`
	Fullname       string    `db:"fullname"`
	Username       string    `db:"username"`
	Password       string    `db:"password"`
	CreatedAt      time.Time `db:"created_at"`
}
