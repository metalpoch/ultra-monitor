package validations

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/metalpoch/ultra-monitor/internal/dto"
)

type StructValidator struct {
	Validator *validator.Validate
}

func (v *StructValidator) Validate(out any) error {
	return v.Validator.Struct(out)
}

func dateFormatValidation(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()
	_, err := time.Parse("2006-01-02", dateStr)
	return err == nil
}

func validateSameLength(fl validator.FieldLevel) bool {
	obj, ok := fl.Parent().Interface().(dto.PrometheusDeviceQuery)
	if !ok {
		return false
	}

	length := len(obj.Shell)
	return length == len(obj.Card) && length == len(obj.Port)
}

func NewValidator() *StructValidator {
	v := validator.New()
	v.RegisterValidation("dateformat", dateFormatValidation)
	v.RegisterValidation("sameLength", validateSameLength)
	return &StructValidator{Validator: v}
}
