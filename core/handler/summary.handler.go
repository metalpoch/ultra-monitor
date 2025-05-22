package handler

import (
	"net/http"
	"net/url"

	"github.com/gofiber/fiber/v3"
	commonModel "github.com/metalpoch/olt-blueprint/common/model"
	"github.com/metalpoch/olt-blueprint/core/model"
	"github.com/metalpoch/olt-blueprint/core/usecase"
)

type SummaryHandler struct {
	Usecase usecase.SummaryUsecase
}

func (hdlr SummaryHandler) UserStatus(c fiber.Ctx) error {
	query := new(commonModel.TranficRangeDate)
	if err := c.Bind().Query(query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.UserStatus(query)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

func (hdlr SummaryHandler) UserStatusByState(c fiber.Ctx) error {
	query := new(model.UserStatusByState)
	if err := c.Bind().Query(query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	res, err := hdlr.Usecase.UserStatusByState(query)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

func (hdlr SummaryHandler) Traffic(c fiber.Ctx) error {
	dates := new(commonModel.TranficRangeDate)
	if err := c.Bind().Query(dates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.Traffic(dates)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

func (hdlr SummaryHandler) TrafficByState(c fiber.Ctx) error {
	state, err := url.QueryUnescape(c.Params("state"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	dates := new(commonModel.TranficRangeDate)
	if err := c.Bind().Query(dates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.TrafficByState(state, dates)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}
