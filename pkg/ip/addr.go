package ip

import (
	"fmt"
	"strconv"
	"strings"
	"unsafe"
)

type Address [4]byte // 32bit

var (
	EmptyAddress     = Address{0x00, 0x00, 0x00, 0x00}
	InvalidAddress   = Address{0x00, 0x00, 0x00, 0x00}
	BroadcastAddress = Address{0xff, 0xff, 0xff, 0xff}
)

func ParseAddress(s string) Address {
	// '.'でパース
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == '.'
	})
	if len(parts) != 4 {
		return InvalidAddress
	}
	ret := Address{}
	for i, part := range parts {
		u, err := strconv.ParseUint(part, 10, 8) // uint64
		if err != nil {
			return InvalidAddress
		}
		ret[i] = uint8(u & 0xff) // ちゃんと8bitで区切るために論理積取ってる?
	}
	return ret
}

func (a Address) Bytes() []byte {
	return a[:]
}

func (a Address) Len() uint8 {
	return uint8(len(a))
}

func (a Address) String() string {
	return fmt.Sprintf("%d.%d.%d.%d", a[0], a[1], a[2], a[3])
}

func (a Address) Uint32() uint32 {
	return *(*uint32)(unsafe.Pointer(&a[0]))
}

func (a Address) IsEmpty() bool {
	return a == EmptyAddress
}