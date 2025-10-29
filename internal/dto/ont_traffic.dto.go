package dto

import "time"

type OntTraffic struct {
	Time        time.Time `json:"time"`
	BpsIn       float64   `json:"bps_in"`
	BpsOut      float64   `json:"bps_out"`
	BytesIn     float64   `json:"bytes_in"`
	BytesOut    float64   `json:"bytes_out"`
	Temperature int32     `json:"temperature"`
	Rx          float64   `json:"rx"`
	Tx          float64   `json:"tx"`
}

type OntTrafficCache struct {
	BytesIn   uint64    `json:"bytes_in"`
	BytesOut  uint64    `json:"bytes_out"`
	LastCheck time.Time `json:"last_check"`
}
