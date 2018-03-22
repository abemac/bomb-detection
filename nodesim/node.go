package nodesim

import (
	"bufio"
	"encoding/json"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/abemac/bomb-detection/constants"
)

//Node represents a node, which is a sensor
type Node struct {
	managerIP   string
	managerPort uint16
	assignedID  uint64
	info        NodeInfo
	muid        int64
}

type NodeInfo struct {
	latitude          float64
	longitude         float64
	batteryPercentage float32
	lock              sync.RWMutex
}

//NewNode creates a new Node
func NewNode(ip string) {
	n := new(Node)
	n.managerIP = ip
	n.managerPort = 12345
	n.assignedID = constants.ID_NOT_ASSIGNED
	n.info.latitude, n.info.longitude = rand.Float64()*180-90, rand.Float64()*180
	n.info.batteryPercentage = rand.Float32()
	go n.mainLoop()
	log.V("New Node created")
}
func CreateNodes(number uint64, ip string) {
	var i uint64
	for i = 0; i < number; i++ {
		NewNode(ip)
	}
	log.I(number, "new nodes created, they are active")
}

func (n *Node) sample() int {
	return rand.Intn(100)
}

func (n *Node) act() {
	for {
		lat := float64(rand.Intn(3) - 1)
		long := float64(rand.Intn(3) - 2)
		n.info.add(lat, long)
		time.Sleep(time.Millisecond * 500)
	}
}

func (n *Node) mainLoop() {
	go n.act()
	n.communicateWithManager()
}
func (n *Node) communicateWithManager() {
	var conn net.Conn
	var connected bool
	for {
		if !connected {
			conn = n.connectToManager()
			connected = true
		}
		resp, err := n.sendInfoAndGetResponse(conn)
		if err != nil {
			connected = false
			conn.Close()
			log.V("Connection to", conn.RemoteAddr(), "closed")
			time.Sleep(time.Second * time.Duration(rand.Int31n(5)+5))
			continue
		}

		err = n.handleResponse(resp, conn)
		if err != nil {
			connected = false
			conn.Close()
			log.V("Connection to", conn.RemoteAddr(), "closed")
			time.Sleep(time.Second * time.Duration(rand.Int31n(5)+5))
			continue
		}

		time.Sleep(time.Second * time.Duration(resp.NextCheckin))
	}
}
func (n *Node) connectToManager() net.Conn {

	for {
		conn, err := net.Dial("tcp", n.managerIP+":"+strconv.Itoa(int(n.managerPort)))
		if err != nil {
			time.Sleep(time.Second * time.Duration(rand.Int31n(5)+5))
		} else {
			log.V("Connection to", conn.RemoteAddr(), "created")
			return conn
		}
	}

}
func (n *Node) sendInfoAndGetResponse(conn net.Conn) (*constants.ManagerToNodeJSON, error) {
	message := new(constants.NodeToManagerJSON)
	message.Latitude, message.Longitude = n.getGPSLoc()
	message.ID = n.assignedID
	message.SampleValid = false
	message.ManagerUID = n.muid
	message.BatteryPercent = n.info.batteryPercentage

	messageJSON, err := json.Marshal(message)
	if err != nil {
		panic(err.Error())
	}
	conn.Write(messageJSON)
	conn.Write(constants.DelimJSON)
	log.D("Sent: ", *message)
	return n.recvFrom(conn)
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

func (n *Node) handleResponse(msg *constants.ManagerToNodeJSON, conn net.Conn) error {
	n.assignedID = msg.AssignedID
	n.muid = msg.ManagerUID
	if msg.PerformSample {
		err := n.sendSample(conn)
		if err != nil {
			return err
		}
	}
	return nil

}
func (n *Node) sendSample(conn net.Conn) error {
	message := new(constants.NodeToManagerJSON)
	message.ID = n.assignedID
	message.SampleValue = n.sample()
	message.SampleValid = true
	message.ManagerUID = n.muid
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return err
	}
	_, err = conn.Write(messageJSON)
	if err != nil {
		return err
	}
	_, err = conn.Write(constants.DelimJSON)
	if err != nil {
		return err
	}
	log.D("Sent: ", *message)
	return nil
}
func (i *NodeInfo) add(latitude float64, longitude float64) {
	i.lock.RLock()
	lat, long := i.latitude, i.longitude
	i.lock.RUnlock()

	lat += latitude
	long += longitude
	if lat > 90 {
		lat -= 180
	}
	if lat < -90 {
		lat += 180
	}
	if long > 180 {
		long -= 360
	}
	if long < -180 {
		long += 360
	}
	i.lock.Lock()
	i.latitude = lat
	i.longitude = long
	i.lock.Unlock()
}

func (n *Node) getGPSLoc() (float64, float64) {
	n.info.lock.RLock()
	defer n.info.lock.RUnlock()
	return n.info.latitude, n.info.longitude
}
