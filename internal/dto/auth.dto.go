package dto

import "time"

type NewUser struct {
	ID              int32  `json:"p00" validate:"required,gte=100000"`
	Username        string `json:"username" validate:"required,min=3,max=15,alphanum"`
	Password        string `json:"password" validate:"required,min=8"`
	PasswordConfirm string `json:"password_confirm" validate:"required,min=8,eqfield=Password"`
	Fullname        string `json:"fullname" validate:"required"`
}
type SignIn struct {
	Username string `json:"username" validate:"required,min=3,max=15,alphanum"`
	Password string `json:"password" validate:"required,min=8"`
}
type ChangePassword struct {
	Password        string `json:"password" validate:"required,min=8"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
	PasswordConfirm string `json:"password_confirm" validate:"required,min=8,eqfield=NewPassword"`
}

type SignInResponse struct {
	ID       int32  `json:"id"`
	Fullname string `json:"fullname"`
	Username string `json:"username"`
	Token    string `json:"token"`
}

type UserResponse struct {
	ID             int32     `json:"id"`
	ChangePassword bool      `json:"change_password"`
	IsAdmin        bool      `json:"is_admin"`
	IsDisabled     bool      `json:"is_disabled"`
	Fullname       string    `json:"fullname"`
	Username       string    `json:"username"`
	CreatedAt      time.Time `json:"created_at"`
}

type User struct {
	ID             int32     `json:"id"`
	ChangePassword bool      `json:"change_password"`
	IsAdmin        bool      `json:"is_admin"`
	IsDisabled     bool      `json:"is_disabled"`
	Fullname       string    `json:"fullname"`
	Username       string    `json:"username"`
	Password       string    `json:"password"`
	CreatedAt      time.Time `json:"created_at"`
}

type UserLogin struct {
	User
	Token string `json:"token"`
}
