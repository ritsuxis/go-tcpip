package ethernet

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"

	"github.com/ritsuxis/go-tcpip/pkg/net"
	"github.com/ritsuxis/go-tcpip/pkg/raw"
)

const (
	headerSize     = 14
	trailerSize    = 0 // FCSの部分だけど、ここではFCSを抜いているので0
	maxPayloadSize = 1500
	minPayloadSize = 14

	maxFrameSize = headerSize + maxPayloadSize + trailerSize
	minFrameSize = headerSize + minPayloadSize + trailerSize
)

type Device struct {
	raw  raw.Device
	addr Address
	mtu  int
}

var _ net.LinkDevice = &Device{} // interface check

func NewDevice(raw raw.Device) (*Device, error) {
	if raw == nil {
		return nil, fmt.Errorf("raw device is required")
	}
	addr := Address{}
	copy(addr[:], raw.Address())
	return &Device{
		raw:  raw,
		addr: addr,
		mtu:  maxPayloadSize,
	}, nil
}

// github.com/ritsuxis/go-tcpip/pkg/net のLinkDeviceのInterfaceに合わせるために必要
func (d *Device) Type() net.HardwareType {
	return net.HardwareTypeEthernet
}

func (d *Device) Name() string {
	return d.raw.Name()
}

func (d *Device) Address() net.HardwareAddress {
	return d.addr
}

func (d *Device) BroadcastAddress() net.HardwareAddress {
	return BroadcastAddress
}

func (d *Device) SetAddress(addr Address) {
	d.addr = addr
}

func (d *Device) MTU() int {
	return d.mtu
}

func (d *Device) HeaderSize() int {
	return headerSize
}

func (d *Device) NeedARP() bool {
	return true
}

func (d *Device) Close() {
	d.raw.Close()
}

func (d *Device) Read(buf []byte) (int, error) {
	return d.raw.Read(buf)
}

func (d *Device) RxHandler(data []byte, callback net.LinkDeviceCallbackHandler) {
	frame, err := parse(data)
	if err != nil {
		log.Println(err)
		return
	}

	// 受信したデータが自分宛(デバイス宛て)かどうか
	if frame.Dst != d.addr {
		if !frame.Dst.isGroupAddress() {
			// other host frame
			return
		}
		if frame.Dst != BroadcastAddress {
			// multicast frame: unsupported
			return
		}
	}
	// if frame.Src == d.addr {
	// 	// loopback frame
	// }

	// 自分宛ならコールバックを呼び出す
	callback(d, frame.Type, frame.payload, frame.Src, frame.Dst)
}

func (d *Device) Tx(Type net.EthernetType, data []byte, dst []byte) error {
	hdr := header{
		Dst:  NewAddress(dst),
		Src:  d.addr,
		Type: Type,
	}
	frame := bytes.NewBuffer(make([]byte, 0))
	binary.Write(frame, binary.BigEndian, hdr)
	binary.Write(frame, binary.BigEndian, data)
	if pad := minFrameSize - frame.Len(); pad > 0 {
		binary.Write(frame, binary.BigEndian, bytes.Repeat([]byte{byte(0)}, pad))
	}
	_, err := d.raw.Write(frame.Bytes())
	return err
}
