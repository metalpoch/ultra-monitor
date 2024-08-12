package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type ElementOLT struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	OLT       string             `bson:"olt"`       // index
	Interface string             `bson:"interface"` //index
	Slot      int8               `bson:"slot"`
	Card      int8               `bson:"card"`
	Port      int8               `bson:"port"`
}
