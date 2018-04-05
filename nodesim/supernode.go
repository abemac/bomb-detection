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

//Supernode represents a node, which is a sensor
type Supernode struct {
	managerIP   string
	managerPort uint16
	assignedID  uint64
	muid        int64
	lat         float64
	long        float64
}

//NewNode creates a new Node
func NewSupernode(ip string, lat float64, long float64) {
	n := new(Supernode)
	n.managerIP = ip
	n.managerPort = 12345
	n.assignedID = constants.ID_NOT_ASSIGNED
	// n.lat, n.long = rand.Float64()*180-90, rand.Float64()*360-180
	n.lat, n.long = lat, long
	go n.mainLoop()
	log.V("New Supernode created")
}
func CreateSupernodes(number uint64, ip string, lat float64, long float64) {
	var i uint64
	for i = 0; i < number; i++ {
		NewSupernode(ip, lat, long)
	}
	log.I(number, "new Supernodes created, they are active")
}
func (n *Supernode) sample() int {
	return rand.Intn(100)
}

func (n *Supernode) mainLoop() {
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

func (n *Supernode) connectToManager() net.Conn {

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

func (n *Supernode) sendInfoAndGetResponse(conn net.Conn) (*constants.ManagerToNodeJSON, error) {
	message := new(constants.NodeToManagerJSON)
	message.Latitude, message.Longitude = n.getGPSLoc()
	message.ID = n.assignedID
	message.SampleValid = false
	message.SuperNode = true
	message.ManagerUID = n.muid
	messageJSON, err := json.Marshal(message)
	if err != nil {
		panic(err.Error())
	}
	conn.Write(messageJSON)
	conn.Write(constants.DelimJSON)
	log.D("Sent: ", *message)

	return n.recvFrom(conn)
}
func (n *Supernode) recvFrom(conn net.Conn) (*constants.ManagerToNodeJSON, error) {
	bytes, err := bufio.NewReader(conn).ReadBytes(constants.DelimJSON[0])
	if err != nil {
		return nil, err
	}
	data := new(constants.ManagerToNodeJSON)
	json.Unmarshal(bytes[:len(bytes)-1], data)
	log.D("Received: ", *data)
	return data, nil
}
func (n *Supernode) handleResponse(msg *constants.ManagerToNodeJSON, conn net.Conn) error {
	n.assignedID = msg.AssignedID
	n.muid = msg.ManagerUID
	n.lat = msg.GoToLat
	n.long = msg.GoToLong
	if msg.PerformSample {
		err := n.sendSample(conn)
		if err != nil {
			return err
		}
	}
	return nil

}
func (n *Supernode) sendSample(conn net.Conn) error {
	message := new(constants.NodeToManagerJSON)
	message.ID = n.assignedID
	message.SampleValue = n.sample()
	message.SampleValid = true
	message.SuperNode = true
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

func (n *Supernode) getGPSLoc() (float64, float64) {
	return n.lat, n.long
}
