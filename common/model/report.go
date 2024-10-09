package model

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	"github.com/metalpoch/olt-blueprint/common/entity"
	"gorm.io/gorm"
)

type Report struct {
	ID               uuid.UUID
	Category         string
	OriginalFilename string
	ContentType      string
	Filepath         string
	User             entity.User
	UserID           uint
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt
}

type NewReport struct {
	Category         string `form:"category"`
	UserID           uint   `form:"user_id"`
	OriginalFilename string
	ContentType      string
	File             *multipart.FileHeader `form:"file"`
}

type FindReports struct {
	UserID   uint   `query:"user_id"`
	Category string `query:"category"`
}

type ReportResponse struct {
	Category         string
	OriginalFilename string
	ContentType      string
	Filepath         string
	User             UserLite
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        time.Time
}
