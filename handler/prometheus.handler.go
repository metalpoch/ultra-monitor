package handler

import (
	"net/url"

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

func (hdlr *PrometheusHandler) GetGponPortsStatusByRegion(c fiber.Ctx) error {
	region, err := url.QueryUnescape(c.Params("region"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.GponPortsStatusByRegion(region)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

func (hdlr *PrometheusHandler) GetGponPortsStatusByState(c fiber.Ctx) error {
	state, err := url.QueryUnescape(c.Params("state"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.GponPortsStatusByState(state)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}
