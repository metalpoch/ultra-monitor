package model

type LocationsDevice struct {
	Devices []uint `redis:"devices"`
}
