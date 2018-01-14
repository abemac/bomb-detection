package manager

import (
	"bufio"
	"container/list"
	"encoding/json"
	"net"

	"github.com/abemac/bomb-detection/manager/constants"
	"github.com/abemac/bomb-detection/manager/logger"
)

var log = logger.NewLogger("Manager.go")

//Manager represents the manager that collects and processes node data
type Manager struct {
	connectedNodes *list.List
}

//NewManager creates a new Manager
func NewManager() *Manager {
	m := new(Manager)
	m.connectedNodes = list.New()
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
func (m *Manager) handleMessage(message *constants.NodeToManagerJSON) *constants.ManagerToNodeJSON {
	resp := new(constants.ManagerToNodeJSON)
	if message.SampleValid {
		resp.PerformSample = false
		//handle new sample data
		log.D("Receieved new sample from ID", message.ID, ", value= ", message.SampleValue)
	} else {
		resp.PerformSample = true
	}
	resp.NextCheckin = 2
	return resp

}
