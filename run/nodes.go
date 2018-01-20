package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/abemac/bomb-detection/constants"
	"github.com/abemac/bomb-detection/logger"

	"github.com/abemac/bomb-detection/nodesim"
)

func main() {
	var loglevel = flag.Int("ll", logger.INFO, "set to vary logging output")
	var ip = flag.String("ip", "127.0.0.1", "manager ip")
	flag.Parse()
	constants.LOG_LEVEL = *loglevel
	handleSIGINT()
	nodesim.CreateNodes(10000, *ip)
	done := make(chan bool)
	<-done
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
