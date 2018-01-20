package main

import (
	"flag"

	"github.com/abemac/bomb-detection/constants"
	"github.com/abemac/bomb-detection/logger"

	"github.com/abemac/bomb-detection/manager"
)

func main() {
	var loglevel = flag.Int("ll", logger.INFO, "set to vary logging output")
	flag.Parse()
	constants.LOG_LEVEL = *loglevel
	m := manager.NewManager()
	m.Run()
}
