package ethernet

import (
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

