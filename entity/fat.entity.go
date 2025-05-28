package entity

import (
	"time"
)

type Fat struct { // Table. TODO: agregar campos de location para eliminar esa tabla
	ID         uint      `db:"id"`
	LocationID uint      `db:"location_id"`
	Splitter   uint8     `db:"splitter"`
	Latitude   float64   `db:"latitude"`
	Longitude  float64   `db:"longitude"`
	Address    string    `db:"address"`
	Fat        string    `db:"fat"`
	ODN        string    `db:"odn"`
	CreatedAt  time.Time `db:"created_at"`
}
