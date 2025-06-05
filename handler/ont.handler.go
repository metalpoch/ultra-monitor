package handler

import (
	"net/url"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/metalpoch/ultra-monitor/internal/dto"
	"github.com/metalpoch/ultra-monitor/usecase"
)

type OntHandler struct {
	Usecase usecase.OntUsecase
}

func (hdlr *OntHandler) OntStatus(c fiber.Ctx) error {
	res, err := hdlr.Usecase.AllOntStatus()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (hdlr *OntHandler) OntStatusByState(c fiber.Ctx) error {
	state, err := url.QueryUnescape(c.Params("state"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.OntStatusByState(state)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (hdlr *OntHandler) OntStatusByOdn(c fiber.Ctx) error {
	state, err := url.QueryUnescape(c.Params("state"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	odn, err := url.QueryUnescape(c.Params("odn"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	dates := new(dto.RangeDate)
	if err := c.Bind().Query(dates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	res, err := hdlr.Usecase.OntStatusByOdn(state, odn, *dates)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (hdlr *OntHandler) OntStatusByOltIP(c fiber.Ctx) error {
	ip, err := url.QueryUnescape(c.Params("ip"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	dates := new(dto.RangeDate)
	if err := c.Bind().Query(dates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.OntStatusByOltIP(ip, *dates)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (hdlr *OntHandler) OntStatusBySysname(c fiber.Ctx) error {
	sysname, err := url.QueryUnescape(c.Params("sysname"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	dates := new(dto.RangeDate)
	if err := c.Bind().Query(dates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.OntStatusBySysname(sysname, *dates)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (hdlr *OntHandler) TrafficOnt(c fiber.Ctx) error {
	ponID, err := fiber.Convert(c.Params("ponID"), strconv.Atoi)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	ontIDX, err := url.QueryUnescape(c.Params("ontIDX"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	dates := new(dto.RangeDate)
	if err := c.Bind().Query(dates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.TrafficOnt(uint64(ponID), ontIDX, *dates)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

func (hdlr *OntHandler) AllOntStatusForecast(c fiber.Ctx) error {
	futureDays, err := strconv.Atoi(c.Query("days", "-1"))
	if err != nil || futureDays < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "days must be a positive integer"})
	}
	dates := new(dto.RangeDate)
	if err := c.Bind().Query(dates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	res, err := hdlr.Usecase.AllOntStatusForecast(*dates, futureDays)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (hdlr *OntHandler) OntStatusByStateForecast(c fiber.Ctx) error {
	state, err := url.QueryUnescape(c.Params("state"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	futureDays, err := strconv.Atoi(c.Query("days", "-1"))
	if err != nil || futureDays < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "days must be a positive integer"})
	}
	dates := new(dto.RangeDate)
	if err := c.Bind().Query(dates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	res, err := hdlr.Usecase.OntStatusByStateForecast(state, *dates, futureDays)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (hdlr *OntHandler) OntStatusByODNForecast(c fiber.Ctx) error {
	state, err := url.QueryUnescape(c.Params("state"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	odn, err := url.QueryUnescape(c.Params("odn"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	futureDays, err := strconv.Atoi(c.Query("days", "-1"))
	if err != nil || futureDays < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "days must be a positive integer"})
	}
	dates := new(dto.RangeDate)
	if err := c.Bind().Query(dates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	res, err := hdlr.Usecase.OntStatusByODNForecast(state, odn, *dates, futureDays)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (h *OntHandler) TrafficOntByDespt(c fiber.Ctx) error {
	despt, err := url.QueryUnescape(c.Params("despt"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	dates := new(dto.RangeDate)
	if err := c.Bind().Query(dates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	traffic, err := h.Usecase.TrafficOntByDespt(despt, *dates)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(traffic)
}
