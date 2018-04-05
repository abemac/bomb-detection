package manager

import (
	"container/heap"
	"fmt"
	"math"
	"sync"
	"time"
)

type supernode struct {
	Value          int
	Latitude       float64
	Longitude      float64
	lastSampleTime int64 //unix timestamp
	mutex          *sync.RWMutex
}

func (m *Manager) newSuperNode() uint64 {
	n := new(supernode)
	n.mutex = new(sync.RWMutex)
	id := m.getNewID()
	m.supernodesMutex.Lock()
	m.supernodes[id] = n
	m.supernodesMutex.Unlock()
	return id
}

func getBestSuperNodeLocation(m *Manager, currentLat float64, currentLong float64) (float64, float64) {
	blocksMutex.Lock()
	defer blocksMutex.Unlock()
	numToConsider := int(math.Min(10, numBlockCols*numBlockRows))
	possibleBlocks := make([]*Block, numToConsider)
	for i := 0; i < numToConsider; i++ {
		possibleBlocks[i] = heap.Pop(&pq).(*Block)
	}
	closest := 0
	dist := math.MaxFloat64
	for i := 0; i < numToConsider; i++ {
		latdist := getBlockLat(possibleBlocks[i].row) - currentLat
		if latdist > 180 {
			latdist -= 180
		} else if latdist < -180 {
			latdist += 180
		}
		longdist := getBlockLong(possibleBlocks[i].col) - currentLong
		if latdist > 360 {
			latdist -= 360
		} else if latdist < -360 {
			latdist += 360
		}
		thisdist := math.Sqrt(latdist*latdist + longdist*longdist)
		if thisdist < dist {
			dist = thisdist
			closest = i
		}
	}
	b := possibleBlocks[closest]
	pq.updateVisitedTime(b, time.Now().Unix())
	for i := 0; i < numToConsider; i++ {
		heap.Push(&pq, possibleBlocks[i])
	}
	return getBlockLat(b.row), getBlockLong(b.col)
}

func (m *Manager) checkIfSupernodeIDValid(id uint64, muid int64) uint64 {
	m.supernodesMutex.RLock()
	_, ok := m.supernodes[id]
	m.supernodesMutex.RUnlock()
	ok = ok && (muid == m.uid)
	if !ok {
		//panic("AHHH") //nodes already assigned ID but manager crashed possibly
		id = m.newSuperNode()
	}
	return id
}
func (m *Manager) updateSuperNodeValue(id uint64, newSample int) {
	m.supernodesMutex.RLock()
	n := m.supernodes[id]
	m.supernodesMutex.RUnlock()
	n.mutex.Lock()
	n.Value = newSample
	n.lastSampleTime = time.Now().Unix()
	n.mutex.Unlock()
}
func (m *Manager) updateSuperNodeLocation(id uint64, latitude float64, longitude float64) {
	m.supernodesMutex.RLock()
	n := m.supernodes[id]
	m.supernodesMutex.RUnlock()
	n.mutex.Lock()
	n.Latitude = latitude
	n.Longitude = longitude
	n.mutex.Unlock()

}
func (m *Manager) removeSuperNode(id uint64) {
	m.supernodesMutex.Lock()
	delete(m.supernodes, id)
	m.supernodesMutex.Unlock()
}
func (m *Manager) printSuperNodes() {
	m.supernodesMutex.RLock()
	for key, val := range m.supernodes {
		val.mutex.RLock()
		fmt.Println(*val, "Id:", key)
		val.mutex.RUnlock()
	}
	m.supernodesMutex.RUnlock()
}

func (m *Manager) periodicallyPrintSuperNodes(sleepTime int) {
	for {
		m.printSuperNodes()
		time.Sleep(time.Second * time.Duration(sleepTime))
	}
}
