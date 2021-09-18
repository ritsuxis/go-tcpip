package tcp_test

import (
	"testing"

	"github.com/ritsuxis/go-tcpip/pkg/ip"
	"github.com/ritsuxis/go-tcpip/pkg/tcp"
)

func TestTcpPacketBuild(t *testing.T) {
	peer := &tcp.Address{
		Addr: ip.ParseAddress("192.0.2.1"),
	}
	bytes := tcp.Build(uint16(1234), uint16(5432), uint32(1), uint32(2), tcp.ControlFlag(tcp.FIN | tcp.ACK), uint16(1240), uint16(0), []byte{byte(0x11), byte(0x22)}, nil)
	packet, _, err := tcp.Parse(bytes, peer.Addr, peer.Addr)
	if err != nil {
		t.Fatal(err)
	}
	packet.Dump()
}