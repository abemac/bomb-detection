package manager

import (
	"bufio"
	"net"

	"github.com/abemac/bomb-detection/constants"
)

type MessageHandler interface {
	handleMessage([]byte) []byte
}

type JSONRequestServer struct {
	msgHandler MessageHandler
}

func NewJSONRequestServer(msgHandler MessageHandler) {
	j := new(JSONRequestServer)
	j.msgHandler = msgHandler
}

func (j *JSONRequestServer) Run() {
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
			go j.handleConnection(conn)
		}

	}
}
func (j *JSONRequestServer) handleConnection(conn net.Conn) {
	defer conn.Close()
	for {
		bytes, err := bufio.NewReader(conn).ReadBytes(constants.DelimJSON[0])
		if err != nil {
			panic("eeek")
		}
		if err != nil {
			log.V("Connection to", conn.RemoteAddr(), "closed")
			break //connection closed
		}
		response := j.msgHandler.handleMessage(bytes)
		if response != nil {
			conn.Write(response)
			conn.Write(constants.DelimJSON)
			log.D("Sent: ", response)
		}
	}

}
