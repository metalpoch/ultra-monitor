package handler

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/metalpoch/olt-blueprint/common/model"
	"github.com/metalpoch/olt-blueprint/core/usecase"
	"github.com/metalpoch/olt-blueprint/core/utils"
)

type InfoHandler struct {
	Usecase usecase.InfoUsecase
}

// Get device
//
//	@Summary		Get data of a device
//	@Description	get data of a device by ID from database
//	@Tags			info
//	@Produce		json
//	@Param			id	path		uint	true	"Device ID"
//	@Success		200	{object}	model.DeviceResponse
//	@Failure		400	{object}	object{message=string}
//	@Failure		500	{object}	object{message=string}
//	@Router			/info/device/{id} [get]
func (hdlr InfoHandler) GetDevice(c fiber.Ctx) error {
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

// Get device by IP
//
//	@Summary		Get data of a device
//	@Description	get data of a device by IP from database
//	@Tags			info
//	@Produce		json
//	@Param			id	path		uint	true	"Device ID"
//	@Success		200	{object}	model.DeviceResponse
//	@Failure		400	{object}	object{message=string}
//	@Failure		500	{object}	object{message=string}
//	@Router			/info/device/ip/{ip} [get]
func (hdlr InfoHandler) GetDeviceByIP(c fiber.Ctx) error {
	ip, err := url.QueryUnescape(c.Params("ip"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	device, err := hdlr.Usecase.GetDeviceByIP(ip)
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
//	@Tags			info
//	@Produce		json
//	@Success		200	{object}	[]model.DeviceLite
//	@Failure		400	{object}	object{message=string}
//	@Failure		500	{object}	object{message=string}
//	@Router			/info/device/ [get]
func (hdlr InfoHandler) GetAllDevice(c fiber.Ctx) error {
	res, err := hdlr.Usecase.GetAllDevice()
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

// Get devices by state
//
//	@Summary		Get data of all devices by state
//	@Description	get data of all devices by state from database
//	@Tags			info
//	@Produce		json
//	@Param			state	path		string					true	"State"
//	@Success		200		{object}	[]model.DeviceResponse
//	@Failure		400		{object}	object{message=string}
//	@Failure		500		{object}	object{message=string}
//	@Router			/info/device/location/{state} [get]
func (hdlr InfoHandler) GetDeviceByState(c fiber.Ctx) error {
	state, err := url.QueryUnescape(c.Params("state"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.GetDeviceByState(state)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

// Get devices by county
//
//	@Summary		Get data of all devices by county
//	@Description	get data of all devices by county from database
//	@Tags			info
//	@Produce		json
//	@Param			state	path		string					true	"State"
//	@Param			county	path		string					true	"County"
//	@Success		200		{object}	[]model.DeviceResponse
//	@Failure		400		{object}	object{message=string}
//	@Failure		500		{object}	object{message=string}
//	@Router			/info/device/location/{state}/{county} [get]
func (hdlr InfoHandler) GetDeviceByCounty(c fiber.Ctx) error {
	state, err := url.QueryUnescape(c.Params("state"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	county, err := url.QueryUnescape(c.Params("county"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.GetDeviceByCounty(state, county)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

// Get devices by municipality
//
//	@Summary		Get data of all devices by municipality
//	@Description	get data of all devices by municipality from database
//	@Tags			info
//	@Produce		json
//	@Param			state			path		string					true	"State"
//	@Param			county			path		string					true	"County"
//	@Param			municipality	path		string					true	"Municipality"
//	@Success		200				{object}	[]model.DeviceResponse
//	@Failure		400				{object}	object{message=string}
//	@Failure		500				{object}	object{message=string}
//	@Router			/info/device/location/{state}/{county}/{municipality} [get]
func (hdlr InfoHandler) GetDeviceByMunicipality(c fiber.Ctx) error {
	state, err := url.QueryUnescape(c.Params("state"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	county, err := url.QueryUnescape(c.Params("county"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	municipality, err := url.QueryUnescape(c.Params("municipality"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.GetDeviceByMunicipality(state, county, municipality)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

// Get interface
//
//	@Summary		Get all data of a interface
//	@Description	get all data of a interface by ID from database
//	@Tags			info
//	@Produce		json
//	@Param			id	path		uint	true	"Interface ID"
//	@Success		200	{object}	model.InterfaceResponse
//	@Failure		400	{object}	object{message=string}
//	@Failure		500	{object}	object{message=string}
//	@Router			/info/interface/{id} [get]
func (hdlr InfoHandler) GetInterface(c fiber.Ctx) error {
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
//	@Tags			info
//	@Produce		json
//	@Param			id	path		uint	true	"Device ID"
//	@Success		200	{object}	[]model.InterfaceWithoutDevice
//	@Failure		400	{object}	object{message=string}
//	@Failure		500	{object}	object{message=string}
//	@Router			/info/interface/device/{id} [get]
func (hdlr InfoHandler) GetInterfacesByDevice(c fiber.Ctx) error {
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

// Get interfaces of a device
//
//	@Summary		Get all interfaces from a device
//	@Description	get all interfaces data by device ID from database
//	@Tags			info
//	@Produce		json
//	@Param			path			uint				true	"Device ID"
//	@Param			shell	query		uint8				true	"Shell GPON"
//	@Param			card	query		uint8					"Card GPON"
//	@Param			port	query		uint8					"PORT GPON"
//	@Success		200	{object}	[]model.InterfaceWithoutDevice
//	@Failure		400	{object}	object{message=string}
//	@Failure		500	{object}	object{message=string}
//	@Router			/info/interface/device/{id}/find [get]
func (hdlr InfoHandler) GetInterfacesByDeviceAndPorts(c fiber.Ctx) error {
	id, err := fiber.Convert(c.Params("id"), strconv.Atoi)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	query := new(model.GponPort)
	if err := c.Bind().Query(query); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	res, err := hdlr.Usecase.GetInterfacesByDeviceAndPorts(uint(id), query.Shell, query.Card, query.Port)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

// Get states
//
//	@Summary		Get all state
//	@Description	get all state  from database
//	@Tags			info
//	@Produce		json
//	@Success		200	{object}	[]string
//	@Failure		400	{object}	object{message=string}
//	@Failure		500	{object}	object{message=string}
//	@Router			/info/location/state [get]
func (hdlr InfoHandler) GetLocationStates(c fiber.Ctx) error {
	res, err := hdlr.Usecase.GetLocationStates()
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

// Get states
//
// @Summary		Get all conties from state
// @Description	get all conties from state from database
// @Tags		info
// @Produce		json
// @Param		state	string		string					true	"State"
// @Success		200		{object}	[]string
// @Failure		400		{object}	object{message=string}
// @Failure		500		{object}	object{message=string}
// @Router		/info/location/{state} [get]
func (hdlr InfoHandler) GetLocationCounties(c fiber.Ctx) error {
	state, err := url.QueryUnescape(c.Params("state"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.GetLocationCounties(state)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

// Get municipalities
//
// @Summary		Get all municipalities from conty & state
// @Description	get all municipalities from conty & state from database
// @Tags		info
// @Produce		json
// @Param		state	string		string					true	"State"
// @Param		county	string		string					true	"County"
// @Success		200		{object}	[]string
// @Failure		400		{object}	object{message=string}
// @Failure		500		{object}	object{message=string}
// @Router		/info/location/{state}/{county} [get]
func (hdlr InfoHandler) GetLocationMunicipalities(c fiber.Ctx) error {
	state, err := url.QueryUnescape(c.Params("state"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	county, err := url.QueryUnescape(c.Params("county"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.GetLocationMunicipalities(state, county)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}
