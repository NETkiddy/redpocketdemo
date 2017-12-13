/*
Package main
The entrance package
*/
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"jcqts/redpocketdemo/service"
)

/*
The entrance func
*/
func main() {
	svc := service.NewService()
	svc.Start()

	// waiting to exit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-c
	log.Printf("Receive ctrl-c")
}
