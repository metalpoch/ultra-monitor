package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/metalpoch/olt-blueprint/auth/usecase"
	"github.com/metalpoch/olt-blueprint/common/model"
)

type UserHandler struct {
	Usecase usecase.UserUsecase
}

func (hdlr UserHandler) Create(ctx fiber.Ctx) error {
	newUser := new(model.NewUser)
	if err := ctx.Bind().JSON(newUser); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := hdlr.Usecase.Create(newUser); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.SendStatus(http.StatusCreated)
}

func (hdlr UserHandler) Login(ctx fiber.Ctx) error {
	credentials := new(model.Login)

	if err := ctx.Bind().JSON(credentials); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	res, err := hdlr.Usecase.Login(credentials.Email, credentials.Password)
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(res)
}

func (hdlr UserHandler) GetOwn(ctx fiber.Ctx) error {
	id := uint(ctx.Locals("id").(float64))
	users, err := hdlr.Usecase.GetUser(id)

	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(users)
}

func (hdlr UserHandler) GetAll(ctx fiber.Ctx) error {
	users, err := hdlr.Usecase.GetAllUsers()
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(users)
}

func (hdlr UserHandler) DeleteUser(ctx fiber.Ctx) error {
	id := uint(ctx.Locals("id").(float64))

	if err := hdlr.Usecase.SoftDelete(id); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.SendStatus(http.StatusOK)
}

func (hdlr UserHandler) ChangePassword(ctx fiber.Ctx) error {
	id := uint(ctx.Locals("id").(float64))
	passwords := new(model.ChangePassword)
	if err := ctx.Bind().JSON(passwords); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err := hdlr.Usecase.ChangePassword(id, passwords)

	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.SendStatus(http.StatusOK)
}
