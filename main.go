package main

import (
	. "github.com/talbright/keds/server"
)

func main() {

	gRPC := NewKedsRPCServer()
	gRPC.Start()

	//TODO load plugins...

	// go gRPC.Start()
	// select {}

}
