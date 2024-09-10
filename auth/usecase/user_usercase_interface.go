package usecase

import "github.com/metalpoch/olt-blueprint/auth/model"

type UserUsecase interface {
	Create(newUser *model.NewUser) (string, error)
	Get() (*model.Users, error)
	GetValue(clave string, valor string) (*model.Users, error)
	DeleteName(name string) (string, error)
	ChangePassword(user *model.NewUser) (string, error)
	Login(email string, password string) (*model.NewUser, error)
}
