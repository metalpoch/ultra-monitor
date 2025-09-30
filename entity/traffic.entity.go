package entity

import "time"

type OltIndex struct {
	IP  string `db:"ip"`
	Idx string `db:"idx"`
}

type SumaryTraffic struct {
	Time     time.Time `db:"time"`
	IP       string    `db:"ip"`
	State    string    `db:"state"`
	Region   string    `db:"region"`
	BpsIn    float64   `db:"bps_in"`
	BpsOut   float64   `db:"bps_out"`
	BytesIn  float64   `db:"bytes_in"`
	BytesOut float64   `db:"bytes_out"`
}
