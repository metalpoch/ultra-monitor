package dto

import "time"

type OntTrafficCache struct {
	BytesIn     uint64    `json:"bytes_in"`
	BytesOut    uint64    `json:"bytes_out"`
	LastCheck   time.Time `json:"last_check"`
}

