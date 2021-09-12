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
	iface, err := ip.CreateInterface(dev, "192.0.2.2", "255.255.255.0", "")
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

	// launch send loop
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		t := time.NewTicker(1 * time.Second)
		defer t.Stop()
		peer := ip.ParseAddress("192.0.2.1")
		data := []byte("1234567890")
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				fmt.Printf("send ICMP Echo to %s\n", peer)
				icmp.EchoRequest(data, peer)
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
