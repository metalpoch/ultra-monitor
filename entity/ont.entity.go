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
	LastCheck   time.Time `db:"last_check"`
	CreatedAt   time.Time `db:"created_at"`
}

type OntTraffic struct {
	OntID       int32     `db:"ont_id"`
	Time        time.Time `db:"time"`
	BpsIn       float64   `db:"bps_in"`
	BpsOut      float64   `db:"bps_out"`
	BytesIn     float64   `db:"bytes_in"`
	BytesOut    float64   `db:"bytes_out"`
	Temperature int32     `db:"temperature"`
	Rx          int32     `db:"rx"`
	Tx          int32     `db:"tx"`
}
