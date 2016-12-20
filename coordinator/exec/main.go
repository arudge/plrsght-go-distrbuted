package main

import (
	"git.target.com/plrsght-go-distrbuted/coordinator"
	"fmt"
)

var dc *coordinator.DatabaseConsumer
var wc *coordinator.WebappConsumer

func main() {
	ea := coordinator.NewEventAggregator()
	dc = coordinator.NewDatabaseConsumer(ea)
	wc = coordinator.NewWebappConsumer(ea)
	ql := coordinator.NewQueueListener(ea)

	go ql.ListenForNewSource()


	fmt.Println("Waiting for messages.....")

	var a string
	fmt.Scanln(&a)
}
