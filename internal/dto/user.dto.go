package dto

type NewUser struct {
	ID              int32  `json:"p00" validate:"required,gte=100000"`
	Username        string `json:"username" validate:"required,min=3,max=15,alphanum"`
	Password        string `json:"password" validate:"required,min=8"`
	PasswordConfirm string `json:"password_confirm" validate:"required,min=8,eqfield=Password"`
	Fullname        string `json:"fullname" validate:"required"`
}
type Login struct {
	Username string `json:"username" validate:"required,min=3,max=15,alphanum"`
	Password string `json:"password" validate:"required,min=8"`
}
type ChangePassword struct {
	Password        string `json:"password" validate:"required,min=8"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
	PasswordConfirm string `json:"password_confirm" validate:"required,min=8,eqfield=NewPassword"`
}
