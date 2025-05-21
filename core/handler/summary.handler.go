package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/metalpoch/olt-blueprint/core/model"
	"github.com/metalpoch/olt-blueprint/core/usecase"
)

type SummaryHandler struct {
	Usecase usecase.SummaryUsecase
}

func (hdlr SummaryHandler) UserStatus(c fiber.Ctx) error {
	query := new(model.UserStatusQuery)
	if err := c.Bind().Query(query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.UserStatus(query)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}
