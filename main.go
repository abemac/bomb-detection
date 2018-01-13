package main

import (
	"github.com/abemac/bomb-detection/manager"
	"github.com/abemac/bomb-detection/nodesim"
)

func main() {
	nodesim.CreateNodes(10)

	m := manager.NewManager()
	m.Run()
}
