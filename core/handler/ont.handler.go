package handler

import (
	"net/http"
	"net/url"

	"github.com/gofiber/fiber/v3"
	commonModel "github.com/metalpoch/olt-blueprint/common/model"
	"github.com/metalpoch/olt-blueprint/core/usecase"
)

type OntHandler struct {
	Usecase usecase.OntUsecase
}

func (hdlr OntHandler) OntStatus(c fiber.Ctx) error {
	dates := new(commonModel.TrafficRangeDate)
	if err := c.Bind().Query(dates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.OntStatus(dates)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

func (hdlr OntHandler) OntStatusByState(c fiber.Ctx) error {
	state, err := url.QueryUnescape(c.Params("state"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	dates := new(commonModel.TrafficRangeDate)
	if err := c.Bind().Query(dates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.OntStatusByState(state, dates)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

func (hdlr OntHandler) OntStatusByOdn(c fiber.Ctx) error {
	state, err := url.QueryUnescape(c.Params("state"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	odn, err := url.QueryUnescape(c.Params("odn"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	dates := new(commonModel.TrafficRangeDate)
	if err := c.Bind().Query(dates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.OntStatusByODN(state, odn, dates)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

func (hdlr OntHandler) TrafficOnt(c fiber.Ctx) error {
	interfaceID, err := url.QueryUnescape(c.Params("interfaceID"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	idx, err := url.QueryUnescape(c.Params("idx"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	dates := new(commonModel.TrafficRangeDate)
	if err := c.Bind().Query(dates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.TrafficOnt(interfaceID, idx, dates)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

/*
 *func (hdlr SummaryHandler) Traffic(c fiber.Ctx) error {
 *        dates := new(commonModel.TranficRangeDate)
 *        if err := c.Bind().Query(dates); err != nil {
 *                return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
 *        }
 *
 *        res, err := hdlr.Usecase.Traffic(dates)
 *        if err != nil {
 *                return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
 *        }
 *
 *        return c.JSON(res)
 *}
 *
 *func (hdlr SummaryHandler) TrafficByState(c fiber.Ctx) error {
 *        state, err := url.QueryUnescape(c.Params("state"))
 *        if err != nil {
 *                return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
 *        }
 *
 *        dates := new(commonModel.TranficRangeDate)
 *        if err := c.Bind().Query(dates); err != nil {
 *                return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
 *        }
 *
 *        res, err := hdlr.Usecase.TrafficByState(state, dates)
 *        if err != nil {
 *                return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
 *        }
 *
 *        return c.JSON(res)
 *}*/
