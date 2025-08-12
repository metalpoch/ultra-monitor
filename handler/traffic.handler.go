package handler

import (
	"net/url"

	"github.com/gofiber/fiber/v3"
	"github.com/metalpoch/ultra-monitor/internal/dto"
	"github.com/metalpoch/ultra-monitor/usecase"
)

type TrafficHandler struct {
	Usecase *usecase.TrafficUsecase
}

func (hdlr *TrafficHandler) DeviceLocation(c fiber.Ctx) error {
	res, err := hdlr.Usecase.DeviceLocation()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

func (hdlr *TrafficHandler) InfoInstance(c fiber.Ctx) error {
	ip, err := url.QueryUnescape(c.Params("ip"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.InfoInstance(ip)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

func (hdlr *TrafficHandler) Total(c fiber.Ctx) error {
	var dates dto.RangeDate
	if err := c.Bind().Query(&dates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.Total(dates.InitDate, dates.FinalDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

func (hdlr *TrafficHandler) Region(c fiber.Ctx) error {
	var dates dto.RangeDate
	if err := c.Bind().Query(&dates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	region, err := url.QueryUnescape(c.Params("region"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.Region(region, dates.InitDate, dates.FinalDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

func (hdlr *TrafficHandler) State(c fiber.Ctx) error {
	var dates dto.RangeDate
	if err := c.Bind().Query(&dates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	state, err := url.QueryUnescape(c.Params("state"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.State(state, dates.InitDate, dates.FinalDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

func (hdlr *TrafficHandler) GroupIP(c fiber.Ctx) error {
	var query dto.GroupedIP
	if err := c.Bind().Query(&query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	res, err := hdlr.Usecase.GroupIP(query.IP, query.InitDate, query.FinalDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

func (hdlr *TrafficHandler) IndexAndIP(c fiber.Ctx) error {
	var dates dto.RangeDate
	if err := c.Bind().Query(&dates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	ip, err := url.QueryUnescape(c.Params("ip"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	index, err := url.QueryUnescape(c.Params("index"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.IndexAndIP(ip, index, dates.InitDate, dates.FinalDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

func (hdlr *TrafficHandler) Regions(c fiber.Ctx) error {
	var dates dto.RangeDate
	if err := c.Bind().Query(&dates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.Regions(dates.InitDate, dates.FinalDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

func (hdlr *TrafficHandler) StatesByRegion(c fiber.Ctx) error {
	var dates dto.RangeDate
	if err := c.Bind().Query(&dates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	region, err := url.QueryUnescape(c.Params("region"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.StatesByRegion(region, dates.InitDate, dates.FinalDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

func (hdlr *TrafficHandler) SysnameByState(c fiber.Ctx) error {
	var dates dto.RangeDate
	if err := c.Bind().Query(&dates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	state, err := url.QueryUnescape(c.Params("state"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.SysnameByState(state, dates.InitDate, dates.FinalDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)

}
