package ioctl

import (
	"bytes"
	"syscall"
	"unsafe"
)

// ここのやつの利用を他のパッケージを頼らずにしている
// https://linuxjm.osdn.jp/html/LDP_man-pages/man7/netdevice.7.html

// socketのインターフェース番号(index)の取得
func SIOCGIFINDEX(name string) (int32, error) {
	// システムコールのsocketを呼び出す
	// socketは通信のためのエンドポイントを作成する
	// AF_INET はIPv4 のアドレスファミリ(プロトコルファミリ)(それで通信する際の約束事を纏めたもの)
	// SOCK_DGRAMはデータグラムのこと(コネクションレス、信頼性無し、固定最大長メッセージ)
	// つまりは、IPv4でデータグラムで通信するためのソケットを作る
	// https://linuxjm.osdn.jp/html/LDP_man-pages/man2/socket.2.html
	// ioctlを使うためにはソケットを開く必要があるためここでは開いている
	soc, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	if err != nil {
		return 0, err
	}
	defer syscall.Close(soc)
	ifreq := struct {
		name  [16]byte
		index int32
		__pad [22]byte
	}{}

	copy(ifreq.name[:syscall.IFNAMSIZ-1], name)
	// uintptrとunsafe.Pointerの説明( https://qiita.com/kitauji/items/291f16f619a939bd7b87 )
	// socのインターフェース番号(index)の取得
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(soc), syscall.SIOCGIFINDEX, uintptr(unsafe.Pointer(&ifreq))); errno != 0 {
		return 0, errno
	}
	return ifreq.index, err
}

// デバイスのアクティブフラグワードを取得する
func SIOCGIFFLAGS(name string) (uint16, error) {
	soc, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	if err != nil {
		return 0, err
	}
	defer syscall.Close(soc)
	ifreq := struct {
		name  [syscall.IFNAMSIZ]byte
		flags uint16
		_pad  [22]byte
	}{}
	copy(ifreq.name[:syscall.IFNAMSIZ-1], name)
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(soc), syscall.SIOCGIFFLAGS, uintptr(unsafe.Pointer(&ifreq))); errno != 0 {
		return 0, errno
	}
	return ifreq.flags, nil
}

// デバイスのアクティブフラグワードを設定する
func SIOCSIFFLAGS(name string, flags uint16) error {
	soc, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	if err != nil {
		return err
	}
	defer syscall.Close(soc)
	ifreq := struct {
		name  [syscall.IFNAMSIZ]byte
		flags uint16
		_pad  [22]byte
	}{}
	copy(ifreq.name[:syscall.IFNAMSIZ-1], name)
	ifreq.flags = flags
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(soc), syscall.SIOCSIFFLAGS, uintptr(unsafe.Pointer(&ifreq))); errno != 0 {
		return errno
	}
	return nil
}

// http://www.coins.tsukuba.ac.jp/~syspro/2012/2012-05-23/sockaddr.html
// https://linuxjm.osdn.jp/html/LDP_man-pages/man2/bind.2.html
type sockaddr struct {
	family uint16   // address family
	addr   [14]byte // 十分に大きい
}

// デバイスのハードウェアアドレスを取得する
func SIOCGIFHWADDR(name string) ([]byte, error) {
	soc, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	if err != nil {
		return nil, err
	}
	defer syscall.Close(soc)
	ifreq := struct {
		name [syscall.IFNAMSIZ]byte
		addr sockaddr
		_pad [8]byte
	}{}
	copy(ifreq.name[:syscall.IFNAMSIZ-1], name)
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(soc), syscall.SIOCGIFHWADDR, uintptr(unsafe.Pointer(&ifreq))); errno != 0 {
		return nil, errno
	}
	return ifreq.addr.addr[:], nil
}

// 仮想NICの作成
func TUNSETIFF(fd uintptr, name string, flags uint16) (string, error) {
	ifreq := struct {
		name  [syscall.IFNAMSIZ]byte
		flags uint16
		_pad  [22]byte
	}{}

	copy(ifreq.name[:syscall.IFNAMSIZ-1], []byte(name))
	// flagsの引数で貰ってるけど使わずに、実際には同じ値をここで再設定している
	ifreq.flags = syscall.IFF_TAP | syscall.IFF_NO_PI
	// 仮想NICの作成
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, syscall.TUNSETIFF, uintptr(unsafe.Pointer(&ifreq))); errno != 0 {
		return "", errno
	}
	// IndexByteは ifreq.nameのスライスの中で最初に0が出てくる所を返す
	// つまりifreq.nameに入っている名前の最後の部分+1を返す(初期化時点では全部0なので)
	return string(ifreq.name[:bytes.IndexByte(ifreq.name[:], 0)]), nil
}
