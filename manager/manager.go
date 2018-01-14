package manager

import (
	"bufio"
	"encoding/json"
	"net"
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
	listener, err := net.Listen("tcp", ":12345")
	if err != nil {
		panic(err.Error())
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.E(err.Error())
		} else {
			log.V("Connection to", conn.RemoteAddr(), "created")
			go m.handleConnection(conn)
		}

	}
}

func (m *Manager) handleConnection(conn net.Conn) {
	defer conn.Close()
	for {
		message, err := m.recvFrom(conn)
		if err != nil {
			log.V("Connection to", conn.RemoteAddr(), "closed")
			break //connection closed
		}
		response := m.handleMessage(message)
		m.send(response, conn)
	}

}

func (m *Manager) handleMessage(message *constants.NodeToManagerJSON) *constants.ManagerToNodeJSON {
	resp := new(constants.ManagerToNodeJSON)
	var id uint64
	if message.ID == constants.ID_NOT_ASSIGNED {
		id = m.newNode()
	} else {
		id = message.ID
	}

	if message.SampleValid {
		resp.PerformSample = false
		m.updateNodeValue(id, message.SampleValue)
		log.D("Receieved new sample from ID", id, ", value= ", message.SampleValue)
	} else {
		m.updateNodeLocation(id, message.Latitude, message.Longitude)
		resp.PerformSample = true
	}
	resp.AssignedID = id
	resp.NextCheckin = 2
	return resp

}

func (m *Manager) send(response *constants.ManagerToNodeJSON, conn net.Conn) {
	messageJSON, err := json.Marshal(response)
	if err != nil {
		panic(err.Error())
	}
	conn.Write(messageJSON)
	conn.Write(constants.DelimJSON)
	log.D("Sent: ", *response)
}
func (m *Manager) recvFrom(conn net.Conn) (*constants.NodeToManagerJSON, error) {
	bytes, err := bufio.NewReader(conn).ReadBytes(constants.DelimJSON[0])
	if err != nil {
		return nil, err
	}

	data := new(constants.NodeToManagerJSON)
	json.Unmarshal(bytes[:len(bytes)-1], data)
	log.D("Received: ", *data)
	return data, nil
}
