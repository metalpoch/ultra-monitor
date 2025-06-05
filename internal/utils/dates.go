package utils

import (
	"time"
)

func DateRangeFromYear() (initDate, endDate time.Time) {
	loc, err := time.LoadLocation("America/Caracas")
	if err != nil {
		loc = time.FixedZone("UTC-4", -4*60*60)
	}

	now := time.Now().In(loc)
	endDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	initDate = endDate.AddDate(-1, 0, 0)

	return
}

func IsDateRangeWithin7Days(start, end time.Time) bool {
	if end.Before(start) {
		return false
	}
	return end.Sub(start).Hours() <= 7*24
}
