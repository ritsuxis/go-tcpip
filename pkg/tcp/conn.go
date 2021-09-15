package tcp

import (
	// "bytes"
	// "encoding/binary"

	// "github.com/ritsuxis/go-tcpip/pkg/net"
)

type Conn struct {
	cb   *cbEntry
	peer *Address
}

// func (conn *Conn) WriteTo(data []byte, peer *Address) error {
// 	hdr := header{}
// 	hdr.SourcePort = conn.cb.Port
// 	hdr.DestinationPort = peer.Port
// 	hdr.OffsetCtrFlag = makeOffsetCtrlFlag() // uint16(int(unsafe.Sizeof(hdr)) + len(data))
// 	datagram := datagram{
// 		header: hdr,
// 		data:   data,
// 	}
// 	datagram.dump()
// 	buf := new(bytes.Buffer)
// 	binary.Write(buf, binary.BigEndian, &hdr)
// 	binary.Write(buf, binary.BigEndian, data)
// 	iface := getAppropriateInterface(conn.cb.Addr, peer.Addr)
// 	b := buf.Bytes()
// 	datagram.Checksum = net.CheckSum16(b, len(b), pseudoHeaderSum(iface.Address(), peer.Addr, len(b)))
// 	binary.BigEndian.PutUint16(b[6:8], datagram.Checksum)
// 	return iface.Tx(net.ProtocolNumberUDP, b, peer.Addr)
// }