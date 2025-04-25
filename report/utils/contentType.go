package utils

import (
	"errors"
	"slices"
)

var allowedTypes []string = []string{
	// office
	"application/msword",
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	"application/vnd.openxmlformats-officedocument.wordprocessingml.template",
	"application/vnd.ms-excel",
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	"application/vnd.openxmlformats-officedocument.spreadsheetml.template",
	"application/vnd.ms-powerpoint",
	"application/vnd.openxmlformats-officedocument.presentationml.presentation",
	"application/vnd.openxmlformats-officedocument.presentationml.template",
	"application/vnd.openxmlformats-officedocument.presentationml.slideshow",
	// libreoffice
	"application/vnd.openxmlformats-officedocument.presentationml.presentation",
	"application/vnd.ms-powerpoint.presentation.macroEnabled.12",
	"application/vnd.oasis.opendocument.text",
	"application/vnd.oasis.opendocument.spreadsheet",
	"application/vnd.oasis.opendocument.presentation",
	"application/vnd.oasis.opendocument.graphics",
	// csv
	"text/csv",
}

func IsValidReport(contentType string) error {
	if !slices.Contains(allowedTypes, contentType) {
		return errors.New("content type not allowed")
	}

	return nil
}
