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
	ethernet "github.com/ritsuxis/go-tcpip/pkg/ether"
	"github.com/ritsuxis/go-tcpip/pkg/icmp"
	"github.com/ritsuxis/go-tcpip/pkg/ip"
	"github.com/ritsuxis/go-tcpip/pkg/net"
	"github.com/ritsuxis/go-tcpip/pkg/raw/tuntap"
	"github.com/ritsuxis/go-tcpip/pkg/tcp"
)

var sig chan os.Signal

func init() {
	icmp.Init()
	arp.Init()
	tcp.Init()
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
	peer := &tcp.Address{
		Addr: ip.ParseAddress("192.0.2.1"),
		Port: 10381,
	}
	icmp.EchoRequest(nil, peer.Addr) // dummy
	conn, err := tcp.Dial(nil, peer)
	if err != nil {
		panic(err)
	}
	// launch send loop
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		t := time.NewTicker(2 * time.Second)
		defer t.Stop()
		data := []byte("hoge\n")
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				// tcp 3way handshake
				switch conn.Cb.State {
				case tcp.Close:
					{
						fmt.Printf("SYN send to %s\n", peer)
						conn.Write(nil, tcp.ControlFlag(tcp.SYN), nil)
					}
				case tcp.SynSent:
					{
						fmt.Printf("ACK send to %s\n", peer)
						conn.Write(nil, tcp.ControlFlag(tcp.ACK), nil)
					}
				case tcp.Established:
					{
						fmt.Printf("ACK send to %s\n", peer)
						conn.Write(nil, tcp.ACK, nil)
					}
				case tcp.Sent:
					{
						fmt.Printf("%d bytes data send to %s\n", len(data), peer)
						conn.Write(data, tcp.ACK, nil)
					}
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
	conn.Close()
	fmt.Println("good bye")
}
