package entity

import (
	"time"

	"github.com/google/uuid"
)

type Report struct {
	ID               uuid.UUID `db:"id"`
	Category         string    `db:"category"`
	OriginalFilename string    `db:"original_filename"`
	ContentType      string    `db:"content_type"`
	Basepath         string    `db:"basepath"`
	Filepath         string    `db:"filepath"`
	UserID           int32     `db:"user_id"`
	CreatedAt        time.Time `db:"created_at"`
}
