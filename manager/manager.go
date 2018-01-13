package manager

import (
	"container/list"
	"fmt"
	"net"
)

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
	m.listenForConnections()
}

func (m *Manager) listenForConnections() {
	listener, err := net.Listen("tcp", ":12345")
	if err != nil {
		panic(err.Error())
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err.Error())
		}
		go func() {
			m.connectedNodes.PushBack(conn)
			fmt.Println("Connected to ", conn.RemoteAddr())
		}()
	}
}
