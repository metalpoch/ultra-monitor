package entity

type TrafficOLT struct {
	Date      int64  `bson:"_id"`
	OLT       string `bson:"olt"`       // index
	Interface string `bson:"interface"` // index
	KbpsIn    int32  `bson:"kbps_in"`
	KbpsOut   int32  `bson:"kbps_out"`
	Bandwidth int16  `bson:"bandwidth"`
}

type TrafficMeasurement struct {
	Date      string `bson:"date"`
	KpbsIn    int32  `bson:"kpbs_in"`
	KpbsOut   int32  `bson:"kbps_out"`
	Bandwidth int16  `bson:"bandwidth"`
}
