package router

import (
	"github.com/gofiber/fiber/v3"
	"github.com/metalpoch/olt-blueprint/common/pkg/tracking"
	"github.com/metalpoch/olt-blueprint/report/handler"
	"github.com/metalpoch/olt-blueprint/report/usecase"
	"gorm.io/gorm"
)

func newReportRouter(server *fiber.App, db *gorm.DB, telegram tracking.SmartModule) {
	hdlr := handler.ReportHandler{
		Usecase: *usecase.NewReportUsecase(db, telegram),
	}

	server.Post("/report/", hdlr.Add)
	server.Get("/reports/", hdlr.GetReports)
	server.Get("/report/categories/", hdlr.GetCategories)
	server.Get("/report/:id", hdlr.Get)
	server.Get("/report/download/:id", hdlr.Download)
}
