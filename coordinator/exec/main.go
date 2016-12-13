package main

import (
	"git.target.com/plrsght-go-distrbuted/coordinator"
	"fmt"
)

func main() {
	ql := coordinator.NewQueueListener()

	go ql.ListenForNewSource()

	fmt.Print("Waiting for messages.....")

	var a string
	fmt.Scanln(&a)
}
