package dto

import (
	"mime/multipart"
	"time"
)

type Fat struct {
	ID           int32     `json:"id"`
	IP           string    `json:"ip"`
	Region       string    `json:"region"`
	State        string    `json:"state"`
	Municipality string    `json:"municipality"`
	County       string    `json:"county"`
	Odn          string    `json:"odn"`
	Fat          string    `json:"fat"`
	Shell        uint8     `json:"shell"`
	Card         uint8     `json:"card"`
	Port         uint8     `json:"port"`
	CreatedAt    time.Time `json:"created_at"`
}

type UpsertFat struct {
	IP                 string    `json:"IP_OLT"`
	Region             string    `json:"REGION"`
	State              string    `json:"ESTADO"`
	Municipality       string    `json:"MUNICIPIO"`
	County             string    `json:"PARROQUIA"`
	Odn                string    `json:"ODN"`
	Fat                string    `json:"FAT"`
	Shell              uint8     `json:"SHELL"`
	Card               uint8     `json:"SLOT"`
	Port               uint8     `json:"PUERTO"`
	Actives            uint32    `json:"ACTIVOS"`
	ProvisionedOffline uint32    `json:"APROVISIONADO_OFFLINE"`
	CutOff             uint32    `json:"CORTADO"`
	InProgress         uint32    `json:"EN_PROCESO"`
	Date               time.Time `json:"date"`
}

type UploadFat struct {
	File *multipart.FileHeader `form:"file" validate:"required"`
	Date string                `form:"date" validate:"required,dateformat"`
}
