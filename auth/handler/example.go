package handler

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/metalpoch/olt-blueprint/auth/model"
	"github.com/metalpoch/olt-blueprint/auth/usecase"
)

type ExampleHandler struct {
	Usecase usecase.ExampleUsecase
}

func (hdlr ExampleHandler) Get(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "id" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "ID invalido"})

	}

	intID, err := strconv.Atoi(id)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.Get(uint8(intID))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(res)
}

// Ejemplo bindeando un modelo
func (hdlr ExampleHandler) Create(ctx fiber.Ctx) error {
	newExample := new(model.NewExample)
	if err := ctx.Bind().JSON(newExample); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.Create(newExample)
	if err != nil {
		return ctx.JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(res)
}
