package tcp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"unsafe"

	"github.com/ritsuxis/go-tcpip/pkg/ip"
	"github.com/ritsuxis/go-tcpip/pkg/net"
)

var winSize uint16

type Conn struct {
	Cb   *cbEntry
	peer *Address
}

func (conn *Conn) Close() {
	// TODO
}

func getAppropriateInterface(local, remote net.ProtocolAddress) net.ProtocolInterface {
	if local.IsEmpty() {
		return ip.GetInterfaceByRemoteAddress(remote)
	}
	return ip.GetInterface(local)
}

func (conn *Conn) Write(data []byte, flag ControlFlag, ops Options) error {
	if conn.peer == nil {
		return fmt.Errorf("this Conn is not dialed")
	}
	return conn.WriteTo(data, conn.peer, flag, ops)
}

func (conn *Conn) WriteTo(data []byte, peer *Address, flag ControlFlag, ops Options) error {
	hdr := header {
		SourcePort: conn.Cb.Port,
		DestinationPort: peer.Port,
		WindowSize: winSize, // TODO: 受信側のサイズに合わせて変える
		Urgent: 0,
	}

	// tcp state
	entry := repo.lookup(conn.Cb.Address)
	if entry == nil {
		return fmt.Errorf("port not found")
	}

	if conn.Cb.State == Close {
		hdr.SequenceNumber = 0
		hdr.ACKNumber = 0
	} else if conn.Cb.State == SynSent {
		sa := <-entry.Number
		hdr.SequenceNumber = sa.Ack
		hdr.ACKNumber = sa.Seq + 1
	} else if conn.Cb.State == Established {
		sa := <-entry.Number
		hdr.SequenceNumber = sa.Ack
		hdr.ACKNumber = sa.Seq + 1
	} else if conn.Cb.State == Sent {
		sa := <-entry.Number
		hdr.SequenceNumber = sa.Ack
		hdr.ACKNumber = sa.Seq
	}

	conn.Cb.State = conn.Cb.State.TransitionSnd(flag)
	log.Printf("Snd: Now TCP state is %s", conn.Cb.State)


	buf := new(bytes.Buffer)

	// option handle
	var opLength = 0 
	if ops != nil { // optionが指定されていた時
		for _, op := range ops {
			opLength += op.Length()
		}

		eolPadding := 4 - (opLength % 4) // 32 bit = 4 byte, 4byteの倍数になるようにパディングする
		if eolPadding != 0 {
			for i := 0; i < eolPadding; i++ { // appendを呼んでるから遅いかなとも思うけど、高々4回しか回さない
				ops = append(ops, NoOperation{})
			}
		}
		hdr.OffsetCtrFlag = makeOffsetCtrlFlag(uint8(int(unsafe.Sizeof(hdr)) + opLength + eolPadding), flag)
		binary.Write(buf, binary.BigEndian, &hdr)
		binary.Write(buf, binary.BigEndian, &ops)
	}else {
		hdr.OffsetCtrFlag = makeOffsetCtrlFlag(uint8(unsafe.Sizeof(hdr)), flag)
		binary.Write(buf, binary.BigEndian, &hdr)
	}

	binary.Write(buf, binary.BigEndian, data)

	packet := packet{
		header: hdr,
		option: ops,
		data: data,
	}

	iface := getAppropriateInterface(conn.Cb.Addr, peer.Addr)
	b := buf.Bytes()
	packet.Checksum = net.CheckSum16(b, len(b), pseudoHeaderSum(iface.Address(), peer.Addr, len(b)))
	binary.BigEndian.PutUint16(b[16:18], packet.Checksum)
	packet.dump()
	return iface.Tx(net.ProtocolNumberTCP, b, peer.Addr)
}

func larger(x uint32, y uint32) uint32 {
	if x < y {
		return y
	}
	return x
}