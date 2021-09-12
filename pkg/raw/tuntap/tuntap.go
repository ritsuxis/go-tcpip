package tuntap

import (
	"io"

	"github.com/ritsuxis/go-tcpip/pkg/ether"
)

type Tap struct {
	io.ReadWriteCloser // Reader, Writer, Closerをinterfaceとしてもっている
	name               string
	buffer             chan *ethernet.Frame
}

const macAddressLength = 6

func NewTap(name string) (*Tap, error) {
	n, f, err := openTap(name)
	if err != nil {
		return nil, err
	}
	return &Tap{
		ReadWriteCloser: f,
		name:            n,
	}, nil
}

func (t Tap) Address() []byte {
	addr, _ := getAddress(t.name)
	return addr[:macAddressLength]
}

func (t Tap) Name() string {
	return t.name
}

func (t Tap) Buffer() chan *ethernet.Frame {
	return t.buffer
}