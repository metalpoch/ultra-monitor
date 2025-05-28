package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/olt-blueprint/entity"
	"github.com/metalpoch/olt-blueprint/internal/dto"
	"github.com/metalpoch/olt-blueprint/internal/jwt"
	"github.com/metalpoch/olt-blueprint/model"
	"github.com/metalpoch/olt-blueprint/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	secret []byte
	repo   repository.UserRepository
}

func NewUserUsecase(db *sqlx.DB, secret []byte) *UserUsecase {
	return &UserUsecase{secret, repository.NewUserRepository(db)}
}

func (uc UserUsecase) Create(newUser *dto.NewUser) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	bytesPsw, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 14)
	if err != nil {
		return err
	}

	if err := uc.repo.Create(ctx, entity.User{
		ID:             newUser.ID,
		Fullname:       newUser.Fullname,
		Email:          newUser.Email,
		Password:       string(bytesPsw),
		ChangePassword: true,
		IsAdmin:        false,
		IsDisabled:     false,
	}); err != nil {
		return err
	}

	return nil
}

func (uc UserUsecase) Login(email string, password string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := uc.repo.UserByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid email or password")
	}
	token, err := jwt.CreateJWT(uc.secret, res.ID, res.IsAdmin)
	if err != nil {
		return nil, err
	}

	return &model.User{
		ID:       res.ID,
		Fullname: res.Fullname,
		Token:    token,
	}, nil
}

func (uc UserUsecase) GetUser(id uint32) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := uc.repo.UserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &model.User{
		ID:             res.ID,
		Email:          res.Email,
		ChangePassword: res.ChangePassword,
		IsDisabled:     res.IsDisabled,
		Fullname:       res.Fullname,
		IsAdmin:        res.IsAdmin,
	}, nil
}

func (uc UserUsecase) Disable(id uint32) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := uc.repo.Disable(ctx, id); err != nil {
		return err
	}
	return nil
}

func (uc UserUsecase) Enable(id uint32) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := uc.repo.Enable(ctx, id); err != nil {
		return err
	}
	return nil
}

func (uc UserUsecase) ChangePassword(id uint32, user *dto.ChangePassword) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := uc.repo.UserByID(ctx, id)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(user.Password)); err != nil {
		return errors.New("invalid password")
	}

	bytesPsw, err := bcrypt.GenerateFromPassword([]byte(user.NewPassword), 14)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}

	if err := uc.repo.ChangePassword(ctx, id, string(bytesPsw)); err != nil {
		return err
	}
	return nil
}
