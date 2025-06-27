package handler

import (
	"net/url"

	"github.com/gofiber/fiber/v3"
	"github.com/metalpoch/ultra-monitor/internal/dto"
	"github.com/metalpoch/ultra-monitor/usecase"
)

type OltHandler struct {
	Usecase usecase.OltUsecase
}

func (hdlr *OltHandler) Add(c fiber.Ctx) error {
	var olt dto.NewOlt
	if err := c.Bind().Body(&olt); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	err := hdlr.Usecase.Add(olt)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (hdlr *OltHandler) GetOlt(c fiber.Ctx) error {
	var param dto.OltIP
	if err := c.Bind().URI(&param); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	olt, err := hdlr.Usecase.Olt(param.IP)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(olt)
}

func (hdlr *OltHandler) DeleteOne(c fiber.Ctx) error {
	var param dto.OltIP
	if err := c.Bind().URI(&param); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err := hdlr.Usecase.Delete(param.IP)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (hdlr *OltHandler) GetAllIP(c fiber.Ctx) error {
	olts, err := hdlr.Usecase.GetAllIP()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(olts)
}

func (hdlr *OltHandler) GetAllSysname(c fiber.Ctx) error {
	olts, err := hdlr.Usecase.GetAllSysname()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(olts)
}

func (hdlr *OltHandler) GetOltsByState(c fiber.Ctx) error {
	state, err := fiber.Convert(c.Params("state"), url.QueryUnescape)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	olts, err := hdlr.Usecase.OltsByState(state)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(olts)
}
