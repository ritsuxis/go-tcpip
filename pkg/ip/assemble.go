package ip

import "github.com/ritsuxis/go-tcpip/pkg/net"

type assembler struct {
	protocol net.ProtocolNumber
	data     []byte
	src      net.ProtocolAddress
	dst      net.ProtocolAddress
	id       uint16
	mtu      int
}

// datagramごとに変わるもの
func newAssembler(protocol net.ProtocolNumber, data []byte, src, dst net.ProtocolAddress, id uint16, mtu int) *assembler {
	return &assembler{
		protocol: protocol,
		data:     data,
		src:      src,
		dst:      dst,
		id:       id,
		mtu:      mtu,
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ipのpayloadの作成
func (a *assembler) assemble() []*datagram {
	ret := []*datagram{}
	var n = len(a.data)
	var slen int
	for done := 0; done < n; done += slen {
		slen = min((n - done), a.mtu) // データサイズがmtuより大きいなら小分けする
		var flag uint16
		if done+slen < n {
			flag = 0x2000
		}
		offset := flag | uint16((done>>3)&0x1ffff)
		var hlen = 20 // optionを指定しないので20で固定する
		var data = a.data[done : done+slen]
		datagram := &datagram{
			header: header{
				VHL:      uint8((4 << 4) | (hlen >> 2)),
				TOS:      0,
				Len:      uint16(hlen + len(data)),
				Id:       a.id,
				Offset:   offset,
				TTL:      0xff,
				Protocol: a.protocol,
				Sum:      0,
				Src:      a.src.(Address),
				Dst:      a.dst.(Address),
			},
			payload: data,
		}
		ret = append(ret, datagram) // 小分けしたdatagramをretに纏める
	}
	return ret
}
