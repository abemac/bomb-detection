package main

import (
	"flag"

	"github.com/abemac/bomb-detection/manager/constants"
	"github.com/abemac/bomb-detection/manager/logger"

	"github.com/abemac/bomb-detection/manager"
	"github.com/abemac/bomb-detection/nodesim"
)

func main() {
	var loglevel = flag.Int("ll", logger.OFF, "set to vary logging output")
	flag.Parse()
	constants.LOG_LEVEL = *loglevel

	nodesim.CreateNodes(10)
	m := manager.NewManager()
	m.Run()
}
