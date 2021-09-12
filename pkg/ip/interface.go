package ip

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"

	"github.com/ritsuxis/go-tcpip/pkg/arp"
	"github.com/ritsuxis/go-tcpip/pkg/net"
)

type Interface struct {
	unicast   Address
	netmask   Address
	broadcast Address
	gateway   Address
	device    *net.Device
}

func newInterface(dev *net.Device, unicast, netmask Address) (*Interface, error) {
	return &Interface{
		unicast: unicast,
		netmask: netmask,
		// network addressの部分はそのままで、host addressの部分を全部1にした値が得られる -> broadcast address
		broadcast: Address{
			unicast[0]&netmask[0] | ^netmask[0],
			unicast[1]&netmask[1] | ^netmask[1],
			unicast[2]&netmask[2] | ^netmask[2],
			unicast[3]&netmask[3] | ^netmask[3],
		},
		device: dev,
	}, nil
}

func CreateInterface(dev *net.Device, unicast, netmask, gateway string) (*Interface, error) {
	addr := ParseAddress(unicast)
	if addr == InvalidAddress {
		return nil, fmt.Errorf("invalid address: %s", unicast)
	}
	mask := ParseAddress(netmask)
	if mask == InvalidAddress { // netmaskが全部0はありえない(全部1はあり得る)
		return nil, fmt.Errorf("invalid address: %s", netmask)
	}
	gw := EmptyAddress
	if gateway != "" {
		gw = ParseAddress(gateway)
		if gw == InvalidAddress {
			return nil, fmt.Errorf("invalid address: %s", gateway)
		}
	}
	net := Address{ // network address
		addr[0] & mask[0],
		addr[1] & mask[1],
		addr[2] & mask[2],
		addr[3] & mask[3],
	}
	iface, err := newInterface(dev, addr, mask)
	if err != nil {
		return nil, err
	}
	repo.add(iface, net, mask, gw)
	return iface, nil
}

func GetInterface(addr net.ProtocolAddress) net.ProtocolInterface {
	// 登録されているデバイスとそのインターフェース群から目的のインターフェースを見つける
	for _, dev := range net.Devices() {
		for _, iface := range dev.Interfaces() {
			if iface.Type() == net.EthernetTypeIP && iface.Address() == addr {
				return iface
			}
		}
	}
	return nil
}

func GetInterfaceByRemoteAddress(remote net.ProtocolAddress) net.ProtocolInterface {
	route := repo.lookup(nil, remote.(Address))
	if route == nil {
		log.Printf("address not found: %s", remote)
		return nil
	}
	log.Printf("address found: %s", remote)
	return route.iface
}

// net.ProtocolInterface に合わせるために必要
func (iface *Interface) Type() net.EthernetType {
	return net.EthernetTypeIP
}

func (iface *Interface) Address() net.ProtocolAddress {
	return iface.unicast
}

func (iface *Interface) Device() *net.Device {
	return iface.device
}

func (iface *Interface) xmit(datagram *datagram, nexthop net.ProtocolAddress) error {
	buf := new(bytes.Buffer)
	// make IPv4 payload
	binary.Write(buf, binary.BigEndian, &datagram.header)
	binary.Write(buf, binary.BigEndian, datagram.payload)
	b := buf.Bytes()
	binary.BigEndian.PutUint16(b[10:12], net.CheckSum16(b, int((datagram.VHL&0x0f)<<2), 0))

	var hardwareAddress []byte
	// arpが必要か見る
	if iface.Device().NeedARP() {
		if nexthop != nil {
			var err error
			// 目的のnexthopのためのarp解決要求を出す
			hardwareAddress, err = arp.Resolve(iface, nexthop.Bytes(), b)
			if err != nil {
				return err
			}
		} else {
			// nexthopがない場合はブロードキャストする
			hardwareAddress = iface.Device().BroadcastAddress().Bytes()
		}
	}
	return iface.Device().Tx(net.EthernetTypeIP, b, hardwareAddress)
}

func (iface *Interface) Tx(protocol net.ProtocolNumber, data []byte, dst net.ProtocolAddress) error {
	var nexthop net.ProtocolAddress
	src := iface.unicast
	if dst.(Address) != BroadcastAddress {
		// ルート検索
		routeEntry := repo.lookup(iface, dst.(Address))
		if routeEntry == nil {
			return fmt.Errorf("route not found")
		}
		iface = routeEntry.iface
		if nexthop = routeEntry.nexthop; nexthop == EmptyAddress {
			nexthop = dst
		}
	}
	assembler := newAssembler(protocol, data, src, dst, idm.next(), iface.Device().MTU())
	for _, datagram := range assembler.assemble() {
		if err := iface.xmit(datagram, nexthop); err != nil {
			return err
		}
	}
	return nil
}
