package entity

import "time"

type Ont struct {
	ID          int32     `db:"id"`
	IP         string    `db:"ip"`
	OntIDX      string    `db:"ont_idx"`
	Serial      string    `db:"serial"`
	Despt       string    `db:"despt"`
	LineProf    string    `db:"line_prof"`
	Description string    `db:"description"`
	Enabled     bool      `db:"enabled"`
	Status      bool      `db:"status"`
	OltDistance int32     `db:"olt_distance"`
	LastCheck   time.Time `db:"last_check"`
	CreatedAt   time.Time `db:"created_at"`
}
