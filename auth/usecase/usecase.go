package usecase

import "github.com/metalpoch/olt-blueprint/auth/model"

type UserUsecase interface {
	Create(newUser *model.NewUser) error
	Login(email string, password string) (*model.User, error)
	GetAll() ([]*model.User, error)
	SoftDelete(id uint) error
	ChangePassword(id uint, password, validatePassowrd string) error
}
