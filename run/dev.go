package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/abemac/bomb-detection/constants"
	"github.com/abemac/bomb-detection/logger"

	"github.com/abemac/bomb-detection/manager"
	"github.com/abemac/bomb-detection/nodesim"
)

func main() {
	var loglevel = flag.Int("ll", logger.INFO, "set to vary logging output")
	flag.Parse()
	constants.LOG_LEVEL = *loglevel
	handleSIGINT()
	nodesim.CreateNodes(1000)
	m := manager.NewManager()
	m.Run()
}

func handleSIGINT() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			fmt.Println(sig)
			os.Exit(0)
		}
	}()
}
