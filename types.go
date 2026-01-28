package capture

type NetworkInterface struct {
	Name        string `json:"name"`
	IpV4Address string `json:"ipv4_address"`
	IpV6Address string `json:"ipv6_address"`
}

type FilterParam struct {
	InterfaceName string // network interface
	Proto         string // tcp or udp
	Host          *struct {
		Src string
		Dst string
		Any string
	} // host
	Port *struct {
		Src int
		Dst int
		Any int
	} // port
}

type OutputPacket struct {
	Timestamp  int64      `json:"timestamp"`
	PacketSize int        `json:"packet_size"`
	SrcIp      string     `json:"src_ip"`
	DstIp      string     `json:"dst_ip"`
	Proto      string     `json:"proto"` // tcp | udp
	TCP        *TCPPacket `json:"tcp"`
	UDP        *UDPPacket `json:"udp"`
}

type TCPOption struct {
	Type       uint8  `json:"type"`
	Length     uint8  `json:"length"`
	OptionData []byte `json:"data"`
}

type TCPPacket struct {
	SrcPort                                    int    `json:"src"`
	DstPort                                    int    `json:"dst"`
	Seq                                        uint32 `json:"seq"`
	Ack                                        uint32 `json:"ack"`
	DataOffset                                 uint8  `json:"data_offset"`
	FIN, SYN, RST, PSH, ACK, URG, ECE, CWR, NS bool
	Window                                     uint16      `json:"window"`
	Checksum                                   uint16      `json:"checksum"`
	Urgent                                     uint16      `json:"urgent"`
	Payload                                    []byte      `json:"payload"`
	Options                                    []TCPOption `json:"options"`
	Padding                                    []byte      `json:"padding"`
}

type UDPPacket struct {
	SrcPort  int    `json:"src"`
	DstPort  int    `json:"dst"`
	Length   uint16 `json:"length"`
	Checksum uint16 `json:"checksum"`
	Payload  []byte `json:"payload"`
}
