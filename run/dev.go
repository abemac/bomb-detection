package main

import (
	"flag"

	"github.com/abemac/bomb-detection/constants"
	"github.com/abemac/bomb-detection/logger"

	"github.com/abemac/bomb-detection/manager"
	"github.com/abemac/bomb-detection/nodesim"
)

func main() {
	var loglevel = flag.Int("ll", logger.INFO, "set to vary logging output")
	flag.Parse()
	constants.LOG_LEVEL = *loglevel

	nodesim.CreateNodes(50, "127.0.0.1")
	m := manager.NewManager()
	m.Run()
}
