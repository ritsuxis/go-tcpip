package tcp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"unsafe"

	"github.com/ritsuxis/go-tcpip/pkg/ip"
	"github.com/ritsuxis/go-tcpip/pkg/net"
)

// "bytes"
// "encoding/binary"

// "github.com/ritsuxis/go-tcpip/pkg/net"

type Conn struct {
	cb   *cbEntry
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
		SourcePort: conn.cb.Port,
		DestinationPort: peer.Port,
		SequenceNumber: 0, // TODO: ctrl seq num
		ACKNumber: 25, // TODO: ctrl ack num
		WindowSize: 1000, // TODO: 受信側のサイズに合わせて変える
		Urgent: 0,
	}

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

	iface := getAppropriateInterface(conn.cb.Addr, peer.Addr)
	b := buf.Bytes()
	packet.Checksum = net.CheckSum16(b, len(b), pseudoHeaderSum(iface.Address(), peer.Addr, len(b)))
	binary.BigEndian.PutUint16(b[16:18], packet.Checksum)
	packet.dump()
	fmt.Println(b)

	return iface.Tx(net.ProtocolNumberTCP, b, peer.Addr)
}