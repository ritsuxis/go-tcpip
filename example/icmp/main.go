package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ritsuxis/go-tcpip/pkg/arp"
	"github.com/ritsuxis/go-tcpip/pkg/ether"
	"github.com/ritsuxis/go-tcpip/pkg/icmp"
	"github.com/ritsuxis/go-tcpip/pkg/ip"
	"github.com/ritsuxis/go-tcpip/pkg/net"
	"github.com/ritsuxis/go-tcpip/pkg/raw/tuntap"
)

var sig chan os.Signal

func init() {
	arp.Init()
	icmp.Init()
}

func setup() (*net.Device, error) {
	// signal handling
	sig = make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	// parse command line params
	name := flag.String("name", "", "device name")
	addr := flag.String("addr", "", "hardware address")
	flag.Parse()
	raw, err := tuntap.NewTap(*name)
	if err != nil {
		return nil, err
	}
	link, err := ethernet.NewDevice(raw)
	if err != nil {
		return nil, err
	}
	if *addr != "" {
		link.SetAddress(ethernet.ParseAddress(*addr))
	}
	dev, err := net.RegisterDevice(link)
	if err != nil {
		return nil, err
	}
	iface, err := ip.CreateInterface(dev, "192.0.2.2", "255.255.255.0", "192.0.2.1")
	if err != nil {
		return nil, err
	}
	dev.RegisterInterface(iface)
	return dev, nil
}

func main() {
	dev, err := setup()
	if err != nil {
		panic(err)
	}
	fmt.Printf("[%s] %s\n", dev.Name(), dev.Address())
	for _, iface := range dev.Interfaces() {
		fmt.Printf("  - %s\n", iface.Address())
	}
	
	// buf := make(chan []byte, dev.HeaderSize()+dev.MTU())
	// go net.Rxloop(dev, buf)
	// buf <- []byte {
	// 	0xce, 0x5c, 0x7f, 0x6b, 0x61, 0xba, 0x22, 0xca, 0xf6, 0xf3, 0xc3, 0x6e, 0x08, 0x00, 0x45, 0x10,
	// 	0x00, 0x3c, 0xfd, 0x46, 0x40, 0x00, 0x40, 0x06, 0xbc, 0x0f, 0xc0, 0xa8, 0x00, 0x02, 0xc0, 0xa8,
	// 	0x00, 0x03, 0xeb, 0x1e, 0x00, 0x17, 0x37, 0x08, 0xb7, 0xf1, 0x00, 0x00, 0x00, 0x00, 0xa0, 0x02,
	// 	0x72, 0x10, 0x81, 0x84, 0x00, 0x00, 0x02, 0x04, 0x05, 0xb4, 0x04, 0x02, 0x08, 0x0a, 0x19, 0x64,
	// 	0x34, 0x32, 0x00, 0x00, 0x00, 0x00, 0x01, 0x03, 0x03, 0x07,
	// }

	dev.Interfaces()

	// launch send loop
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		t := time.NewTicker(1 * time.Second)
		defer t.Stop()
		peer := ip.ParseAddress("172.29.18.100")
		data := []byte("1234567890")
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				if ip.GetInterfaceByRemoteAddress(peer) == nil {
					peer := ip.ParseAddress("192.0.2.1")
					fmt.Printf("send ICMP Echo to %s\n", peer)
					icmp.EchoRequest(data, peer)
				} else {
					fmt.Printf("send ICMP Echo to %s\n", peer)
					icmp.EchoRequest(data, peer)
				}
			}
		}
	}()
	select {
	case s := <-sig:
		fmt.Printf("sig: %s\n", s)
		cancel()
	}
	wg.Wait()
	dev.Shutdown()
	fmt.Println("good bye")
}
