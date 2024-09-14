package usecase

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/metalpoch/olt-blueprint/auth/entity"
	"github.com/metalpoch/olt-blueprint/auth/model"
	"github.com/metalpoch/olt-blueprint/auth/repository"
	"github.com/metalpoch/olt-blueprint/auth/utils"
)

type userUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) *userUsecase {
	return &userUsecase{repo}
}

func (use userUsecase) Create(newUser *model.NewUser) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if newUser.Password != newUser.PasswordConfirm {
		return errors.New("passwords do not match")
	}
	password, err := utils.HashPassword(newUser.Password)
	if err != nil {
		return err
	}

	user := &entity.User{
		ID:               newUser.ID,
		Fullname:         newUser.Fullname,
		Email:            newUser.Email,
		Password:         password,
		ChangePassword:   true,
		IsAdmin:          false,
		StatesPermission: []string{},
	}

	if err := use.repo.Create(ctx, user); err != nil {
		return err
	}

	return nil
}

func (use userUsecase) Login(email string, password string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := use.repo.GetUserByEmail(ctx, email)
	if err != nil {
		log.Println("An error has occurred:", err.Error())
		return nil, errors.New("invalid email or password")
	}

	if err := utils.CheckPasswordHash(password, res.Password); err != nil {
		log.Println("An error has occurred:", err.Error())
		return nil, errors.New("invalid email or password")
	}

	return &model.User{
		ID:               res.ID,
		Email:            res.Email,
		ChangePassword:   res.ChangePassword,
		Fullname:         res.Fullname,
		IsAdmin:          res.IsAdmin,
		StatesPermission: res.StatesPermission,
	}, nil

}

func (use userUsecase) GetAll() ([]*model.User, error) {
	users := []*model.User{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := use.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	for _, e := range res {
		users = append(users, &model.User{
			ID:               e.ID,
			Email:            e.Email,
			ChangePassword:   e.ChangePassword,
			Fullname:         e.Fullname,
			IsAdmin:          e.IsAdmin,
			StatesPermission: e.StatesPermission,
		})
	}

	return users, nil
}

func (use userUsecase) SoftDelete(id uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := use.repo.SoftDelete(ctx, id); err != nil {
		return err
	}
	return nil
}

func (use userUsecase) ChangePassword(id uint, password, validatePassowrd string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if password != validatePassowrd {
		return errors.New("passwords do not match")
	}

	if err := use.repo.ChangePassword(ctx, id, password); err != nil {
		log.Println("An error has occurred:", err.Error())
		return err
	}
	return nil
}
