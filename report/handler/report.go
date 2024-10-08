package handler

import (
	"path"

	"github.com/gofiber/fiber/v3"
	"github.com/metalpoch/olt-blueprint/common/model"
	"github.com/metalpoch/olt-blueprint/report/usecase"
	"github.com/metalpoch/olt-blueprint/report/utils"
)

type ReportHandler struct {
	Usecase usecase.ReportUsecase
}

// Save Fats
//
//	@Summary		Save a file report
//	@Description	Save a file report into the  database
//	@Tags			fat
//	@Produce		json
//	@Success		201
//	@Failure		400	{object}	object{message=string}
//	@Failure		500	{object}	object{message=string}
//	@Router			/report/ [post]
func (hdlr ReportHandler) Add(c fiber.Ctx) error {
	newReport := new(model.NewReport)
	if err := c.Bind().MultipartForm(newReport); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	f, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	id, err := hdlr.Usecase.Add(newReport)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := c.SaveFile(f, path.Join("./report/data/", id)); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusCreated)
}

// Get a Fat
//
//	@Summary		Get report by id
//	@Description	get report by id from database
//	@Tags			fat
//	@Produce		json
//	@Success		200	{object}	model.ReportResponse
//	@Failure		400	{object}	object{message=string}
//	@Failure		500	{object}	object{message=string}
//	@Router			/report/:id [get]
func (hdlr ReportHandler) Get(c fiber.Ctx) error {
	id := c.Params("id")
	rp, err := hdlr.Usecase.Get(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(utils.ReportResponse(rp))
}
