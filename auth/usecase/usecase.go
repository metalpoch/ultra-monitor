package usecase

import "github.com/metalpoch/olt-blueprint/auth/model"

type UserUsecase interface {
	Create(newUser *model.NewUser) error
	Login(email string, password string) (*model.LoginResponse, error)
	GetUser(id uint) (*model.User, error)
	GetAllUsers() ([]*model.User, error)
	SoftDelete(id string) error
	ChangePassword(id uint, passwords *model.ChangePassword) error
}
