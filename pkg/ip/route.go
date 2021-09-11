package ip

import (
	"encoding/binary"
	"log"
	"sync"
)

// ルーティングテーブルを作るときも排他制御が必要
// https://www.infraexpert.com/study/routing2.html

type routeEntry struct {
	network Address
	netmask Address
	nexthop Address
	iface   *Interface
}

type routeTable struct {
	storage []*routeEntry
	mutex   sync.RWMutex
}

var repo *routeTable

func newRouteTable() *routeTable {
	return &routeTable{
		storage: make([]*routeEntry, 0, 1024),
	}
}

func (tbl *routeTable) add(iface *Interface, network, netmask, nexthop Address) {
	tbl.mutex.Lock()
	tbl.storage = append(tbl.storage, &routeEntry{network, netmask, nexthop, iface})
	tbl.mutex.Unlock()
}

func (tbl *routeTable) del(iface *Interface) {
	// 新しいrouteTableを用意して、引数に渡されたiface以外を登録する
	newStorage := make([]*routeEntry, 0, cap(tbl.storage))
	tbl.mutex.Lock()
	for _, entry := range tbl.storage {
		if entry.iface != iface {
			newStorage = append(newStorage, entry)
		}
	}
	tbl.storage = newStorage
	tbl.mutex.Unlock()
}

func (tbl *routeTable) lookup(iface *Interface, dst Address) *routeEntry {
	var candidate *routeEntry
	tbl.mutex.RLock()
	for _, entry := range tbl.storage {
		// network addressが同じかつ、同じインターフェース(もしくはインターフェースがない)なら
		if dst.Uint32()&entry.netmask.Uint32() == entry.network.Uint32() && (iface == nil || entry.iface == iface) {
			// candidateがnilまたはentryがcandidateよりもより大きいnetmaskを持っているなら(つまり、より小さいホスト部を持つなら)
			if candidate == nil || ntoh32(candidate.netmask.Bytes()) < ntoh32(entry.netmask.Bytes()) {
				candidate = entry
			}
		}
	}
	tbl.mutex.RUnlock()
	return candidate
}

func (tbl *routeTable) length() int {
	tbl.mutex.RLock()
	defer tbl.mutex.RUnlock()
	return len(tbl.storage)
}

func (tbl *routeTable) dump() {
	tbl.mutex.RLock()
	defer tbl.mutex.RUnlock()
	log.Printf("route table dump: %d entries\n", len(tbl.storage))
	for _, entry := range tbl.storage {
		log.Printf("network=%s, netmask=%s, nexthop=%s, iface=%s\n", entry.network, entry.netmask, entry.nexthop, entry.iface.Device().Name())
	}
}

func ntoh32(i []byte) uint32 {
	return binary.BigEndian.Uint32(i)
}