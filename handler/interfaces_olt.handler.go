package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/metalpoch/ultra-monitor/internal/dto"
	"github.com/metalpoch/ultra-monitor/usecase"
)

type InterfacesOltHandler struct {
	Usecase *usecase.InterfaceOltUsecase
}

func (hdlr *InterfacesOltHandler) GetAll(c fiber.Ctx) error {
	olts, err := hdlr.Usecase.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(olts)
}

func (hdlr *InterfacesOltHandler) Update(c fiber.Ctx) error {
	var newValue dto.InterfacesOlt
	if err := c.Bind().Body(&newValue); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := hdlr.Usecase.Update(newValue); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusOK)
}
