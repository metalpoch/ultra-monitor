package router

// import (
// 	"github.com/gofiber/fiber/v3"
// 	"github.com/metalpoch/olt-blueprint/traffic/handler"
// 	"github.com/metalpoch/olt-blueprint/traffic/usecase"
// 	"gorm.io/gorm"
// )

// func newDeviceRouter(server *fiber.App, db *gorm.DB, secret []byte) {
// 	hdlr := handler.TrafficHandler{
// 		Usecase: *usecase.NewTrafficUsecase(db, secret),
// 	}

// 	server.Get("/interface/:id", hdlr.GetByInterface)
// 	server.Get("/device/:id", hdlr.GetByDevice)
// }
