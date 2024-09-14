package model

type NewUser struct {
	ID              uint   `json:"p00" validate:"required,gte=100000"`
	Email           string `json:"email" validate:"required"`
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

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
