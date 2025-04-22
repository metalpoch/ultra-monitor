package model

import (
	"time"

	"github.com/metalpoch/olt-blueprint/common/entity"
)

type Measurement struct {
	Interface   entity.Interface
	InterfaceID uint
	Date        time.Time
	Bandwidth   uint64
	In          uint64
	Out         uint64
}
