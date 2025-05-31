package handler

import (
	"net/url"
	"path"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/metalpoch/ultra-monitor/internal/dto"
	"github.com/metalpoch/ultra-monitor/internal/utils"
	"github.com/metalpoch/ultra-monitor/usecase"
)

type ReportHandler struct {
	Usecase         *usecase.ReportUsecase
	ReportDirectory string
}

func (hdlr *ReportHandler) Add(c fiber.Ctx) error {
	var newReport dto.NewReport
	if err := c.Bind().Form(newReport); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	f, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := utils.IsValidReport(f.Header.Get("Content-Type")); err != nil {
		return c.Status(fiber.StatusUnsupportedMediaType).JSON(fiber.Map{"error": err.Error()})
	}

	newReport.Basepath = hdlr.ReportDirectory
	newReport.File = f

	id, err := hdlr.Usecase.Add(newReport)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := c.SaveFile(f, path.Join(hdlr.ReportDirectory, id)); err != nil {
		hdlr.Usecase.Delete(id)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (hdlr *ReportHandler) Get(c fiber.Ctx) error {
	id := c.Params("id")
	res, err := hdlr.Usecase.Get(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

func (hdlr *ReportHandler) GetCategories(c fiber.Ctx) error {
	rp, err := hdlr.Usecase.GetCategories()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(rp)
}

func (hdlr ReportHandler) GetReportsByUser(c fiber.Ctx) error {
	userID, err := fiber.Convert(c.Params("userID"), strconv.Atoi)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var pag dto.Pagination
	if err := c.Bind().Query(pag); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.GetReportsByUser(userID, pag)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (hdlr ReportHandler) GetReportsByCategory(c fiber.Ctx) error {
	category, err := url.QueryUnescape(c.Params("category"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var pag dto.Pagination
	if err := c.Bind().Query(pag); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.GetReportsByCategory(category, pag)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}
