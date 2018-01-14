package nodesim

import (
	"bufio"
	"encoding/json"
	"math/rand"
	"net"
	"strconv"
	"time"

	"github.com/abemac/bomb-detection/constants"
)

//Node represents a node, which is a sensor
type Node struct {
	managerIP   string
	managerPort uint16
	assignedID  uint64
}

//NewNode creates a new Node
func NewNode() {
	n := new(Node)
	n.managerIP = "127.0.0.1"
	n.managerPort = 12345
	n.assignedID = constants.ID_NOT_ASSIGNED
	go n.mainLoop()
	log.V("New Node created")
}
func CreateNodes(number uint64) {
	var i uint64
	for i = 0; i < number; i++ {
		NewNode()
	}
	log.I(number, "new nodes created, they are active")
}

func (n *Node) mainLoop() {
	for {
		conn, err := n.connectToManager()
		if err != nil {
			panic(err.Error())
		}
		n.sendInfo(conn)
		resp, err := n.recvFrom(conn)
		if err != nil {
			panic(err.Error())
		}
		shouldSample := n.handleResponse(resp)
		if shouldSample {
			n.sendSample(conn)
		}
		if resp.NextCheckin > 0 {
			conn.Close()
			log.V("Connection to", conn.RemoteAddr(), "closed")
		}

		time.Sleep(time.Second * time.Duration(resp.NextCheckin))
	}
}

func (n *Node) connectToManager() (net.Conn, error) {

	for {
		conn, err := net.Dial("tcp", n.managerIP+":"+strconv.Itoa(int(n.managerPort)))
		if err != nil {
			time.Sleep(time.Second * 5)
		} else {
			log.V("Connection to", conn.RemoteAddr(), "created")
			return conn, err
		}
	}

}

func (n *Node) sendInfo(conn net.Conn) {
	message := new(constants.NodeToManagerJSON)
	message.Latitude, message.Longitude = n.getGPSLoc()
	message.ID = n.assignedID
	message.SampleValid = false
	messageJSON, err := json.Marshal(message)
	if err != nil {
		panic(err.Error())
	}
	conn.Write(messageJSON)
	conn.Write(constants.DelimJSON)
	log.D("Sent: ", *message)
}

func (n *Node) sendSample(conn net.Conn) {
	message := new(constants.NodeToManagerJSON)
	message.ID = n.assignedID
	message.SampleValue = n.sample()
	message.SampleValid = true
	messageJSON, err := json.Marshal(message)
	if err != nil {
		panic(err.Error())
	}
	conn.Write(messageJSON)
	conn.Write(constants.DelimJSON)
	log.D("Sent: ", *message)
}

func (n *Node) recvFrom(conn net.Conn) (*constants.ManagerToNodeJSON, error) {
	bytes, err := bufio.NewReader(conn).ReadBytes(constants.DelimJSON[0])
	if err != nil {
		return nil, err
	}
	data := new(constants.ManagerToNodeJSON)
	json.Unmarshal(bytes[:len(bytes)-1], data)
	log.D("Received: ", *data)
	return data, nil
}

func (n *Node) handleResponse(msg *constants.ManagerToNodeJSON) bool {
	n.assignedID = msg.AssignedID
	return msg.PerformSample

}

func (n *Node) sample() int {
	return rand.Intn(100)
}
func (n *Node) getGPSLoc() (float64, float64) {
	return rand.Float64(), rand.Float64()
}
