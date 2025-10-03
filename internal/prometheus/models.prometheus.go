package prometheus

import "time"

type dataProm struct {
	Labels map[string]string
	Value  float64
	Time   time.Time
}

type TrafficQuery struct {
	InitDate  time.Time
	FinalDate time.Time
	Region    string
	State     string
	Instance  string
	PortIndex string
}

type Traffic struct {
	Time      time.Time
	BpsIn     float64
	BpsOut    float64
	Bandwidth float64
	BytesIn   float64
	BytesOut  float64
}

type InfoDevice struct {
	Region       string
	State        string
	IP           string
	IfName       string
	IfIndex      int64
	IfOperStatus int8
}

type DeviceLocation struct {
	Region  string
	State   string
	IP      string
	SysName string
}

type trafficStats struct {
	Port      string
	IfName    string
	IfSpeed   float64
	MaxInBps  float64
	AvgInBps  float64
	MaxOutBps float64
	AvgOutBps float64
	UsageIn   float64
	UsageOut  float64
	Samples   int
	Instance  string
}

type GponStats struct {
	Port      string
	IfName    string
	IfSpeed   float64
	MaxInBps  float64
	AvgInBps  float64
	MaxOutBps float64
	AvgOutBps float64
	UsageIn   float64
	UsageOut  float64
}

type RegionStats struct {
	State     string
	IfSpeed   float64
	MaxInBps  float64
	AvgInBps  float64
	MaxOutBps float64
	AvgOutBps float64
	UsageIn   float64
	UsageOut  float64
}

type OltStats struct {
	Instance  string
	SysName   string
	IfSpeed   float64
	MaxInBps  float64
	AvgInBps  float64
	MaxOutBps float64
	AvgOutBps float64
	UsageIn   float64
	UsageOut  float64
}

type TrafficByInstance struct {
	IP        string
	State     string
	Region    string
	SysName   string
	Time      time.Time
	BpsIn     float64
	BpsOut    float64
	Bandwidth float64
	BytesIn   float64
	BytesOut  float64
}
