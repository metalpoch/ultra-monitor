package dto

type Prometheus struct {
	Region string `json:"region"`
	State  string `json:"state"`
	IP     string `json:"ip"`
	IDX    int64  `json:"idx"`
	Shell  uint8  `json:"shell"`
	Card   uint8  `json:"card"`
	Port   uint8  `json:"port"`
}
