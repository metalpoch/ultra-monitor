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
