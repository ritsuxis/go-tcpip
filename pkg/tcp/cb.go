package tcp

import (
	"container/list"
	"sync"

	"github.com/ritsuxis/go-tcpip/pkg/net"
)

type queueEntry struct {
	addr net.ProtocolAddress
	port uint16
	data []byte
}

// callback entry
type cbEntry struct {
	*Address
	rxQueue chan *queueEntry
}

type cbRepository struct {
	list  *list.List
	mutex sync.RWMutex
}

var repo *cbRepository

func newCbRepository() *cbRepository {
	return &cbRepository{
		list: list.New(),
	}
}