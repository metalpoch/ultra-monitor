package usecase

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/metalpoch/olt-blueprint/auth/entity"
	"github.com/metalpoch/olt-blueprint/auth/model"
	"github.com/metalpoch/olt-blueprint/auth/repository"
	"github.com/metalpoch/olt-blueprint/auth/utils"
)

type userUsecase struct {
	secret []byte
	repo   repository.UserRepository
}

func NewUserUsecase(db *sql.DB, secret []byte) *userUsecase {
	return &userUsecase{secret, repository.NewUserRepository(db)}
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

func (use userUsecase) Login(email string, password string) (*model.LoginResponse, error) {
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

	token, err := utils.CreateJWT(use.secret, res.ID, res.IsAdmin)
	if err != nil {
		log.Println("An error has occurred:", err.Error())
		return nil, err
	}

	return &model.LoginResponse{
		Token: token,
		User: model.User{
			ID:               res.ID,
			Email:            res.Email,
			ChangePassword:   res.ChangePassword,
			Fullname:         res.Fullname,
			IsAdmin:          res.IsAdmin,
			StatesPermission: res.StatesPermission,
		},
	}, nil
}

func (use userUsecase) GetUser(id uint) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := use.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
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

func (use userUsecase) GetAllUsers() ([]*model.User, error) {
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

func (use userUsecase) SoftDelete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	p00, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	if err := use.repo.SoftDelete(ctx, uint(p00)); err != nil {
		return err
	}
	return nil
}

func (use userUsecase) ChangePassword(id uint, user *model.ChangePassword) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := use.repo.GetUserByID(ctx, id)
	if err != nil {
		return err
	}

	if err := utils.CheckPasswordHash(user.Password, res.Password); err != nil {
		return errors.New("invalid password")
	}

	if user.NewPassword != user.NewPasswordConfirm {
		return errors.New("passwords do not match")
	}

	if err := use.repo.ChangePassword(ctx, id, user.NewPassword); err != nil {
		log.Println("An error has occurred:", err.Error())
		return err
	}
	return nil
}
