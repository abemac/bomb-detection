package manager

import (
	"fmt"
	"sync"
	"time"
)

type node struct {
	Value             int
	Latitude          float64
	Longitude         float64
	BatteryPercentage float32
	lastSampleTime    int64 //unix timestamp
	mutex             *sync.RWMutex
	row               int
	col               int
}

func (m *Manager) newNode() uint64 {
	n := new(node)
	n.mutex = new(sync.RWMutex)
	n.row = -1
	n.col = -1
	id := m.getNewID()
	m.nodesMutex.Lock()
	m.nodes[id] = n
	m.nodesMutex.Unlock()
	return id
}
func (m *Manager) checkIfNodeIDValid(id uint64, muid int64) uint64 {
	m.nodesMutex.RLock()
	_, ok := m.nodes[id]
	m.nodesMutex.RUnlock()
	ok = ok && (muid == m.uid)

	if !ok {
		//panic("AHHH") //nodes already assigned ID but manager crashed possibly
		id = m.newNode()
	}
	return id
}
func (m *Manager) updateNodeValue(id uint64, newSample int) {
	m.nodesMutex.RLock()
	n := m.nodes[id]
	m.nodesMutex.RUnlock()
	n.mutex.Lock()
	n.Value = newSample
	n.lastSampleTime = time.Now().Unix()
	n.mutex.Unlock()
}
func (m *Manager) updateNodeBatteryLevel(id uint64, batteryLevel float32) {
	m.nodesMutex.RLock()
	n := m.nodes[id]
	m.nodesMutex.RUnlock()
	n.mutex.Lock()
	n.BatteryPercentage = batteryLevel
	n.mutex.Unlock()
}
func (m *Manager) updateNodeLocation(id uint64, latitude float64, longitude float64) {
	m.nodesMutex.RLock()
	n := m.nodes[id]
	m.nodesMutex.RUnlock()
	n.mutex.Lock()
	newRow := getRowFromLat(latitude)
	newCol := getColFromLong(longitude)
	if newRow != n.row || newCol != n.col {
		blocksMutex.Lock()
		if n.row != -1 {
			pq.updateCount(blocks[n.row][n.col], blocks[n.row][n.col].count-1)
		}
		pq.updateCount(blocks[newRow][newCol], blocks[newRow][newCol].count+1)
		pq.updateVisitedTime(blocks[newRow][newCol], time.Now().Unix())
		n.row = newRow
		n.col = newCol
		blocksMutex.Unlock()
	}
	n.Latitude = latitude
	n.Longitude = longitude
	n.mutex.Unlock()
}
func (m *Manager) removeNode(id uint64) {
	m.nodesMutex.Lock()
	delete(m.nodes, id)
	m.nodesMutex.Unlock()
}
func (m *Manager) periodicallyPrintNodes(sleepTime int) {
	for {
		m.printNodes()
		time.Sleep(time.Second * time.Duration(sleepTime))
	}
}
func (m *Manager) printNodes() {
	m.nodesMutex.RLock()
	for key, val := range m.nodes {
		val.mutex.RLock()
		fmt.Println(*val, "Id:", key)
		val.mutex.RUnlock()
	}
	m.nodesMutex.RUnlock()
}
