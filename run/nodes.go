package main

import (
	"flag"
	"sync"

	"github.com/abemac/bomb-detection/constants"
	"github.com/abemac/bomb-detection/logger"

	"github.com/abemac/bomb-detection/nodesim"
)

func main() {
	var loglevel = flag.Int("ll", logger.INFO, "set to vary logging output")
	var ip = flag.String("ip", "127.0.0.1", "manager ip")
	var nodeConfigFile = flag.String("nodeConfigFile", "", "JSON configuration for nodes")
	flag.Parse()
	constants.LOG_LEVEL = *loglevel
	// nodesim.CreateNodes(9900, *ip)
	// nodesim.CreateSupernodes(100, *ip)
	nodesim.ExecConfigFile(*nodeConfigFile, *ip)
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()

}
