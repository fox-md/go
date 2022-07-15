package main

import (
	"flag"
	"fmt"
	"log"
	"time"
)

func main() {

	var shutdownDelay int
	currentTime := time.Now()

	flag.IntVar(&shutdownDelay, "timeout", 0, "Shutdown delay. Default: 0 sec")
	flag.Parse()

	log.Print("Delay for ", shutdownDelay, " seconds")
	fmt.Println(currentTime.Format("2006/01/02 15:04:05"), "Delay for", shutdownDelay, "seconds")

	time.Sleep(time.Duration(shutdownDelay) * time.Second)
}
