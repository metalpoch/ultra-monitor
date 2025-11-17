package dto

type InterfacesDetailedOlt struct {
	Region     string `json:"region"`
	State      string `json:"state"`
	Sysname    string `json:"sysname"`
	OltIP      string `json:"olt_ip"`
	OltVerbose string `json:"olt_verbose"`
}

type InterfacesOlt struct {
	OltVerbose string `json:"olt_verbose"`
	OltIP      string `json:"olt_ip"`
}
