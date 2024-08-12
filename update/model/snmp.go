package model

type Snmp struct {
	IfName    map[int]string
	ByteIn    map[int]int64
	ByteOut   map[int]int64
	Bandwidth map[int]int16
}
