package main

import "net/http"

func main() {
	println("NODE STARTED")
	go http.ListenAndServe(":8080", nil)
	select {} // Stops the node from closing down.
}
