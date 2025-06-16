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

func (hdlr *OntHandler) OntStatusIPSummary(c fiber.Ctx) error {
	ip, err := url.QueryUnescape(c.Params("ip"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var dates dto.RangeDate
	if err := c.Bind().Query(&dates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.GetStatusIPSummary(ip, dates)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (hdlr *OntHandler) OntStatusStateSummary(c fiber.Ctx) error {
	var dates dto.RangeDate
	if err := c.Bind().Query(&dates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.GetStatusStateSummary(dates)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (hdlr *OntHandler) OntStatusByStateSummary(c fiber.Ctx) error {
	state, err := url.QueryUnescape(c.Params("state"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var dates dto.RangeDate
	if err := c.Bind().Query(&dates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.GetStatusByStateSummary(state, dates)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (hdlr *OntHandler) OntStatusMunicipalitySummary(c fiber.Ctx) error {
	state, err := url.QueryUnescape(c.Params("state"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var dates dto.RangeDate
	if err := c.Bind().Query(&dates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.GetStatusMunicipalitySummary(state, dates)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (hdlr *OntHandler) OntStatusCountySummary(c fiber.Ctx) error {
	state, err := url.QueryUnescape(c.Params("state"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	municipality, err := url.QueryUnescape(c.Params("municipality"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var dates dto.RangeDate
	if err := c.Bind().Query(&dates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.GetStatusCountySummary(state, municipality, dates)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (hdlr *OntHandler) OntStatusOdnSummary(c fiber.Ctx) error {
	state, err := url.QueryUnescape(c.Params("state"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	municipality, err := url.QueryUnescape(c.Params("municipality"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	county, err := url.QueryUnescape(c.Params("county"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var dates dto.RangeDate
	if err := c.Bind().Query(&dates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.GetStatusOdnSummary(state, municipality, county, dates)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (hdlr *OntHandler) TrafficOnt(c fiber.Ctx) error {
	ponID, err := fiber.Convert(c.Params("ponID"), strconv.Atoi)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	ontIDX, err := fiber.Convert(c.Params("idx"), strconv.Atoi)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var dates dto.RangeDate
	if err := c.Bind().Query(&dates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.TrafficOnt(ponID, int64(ontIDX), dates)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

func (hdlr *OntHandler) TrafficOntByDespt(c fiber.Ctx) error {
	despt, err := url.QueryUnescape(c.Params("despt"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var dates dto.RangeDate
	if err := c.Bind().Query(&dates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	traffic, err := hdlr.Usecase.TrafficOntByDespt(despt, dates)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(traffic)
}

func (hdlr *OntHandler) GetStatusSummaryForecast(c fiber.Ctx) error {
	futureDays, err := strconv.Atoi(c.Query("days", "-1"))
	if err != nil || futureDays < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "days must be a positive integer"})
	}
	res, err := hdlr.Usecase.GetStatusSummaryForecast(futureDays)
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
	res, err := hdlr.Usecase.GetStatusByStateForecast(state, futureDays)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}
