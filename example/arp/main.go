package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ritsuxis/go-tcpip/pkg/arp"
	ethernet "github.com/ritsuxis/go-tcpip/pkg/ether"
	"github.com/ritsuxis/go-tcpip/pkg/ip"
	"github.com/ritsuxis/go-tcpip/pkg/net"
	"github.com/ritsuxis/go-tcpip/pkg/raw/tuntap"
)

var sig chan os.Signal

const (
	etherTapName    = "tap0"
	etherTapHwAddr  = "00:00:5e:00:53:01"
	etherTapIPAddr  = "192.0.2.2"
	etherTapNetMask = "255.255.255.0"
	etherNextHop    = "192.0.2.1"
)

func init() {
	arp.Init()
}

func setup() (*net.Device, error) {
	// signal handling
	sig = make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	// parse command line params
	// name := flag.String("name", "", "device name")
	// addr := flag.String("addr", "", "hardware address")
	name := etherTapName
	addr := etherTapHwAddr
	flag.Parse()
	raw, err := tuntap.NewTap(name)
	if err != nil {
		return nil, err
	}
	link, err := ethernet.NewDevice(raw)
	if err != nil {
		return nil, err
	}
	if addr != "" {
		link.SetAddress(ethernet.ParseAddress(addr))
	}
	dev, err := net.RegisterDevice(link)
	if err != nil {
		return nil, err
	}
	iface, err := ip.CreateInterface(dev, etherTapIPAddr, etherTapNetMask, "")
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

	for {
		
		break
	}
	select {
	case s := <-sig:
		fmt.Printf("sig: %s\n", s)
		dev.Shutdown()
	}
	fmt.Println("good bye")
}
