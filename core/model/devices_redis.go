package model

type RedisDevice struct {
	Devices []uint `redis:"devices"`
}
