package manager

import (
	"fmt"
	"sync"
	"time"
)

type node struct {
	Value          int
	Latitude       float64
	Longitude      float64
	lastSampleTime int64 //unix timestamp
	mutex          *sync.RWMutex
}

func (m *Manager) newNode() uint64 {
	n := new(node)
	n.mutex = new(sync.RWMutex)
	id := m.getNewID()
	m.nodesMutex.Lock()
	m.nodes[id] = n
	m.nodesMutex.Unlock()
	return id
}
func (m *Manager) removeNode(id uint64) {
	m.nodesMutex.Lock()
	delete(m.nodes, id)
	m.nodesMutex.Unlock()
}
func (m *Manager) checkIfNodeIDValid(id uint64) uint64 {
	m.nodesMutex.RLock()
	_, ok := m.nodes[id]
	m.nodesMutex.RUnlock()
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
	n.updateValue(newSample)
}
func (m *Manager) updateNodeLocation(id uint64, latitude float64, longitude float64) {
	m.nodesMutex.RLock()
	n := m.nodes[id]
	m.nodesMutex.RUnlock()
	n.mutex.Lock()
	n.Latitude = latitude
	n.Longitude = longitude
	n.mutex.Unlock()
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

func (m *Manager) periodicallyPrintNodes(sleepTime int) {
	for {
		m.printNodes()
		time.Sleep(time.Second * time.Duration(sleepTime))
	}
}

func (n *node) updateValue(newSample int) {
	n.mutex.Lock()
	n.Value = newSample
	n.lastSampleTime = time.Now().Unix()
	n.mutex.Unlock()
}
