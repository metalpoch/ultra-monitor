package model

type NewUser struct {
	Id               string   `json:"_id"`
	Email            string   `json:"email"`
	Password         string   `json:"password"`
	ChangePassw      bool     `json:"change_password"`
	Fullname         string   `json:"fullname"`
	P00              uint     `json:"p00"`
	IsAdmin          bool     `json:"is_admin"`
	StatesPermission []string `json:"states_permission"`
}

type Users []*NewUser

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
