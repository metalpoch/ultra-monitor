package handler

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/metalpoch/olt-blueprint/traffic/model"
	"github.com/metalpoch/olt-blueprint/traffic/usecase"
)

type TrafficHandler struct {
	Usecase usecase.TrafficUsecase
}

func (hdlr TrafficHandler) GetTrafficByInterface(c fiber.Ctx) error {
	param := c.Params("id", "")
	if param == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "interface id required"})
	}

	id, err := strconv.Atoi(param)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	query := new(model.RangeDate)
	if err := c.Bind().Query(query); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.GetTrafficByInterface(uint(id), query)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

func (hdlr TrafficHandler) GetTrafficByDevice(c fiber.Ctx) error {
	param := c.Params("id", "")
	if param == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "device id required"})
	}

	id, err := strconv.Atoi(param)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	query := new(model.RangeDate)
	if err := c.Bind().Query(query); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.GetTrafficByDevice(uint(id), query)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}
