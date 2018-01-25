package manager

import (
	"encoding/json"
	"sync"
	"sync/atomic"

	"github.com/abemac/bomb-detection/constants"
)

//Manager represents the manager that collects and processes node data
type Manager struct {
	nodes              map[uint64]*node
	supernodes         map[uint64]*supernode
	lastAssignedNodeID uint64
	nodesMutex         *sync.RWMutex
	supernodesMutex    *sync.RWMutex
}

//NewManager creates a new Manager
func NewManager() *Manager {
	m := new(Manager)
	m.nodes = make(map[uint64]*node)
	m.supernodes = make(map[uint64]*supernode)
	m.lastAssignedNodeID = 0
	m.nodesMutex = new(sync.RWMutex)
	m.supernodesMutex = new(sync.RWMutex)
	return m
}

func (m *Manager) Run() {
	log.I("Manager Started")
	go NewWebUI(m).Run()
	NewJSONRequestServer(m, 12345).Run()
}

func (m *Manager) handleMessage(bytes []byte) []byte {
	message := new(constants.NodeToManagerJSON)
	json.Unmarshal(bytes[:len(bytes)-1], message)
	log.D("Received: ", *message)

	var id uint64

	if message.SuperNode {
		if message.ID == constants.ID_NOT_ASSIGNED {
			id = m.newSuperNode()
		} else {
			id = m.checkIfSupernodeIDValid(message.ID)
		}

		if message.SampleValid {
			m.updateSuperNodeValue(id, message.SampleValue)
			log.D("Receieved new sample from ID", id, ", value= ", message.SampleValue)
			return nil
		}
		response := new(constants.ManagerToNodeJSON)
		m.updateSuperNodeLocation(id, message.Latitude, message.Longitude)
		response.PerformSample = true
		response.AssignedID = id
		response.NextCheckin = 2
		responseBytes, err := json.Marshal(response)
		if err != nil {
			panic(err.Error())
		}
		log.D("Sending: ", *response)
		return responseBytes
	}

	if message.ID == constants.ID_NOT_ASSIGNED {
		id = m.newNode()
	} else {
		id = m.checkIfNodeIDValid(message.ID)
	}

	if message.SampleValid {
		m.updateNodeValue(id, message.SampleValue)
		log.D("Receieved new sample from ID", id, ", value= ", message.SampleValue)
		return nil
	}
	response := new(constants.ManagerToNodeJSON)
	m.updateNodeLocation(id, message.Latitude, message.Longitude)
	response.PerformSample = true
	response.AssignedID = id
	response.NextCheckin = 2
	responseBytes, err := json.Marshal(response)
	if err != nil {
		panic(err.Error())
	}
	log.D("Sending: ", *response)
	return responseBytes

}
func (m *Manager) getNewID() uint64 {
	//18,446,744,073,709,551,615 possible ids (18 quintillion)
	//if everyone on the planet carried around 2 million smartphones there would still be ids to give
	id := atomic.AddUint64(&m.lastAssignedNodeID, 1)
	return id
}
