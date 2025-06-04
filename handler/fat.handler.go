package handler

import (
	"net/url"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/metalpoch/ultra-monitor/internal/dto"
	"github.com/metalpoch/ultra-monitor/model"
	"github.com/metalpoch/ultra-monitor/usecase"
)

type FatHandler struct {
	Usecase *usecase.FatUsecase
}

func (hdlr *FatHandler) AddFat(c fiber.Ctx) error {
	fat := new(model.Fat)
	if err := c.Bind().Body(&fat); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := hdlr.Usecase.AddFat(*fat); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusCreated)
}

func (hdlr *FatHandler) DeleteOne(c fiber.Ctx) error {
	id, err := fiber.Convert(c.Params("id"), strconv.Atoi)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	if err := hdlr.Usecase.DeleteOne(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (hdlr *FatHandler) GetAll(c fiber.Ctx) error {
	pag := new(dto.Pagination)
	if err := c.Bind().Query(pag); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.GetAll(*pag)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (hdlr *FatHandler) GetByID(c fiber.Ctx) error {
	id, err := fiber.Convert(c.Params("id"), strconv.Atoi)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	tao, err := hdlr.Usecase.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(tao)
}

func (hdlr *FatHandler) GetByFat(c fiber.Ctx) error {
	tao, err := url.QueryUnescape(c.Params("fat"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.GetByFat(tao)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (hdlr *FatHandler) GetByOdn(c fiber.Ctx) error {
	state, err := url.QueryUnescape(c.Params("state"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	odn, err := url.QueryUnescape(c.Params("odn"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.GetByOdn(state, odn)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (hdlr *FatHandler) GetOdnStates(c fiber.Ctx) error {
	state := c.Query("state")
	res, err := hdlr.Usecase.GetOdnStates(state)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (hdlr *FatHandler) GetOdnCounty(c fiber.Ctx) error {
	state, err := url.QueryUnescape(c.Params("state"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	county, err := url.QueryUnescape(c.Params("county"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.GetOdnCounty(state, county)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (hdlr *FatHandler) GetOdnMunicipality(c fiber.Ctx) error {
	state, err := url.QueryUnescape(c.Params("state"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	county, err := url.QueryUnescape(c.Params("county"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	municipality, err := url.QueryUnescape(c.Params("municipality"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.GetOdnMunicipality(state, county, municipality)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (hdlr *FatHandler) GetOdnByOlt(c fiber.Ctx) error {
	oltIP, err := url.QueryUnescape(c.Params("oltIP"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if oltIP == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "olt_ip query parameter is required"})
	}

	res, err := hdlr.Usecase.GetOdnByOlt(oltIP)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (hdlr *FatHandler) GetOdnOltPort(c fiber.Ctx) error {
	oltIP, err := url.QueryUnescape(c.Params("oltIP"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	shell, err := fiber.Convert(c.Params("shell"), strconv.Atoi)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	card, err := fiber.Convert(c.Params("card"), strconv.Atoi)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	port, err := fiber.Convert(c.Params("port"), strconv.Atoi)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.GetOdnOltPort(oltIP, shell, card, port)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}
