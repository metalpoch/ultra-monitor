package handler

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/metalpoch/ultra-monitor/internal/dto"
	"github.com/metalpoch/ultra-monitor/usecase"
)

type UserHandler struct {
	Usecase usecase.UserUsecase
}

func (hdlr *UserHandler) Create(c fiber.Ctx) error {
	newUser := new(dto.NewUser)
	if err := c.Bind().JSON(newUser); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := hdlr.Usecase.Create(newUser); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(http.StatusCreated)
}

func (hdlr *UserHandler) Login(c fiber.Ctx) error {
	credentials := new(dto.SignIn)

	if err := c.Bind().JSON(credentials); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	user, err := hdlr.Usecase.Login(credentials.Username, credentials.Password)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	if user.ChangePassword {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "you must change your password", "token": user.Token})
	}

	if user.IsDisabled {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "your account is disabled"})
	}

	return c.JSON(dto.SignInResponse{
		ID:       user.ID,
		Fullname: user.Fullname,
		Username: user.Username,
		Token:    user.Token,
	})
}

func (hdlr *UserHandler) AllUsers(c fiber.Ctx) error {
	users, err := hdlr.Usecase.AllUsers()

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(users)
}

func (hdlr *UserHandler) GetOwn(c fiber.Ctx) error {
	id, ok := c.Locals("id").(int32)
	if !ok {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	users, err := hdlr.Usecase.GetUser(int(id))

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(users)
}

func (hdlr *UserHandler) Disable(c fiber.Ctx) error {
	id, err := fiber.Convert(c.Params("id"), strconv.Atoi)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := hdlr.Usecase.Disable(id); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(http.StatusNoContent)
}

func (hdlr *UserHandler) Enable(c fiber.Ctx) error {
	id, err := fiber.Convert(c.Params("id"), strconv.Atoi)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := hdlr.Usecase.Enable(id); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(http.StatusOK)
}

func (hdlr *UserHandler) TemporalPassword(c fiber.Ctx) error {
	id, err := fiber.Convert(c.Params("id"), strconv.Atoi)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

  var password dto.CreateTemporalPassword
	if err := c.Bind().JSON(&password); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := hdlr.Usecase.EnableChangePassword(id, password); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(http.StatusOK)
}

func (hdlr *UserHandler) ChangePassword(c fiber.Ctx) error {
	id, ok := c.Locals("id").(int32)
	if !ok {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	passwords := new(dto.ChangePassword)
	if err := c.Bind().JSON(passwords); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err := hdlr.Usecase.ChangePassword(int(id), passwords)

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(http.StatusOK)
}
