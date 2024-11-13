package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/metalpoch/olt-blueprint/common/model"
	"github.com/metalpoch/olt-blueprint/report/usecase"
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
//	@Failure		400	{object}	object{error=string}
//	@Failure		500	{object}	object{error=string}
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
//	@Success		200	{object}	model.FatResponse
//	@Failure		400	{object}	object{error=string}
//	@Failure		404	{object}	object{error=string}
//	@Failure		500	{object}	object{error=string}
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

	if fat.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{})
	}

	return c.JSON(fat)
}

// Get all Fats
//
//	@Summary		Get all fats
//	@Description	get all fats from database
//	@Tags			fat
//	@Produce		json
//	@Success		200	{object}	[]model.FatResponse
//	@Failure		400	{object}	object{error=string}
//	@Failure		500	{object}	object{error=string}
//	@Router			/fat/ [get]
func (hdlr FatHandler) GetAll(c fiber.Ctx) error {
	page := new(model.Page)
	if err := c.Bind().Query(page); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	fats, err := hdlr.Usecase.GetAll(page)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fats)
}

// Delete Fat
//
//	@Summary		delete a fat id
//	@Description	delete a fat by id from database
//	@Tags			fat
//	@Produce		json
//	@Success		200
//	@Failure		400	{object}	object{error=string}
//	@Failure		500	{object}	object{error=string}
//	@Router			/fat/:id [delete]
func (hdlr FatHandler) Delete(c fiber.Ctx) error {
	id, err := fiber.Convert(c.Params("id"), strconv.Atoi)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := hdlr.Usecase.Delete(uint(id)); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusOK)
}
