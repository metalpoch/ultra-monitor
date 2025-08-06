package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/metalpoch/ultra-monitor/usecase"
)

type PrometheusHandler struct {
	Usecase *usecase.PrometheusUsecase
}

func (hdlr *PrometheusHandler) GetGponPortsStatus(c fiber.Ctx) error {
	res, err := hdlr.Usecase.GponPortsStatus()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}
