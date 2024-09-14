package model

type NewUser struct {
	ID              uint   `json:"p00" validate:"required,gte=100000"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8"`
	PasswordConfirm string `json:"password_confirm" validate:"required,min=8"`
	Fullname        string `json:"fullname" validate:"required"`
}

type User struct {
	ID               uint     `json:"p00"`
	Email            string   `json:"email"`
	ChangePassword   bool     `json:"change_password"`
	Fullname         string   `json:"fullname"`
	IsAdmin          bool     `json:"is_admin"`
	StatesPermission []string `json:"states_permission"`
}

type ChangePassword struct {
	Password           string `json:"password" validate:"required,min=8"`
	NewPassword        string `json:"new_password" validate:"required,min=8"`
	NewPasswordConfirm string `json:"new_password_confirm" validate:"required,min=8"`
}

type Login struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginResponse struct {
	User
	Token string `json:"token"`
}
