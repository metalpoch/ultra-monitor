package model

import (
	"time"

	"github.com/google/uuid"
)

type Report struct {
	ID               uuid.UUID `json:"id"`
	Category         string    `json:"category"`
	OriginalFilename string    `json:"original_filename"`
	ContentType      string    `json:"content_type"`
	Basepath         string    `json:"basepath"`
	Filepath         string    `json:"filepath"`
	UserID           int32     `json:"user_id"`
	CreatedAt        time.Time `json:"created_at"`
}
