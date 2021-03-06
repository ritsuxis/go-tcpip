package tcp

import (
	"fmt"

	"github.com/ritsuxis/go-tcpip/pkg/ip"
)

func Dial(local, remote *Address) (*Conn, error) {
	if local == nil {
		iface := ip.GetInterfaceByRemoteAddress(remote.Addr)
		if iface == nil {
			return nil, fmt.Errorf("dial failure")
		}
		local = &Address{
			Addr: iface.Address(),
		}
	}
	entry := repo.add(local)
	if entry == nil {
		return nil, fmt.Errorf("dial failure")
	}
	return &Conn{
		Cb:   entry,
		peer: remote,
	}, nil
}