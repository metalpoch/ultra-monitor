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

type PonPort struct {
	IDX    int64 `json:"idx"`
	Shell  uint8 `json:"shell"`
	Card   uint8 `json:"card"`
	Port   uint8 `json:"port"`
	Status int8  `json:"status"`
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
	Switch       string `json:"switch"`
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

type Volume struct {
	BytesIn  string `json:"in"`
	BytesOut string `json:"out"`
}

type GponStats struct {
	Port      string  `json:"port"`
	IfName    string  `json:"if_name"`
	IfSpeed   float64 `json:"if_speed"`
	MaxInBps  float64 `json:"max_in_bps"`
	AvgInBps  float64 `json:"avg_in_bps"`
	MaxOutBps float64 `json:"max_out_bps"`
	AvgOutBps float64 `json:"avg_out_bps"`
	UsageIn   float64 `json:"usage_in"`
	UsageOut  float64 `json:"usage_out"`
}

type StateStats struct {
	State     string  `json:"state"`
	IfSpeed   float64 `json:"if_speed"`
	MaxInBps  float64 `json:"max_in_bps"`
	AvgInBps  float64 `json:"avg_in_bps"`
	MaxOutBps float64 `json:"max_out_bps"`
	AvgOutBps float64 `json:"avg_out_bps"`
	UsageIn   float64 `json:"usage_in"`
	UsageOut  float64 `json:"usage_out"`
}

type OltStats struct {
	Instance  string  `json:"ip"`
	SysName   string  `json:"sys_name"`
	Switch    string  `json:"switch"`
	IfSpeed   float64 `json:"if_speed"`
	MaxInBps  float64 `json:"max_in_bps"`
	AvgInBps  float64 `json:"avg_in_bps"`
	MaxOutBps float64 `json:"max_out_bps"`
	AvgOutBps float64 `json:"avg_out_bps"`
	UsageIn   float64 `json:"usage_in"`
	UsageOut  float64 `json:"usage_out"`
}
