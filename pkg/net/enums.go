package net

type HardwareType uint16

const (
	// arpのヘッダにあるハードウェアタイプの欄
	// 今日だとイーサネットばっかり使われているが、他にもある
	// http://www.iana.org/assignments/arp-parameters/arp-parameters.xhtml
	HardwareTypeLoopBack = 0x0000
	HardwareTypeEthernet = 0x0001

	// イーサネットフレームのヘッダにある上位プロトコルの番号
	// https://www.iana.org/assignments/ieee-802-numbers/ieee-802-numbers.txt
	EthernetTypeIP   EthernetType = 0x0800
	EthernetTypeARP  EthernetType = 0x0806
	EthernetTypeIPv6 EthernetType = 0x86dd

	// IPのヘッダにある上位プロトコルの番号
	ProtocolNumberICMP ProtocolNumber = 0x01
	ProtocolNumberTCP  ProtocolNumber = 0x06
	ProtocolNumberUDP  ProtocolNumber = 0x11
)

func (t HardwareType) String() string {
	switch t {
	case HardwareTypeLoopBack:
		return "LoopBack"
	case HardwareTypeEthernet:
		return "Ethernet"
	default:
		return "Unknown"
	}
}

type EthernetType uint16

func (t EthernetType) String() string {
	switch t {
	case EthernetTypeIP:
		return "IP"
	case EthernetTypeARP:
		return "ARP"
	case EthernetTypeIPv6:
		return "IPv6"
	default:
		return "Unknown"
	}
}

type ProtocolNumber uint8

func (t ProtocolNumber) String() string {
	switch t {
	case ProtocolNumberICMP:
		return "ICMP"
	case ProtocolNumberTCP:
		return "TCP"
	case ProtocolNumberUDP:
		return "UDP"
	default:
		return "Unknown"
	}
}
