package manager

import (
	"sync"
	"time"
)

type node struct {
	value          int
	latitude       float64
	longitude      float64
	superNode      bool
	lastSampleTime int64 //unix timestamp
	mutex          *sync.Mutex
}

func (m *Manager) newNode() uint64 {
	n := new(node)
	n.mutex = new(sync.Mutex)
	id := m.getNewNodeID()
	m.mapMutex.Lock()
	m.nodes[id] = n
	m.mapMutex.Unlock()
	return id
}
func (m *Manager) removeNode(id uint64) {
	m.mapMutex.Lock()
	delete(m.nodes, id)
	m.mapMutex.Unlock()
}
func (m *Manager) getNewNodeID() uint64 {
	//18,446,744,073,709,551,615 possible ids (18 quintillion)
	//if everyone on the planet carried around 2 million smartphones there would still be ids to give
	m.idMutex.Lock()
	defer m.idMutex.Unlock()
	m.lastAssignedNodeID++
	return m.lastAssignedNodeID
}
func (m *Manager) updateNodeValue(id uint64, newSample int) {
	m.mapMutex.Lock()
	n := m.nodes[id]
	m.mapMutex.Unlock()
	n.updateValue(newSample)
}
func (m *Manager) updateNodeLocation(id uint64, latitude float64, longitude float64) {
	m.mapMutex.Lock()
	n, ok := m.nodes[id]
	if !ok {
		panic("AHHH")
	}
	m.mapMutex.Unlock()
	n.mutex.Lock()
	n.latitude = latitude
	n.longitude = longitude
	n.mutex.Unlock()

}
func (n *node) updateValue(newSample int) {
	n.mutex.Lock()
	n.value = newSample
	n.lastSampleTime = time.Now().Unix()
	n.mutex.Unlock()
}
