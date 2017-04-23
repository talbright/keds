package main

import (
	"github.com/spf13/viper"
	. "github.com/talbright/keds/server"
	_ "github.com/talbright/keds/utils/config"

	"log"
)

func main() {
	log.Printf("config: %v", viper.AllSettings())
	gRPC := NewKedsRPCServer()
	gRPC.Start()
}
