package model

type Traffic struct {
	Date      int64  `json:"date"`
	OLT       string `json:"olt"`       // index
	Interface string `json:"interface"` // index
	KpbsIn    int32  `json:"kpbs_in"`
	KpbsOut   int32  `json:"kbps_out"`
	Bandwidth int16  `json:"bandwidth"`
}

type QueryTraffic struct {
	Sysname    string   `form:"sysname"`
	ElementsID []string `form:"id"`
	Firstday   int      `form:"firstday"`
	Lastday    int      `form:"lastday"`
}

type TrafficResponse struct {
	Date      string `json:"date"`
	KpbsIn    int32  `json:"kpbs_in"`
	KpbsOut   int32  `json:"kbps_out"`
	Bandwidth int16  `json:"bandwidth"`
}

// type TrafficOLT struct {
// 	Date      int64  `bson:"_id"`
// 	OLT       string `bson:"olt"`       // index
// 	Interface string `bson:"interface"` // index
// 	KbpsIn    int16  `bson:"kbps_in"`
// 	KbpsOut   int16  `bson:"kbps_out"`
// 	Bandwidth int16  `bson:"bandwidth"`
// }

// type TrafficMeasurement struct {
// 	Date      string `bson:"date"`
// 	KpbsIn    int16  `bson:"kpbs_in"`
// 	KpbsOut   int16  `bson:"kbps_out"`
// 	Bandwidth int16  `bson:"bandwidth"`
// }
