package handler

import (
	"net/url"
	"slices"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/metalpoch/ultra-monitor/internal/constants"
	"github.com/metalpoch/ultra-monitor/internal/dto"
	"github.com/metalpoch/ultra-monitor/usecase"
)

type FatHandler struct {
	Usecase *usecase.FatUsecase
}

func (hdlr *FatHandler) GetAll(c fiber.Ctx) error {
	var query dto.Find
	if err := c.Bind().Query(&query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.GetAll(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (hdlr *FatHandler) UpsertFats(c fiber.Ctx) error {
	var form dto.UploadFat
	if err := c.Bind().Form(&form); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if form.File.Header.Get("Content-Type") != "text/csv" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "content type must csv"})
	}

	f, err := form.File.Open()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	defer f.Close()

	date, _ := time.Parse("2006-01-02", form.Date)
	res, err := hdlr.Usecase.UpsertFats(f, date)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

func (hdlr *FatHandler) GetByID(c fiber.Ctx) error {
	id, err := fiber.Convert(c.Params("id"), strconv.Atoi)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	res, err := hdlr.Usecase.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (hdlr *FatHandler) GetAllByIP(c fiber.Ctx) error {
	ip, err := url.QueryUnescape(c.Params("ip"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.GetAllByIP(ip)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

func (hdlr *FatHandler) GetRegions(c fiber.Ctx) error {
	res, err := hdlr.Usecase.GetRegions()
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

func (hdlr *FatHandler) GetStates(c fiber.Ctx) error {
	res, err := hdlr.Usecase.GetStates()
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

func (hdlr *FatHandler) GetMunicipalities(c fiber.Ctx) error {
	state, err := url.QueryUnescape(c.Params("state"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.GetMunicipalities(state)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

func (hdlr *FatHandler) GetCounties(c fiber.Ctx) error {
	state, err := url.QueryUnescape(c.Params("state"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	municipality, err := url.QueryUnescape(c.Params("municipality"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.GetCounties(state, municipality)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

func (hdlr *FatHandler) GetODN(c fiber.Ctx) error {
	state, err := url.QueryUnescape(c.Params("state"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	municipality, err := url.QueryUnescape(c.Params("municipality"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	county, err := url.QueryUnescape(c.Params("county"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.GetODN(state, municipality, county)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

func (hdlr *FatHandler) FindByStates(c fiber.Ctx) error {
	state, err := url.QueryUnescape(c.Params("state"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var pag dto.Pagination
	if err := c.Bind().Query(&pag); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.FindByStates(state, pag)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

func (hdlr *FatHandler) FindByMunicipality(c fiber.Ctx) error {
	state, err := url.QueryUnescape(c.Params("state"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	municipality, err := url.QueryUnescape(c.Params("municipality"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var pag dto.Pagination
	if err := c.Bind().Query(&pag); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.FindByMunicipality(state, municipality, pag)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

func (hdlr *FatHandler) FindByCounty(c fiber.Ctx) error {
	state, err := url.QueryUnescape(c.Params("state"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	municipality, err := url.QueryUnescape(c.Params("municipality"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	county, err := url.QueryUnescape(c.Params("county"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var pag dto.Pagination
	if err := c.Bind().Query(&pag); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.FindByCounty(state, municipality, county, pag)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

func (hdlr *FatHandler) FindBytOdn(c fiber.Ctx) error {
	state, err := url.QueryUnescape(c.Params("state"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	municipality, err := url.QueryUnescape(c.Params("municipality"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	county, err := url.QueryUnescape(c.Params("county"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	odn, err := url.QueryUnescape(c.Params("odn"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var pag dto.Pagination
	if err := c.Bind().Query(&pag); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.FindBytOdn(state, municipality, county, odn, pag)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

func (hdlr *FatHandler) GetAllFatStatus(c fiber.Ctx) error {
	res, err := hdlr.Usecase.GetAllFatStatus()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (hdlr *FatHandler) GetAllFatStatusByRegion(c fiber.Ctx) error {
	region, err := url.QueryUnescape(c.Params("region"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	regions := make([]string, 0, len(constants.STATES_BY_REGION))
	for key := range constants.STATES_BY_REGION {
		regions = append(regions, key)
	}

	if !slices.Contains(regions, region) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "region invalid"})
	}

	res, err := hdlr.Usecase.GetAllFatStatusByRegion(region)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (hdlr *FatHandler) GetAllFatStatusByState(c fiber.Ctx) error {
	state, err := url.QueryUnescape(c.Params("state"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.GetAllFatStatusByState(state)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (hdlr *FatHandler) GetAllFatStatusByMunicipality(c fiber.Ctx) error {
	state, err := url.QueryUnescape(c.Params("state"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	municipality, err := url.QueryUnescape(c.Params("municipality"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.GetAllFatStatusByMunicipality(state, municipality)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (hdlr *FatHandler) GetAllFatStatusByCounty(c fiber.Ctx) error {
	state, err := url.QueryUnescape(c.Params("state"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	municipality, err := url.QueryUnescape(c.Params("municipality"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	county, err := url.QueryUnescape(c.Params("county"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.GetAllFatStatusByCounty(state, municipality, county)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (hdlr *FatHandler) GetAllFatStatusByODN(c fiber.Ctx) error {
	state, err := url.QueryUnescape(c.Params("state"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	municipality, err := url.QueryUnescape(c.Params("municipality"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	county, err := url.QueryUnescape(c.Params("county"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	odn, err := url.QueryUnescape(c.Params("odn"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.GetAllFatStatusByODN(state, municipality, county, odn)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (hdlr *FatHandler) GetAllFatStatusByFat(c fiber.Ctx) error {
	state, err := url.QueryUnescape(c.Params("state"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	municipality, err := url.QueryUnescape(c.Params("municipality"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	county, err := url.QueryUnescape(c.Params("county"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	odn, err := url.QueryUnescape(c.Params("odn"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	fat, err := url.QueryUnescape(c.Params("fat"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.GetAllFatStatusByFat(state, municipality, county, odn, fat)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (hdlr *FatHandler) GetFatStatusStateByRegion(c fiber.Ctx) error {
	region, err := url.QueryUnescape(c.Params("region"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.GetFatStatusStateByRegion(region)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}
func (hdlr *FatHandler) GetFatStatusOltByState(c fiber.Ctx) error {
	state, err := url.QueryUnescape(c.Params("state"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.GetFatStatusOltByState(state)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

func (hdlr *FatHandler) GetFatStatusGponByOlt(c fiber.Ctx) error {
	ip, err := url.QueryUnescape(c.Params("ip"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.GetFatStatusGponByOlt(ip)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

func (hdlr *FatHandler) GetFieldsOptions(c fiber.Ctx) error {
	field, err := url.QueryUnescape(c.Params("field"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := hdlr.Usecase.GetFieldsOptions(field)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

func (hdlr *FatHandler) DeleteByDate(c fiber.Ctx) error {
	dateStr := c.Params("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid date format, expected YYYY-MM-DD"})
	}

	err = hdlr.Usecase.DeleteByDate(date)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusOK)
}


