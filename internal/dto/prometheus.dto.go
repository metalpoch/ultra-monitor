package dto

import "time"

type Prometheus struct {
	Region string `json:"region"`
	State  string `json:"state"`
	IP     string `json:"ip"`
	IDX    int64  `json:"idx"`
	Shell  uint8  `json:"shell"`
	Card   uint8  `json:"card"`
	Port   uint8  `json:"port"`
}

type DeviceLocation struct {
	Region  string `json:"region"`
	State   string `json:"state"`
	IP      string `json:"ip"`
	SysName string `json:"sysName"`
}

type InfoDevice struct {
	Region  string `json:"region"`
	State   string `json:"state"`
	IP      string `json:"ip"`
	IfName  string `json:"if_name"`
	IfIndex int64  `json:"if_index"`
}

type traffic struct {
	Time        time.Time `json:"time"`
	Description string    `json:"description"`
	BpsIn       float64   `json:"bps_in"`
	BpsOut      float64   `json:"bps_out"`
	Bandwidth   float64   `json:"bw_mbps"`
	BytesIn     float64   `json:"bytes_in"`
	BytesOut    float64   `json:"bytes_out"`
}

type Traffic map[string][]*traffic

type GroupedIP struct {
	IP []string `query:"ip" validate:"required"`
	RangeDate
}
