package model

type Page struct {
	Num  int `query:"page" validate:"gt=0"`
	Size int `query:"page_size" validate:"gt=0"`
}
