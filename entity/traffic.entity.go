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
	Sysname  string    `db:"sysname"`
	BpsIn    float64   `db:"bps_in"`
	BpsOut   float64   `db:"bps_out"`
	BytesIn  float64   `db:"bytes_in"`
	BytesOut float64   `db:"bytes_out"`
}

type TrafficSummary struct {
	Time          time.Time `db:"time"`
	Sysname       string    `db:"sysname"`
	TotalBpsIn    float64   `db:"total_bps_in"`
	TotalBpsOut   float64   `db:"total_bps_out"`
	TotalBytesIn  float64   `db:"total_bytes_in"`
	TotalBytesOut float64   `db:"total_bytes_out"`
}

type TrafficByRegion struct {
	Region        string    `db:"region"`
	Time          time.Time `db:"time"`
	TotalBpsIn    float64   `db:"total_bps_in"`
	TotalBpsOut   float64   `db:"total_bps_out"`
	TotalBytesIn  float64   `db:"total_bytes_in"`
	TotalBytesOut float64   `db:"total_bytes_out"`
}

type TrafficByState struct {
	State         string    `db:"state"`
	Time          time.Time `db:"time"`
	TotalBpsIn    float64   `db:"total_bps_in"`
	TotalBpsOut   float64   `db:"total_bps_out"`
	TotalBytesIn  float64   `db:"total_bytes_in"`
	TotalBytesOut float64   `db:"total_bytes_out"`
}

type TrafficBySysname struct {
	Sysname       string    `db:"sysname"`
	Time          time.Time `db:"time"`
	TotalBpsIn    float64   `db:"total_bps_in"`
	TotalBpsOut   float64   `db:"total_bps_out"`
	TotalBytesIn  float64   `db:"total_bytes_in"`
	TotalBytesOut float64   `db:"total_bytes_out"`
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
