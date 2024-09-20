package entity

import "time"

type Measurement struct {
	Date        time.Time
	Bw          uint
	In          uint
	Out         uint
	InterfaceID uint
}
