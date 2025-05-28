package model

import "time"

type User struct { // Table
	ID             uint32    `json:"p00"`
	ChangePassword bool      `json:"change_password,omitempty"`
	IsAdmin        bool      `json:"is_admin,omitempty"`
	IsDisabled     bool      `json:"is_disabled,omitempty"`
	Fullname       string    `json:"fullname"`
	Email          string    `json:"email"`
	Password       string    `json:"password,omitempty"`
	Token          string    `json:"token,omitempty"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
}
