package constants

const (
	SYSNAME_OID        string = "1.3.6.1.2.1.1.5.0"
	SYSLOCATION_OID    string = "1.3.6.1.2.1.1.6.0"
	IF_DESCR_OID       string = "1.3.6.1.2.1.2.2.1.2"
	IF_ALIAS_OID       string = "1.3.6.1.2.1.31.1.1.1.18"
	IF_NAME_OID        string = "1.3.6.1.2.1.31.1.1.1.1"
	IF_HIGH_SPEED_OID  string = "1.3.6.1.2.1.31.1.1.1.15" // bandwidth as Mbps
	SYSNAME            string = "sysname"
	IF_ALIAS           string = "ifalias"
	IF_NAME            string = "ifname"
	IF_DESCR           string = "ifdescr"
	BANDWIDTH          string = "bw"
	IN                 string = "in"
	OUT                string = "out"
	OVERFLOW_COUNTER64 uint64 = 1<<64 - 1
)
