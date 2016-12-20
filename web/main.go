package main

import (
	"net/http"
	"git.target.com/plrsght-go-distrbuted/web/controller"
)

func main() {
	controller.Initialize()

	http.ListenAndServe(":3000", nil)
}
