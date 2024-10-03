package handler

import (
	"github.com/metalpoch/olt-blueprint/common/usecase"
)

type DeviceHandler struct {
	Usecase usecase.DeviceUsecase
}

// // All Device
// //
// //	@Summary		Get all device
// //	@Description	get all device from database
// //	@Tags			device
// //	@Produce		json
// //	@Success		200	{object}	[]model.DeviceResponse
// //	@Failure		400	{object}	object{message=string}
// //	@Failure		500	{object}	object{message=string}
// //	@Router			/device/ [get]
// func (hdlr DeviceHandler) GetAll(c fiber.Ctx) error {
// 	devices, err := hdlr.Usecase.GetAll()
// 	if err != nil {
// 		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
// 	}
// 	var res []model.DeviceResponse
// 	for _, d := range devices {
// 		res = append(res, utils.DeviceResponse(d))
// 	}
// 	return c.JSON(res)
// }

// // Get a device
// //
// //	@Summary		Get a device
// //	@Description	get a device by ID from database
// //	@Tags			device
// //	@Produce		json
// //	@Param			id	path		uint	true	"Device ID"
// //	@Success		200	{object}	model.DeviceResponse
// //	@Failure		400	{object}	object{message=string}
// //	@Failure		500	{object}	object{message=string}
// //	@Router			/device/{id} [get]
// func (hdlr DeviceHandler) GetByID(c fiber.Ctx) error {
// 	id, err := fiber.Convert(c.Params("id"), strconv.Atoi)
// 	if err != nil {
// 		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
// 	}

// 	device, err := hdlr.Usecase.GetByID(uint(id))
// 	if err != nil {
// 		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
// 	}

// 	res := utils.DeviceResponse(device)

// 	return c.JSON(res)
// }
