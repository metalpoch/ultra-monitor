package entity

type Count struct {
	Date      int64  `bson:"date"`
	OLT       string `bson:"olt"`       // index
	Interface string `bson:"interface"` // index
	BytesIn   int64  `bson:"bytes_in"`
	BytesOut  int64  `bson:"bytes_out"`
	Bandwidth int16  `bson:"bandwidth"`
}
