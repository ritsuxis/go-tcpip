package tcp

import (
	"container/list"
	"fmt"
	"log"
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
	State   State
	Number chan SeqAck 
}

type SeqAck struct {
	Seq uint32
	Ack uint32
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

func (repo *cbRepository) lookupUnlocked(addr *Address) *cbEntry {
	for elem := repo.list.Front(); elem != nil; elem = elem.Next() {
		entry := elem.Value.(*cbEntry)
		if entry.Port == addr.Port && (entry.Addr.IsEmpty() || entry.Addr == addr.Addr) {
			return entry
		}
	}
	return nil
}

func (repo *cbRepository) lookup(addr *Address) *cbEntry {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()
	return repo.lookupUnlocked(addr)
}

func (repo *cbRepository) getAvailablePort(addr net.ProtocolAddress) uint16 {
	var port uint16
	for port = 40000; port <= 65535; port++ {
		var elem *list.Element
		for elem = repo.list.Front(); elem != nil; elem = elem.Next() {
			entry := elem.Value.(*cbEntry)
			if entry.Port == port && (entry.Addr.IsEmpty() || entry.Addr == addr) {
				break
			}
		}
		if elem == nil {
			return port
		}
	}
	return 0
}

func (repo *cbRepository) add(addr *Address) *cbEntry {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	if addr.Port == 0 {
		addr.Port = repo.getAvailablePort(addr.Addr)
		if addr.Port == 0 {
			return nil
		}
	} else {
		if repo.lookupUnlocked(addr) != nil {
			fmt.Println("entry exists")
		}
	}
	entry := &cbEntry{
		Address: addr,
		rxQueue: make(chan *queueEntry),
		State: NewTcpState(),
		Number: make(chan SeqAck),
	}
	log.Printf("Add to cbRepository: %s", entry)
	repo.list.PushBack(entry)
	return entry
}

func (repo *cbRepository) del(entry *cbEntry) bool {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	for elem := repo.list.Front(); elem != nil; elem = elem.Next() {
		if elem.Value == entry {
			repo.list.Remove(elem)
			return true
		}
	}
	return false
}
