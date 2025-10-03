package dto

import "time"

type RangeDate struct {
	InitDate  time.Time `query:"initDate" validate:"required"`
	FinalDate time.Time `query:"finalDate" validate:"required"`
}

type Pagination struct {
	Page  uint16 `query:"page" validate:"required"`
	Limit uint16 `query:"limit" validate:"required"`
}

type Find struct {
	Field string `query:"field" validate:"required_with=Value"`
	Value string `query:"value" validate:"required_with=Field"`
	Pagination
}

type LocationHierarchy struct {
	Regions []string          `json:"regions"`
	States  map[string][]string `json:"states"`
	Olts    map[string][]OltInfo `json:"olts"`
}

type OltInfo struct {
	IP      string `json:"ip"`
	SysName string `json:"sys_name"`
}
