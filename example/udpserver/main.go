package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ritsuxis/go-tcpip/pkg/ether"
	"github.com/ritsuxis/go-tcpip/pkg/icmp"
	"github.com/ritsuxis/go-tcpip/pkg/ip"
	"github.com/ritsuxis/go-tcpip/pkg/net"
	"github.com/ritsuxis/go-tcpip/pkg/raw/tuntap"
	"github.com/ritsuxis/go-tcpip/pkg/udp"
)

var sig chan os.Signal

func init() {
	icmp.Init()
}

func setup() error {
	// signal handling
	sig = make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	// parse command line params
	name := flag.String("name", "", "device name")
	addr := flag.String("addr", "", "hardware address")
	flag.Parse()
	raw, err := tuntap.NewTap(*name)
	if err != nil {
		return err
	}
	link, err := ethernet.NewDevice(raw)
	if err != nil {
		return err
	}
	if *addr != "" {
		link.SetAddress(ethernet.ParseAddress(*addr))
	}
	dev, err := net.RegisterDevice(link)
	if err != nil {
		return err
	}
	iface, err := ip.CreateInterface(dev, "192.0.2.2", "255.255.255.0", "")
	if err != nil {
		return err
	}
	dev.RegisterInterface(iface)
	return nil
}

func main() {
	err := setup()
	if err != nil {
		panic(err)
	}
	conn, err := udp.Listen(
		&udp.Address{
			Addr: ip.EmptyAddress,
			Port: 7,
		},
	)
	if err != nil {
		panic(err)
	}
	// launch receive loop
	go func() {
		buf := make([]byte, 65536)
		for {
			n, peer, err := conn.ReadFrom(buf)
			if n > 0 {
				log.Printf("main: receive %d bytes data from %s", n, peer)
				conn.WriteTo(buf[:n], peer)
			}
			if err != nil {
				if err != io.EOF {
					fmt.Println(err)
				}
				return
			}
		}
	}()
	select {
	case s := <-sig:
		fmt.Printf("sig: %s\n", s)
		conn.Close()
	}
	fmt.Println("good bye")
}
