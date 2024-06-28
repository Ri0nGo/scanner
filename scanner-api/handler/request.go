package handler

type ReqScanResponse struct {
	ScanResult any            `json:"scan_result"`
	ScanStatus map[string]any `json:"scan_status"`
}

type ScanDetail struct {
	Type  string `json:"type"`  // "TCP" / "PING"
	Ports string `json:"ports"` // 端口范围
}

type Address struct {
	Type      string   `json:"type"`       // "IP" / "NETWORK"
	IPS       []string `json:"ips"`        // IP地址列表
	IPNetwork string   `json:"ip_network"` // IP网段
}

type ReqScan struct {
	UUID     string       `json:"uuid"`
	ScanType []ScanDetail `json:"scan_type"`
	Address  Address      `json:"address"`
}
