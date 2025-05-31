package dto

import "mime/multipart"

type NewReport struct {
	Category         string `form:"category" validate:"required"`
	UserID           int32  `form:"user_id" validate:"required"`
	Basepath         string
	OriginalFilename string
	ContentType      string
	File             *multipart.FileHeader `form:"file" validate:"required"`
}
