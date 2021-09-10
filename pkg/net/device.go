package net

import (
	"fmt"
	"sync"
)

// 後で見る
type LinkDeviceCallbackHandler func(link LinkDevice, protocol EthernetType, payload []byte, src, dst HardwareAddress)

type LinkDevice interface {
	Type() HardwareType
	Name() string
	Address() HardwareAddress
	BroadcastAddress() HardwareAddress
	MTU() int
	HeaderSize() int
	NeedARP() bool
	Close()
	Read(data []byte) (int, error)
	RxHandler(frame []byte, callback LinkDeviceCallbackHandler)
	Tx(proto EthernetType, data []byte, dst []byte) error
}

type Device struct {
	LinkDevice
	errors     chan error
	interfaces []ProtocolInterface
	sync.RWMutex // 排他制御
}

var devices = sync.Map{}

// デバイスを登録
func RegisterDevice(link LinkDevice) (*Device, error) {
	// すでに登録済みかどうか確認する
	if _, exists := devices.Load(link); exists {
		return nil, fmt.Errorf("link device '%s' is already registered", link.Name())
	}
	dev := &Device{
		LinkDevice: link,
		errors:     make(chan error),
	}

	// launch rx loop
	go func() {
		var buf = make([]byte, dev.HeaderSize()+dev.MTU())
		for {
			n, err := dev.Read(buf)
			if n > 0 {
				dev.RxHandler(buf[:n], rxHander)
			}
			if err != nil {
				dev.errors <- err
				break
			}
		}
		close(dev.errors)
	}()

	// 登録
	devices.Store(link, dev)
	return dev, nil
}

func rxHander(link LinkDevice, protocol EthernetType, payload []byte, src, dst HardwareAddress){
	protocols.Range(func(key, value interface{}) bool {
		var (
			Type = key.(EthernetType)
			entry = value.(*entry)
		)
		if Type == EthernetType(protocol) {
			dev, ok := devices.Load(link)
			if !ok {
				panic("device not found")
			}
			entry.rxQueue <- &packet{
				dev: dev.(*Device),
				data: payload,
				src: src,
				dst: dst,
			}
			return false
		}
		return true
	})
}

func Devices() []*Device {
	ret := []*Device{}
	devices.Range(func(_, value interface{}) bool {
		ret = append(ret, value.(*Device))
		return true
	})
	return ret
}

// 登録されているinterfaceを返す
func (d *Device) Interfaces() []ProtocolInterface {
	d.RLock() // Read onlyにロック
	ret := make([]ProtocolInterface, len(d.interfaces))
	for i, iface := range d.interfaces {
		ret[i] = iface
	}
	d.RUnlock() // ロック解除
	return ret
}
