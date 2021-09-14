package tcp_test

import (
	"testing"

	"github.com/ritsuxis/go-tcpip/pkg/tcp"
)

func TestFlagCheck(t *testing.T) {
	f := (tcp.FIN | tcp.ACK | tcp.URG)

	if f.String() != "UA---F" {
		t.Fatal("not match")
	}
}
