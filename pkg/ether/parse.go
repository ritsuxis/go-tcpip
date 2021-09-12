package ethernet

import (
	"bytes"
	"encoding/binary"

	"github.com/ritsuxis/go-tcpip/pkg/net"
)

// Ethernet header
type header struct {
	Dst  Address
	Src  Address
	Type net.EthernetType
}

// Ethernet frame
type Frame struct {
	header
	payload []byte
}

func parse(data []byte) (*Frame, error) {
	frame := Frame{}
	buf := bytes.NewBuffer(data)
	// 来たデータはネットワークの標準でビッグエンディアンになっている
	// bufに格納したdataをheaderの分だけ読みだして構造体に埋める
	if err := binary.Read(buf, binary.BigEndian, &frame.header); err != nil {
		return nil, err
	}
	// 残りの部分は全部payloadに読みだす
	frame.payload = buf.Bytes()
	return &frame, nil
}
