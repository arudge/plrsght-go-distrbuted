package main

import (
	"git.target.com/plrsght-go-distrbuted/coordinator"
	"fmt"
)

var dc *coordinator.DatabaseConsumer

func main() {
	ea := coordinator.NewEventAggregator()
	dc = coordinator.NewDatabaseConsumer(ea)
	ql := coordinator.NewQueueListener(ea)

	go ql.ListenForNewSource()


	fmt.Print("Waiting for messages.....")

	var a string
	fmt.Scanln(&a)
}
