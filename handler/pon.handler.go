package handler

import (
	"net/url"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/metalpoch/ultra-monitor/usecase"
)

type PonHandler struct {
	Usecase *usecase.PonUsecase
}

func NewPonHandler(uc *usecase.PonUsecase) *PonHandler {
	return &PonHandler{uc}
}

func (hdlr *PonHandler) GetAllByDevice(c fiber.Ctx) error {
	sysname, err := url.QueryUnescape(c.Params("sysname"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.GetAllByDevice(sysname)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

func (hdlr *PonHandler) GetByOltAndPort(c fiber.Ctx) error {
	sysname, err := url.QueryUnescape(c.Params("sysname"))
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

	res, err := hdlr.Usecase.PonByOltAndPort(sysname, shell, card, port)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}
