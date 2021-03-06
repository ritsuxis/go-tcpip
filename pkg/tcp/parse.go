package tcp

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"unsafe"

	"github.com/ritsuxis/go-tcpip/pkg/net"
)

/*
TCP Header Format

    0                   1                   2                   3
    0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |          Source Port          |       Destination Port        |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                        Sequence Number                        |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                    Acknowledgment Number                      |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |  Data |           |U|A|P|R|S|F|                               |
   | Offset| Reserved  |R|C|S|S|Y|I|            Window             |
   |       |           |G|K|H|T|N|N|                               |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |           Checksum            |         Urgent Pointer        |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                    Options                    |    Padding    |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                             data                              |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
*/

type header struct {
	SourcePort      uint16
	DestinationPort uint16
	SequenceNumber  uint32
	ACKNumber       uint32
	OffsetCtrFlag
	WindowSize uint16
	Checksum   uint16
	Urgent     uint16 // 緊急ポインタに関する実装はしない
}

type packet struct {
	header
	option Options
	data   []byte
}

type OffsetCtrFlag uint16

func (p packet) Dump() {
	p.dump()
}

func (p packet) dump() {
	log.Printf(">>>>>>>>>>TCP>>>>>>>>>>>")
	log.Printf("     src port: %d\n", p.SourcePort)
	log.Printf("     dst port: %d\n", p.DestinationPort)
	log.Printf("   seq number: %d\n", p.SequenceNumber)
	log.Printf("   ack number: %d\n", p.ACKNumber)
	log.Printf("       offset: %d\n", p.OffsetCtrFlag.Offset())
	log.Printf("         flag: %s\n", p.OffsetCtrFlag.ControlFlag().String())
	log.Printf("  window size: %d\n", p.WindowSize)
	log.Printf("     checksum: 0x%04x\n", p.Checksum)
	log.Printf("urgent number: %d\n", p.Urgent)
	fmt.Println(hex.Dump(p.data))
	log.Printf(">>>>>>>>>>>>>>>>>>>>>>>>")
}

func pseudoHeaderSum(src, dst net.ProtocolAddress, n int) uint32 {
	pseudo := new(bytes.Buffer)
	binary.Write(pseudo, binary.BigEndian, src.Bytes())
	binary.Write(pseudo, binary.BigEndian, dst.Bytes())
	binary.Write(pseudo, binary.BigEndian, uint16(net.ProtocolNumberTCP))
	binary.Write(pseudo, binary.BigEndian, uint16(n))
	return uint32(^net.CheckSum16(pseudo.Bytes(), pseudo.Len(), 0))
}

func Parse(data []byte, src, dst net.ProtocolAddress) (*packet, ControlFlag, error) {
	hdr := header{}
	if len(data) < int(unsafe.Sizeof(hdr)) {
		return nil, 0, fmt.Errorf("message is too short: %d", len(data))
	}
	buf := bytes.NewBuffer(data)
	if err := binary.Read(buf, binary.BigEndian, &hdr); err != nil {
		return nil, 0, err
	}

	sum := net.CheckSum16(data, len(data), pseudoHeaderSum(src, dst, len(data)))
	if sum != 0 {
		return nil, 0, fmt.Errorf("tcp checksum failure: 0x%04x", sum)
	}

	// // parseする際にoffsetsizeを読めばどこからがdataがわかる
	// // headerはoption抜きだと20byteのため、その差分がheader
	// if opsize := hdr.OffsetCtrFlag.Offset() - 20; opsize > 0 {
	// 	log.Printf("Offset size: %d byte", hdr.OffsetCtrFlag.Offset())
	// 	log.Printf("Option size: %d byte", opsize)
	// 	ops := make(Options, opsize)
	// 	if err := binary.Read(buf, binary.BigEndian, &ops); err != nil {
	// 		return nil, 0, err
	// 	}
	// 	return &packet{
	// 		header: hdr,
	// 		option: ops,
	// 		data:   buf.Bytes(),
	// 	}, hdr.OffsetCtrFlag.ControlFlag() ,nil
	// } else {
	// 	return &packet{
	// 		header: hdr,
	// 		data:   buf.Bytes(),
	// 	}, hdr.OffsetCtrFlag.ControlFlag(), nil
	// }
	pk := packet{
		header: hdr,
		data:   buf.Bytes(),
	}
	pk.dump()

	return &pk, hdr.OffsetCtrFlag.ControlFlag(), nil
}

func makeOffsetCtrlFlag(offset uint8, flag ControlFlag) OffsetCtrFlag {
	return OffsetCtrFlag(uint16(offset/4)<<12 | uint16(flag)) // offsetは32bit word
}

// TCP header Length
func (ocf OffsetCtrFlag) Offset() int {
	ocf8 := uint(ocf >> 8)
	return 4 * int(ocf8>>4)
}

func (ocf OffsetCtrFlag) ControlFlag() ControlFlag {
	return ControlFlag(uint8(ocf))
}
