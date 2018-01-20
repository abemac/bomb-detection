package main

import (
	"flag"

	"github.com/abemac/bomb-detection/constants"
	"github.com/abemac/bomb-detection/logger"

	"github.com/abemac/bomb-detection/nodesim"
)

func main() {
	var loglevel = flag.Int("ll", logger.INFO, "set to vary logging output")
	var ip = flag.String("ip", "127.0.0.1", "manager ip")
	flag.Parse()
	constants.LOG_LEVEL = *loglevel
	nodesim.CreateNodes(10000, *ip)
	done := make(chan bool)
	<-done
}
