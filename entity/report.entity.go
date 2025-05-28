package entity

import (
	"time"

	"github.com/google/uuid"
)

type Report struct { // Table
	ID               uuid.UUID `db:"id"`
	Category         string    `db:"category"`
	OriginalFilename string    `db:"original_filename"`
	ContentType      string    `db:"content_type"`
	Basepath         string    `db:"basepath"`
	Filepath         string    `db:"filepath"`
	UserID           uint32    `db:"user_id"`
	CreatedAt        time.Time `db:"created_at"`
}

// func (rp *Report) BeforeCreate(tx *gorm.DB) (err error) {
// 	rp.ID = uuid.New()
// 	rp.Filepath = path.Join(rp.Basepath, rp.ID.String())
// 	return
// }
