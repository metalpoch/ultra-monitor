package usecase

import (
	"context"
	"time"

	"github.com/metalpoch/olt-blueprint/auth/entity"
	"github.com/metalpoch/olt-blueprint/auth/model"
	"github.com/metalpoch/olt-blueprint/auth/repository"
)

type userUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) *userUsecase {
	return &userUsecase{repo}
}

func (use userUsecase) Create(newUser *model.NewUser) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	data := entity.User{
		Email:            newUser.Email,
		Password:         newUser.Password,
		ChangePassw:      newUser.ChangePassw,
		Fullname:         newUser.Fullname,
		P00:              newUser.P00,
		IsAdmin:          newUser.IsAdmin,
		StatesPermission: newUser.StatesPermission,
	}

	res, err := use.repo.Create(ctx, &data)
	if err != nil {
		return "Error 2", err
	}
	return res, nil

}

func (use userUsecase) Get() (*model.Users, error) {
	var modelUsers model.Users
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := use.repo.Get(ctx)
	if err != nil {
		return nil, err
	}

	for _, userEntity := range *res {
		newModel := model.NewUser{
			Id:               userEntity.ID.String(),
			Email:            userEntity.Email,
			Password:         userEntity.Password,
			ChangePassw:      userEntity.ChangePassw,
			Fullname:         userEntity.Fullname,
			P00:              userEntity.P00,
			IsAdmin:          userEntity.IsAdmin,
			StatesPermission: userEntity.StatesPermission,
		}
		temp := &newModel
		modelUsers = append(modelUsers, temp)
	}

	return &modelUsers, nil
}

func (use userUsecase) GetValue(clave string, valor string) (*model.Users, error) {
	var modelUsers model.Users
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := use.repo.GetValue(ctx, clave, valor)
	if err != nil {
		return nil, err
	}

	for _, userEntity := range *res {
		newModel := model.NewUser{
			Id:               userEntity.ID.String(),
			Email:            userEntity.Email,
			Password:         userEntity.Password,
			ChangePassw:      userEntity.ChangePassw,
			Fullname:         userEntity.Fullname,
			P00:              userEntity.P00,
			IsAdmin:          userEntity.IsAdmin,
			StatesPermission: userEntity.StatesPermission,
		}
		temp := &newModel
		modelUsers = append(modelUsers, temp)
	}

	return &modelUsers, nil
}

func (use userUsecase) DeleteName(name string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := use.repo.DeleteName(ctx, name)

	if err != nil {
		return "Ocurrio un error 2", err
	}

	return res, nil
}

func (use userUsecase) ChangePassword(user *model.NewUser) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	data := entity.User{
		Email:            user.Email,
		Password:         user.Password,
		ChangePassw:      user.ChangePassw,
		Fullname:         user.Fullname,
		P00:              user.P00,
		IsAdmin:          user.IsAdmin,
		StatesPermission: user.StatesPermission,
	}

	res, err := use.repo.ChangePassword(ctx, &data)
	if err != nil {
		return "ocurrio un error 2", err
	}

	return res, nil
}

func (use userUsecase) Login(email string, password string) (string, *model.NewUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, user, err := use.repo.Login(ctx, email, password)
	if err != nil {
		return "ocurrio un error en 2", nil, err
	}

	newModel := model.NewUser{
		Id:               user.ID.String(),
		Email:            user.Email,
		Password:         user.Password,
		ChangePassw:      user.ChangePassw,
		Fullname:         user.Fullname,
		P00:              user.P00,
		IsAdmin:          user.IsAdmin,
		StatesPermission: user.StatesPermission,
	}

	return res, &newModel, nil
}
