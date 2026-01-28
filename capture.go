package capture

import "C"
import (
	"sync/atomic"

	"github.com/gopacket/gopacket"
	"github.com/gopacket/gopacket/layers"
	"github.com/gopacket/gopacket/pcap"
)

type Listener interface {
	OnPkg(pkg *OutputPacket)
	OnClose()
}

var _captureHandle atomic.Value

func StartCapture(filter FilterParam, listener Listener) error {
	StopCapture()
	handle, err := pcap.OpenLive(filter.InterfaceName,
		65536,
		false,
		-1,
	)
	if err != nil {
		return err
	}
	_captureHandle.Store(handle)
	defer handle.Close()
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	defer func() {
		listener.OnClose()
	}()
	pkgChan := packetSource.Packets()
	for {
		packet, ok := <-pkgChan
		if !ok {
			break
		}
		// IP
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
		if filter.Host != nil {
			if srcIp == "" && dstIp == "" {
				continue
			}
			if filter.Host.Any != "" {
				if filter.Host.Any != srcIp || filter.Host.Any != dstIp {
					continue
				}
			} else if filter.Host.Src != "" {
				if filter.Host.Any != srcIp {
					continue
				}
			} else if filter.Host.Dst != "" {
				if filter.Host.Any != dstIp {
					continue
				}
			}
		}
		// TCP/UDP 、 Port
		tcpLayer := packet.Layer(layers.LayerTypeTCP)
		udpLayer := packet.Layer(layers.LayerTypeUDP)
		var srcPort, dstPort int
		if tcpLayer != nil {
			tcp := tcpLayer.(*layers.TCP)
			srcPort = int(tcp.SrcPort)
			dstPort = int(tcp.DstPort)
		} else if udpLayer != nil {
			tcp := udpLayer.(*layers.UDP)
			srcPort = int(tcp.SrcPort)
			dstPort = int(tcp.DstPort)
		}
		if filter.Proto != "" {
			if tcpLayer == nil && udpLayer == nil {
				continue
			}
			if filter.Proto == "tcp" && tcpLayer == nil {
				continue
			} else if filter.Proto == "udp" && udpLayer == nil {
				continue
			}
		}
		if filter.Port != nil {
			if srcPort == 0 || dstPort == 0 {
				continue
			}
			if filter.Port.Any > 0 {
				if filter.Port.Any != srcPort && filter.Port.Any != dstPort {
					continue
				}
			} else if filter.Port.Src > 0 {
				if filter.Port.Src != srcPort {
					continue
				}
			} else if filter.Port.Dst > 0 {
				if filter.Port.Dst != dstPort {
					continue
				}
			}
		}
		// Output
		output := convPackage(packet)
		if output != nil {
			listener.OnPkg(output)
		}
	}
	return nil
}

func StopCapture() {
	if _captureHandle.Load() == nil {
		return
	}
	handler := _captureHandle.Load().(*pcap.Handle)
	handler.Close()
}
