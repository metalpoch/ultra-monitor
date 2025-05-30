package model

import "time"

type User struct {
	ID             int32     `json:"p00"`
	ChangePassword bool      `json:"change_password,omitempty"`
	IsAdmin        bool      `json:"is_admin,omitempty"`
	IsDisabled     bool      `json:"is_disabled,omitempty"`
	Fullname       string    `json:"fullname"`
	Username       string    `json:"username"`
	Password       string    `json:"password,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

type Login struct {
	ID       int32  `json:"p00"`
	Fullname string `json:"fullname"`
	Username string `json:"username"`
	Token    string `json:"token"`
}
