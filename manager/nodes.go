package manager

import (
	"fmt"
	"sync"
	"time"
)

type node struct {
	Value          int     `json:"value"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	SuperNode      bool    `json:"superNode"`
	lastSampleTime int64   //unix timestamp
	mutex          *sync.RWMutex
}

func (m *Manager) newNode() uint64 {
	n := new(node)
	n.mutex = new(sync.RWMutex)
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
	m.mapMutex.RLock()
	n := m.nodes[id]
	m.mapMutex.RUnlock()
	n.updateValue(newSample)
}
func (m *Manager) updateNodeLocation(id uint64, latitude float64, longitude float64) {
	m.mapMutex.RLock()
	n, ok := m.nodes[id]
	if !ok {
		panic("AHHH")
	}
	m.mapMutex.RUnlock()
	n.mutex.Lock()
	n.Latitude = latitude
	n.Longitude = longitude
	n.mutex.Unlock()

}
func (m *Manager) printNodes() {
	m.mapMutex.RLock()
	for key, val := range m.nodes {
		val.mutex.RLock()
		fmt.Println(*val, "Id:", key)
		val.mutex.RUnlock()
	}
	m.mapMutex.RUnlock()
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
