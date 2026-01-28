package capture

import (
	"github.com/gopacket/gopacket/pcap"
	"net"
)

// ListNetworkInterfaces all ipv4 and ipv6 network interfaces
func ListNetworkInterfaces() ([]NetworkInterface, error) {
	var interfaces []NetworkInterface
	devices, err := pcap.FindAllDevs()
	if err != nil {
		return nil, err
	}
	for _, device := range devices {
		i := NetworkInterface{
			Name: device.Name,
		}
		for _, address := range device.Addresses {
			if len(address.IP) == net.IPv4len {
				i.IpV4Address = address.IP.String()
			} else if len(address.IP) == net.IPv6len {
				i.IpV6Address = address.IP.String()
			}
		}
		if i.IpV4Address == "" || i.IpV6Address == "" {
			continue
		}
		interfaces = append(interfaces, i)
	}
	return interfaces, nil
}
