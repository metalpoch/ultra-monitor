package validations

import (
	"time"

	"github.com/go-playground/validator/v10"
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

func NewValidator() *StructValidator {
	v := validator.New()
	v.RegisterValidation("dateformat", dateFormatValidation)
	return &StructValidator{Validator: v}
}
