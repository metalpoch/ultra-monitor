package usecase

import "github.com/metalpoch/olt-blueprint/common/model"

type UserUsecase interface {
	Create(newUser *model.NewUser) error
	Login(email string, password string) (*model.LoginResponse, error)
	GetUser(id uint) (*model.User, error)
	GetAllUsers() ([]*model.FullUser, error)
	SoftDelete(id uint) error
	ChangePassword(id uint, passwords *model.ChangePassword) error
}
