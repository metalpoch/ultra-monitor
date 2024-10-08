package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Report struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	Category  string
	Filename  string
	Filepath  string
	User      User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UserID    uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (rp *Report) BeforeCreate(tx *gorm.DB) (err error) {
	rp.ID = uuid.New()
	return
}
