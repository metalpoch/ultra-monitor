package model

type Count struct {
	OLT       string `json:"olt"`       // index
	Interface string `json:"interface"` // index
	Date      int64  `json:"date"`
	BytesIn   int64  `json:"bytes_in"`
	BytesOut  int64  `json:"bytes_out"`
	Bandwidth int16  `json:"bandwidth"`
}

type CountDiff struct {
	OLT           string `json:"olt"`       // index
	Interface     string `json:"interface"` // index
	PrevDate      int64  `json:"prev_date"`
	PrevBytesIn   int64  `json:"prev_bytes_in"`
	PrevBytesOut  int64  `json:"prev_bytes_out"`
	CurrDate      int64  `json:"curr_date"`
	CurrBytesIn   int64  `json:"curr_bytes_in"`
	CurrBytesOut  int64  `json:"curr_bytes_out"`
	CurrBandwidth int16  `json:"bandwidth"`
}
