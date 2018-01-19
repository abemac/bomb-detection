package manager

import (
	"encoding/json"
	"sync"

	"github.com/abemac/bomb-detection/constants"
)

//Manager represents the manager that collects and processes node data
type Manager struct {
	nodes              map[uint64]*node
	lastAssignedNodeID uint64
	mapMutex           *sync.RWMutex
	idMutex            *sync.Mutex
}

//NewManager creates a new Manager
func NewManager() *Manager {
	m := new(Manager)
	m.nodes = make(map[uint64]*node)
	m.lastAssignedNodeID = 0
	m.mapMutex = new(sync.RWMutex)
	m.idMutex = new(sync.Mutex)
	return m
}

func (m *Manager) Run() {
	log.I("Manager Started")
	go NewWebAPI(m).Run()
	NewJSONRequestServer(m, 12345).Run()
}

func (m *Manager) handleMessage(bytes []byte) []byte {
	message := new(constants.NodeToManagerJSON)
	json.Unmarshal(bytes[:len(bytes)-1], message)
	log.D("Received: ", *message)

	var id uint64
	if message.ID == constants.ID_NOT_ASSIGNED {
		id = m.newNode()
	} else {
		id = message.ID
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
