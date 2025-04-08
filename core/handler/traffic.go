package handler

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/metalpoch/olt-blueprint/common/model"
	"github.com/metalpoch/olt-blueprint/core/usecase"
)

type TrafficHandler struct {
	Usecase usecase.TrafficUsecase
}

// Traffic by interface
//
// @Summary		Get a array of interface traffic in a range date
// @Description	get traffic by interface ID
// @Tags		traffic
// @Produce		json
// @Param		id		path		uint					true	"Interface ID"
// @Param		init_date	query		time.Time				true	"Start date for traffic range"
// @Param		end_date	query		time.Time				true	"End date for traffic range"
// @Success		200			{object}	[]model.Traffic
// @Failure		400			{object}	object{message=string}
// @Failure		500			{object}	object{message=string}
// @Router		/traffic/interface/{id} [get]
func (hdlr TrafficHandler) GetByInterface(c fiber.Ctx) error {
	id, err := fiber.Convert(c.Params("id"), strconv.Atoi)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	query := new(model.TranficRangeDate)
	if err := c.Bind().Query(query); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	res, err := hdlr.Usecase.GetTrafficByInterface(uint(id), query)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

// Traffic by device
//
// @Summary		Get a array of device traffic in a range date
// @Description	get traffic by device ID
// @Tags		traffic
// @Produce		json
// @Param		id			path		uint					true	"Device ID"
// @Param		init_date	query		time.Time				true	"Start date for traffic range"
// @Param		end_date	query		time.Time				true	"End date for traffic range"
// @Success		200			{object}	[]model.Traffic
// @Failure		400			{object}	object{message=string}
// @Failure		500			{object}	object{message=string}
// @Router		/traffic/device/{id} [get]
func (hdlr TrafficHandler) GetByDevice(c fiber.Ctx) error {
	id, err := fiber.Convert(c.Params("id"), strconv.Atoi)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	query := new(model.TranficRangeDate)
	if err := c.Bind().Query(query); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.GetTrafficByDevice(uint(id), query)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

// Traffic by fat
//
// @Summary		Get a array of fat traffic in a range date
// @Description	get traffic by fat ID
// @Tags		traffic
// @Produce		json
// @Param		id			path		uint					true	"Fat ID"
// @Param		init_date	query		time.Time				true	"Start date for traffic range"
// @Param		end_date	query		time.Time				true	"End date for traffic range"
// @Success		200			{object}	[]model.Traffic
// @Failure		400			{object}	object{message=string}
// @Failure		500			{object}	object{message=string}
// @Router		/traffic/fat/{id} [get]
func (hdlr TrafficHandler) GetByFat(c fiber.Ctx) error {
	id, err := fiber.Convert(c.Params("id"), strconv.Atoi)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	query := new(model.TranficRangeDate)
	if err := c.Bind().Query(query); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.GetTrafficByFat(uint(id), query)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

// Traffic by state
//
// @Summary		Get a array of state traffic in a range date
// @Description	get traffic by state
// @Tags		traffic
// @Produce		json
// @Param		state		string		string					true	"State"
// @Param		init_date	query		time.Time				true	"Start date for traffic range"
// @Param		end_date	query		time.Time				true	"End date for traffic range"
// @Success		200			{object}	[]model.Traffic
// @Failure		400			{object}	object{message=string}
// @Failure		500			{object}	object{message=string}
// @Router		/traffic/location/{state} [get]
func (hdlr TrafficHandler) GetByState(c fiber.Ctx) error {
	state, err := url.QueryUnescape(c.Params("state"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	query := new(model.TranficRangeDate)
	if err := c.Bind().Query(query); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.GetTrafficByState(state, query)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

// Traffic by county
//
// @Summary		Get a array of county traffic in a range date
// @Description	get traffic by county
// @Tags		traffic
// @Produce		json
// @Param		state		string		string					true	"State"
// @Param		county		path		string					true	"County"
// @Param		init_date	query		time.Time				true	"Start date for traffic range"
// @Param		end_date	query		time.Time				true	"End date for traffic range"
// @Success		200			{object}	[]model.Traffic
// @Failure		400			{object}	object{message=string}
// @Failure		500			{object}	object{message=string}
// @Router		/traffic/location/{state}/{county} [get]
func (hdlr TrafficHandler) GetByCounty(c fiber.Ctx) error {
	state, err := url.QueryUnescape(c.Params("state"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	county, err := url.QueryUnescape(c.Params("county"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	query := new(model.TranficRangeDate)
	if err := c.Bind().Query(query); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.GetTrafficByCounty(state, county, query)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

// Traffic by municipality
//
// @Summary		Get a array of municipality traffic in a range date
// @Description	get traffic by municipality
// @Tags		traffic
// @Produce		json
// @Param		id			path		uint					true	"Location ID"
// @Param		init_date	query		time.Time				true	"Start date for traffic range"
// @Param		end_date	query		time.Time				true	"End date for traffic range"
// @Success		200			{object}	[]model.Traffic
// @Failure		400			{object}	object{message=string}
// @Failure		500			{object}	object{message=string}
// @Router		/traffic/location/{state}/{county}/{municipality} [get]
func (hdlr TrafficHandler) GetByMunicipaly(c fiber.Ctx) error {
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

	query := new(model.TranficRangeDate)
	if err := c.Bind().Query(query); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.GetTrafficByMunicipality(state, county, municipality, query)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

// Traffic by ODN
//
// @Summary		Get a array of ODN traffic in a range date
// @Description	get traffic by ODN name
// @Tags		traffic
// @Produce		json
// @Param		odn		path		string					true	"ODN"
// @Param		init_date	query		time.Time				true	"Start date for traffic range"
// @Param		end_date	query		time.Time				true	"End date for traffic range"
// @Success		200			{object}	[]model.Traffic
// @Failure		400			{object}	object{message=string}
// @Failure		500			{object}	object{message=string}
// @Router		/traffic/interface/{id} [get]
func (hdlr TrafficHandler) GetTrafficByODN(c fiber.Ctx) error {
	odn, err := url.QueryUnescape(c.Params("odn"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	query := new(model.TranficRangeDate)
	if err := c.Bind().Query(query); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	res, err := hdlr.Usecase.GetTrafficByODN(odn, query)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

// Traffic by State
//
// @Summary		Get a array of state traffic per device
// @Description	get traffic by state per device
// @Tags		traffic
// @Produce		json
// @Success		200			{object}	[]model.TrafficState
// @Failure		400			{object}	object{message=string}
// @Failure		500			{object}	object{message=string}
// @Router		/traffic/interface/{id} [get]
func (hdlr TrafficHandler) GetTotalTrafficByState(c fiber.Ctx) error {
	month, err := url.QueryUnescape(c.Params("month"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	res, err := hdlr.Usecase.GetTotalTrafficByState(month)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

// Traffic by State
//
// @Summary		Get a array of state traffic per device
// @Description	get traffic by state per device
// @Tags		traffic
// @Produce		json
// @Success		200			{object}	[]model.TrafficState
// @Failure		400			{object}	object{message=string}
// @Failure		500			{object}	object{message=string}
// @Router		/traffic/interface/{id} [get]
func (hdlr TrafficHandler) GetTotalTrafficByState_N(c fiber.Ctx) error {
	n, err := url.QueryUnescape(c.Params("n"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	n_int, err := strconv.ParseInt(n, 10, 8)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	query := new(model.TranficRangeDate)
	if err := c.Bind().Query(query); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	res, err := hdlr.Usecase.GetTotalTrafficByState_N(query, int8(n_int))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}
