package prometheus

import "time"

type dataProm struct {
	Labels map[string]string
	Value  float64
	Time   time.Time
}

type TrafficResult struct {
	OltIP       string
	OltRegion   string
	OltState    string
	SysLocation string
	SysName     string
	IfIndex     int64
	IfName      string
	IfDescr     string
	IfAlias     string
	IfSpeed     float64
	BpsIn       float64
	BpsOut      float64
	Bandwidth   float64
	BytesIn     float64
	BytesOut    float64
	Time        time.Time
}

type InfoDevice struct {
	Region  string
	State   string
	IP      string
	IfName  string
	IfIndex int64
}
