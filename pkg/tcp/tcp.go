package tcp

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	"github.com/ritsuxis/go-tcpip/pkg/ip"
	"github.com/ritsuxis/go-tcpip/pkg/net"
)

func init() {
	ip.RegisterProtocol(net.ProtocolNumberTCP, rxHandler)
	repo = newCbRepository()
}

func Init() {
	// do nothing
}

func rxHandler(iface net.ProtocolInterface, data []byte, src, dst net.ProtocolAddress) error {
	// Proto
	return nil
}

// Used for debugging purposes only
func Build(src, dst uint16, seq, ack uint32, flag ControlFlag, ws, urgent uint16, data []byte, ops Options) []byte {
	hdr := header {
		SourcePort: src,
		DestinationPort: dst,
		SequenceNumber: seq,
		ACKNumber: ack,
		WindowSize: ws,
		Urgent: urgent,
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
	b := buf.Bytes()
	peer := &Address{
		Addr: ip.ParseAddress("192.0.2.1"),
	}
	packet.Checksum = net.CheckSum16(b, len(b), pseudoHeaderSum(peer.Addr, peer.Addr, len(b)))
	binary.BigEndian.PutUint16(b[16:18], packet.Checksum)
	packet.dump()
	return b
}