package handler

import (
	"net/http"
	//"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/metalpoch/olt-blueprint/auth/model"
	"github.com/metalpoch/olt-blueprint/auth/usecase"
)

type UserHandler struct {
	Usecase usecase.UserUsecase
}

func (hdlr UserHandler) Get(ctx fiber.Ctx) error {

	res, err := hdlr.Usecase.Get()
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(res)
}

func (hdlr UserHandler) Create(ctx fiber.Ctx) error {
	newUser := new(model.NewUser)
	if err := ctx.Bind().JSON(newUser); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error al insertar": err.Error()})
	}

	res, err := hdlr.Usecase.Create(newUser)
	if err != nil {
		return ctx.JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(res)
}

func (hdlr UserHandler) GetValue(ctx fiber.Ctx) error {
	clave := ctx.Params("clave")
	valor := ctx.Params("valor")
	if clave == "clave" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "clave invalida"})
	}
	if valor == "valor" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "valor invalida"})
	}

	res, err := hdlr.Usecase.GetValue(clave, valor)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error al obtner los valores": err.Error()})
	}
	return ctx.JSON(res)
}

func (hdlr UserHandler) DeleteName(ctx fiber.Ctx) error {
	name := ctx.Params("p00")
	if name == "p00" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "fullname invalida"})
	}

	res, err := hdlr.Usecase.DeleteName(name)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error al obtner los valores": err.Error()})
	}
	return ctx.JSON(res)
}
func (hdlr UserHandler) ChangePassword(ctx fiber.Ctx) error {
	newUser := new(model.NewUser)
	if err := ctx.Bind().JSON(newUser); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.ChangePassword(newUser)

	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error al cambiar contraseña": err.Error()})
	}
	return ctx.JSON(res)
}
func (hdlr UserHandler) Login(ctx fiber.Ctx) error {
	email := ctx.Params("email")
	password := ctx.Params("password")
	if email == "email" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "email invalida"})
	}
	if password == "password" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "password invalida"})
	}

	_, user, err := hdlr.Usecase.Login(email, password)

	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error al inicar sesion": err.Error()})
	}
	return ctx.JSON(user)
}
