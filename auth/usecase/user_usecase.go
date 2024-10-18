package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/metalpoch/olt-blueprint/auth/repository"
	"github.com/metalpoch/olt-blueprint/auth/utils"
	"github.com/metalpoch/olt-blueprint/common/constants"
	"github.com/metalpoch/olt-blueprint/common/entity"
	"github.com/metalpoch/olt-blueprint/common/model"
	"github.com/metalpoch/olt-blueprint/common/pkg/tracking"
	"gorm.io/gorm"
)

type userUsecase struct {
	secret   []byte
	telegram tracking.SmartModule
	repo     repository.UserRepository
}

func NewUserUsecase(db *gorm.DB, secret []byte, telegram tracking.SmartModule) *userUsecase {
	return &userUsecase{secret, telegram, repository.NewUserRepository(db)}
}

func (use userUsecase) Create(newUser *model.NewUser) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if newUser.Password != newUser.PasswordConfirm {
		return errors.New("passwords do not match")
	}
	password, err := utils.HashPassword(newUser.Password)
	if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_AUTH,
			constants.CATEGORY_UTILS,
			fmt.Sprintf("(userUsecase).Create - utils.HashPassword(%s)", newUser.Password),
			err,
		)
		return err
	}

	user := &entity.User{
		ID:             newUser.ID,
		Fullname:       newUser.Fullname,
		Email:          newUser.Email,
		Password:       password,
		ChangePassword: true,
		IsAdmin:        false,
	}

	if err := use.repo.Create(ctx, user); err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_AUTH,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(userUsecase).Create - use.repo.Create(ctx, %v)", user),
			err,
		)
		return err
	}

	return nil
}

func (use userUsecase) Login(email string, password string) (*model.LoginResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := use.repo.GetUserByEmail(ctx, email)
	if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_AUTH,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(userUsecase).Login - use.repo.GetUserByEmail(ctx, %s)", email),
			err,
		)
		return nil, errors.New("invalid email or password")
	}

	if err := utils.CheckPasswordHash(password, res.Password); err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_AUTH,
			constants.CATEGORY_UTILS,
			fmt.Sprintf("(userUsecase).Login - utils.CheckPasswordHash(%s, %s)", password, res.Password),
			err,
		)
		return nil, errors.New("invalid email or password")
	}

	token, err := utils.CreateJWT(use.secret, res.ID, res.IsAdmin)
	if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_AUTH,
			constants.CATEGORY_UTILS,
			fmt.Sprintf("(userUsecase).Login - utils.CreateJWT(%x, %d, %t)", use.secret, res.ID, res.IsAdmin),
			err,
		)
		return nil, err
	}

	return &model.LoginResponse{
		Token: token,
		User: model.User{
			ID:             res.ID,
			Email:          res.Email,
			ChangePassword: res.ChangePassword,
			Fullname:       res.Fullname,
			IsAdmin:        res.IsAdmin,
		},
	}, nil
}

func (use userUsecase) GetUser(id uint) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := use.repo.GetUserByID(ctx, id)
	if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_AUTH,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(userUsecase).GetUser - use.repo.GetUserByID(ctx, %d)", id),
			err,
		)
		return nil, err
	}

	return &model.User{
		ID:             res.ID,
		Email:          res.Email,
		ChangePassword: res.ChangePassword,
		Fullname:       res.Fullname,
		IsAdmin:        res.IsAdmin,
	}, nil
}

func (use userUsecase) GetAllUsers() ([]*model.FullUser, error) {
	users := []*model.FullUser{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := use.repo.GetAll(ctx)
	if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_AUTH,
			constants.CATEGORY_DATABASE,
			"(userUsecase).GetAllUsers - use.repo.GetAll(ctx)",
			err,
		)
		return nil, err
	}

	for _, e := range res {
		users = append(users, &model.FullUser{
			ID:             e.ID,
			Email:          e.Email,
			ChangePassword: e.ChangePassword,
			Fullname:       e.Fullname,
			IsAdmin:        e.IsAdmin,
			CreatedAt:      e.CreatedAt,
			UpdatedAt:      e.UpdatedAt,
		})
	}

	return users, nil
}

func (use userUsecase) SoftDelete(id uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := use.repo.SoftDelete(ctx, id); err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_AUTH,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(userUsecase).SoftDelete - use.repo.SoftDelete(ctx, %d)", id),
			err,
		)
		return err
	}
	return nil
}

func (use userUsecase) ChangePassword(id uint, user *model.ChangePassword) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := use.repo.GetUserByID(ctx, id)
	if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_AUTH,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(userUsecase).ChangePassword - use.repo.GetUserByID(ctx, %d)", id),
			err,
		)
		return err
	}

	if err := utils.CheckPasswordHash(user.Password, res.Password); err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_AUTH,
			constants.CATEGORY_UTILS,
			fmt.Sprintf("(userUsecase).ChangePassword - utils.CheckPasswordHash(%s, %s)", user.Password, res.Password),
			err,
		)
		return errors.New("invalid password")
	}

	if user.NewPassword != user.NewPasswordConfirm {
		return errors.New("passwords do not match")
	}

	password, err := utils.HashPassword(user.NewPassword)
	if err != nil {
		return err
	}

	if err := use.repo.ChangePassword(ctx, id, password); err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_AUTH,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(userUsecase).ChangePassword - use.repo.ChangePassword(ctx, %d, %s)", id, password),
			err,
		)
		return err
	}
	return nil
}
