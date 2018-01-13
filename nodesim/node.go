package nodesim

import (
	"math/rand"
	"net"
	"strconv"
	"time"
)

//Node represents a node, which is a sensor
type Node struct {
	managerIP   string
	managerPort uint16
	conn        net.Conn
}

//NewNode creates a new Node
func NewNode() {
	n := new(Node)
	n.managerIP = "127.0.0.1"
	n.managerPort = 12345
	go n.mainLoop()
}

func (n *Node) mainLoop() {
	n.connectToManager()

}
func (n *Node) sample() int {
	return rand.Intn(100)
}

func (n *Node) connectToManager() {

	for {
		conn, err := net.Dial("tcp", n.managerIP+":"+strconv.Itoa(int(n.managerPort)))
		if err != nil {
			time.Sleep(time.Second)
		} else {
			n.conn = conn
			break
		}
	}

}
func (n *Node) disconnectFromManager() {
	n.conn.Close()
}

func (n *Node) getGPSLoc() (float64, float64) {
	return rand.Float64(), rand.Float64()
}

func CreateNodes(number uint) {
	var i uint
	for i = 0; i < number; i++ {
		NewNode()
	}
}
