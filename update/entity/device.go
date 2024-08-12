package entity

type Device struct {
	Sysname   string `bson:"_id"`
	IP        string `bson:"ip"`
	Community string `bson:"community"`
}
