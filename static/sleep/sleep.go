package main

import (
	"flag"
	"log"
	"time"
)

func main() {

	var shutdownDelay int

	flag.IntVar(&shutdownDelay, "timeout", 0, "Shutdown delay. Default: 0 sec")
	flag.Parse()

	log.Print("Delay for ", shutdownDelay, " seconds")
	time.Sleep(time.Duration(shutdownDelay) * time.Second)
}
