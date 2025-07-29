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
	Time        time.Time
	Description string
	BpsIn       float64
	BpsOut      float64
	Bandwidth   float64
	BytesIn     float64
	BytesOut    float64
}

type InfoDevice struct {
	Region  string
	State   string
	IP      string
	IfName  string
	IfIndex int64
}
