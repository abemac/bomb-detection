package main

import (
	"github.com/abemac/bomb-detection/manager"
	"github.com/abemac/bomb-detection/nodesim"
)

func main() {
	m := manager.NewManager()
	for i := 0; i < 20; i++ {
		nodesim.NewNode()
	}

	m.Run()
}
