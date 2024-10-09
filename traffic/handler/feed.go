package handler

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/metalpoch/olt-blueprint/traffic/usecase"
	"github.com/metalpoch/olt-blueprint/traffic/utils"
)

type FeedHandler struct {
	Usecase usecase.FeedUsecase
}

// Get device
//
//	@Summary		Get data of a device
//	@Description	get data of a device by ID from database
//	@Tags			feed
//	@Produce		json
//	@Param			id	path		uint	true	"Device ID"
//	@Success		200	{object}	model.DeviceResponse
//	@Failure		400	{object}	object{message=string}
//	@Failure		500	{object}	object{message=string}
//	@Router			/feed/device/{id} [get]
func (hdlr FeedHandler) GetDevice(c fiber.Ctx) error {
	id, err := fiber.Convert(c.Params("id"), strconv.Atoi)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	device, err := hdlr.Usecase.GetDevice(uint(id))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res := utils.DeviceResponse(device)

	return c.JSON(res)
}

// Get all devices
//
//	@Summary		Get all devices
//	@Description	get all device ID from database
//	@Tags			feed
//	@Produce		json
//	@Success		200	{object}	[]model.DeviceLite
//	@Failure		400	{object}	object{message=string}
//	@Failure		500	{object}	object{message=string}
//	@Router			/feed/device/ [get]
func (hdlr FeedHandler) GetAllDevice(c fiber.Ctx) error {
	res, err := hdlr.Usecase.GetAllDevice()
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

// Get interface
//
//	@Summary		Get all data of a interface
//	@Description	get all data of a interface by ID from database
//	@Tags			feed
//	@Produce		json
//	@Param			id	path		uint	true	"Interface ID"
//	@Success		200	{object}	model.InterfaceResponse
//	@Failure		400	{object}	object{message=string}
//	@Failure		500	{object}	object{message=string}
//	@Router			/feed/interface/{id} [get]
func (hdlr FeedHandler) GetInterface(c fiber.Ctx) error {
	id, err := fiber.Convert(c.Params("id"), strconv.Atoi)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	iface, err := hdlr.Usecase.GetInterface(uint(id))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	res := utils.InterfaceResponse(iface)

	return c.JSON(res)
}

// Get interfaces of a device
//
//	@Summary		Get all interfaces from a device
//	@Description	get all interfaces data by device ID from database
//	@Tags			feed
//	@Produce		json
//	@Param			id	path		uint	true	"Device ID"
//	@Success		200	{object}	[]model.InterfaceWithoutDevice
//	@Failure		400	{object}	object{message=string}
//	@Failure		500	{object}	object{message=string}
//	@Router			/feed/interface/device/{id} [get]
func (hdlr FeedHandler) GetInterfacesByDevice(c fiber.Ctx) error {
	id, err := fiber.Convert(c.Params("id"), strconv.Atoi)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	res, err := hdlr.Usecase.GetInterfacesByDevice(uint(id))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}
