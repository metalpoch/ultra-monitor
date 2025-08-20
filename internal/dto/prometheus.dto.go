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
	Status int8   `json:"status"`
}

type PrometheusDeviceQuery struct {
	IP    string   `json:"ip"`
	Shell []uint32 `json:"shell" validate:"required,min=1,sameLength"`
	Card  []uint32 `json:"card"`
	Port  []uint32 `json:"port"`
}

type PrometheusPortStatus struct {
	Olts          uint32 `json:"olts"`
	Cards         uint32 `json:"cards"`
	GponActives   uint32 `json:"gpon_actives"`
	GponInactives uint32 `json:"gpon_inactives"`
	TotalGpon     uint32 `json:"total_gpon"`
}

type DeviceLocation struct {
	Region  string `json:"region"`
	State   string `json:"state"`
	IP      string `json:"ip"`
	SysName string `json:"sysName"`
}

type InfoDevice struct {
	Region       string `json:"region"`
	State        string `json:"state"`
	IP           string `json:"ip"`
	IfName       string `json:"if_name"`
	IfIndex      int64  `json:"if_index"`
	IfOperStatus int8   `json:"if_oper_status"`
}

type Traffic struct {
	Time      time.Time `json:"time"`
	BpsIn     float64   `json:"bps_in"`
	BpsOut    float64   `json:"bps_out"`
	Bandwidth float64   `json:"bw_mbps"`
	BytesIn   float64   `json:"bytes_in"`
	BytesOut  float64   `json:"bytes_out"`
}

type TrafficByLabel map[string][]Traffic

type GroupedIP struct {
	IP []string `query:"ip" validate:"required"`
	RangeDate
}
