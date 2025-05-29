package handler

import (
	"net/url"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/metalpoch/olt-blueprint/internal/dto"
	"github.com/metalpoch/olt-blueprint/usecase"
)

type OltHandler struct {
	Usecase usecase.OltUsecase
}

func (hdlr OltHandler) Add(c fiber.Ctx) error {
	olt := new(dto.NewOlt)
	if err := c.Bind().Query(olt); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	err := hdlr.Usecase.Add(*olt)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (hdlr OltHandler) UpdateOne(c fiber.Ctx) error {
	id, err := fiber.Convert(c.Params("id"), strconv.Atoi)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	olt := new(dto.NewOlt)
	if err := c.Bind().Query(olt); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err = hdlr.Usecase.Update(uint64(id), *olt)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (hdlr OltHandler) DeleteOne(c fiber.Ctx) error {
	id, err := fiber.Convert(c.Params("id"), strconv.Atoi)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	err = hdlr.Usecase.Delete(uint64(id))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (hdlr OltHandler) GetOlt(c fiber.Ctx) error {
	id, err := fiber.Convert(c.Params("id"), strconv.Atoi)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	olt, err := hdlr.Usecase.Olt(uint64(id))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(olt)
}

func (hdlr OltHandler) GetOlts(c fiber.Ctx) error {
	pag := new(dto.Pagination)
	if err := c.Bind().Query(pag); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	olts, err := hdlr.Usecase.Olts(pag.Page, pag.Limit)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(olts)
}

func (hdlr OltHandler) GetOltsByState(c fiber.Ctx) error {
	state, err := fiber.Convert(c.Params("state"), url.QueryUnescape)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	pag := new(dto.Pagination)
	if err := c.Bind().Query(pag); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	olts, err := hdlr.Usecase.OltsByState(state, pag.Page, pag.Limit)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(olts)
}
func (hdlr OltHandler) GetOltsByCounty(c fiber.Ctx) error {
	county, err := fiber.Convert(c.Params("county"), url.QueryUnescape)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	state, err := fiber.Convert(c.Params("state"), url.QueryUnescape)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	pag := new(dto.Pagination)
	if err := c.Bind().Query(pag); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	olts, err := hdlr.Usecase.OltsByCounty(state, county, pag.Page, pag.Limit)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(olts)
}
func (hdlr OltHandler) GetOltsByMunicipality(c fiber.Ctx) error {
	municipality, err := fiber.Convert(c.Params("municipality"), url.QueryUnescape)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	county, err := fiber.Convert(c.Params("county"), url.QueryUnescape)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	state, err := fiber.Convert(c.Params("state"), url.QueryUnescape)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	pag := new(dto.Pagination)
	if err := c.Bind().Query(pag); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	olts, err := hdlr.Usecase.OltsByMunicipality(state, county, municipality, pag.Page, pag.Limit)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(olts)
}
