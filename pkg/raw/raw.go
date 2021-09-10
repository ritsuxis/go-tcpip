package raw

import "io"

// rawフォルダ内のtuntap.go等で作成されるデバイスを操作するためのInterface
type Device interface {
	io.ReadWriteCloser
	Name() string
	Address() []byte
}