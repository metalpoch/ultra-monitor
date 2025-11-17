package entity

type InterfacesDetailedOlt struct {
	Region     string `db:"region"`
	State      string `db:"state"`
	Sysname    string `db:"sysname"`
	OltIP      string `db:"olt_ip"`
	OltVerbose string `db:"olt_verbose"`
}

type InterfacesOlt struct {
	OltVerbose string `db:"olt_verbose"`
	OltIP      string `db:"olt_ip"`
}
