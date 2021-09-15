package tcp

import (
	"fmt"

	"github.com/ritsuxis/go-tcpip/pkg/net"
)

type Address struct {
	Addr net.ProtocolAddress
	Port uint16
}

func (a Address) Network() string {
	return "tcp"
}

func (a Address) String() string {
	return fmt.Sprintf("%s:%d", a.Addr, a.Port)
}