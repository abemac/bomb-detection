package nodesim

import (
	"math/rand"
	"net"
	"strconv"
)

//Node represents a node, which is a sensor
type Node struct {
	managerIP   string
	managerPort uint16
	conn        net.Conn
}

//NewNode creates a new Node
func NewNode() *Node {
	n := new(Node)
	n.managerIP = "127.0.0.1"
	n.managerPort = 12345
	n.connectToManager()
	return n
}

func (n *Node) sample() int {
	return rand.Intn(100)
}

func (n *Node) connectToManager() {
	conn, err := net.Dial("tcp", n.managerIP+":"+strconv.Itoa(int(n.managerPort)))
	if err != nil {
		panic(err.Error())
	}
	n.conn = conn
}
func (n *Node) disconnectFromManager() {
	n.conn.Close()
}
