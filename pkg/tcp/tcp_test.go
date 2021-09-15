package tcp_test

import (
	"testing"

	// "github.com/ritsuxis/go-tcpip/pkg/ip"
	"github.com/ritsuxis/go-tcpip/pkg/tcp"
)

func TestTcpPacketBuild(t *testing.T) {
	tcp.Build(1234, 5432, 1, 2, tcp.ControlFlag(tcp.FIN | tcp.ACK), 1240, 20, []byte{byte(0x11), byte(0x22)}, tcp.Options{tcp.EndOptionList{}})
}