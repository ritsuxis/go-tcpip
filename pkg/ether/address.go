package ethernet

import (
	// "fmt"
	"fmt"
	"strconv"
	"strings"
)

const AddressLength = 6

type Address [AddressLength]byte // 6byte

var (
	EmptyAddress = Address{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	InvalidAddress = Address{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	BroadcastAddress = Address{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
)

// 後で見る
// 多分一々登録しなくても済むようにしてる？
func NewAddress(b []byte) Address {
	var ret Address
	copy(ret[:], b)
	return ret
}

func ParseAddress(s string) Address {
	// :か-を区切り文字としてパース
	// e.g.  00:00:00:00:00:00
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == ':' || r == '-'
	})
	ret := Address{}
	if len(parts) != AddressLength {
		return InvalidAddress
	}
	for i, part := range parts {
		u, err := strconv.ParseUint(part, 16, 8) // 文字列を16進数と見て8bit長でパース
		if err != nil {
			return InvalidAddress
		}
		ret[i] = byte(u)
	}
	return ret
}

func (a Address) isGroupAddress() bool {
	return (a[0] & 0x01) != 0
}

func (a Address) Bytes() []byte {
	return a[:]
}

func (a Address) Len() uint8 {
	return uint8(len(a))
}

func (a Address) String() string {
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", a[0], a[1], a[2], a[3], a[4], a[5])
}