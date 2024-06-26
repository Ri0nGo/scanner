package scanner

const (
	ScanTypePING string = "PING"
	ScanTypeTCP  string = "TCP"

	AddressTypeIP      string = "IP"
	AddressTypeNetwork string = "NETWORK"

	PingNoUse   = 0
	PingSuccess = 1
	PingFailed  = 2

	PortOpen  = 1
	PortClose = 2

	IPStatusIdle    = 0
	IPStatusRunning = 1
	IPStatusFinish  = 2
)
