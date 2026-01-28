package capture

import (
	"github.com/gopacket/gopacket"
	"github.com/gopacket/gopacket/layers"
)

func convPackage(packet gopacket.Packet) *OutputPacket {
	var srcIp, dstIp string
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer != nil {
		srcIp = ipLayer.(*layers.IPv4).SrcIP.String()
		dstIp = ipLayer.(*layers.IPv4).DstIP.String()
	} else {
		ipLayer = packet.Layer(layers.LayerTypeIPv6)
		if ipLayer != nil {
			srcIp = ipLayer.(*layers.IPv6).SrcIP.String()
			dstIp = ipLayer.(*layers.IPv6).DstIP.String()
		}
	}
	pkg := &OutputPacket{
		Timestamp:  packet.Metadata().Timestamp.UnixMilli(),
		PacketSize: packet.Metadata().Length,
		SrcIp:      srcIp,
		DstIp:      dstIp,
	}
	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if tcpLayer != nil {
		tcp, _ := tcpLayer.(*layers.TCP)
		srcPort := int(tcp.SrcPort)
		dstPort := int(tcp.DstPort)
		pkg.Proto = "tcp"
		pkg.TCP = &TCPPacket{
			SrcPort:    srcPort,
			DstPort:    dstPort,
			Seq:        tcp.Seq,
			Ack:        tcp.Ack,
			DataOffset: tcp.DataOffset,
			FIN:        tcp.FIN,
			SYN:        tcp.SYN,
			RST:        tcp.RST,
			PSH:        tcp.PSH,
			ACK:        tcp.ACK,
			URG:        tcp.URG,
			ECE:        tcp.ECE,
			CWR:        tcp.CWR,
			NS:         tcp.NS,
			Window:     tcp.Window,
			Checksum:   tcp.Checksum,
			Urgent:     tcp.Urgent,
			Payload:    tcp.Payload,
			Options:    convTCPOption(tcp.Options),
			Padding:    tcp.Padding,
		}
		return pkg
	}
	udpLayer := packet.Layer(layers.LayerTypeUDP)
	if udpLayer != nil {
		udp, _ := udpLayer.(*layers.UDP)
		srcPort := int(udp.SrcPort)
		dstPort := int(udp.DstPort)
		pkg.Proto = "udp"
		pkg.UDP = &UDPPacket{
			SrcPort:  srcPort,
			DstPort:  dstPort,
			Length:   udp.Length,
			Checksum: udp.Checksum,
			Payload:  udp.Payload,
		}
		return pkg
	}
	return nil
}

func convTCPOption(ops []layers.TCPOption) []TCPOption {
	var result []TCPOption
	for _, op := range ops {
		result = append(result, TCPOption{
			Type:       uint8(op.OptionType),
			Length:     op.OptionLength,
			OptionData: op.OptionData,
		})
	}
	return result
}
