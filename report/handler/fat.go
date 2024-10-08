package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/metalpoch/olt-blueprint/common/model"
	"github.com/metalpoch/olt-blueprint/report/usecase"
	"github.com/metalpoch/olt-blueprint/report/utils"
)

type FatHandler struct {
	Usecase usecase.FatUsecase
}

// Save Fats
//
//	@Summary		Save a fat
//	@Description	Save a fat into the  database
//	@Tags			fat
//	@Produce		json
//	@Success		201
//	@Failure		400	{object}	object{message=string}
//	@Failure		500	{object}	object{message=string}
//	@Router			/fat/ [post]
func (hdlr FatHandler) Add(c fiber.Ctx) error {
	newFat := new(model.NewFat)
	if err := c.Bind().JSON(newFat); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := hdlr.Usecase.Add(newFat); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusCreated)
}

// Get a Fat
//
//	@Summary		Get a fat id
//	@Description	get a fat by id from database
//	@Tags			fat
//	@Produce		json
//	@Success		200	{object}	[]model.FatResponse
//	@Failure		400	{object}	object{message=string}
//	@Failure		500	{object}	object{message=string}
//	@Router			/fat/:id [get]
func (hdlr FatHandler) Get(c fiber.Ctx) error {
	id, err := fiber.Convert(c.Params("id"), strconv.Atoi)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	fat, err := hdlr.Usecase.Get(uint(id))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(utils.FatResponse(fat))
}
