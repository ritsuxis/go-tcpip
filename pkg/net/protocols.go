package net

import (
	"fmt"
	"log"
	"sync"
)

// RXは多分受信の意味
type ProtocolRxHandler func(dev *Device, data []byte, src, dst HardwareAddress) error

type packet struct {
	dev *Device
	data []byte
	src HardwareAddress
	dst HardwareAddress
}

type entry struct {
	Type EthernetType
	rxHandler ProtocolRxHandler
	rxQueue chan *packet // チャネル
}

var protocols = sync.Map{} // キャッシュみたいなやつ

func RegisterProtocol(Type EthernetType, rxHandler ProtocolRxHandler) error {
	// 最初に登録しようとしているプロトコルがすでに登録済みか確認する
	if _, exists := protocols.Load(Type); exists {
		return fmt.Errorf("protocol `%s` is already registered", Type)
	}
	entry := &entry{
		Type: Type,
		rxHandler: rxHandler,
		rxQueue: make(chan *packet),
	}

	// launch rx loop
	// goroutine
	go func() {
		for packet := range entry.rxQueue {
			if err := entry.rxHandler(packet.dev, packet.data, packet.src, packet.dst); err != nil {
				log.Println(err)
			}
		}
	}()

	// 登録
	protocols.Store(Type, entry)
	return nil
}