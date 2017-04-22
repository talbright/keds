package main

import (
	"github.com/spf13/viper"
	. "github.com/talbright/keds/server"
	_ "github.com/talbright/keds/utils/config"

	"log"
	"strings"
)

func main() {

	log.Printf("config: %v", viper.AllSettings())
	pp := viper.GetStringSlice("plugin_path")
	log.Printf("plugin path: %s", strings.Join(pp, ":"))
	gRPC := NewKedsRPCServer()
	gRPC.Start()

}
