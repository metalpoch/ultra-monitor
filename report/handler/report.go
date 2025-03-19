package handler

import (
	"fmt"
	"path"
	"path/filepath"

	"github.com/gofiber/fiber/v3"
	"github.com/metalpoch/olt-blueprint/common/constants"
	"github.com/metalpoch/olt-blueprint/common/model"
	"github.com/metalpoch/olt-blueprint/report/usecase"
	"github.com/metalpoch/olt-blueprint/report/utils"
)

type ReportHandler struct {
	Usecase usecase.ReportUsecase
}

// Save Report
//
//	@Summary		Save a file report
//	@Description	Save a file report into the database
//	@Tags			report
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

	if err := utils.IsValidReport(f.Header.Get("Content-Type")); err != nil {
		return c.Status(fiber.StatusUnsupportedMediaType).JSON(fiber.Map{"error": err.Error()})
	}
	newReport.File = f
	id, err := hdlr.Usecase.Add(newReport)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := c.SaveFile(f, path.Join(constants.BASE_FILEPATH, id)); err != nil {
		hdlr.Usecase.Delete(id)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusCreated)
}

// Get a Report
//
//	@Summary		Get report by id
//	@Description	get report by id from database
//	@Tags			report
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

// Get all categories
//
//	@Summary		Get all report categories
//	@Description	get all report categories from database
//	@Tags			report
//	@Produce		json
//	@Success		200	{object}	[]string
//	@Failure		400	{object}	object{message=string}
//	@Failure		500	{object}	object{message=string}
//	@Router			/report/categories [get]
func (hdlr ReportHandler) GetCategories(c fiber.Ctx) error {
	rp, err := hdlr.Usecase.GetCategories()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(rp)
}

// Get reports
//
//	@Summary		Get reports by category and/or user id
//	@Description	get reports by category and/or user id from database
//	@Tags			report
//	@Produce		json
//	@Success		200	{object}	[]model.ReportResponse
//	@Failure		400	{object}	object{message=string}
//	@Failure		500	{object}	object{message=string}
//	@Router			/reports/ [get]
func (hdlr ReportHandler) GetReports(c fiber.Ctx) error {
	query := new(model.FindReports)
	if err := c.Bind().Query(query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if query.Category == "" && query.UserID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "missing fields"})
	}

	res, err := hdlr.Usecase.GetReports(query)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

// Download Report
//
//	@Summary		Download report by id
//	@Description	download report by id from database
//	@Tags			report
//	@Produce		json
//	@Success		200	{object}	model.ReportResponse
//	@Failure		400	{object}	object{message=string}
//	@Failure		500	{object}	object{message=string}
//	@Router			/report/download/:id [get]
func (hdlr ReportHandler) Download(c fiber.Ctx) error {
	id := c.Params("id")
	rp, err := hdlr.Usecase.Get(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	filename := fmt.Sprintf("%s-%v%s", rp.Category, rp.CreatedAt.Format("20060102150405"), filepath.Ext(rp.OriginalFilename))
	return c.Download(rp.Filepath, filename)
}
