package main

import (
	"flag"

	"github.com/abemac/bomb-detection/constants"
	"github.com/abemac/bomb-detection/logger"

	"github.com/abemac/bomb-detection/manager"
)

func main() {
	var loglevel = flag.Int("ll", logger.INFO, "set to vary logging output")
	var distPath = flag.String("DIST_PATH", "", "Set the location of the directory to serve as http server")
	flag.Parse()
	constants.LOG_LEVEL = *loglevel
	constants.DIST_PATH = *distPath
	m := manager.NewManager()
	m.Run()
}
