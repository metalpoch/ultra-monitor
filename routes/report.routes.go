package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/handler"
	"github.com/metalpoch/ultra-monitor/internal/cache"
	"github.com/metalpoch/ultra-monitor/usecase"
)

func NewReportRoutes(app *fiber.App, db *sqlx.DB, cache *cache.Redis, reportsDir string) {
	hdlr := &handler.ReportHandler{
		Usecase:         usecase.NewReportUsecase(db),
		ReportDirectory: reportsDir,
	}

	route := app.Group("/api/report")
	route.Post("/", hdlr.Add)
	route.Get("/file/:id", hdlr.Get)
	route.Delete("/file/:id", hdlr.Delete)

	route.Get("/categories", hdlr.GetCategories)
	route.Get("/category/:category", hdlr.GetReportsByCategory)

	route.Get("/user/:id", hdlr.GetReportsByUser)
}
