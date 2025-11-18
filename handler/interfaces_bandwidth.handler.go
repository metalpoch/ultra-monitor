package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/metalpoch/ultra-monitor/usecase"
)

type InterfacesBandwidthHandler struct {
	Usecase *usecase.InterfaceBandwidthUsecase
}

func (hdlr *InterfacesBandwidthHandler) GetAll(c fiber.Ctx) error {
	olts, err := hdlr.Usecase.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(olts)
}
